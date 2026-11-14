package translate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/platform/db"
	"phsar/internal/seller"
)

// Worker memproses antrean translation_jobs. Semua job untuk satu entitas (mis. dua bahasa
// target untuk satu produk) dikerjakan dalam satu panggilan AI.
type Worker struct {
	Pool     *pgxpool.Pool
	T        Translator
	Log      *slog.Logger
	Interval time.Duration
}

const maxAttempts = 5

func (w *Worker) Run(ctx context.Context) {
	if w.T == nil {
		w.Log.Warn("terjemahan AI nonaktif: isi ANTHROPIC_API_KEY (atau GOOGLE_TRANSLATE_API_KEY) di .env; job tetap menunggu di antrean")
		return
	}
	if w.Interval == 0 {
		w.Interval = 3 * time.Second
	}
	// Job "running" yang tertinggal (server mati di tengah jalan) dikembalikan ke antrean.
	_, _ = w.Pool.Exec(ctx, `UPDATE translation_jobs SET status = 'queued' WHERE status = 'running'`)
	w.Log.Info("worker terjemahan berjalan")
	t := time.NewTicker(w.Interval)
	defer t.Stop()
	for {
		for {
			did, err := w.RunOnce(ctx)
			if err != nil {
				w.Log.Error("worker terjemahan", "err", err)
				break
			}
			if !did || ctx.Err() != nil {
				break
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

type job struct {
	id             int64
	table, entity  string
	source, target string
	attempts       int
}

// RunOnce mengambil satu entitas dari antrean dan memprosesnya. false = antrean kosong.
func (w *Worker) RunOnce(ctx context.Context) (bool, error) {
	rows, err := w.Pool.Query(ctx, `
		WITH first AS (
			SELECT entity_table, entity_id FROM translation_jobs
			WHERE status = 'queued' AND (run_after IS NULL OR run_after <= now())
			ORDER BY priority DESC, created_at LIMIT 1 FOR UPDATE SKIP LOCKED)
		UPDATE translation_jobs j SET status = 'running', attempts = j.attempts + 1
		FROM first
		WHERE j.entity_table = first.entity_table AND j.entity_id = first.entity_id
		  AND j.status = 'queued' AND (j.run_after IS NULL OR j.run_after <= now())
		RETURNING j.id, j.entity_table, j.entity_id, j.source_lang::text, j.target_lang::text, j.attempts`)
	if err != nil {
		return false, err
	}
	jobs, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (job, error) {
		var j job
		err := r.Scan(&j.id, &j.table, &j.entity, &j.source, &j.target, &j.attempts)
		return j, err
	})
	if err != nil || len(jobs) == 0 {
		return false, err
	}

	var model string
	switch jobs[0].table {
	case "store_translations":
		model, err = w.store(ctx, jobs)
	case "product_translations":
		model, err = w.product(ctx, jobs)
	default:
		err = errSkip
	}
	ids := make([]int64, len(jobs))
	for i, j := range jobs {
		ids[i] = j.id
	}
	switch {
	case errors.Is(err, errSkip):
		_, err = w.Pool.Exec(ctx, `UPDATE translation_jobs SET status = 'skipped', done_at = now() WHERE id = ANY($1)`, ids)
	case err != nil:
		w.Log.Warn("terjemahan gagal, dicoba lagi nanti", "entity", jobs[0].table+"/"+jobs[0].entity, "attempt", jobs[0].attempts, "err", err)
		_, err = w.Pool.Exec(ctx, `UPDATE translation_jobs SET
			status = CASE WHEN attempts >= $2 THEN 'failed' ELSE 'queued' END,
			error = $3, run_after = now() + make_interval(mins => attempts * attempts)
			WHERE id = ANY($1)`, ids, maxAttempts, truncate(err.Error(), 1000))
	default:
		w.Log.Info("terjemahan selesai", "entity", jobs[0].table+"/"+jobs[0].entity, "model", model)
		_, err = w.Pool.Exec(ctx, `UPDATE translation_jobs SET status = 'done', model = $2, error = NULL, done_at = now() WHERE id = ANY($1)`, ids, model)
	}
	return true, err
}

var errSkip = errors.New("tidak ada yang perlu diterjemahkan")

func targetsOf(jobs []job, source string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, j := range jobs {
		if j.target != source && !seen[j.target] {
			seen[j.target] = true
			out = append(out, j.target)
		}
	}
	return out
}

// ---- Deskripsi toko ----

func (w *Worker) store(ctx context.Context, jobs []job) (string, error) {
	storeID := jobs[0].entity
	var slug, source, about string
	err := w.Pool.QueryRow(ctx, `
		SELECT st.slug, t.lang::text, t.about FROM stores st
		JOIN store_translations t ON t.store_id = st.id AND t.source = 'seller'
		WHERE st.id = $1::uuid AND st.deleted_at IS NULL
		ORDER BY t.updated_at DESC LIMIT 1`, storeID).Scan(&slug, &source, &about)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && strings.TrimSpace(about) == "") {
		return "", errSkip
	}
	if err != nil {
		return "", err
	}
	targets := targetsOf(jobs, source)
	if len(targets) == 0 {
		return "", errSkip
	}
	res, err := w.T.Translate(ctx, Request{Kind: "store description", Source: source, Targets: targets, Fields: map[string]string{"about": about}})
	if err != nil {
		return "", err
	}
	return res.Model, db.WithActor(ctx, w.Pool, "", "system", "", func(tx pgx.Tx) error {
		for _, l := range targets {
			if _, err := tx.Exec(ctx, `
				INSERT INTO store_translations (store_id, lang, about, slug, source) VALUES ($1::uuid, $2::lang_code, $3, $4, 'machine')
				ON CONFLICT (store_id, lang) DO UPDATE SET about = EXCLUDED.about, source = 'machine'
				WHERE store_translations.source <> 'seller'`, storeID, l, res.Texts[l]["about"], slug); err != nil {
				return err
			}
		}
		return nil
	})
}

// ---- Produk: nama, deskripsi, nama opsi & nilai varian ----

func (w *Worker) product(ctx context.Context, jobs []job) (string, error) {
	productID := jobs[0].entity
	var source, name, desc string
	err := w.Pool.QueryRow(ctx, `
		SELECT p.source_lang::text, t.name, t.description FROM products p
		JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE p.id = $1::uuid AND p.deleted_at IS NULL`, productID).Scan(&source, &name, &desc)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errSkip
	}
	if err != nil {
		return "", err
	}
	// Bahasa yang terjemahannya sudah diedit sendiri oleh seller tidak ditimpa.
	targets := []string{}
	for _, l := range targetsOf(jobs, source) {
		var edited bool
		_ = w.Pool.QueryRow(ctx, `SELECT seller_edited FROM product_translations WHERE product_id = $1::uuid AND lang = $2::lang_code`,
			productID, l).Scan(&edited)
		if !edited {
			targets = append(targets, l)
		}
	}
	if len(targets) == 0 {
		return "", errSkip
	}

	fields := map[string]string{"name": name}
	if strings.TrimSpace(desc) != "" {
		fields["description"] = desc
	}
	// Nama opsi ("Warna") dan nilainya ("Merah", "Hitam") ikut diterjemahkan.
	type ref struct{ table, id string }
	refs := map[string]ref{}
	rows, err := w.Pool.Query(ctx, `
		SELECT 'o', o.id::text, o.position::int, 0, coalesce(o.name_i18n->>$2, '') FROM product_options o WHERE o.product_id = $1::uuid
		UNION ALL
		SELECT 'v', v.id::text, o.position::int, v.position::int, coalesce(v.value_i18n->>$2, '')
		FROM product_option_values v JOIN product_options o ON o.id = v.option_id WHERE o.product_id = $1::uuid`, productID, source)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		var kind, id, text string
		var op, vp int
		if err := rows.Scan(&kind, &id, &op, &vp, &text); err != nil {
			rows.Close()
			return "", err
		}
		if strings.TrimSpace(text) == "" {
			continue
		}
		key := fmt.Sprintf("option_%d", op)
		table := "product_options"
		if kind == "v" {
			key = fmt.Sprintf("value_%d_%d", op, vp)
			table = "product_option_values"
		}
		fields[key] = text
		refs[key] = ref{table, id}
	}
	rows.Close()

	res, err := w.T.Translate(ctx, Request{Kind: "product", Source: source, Targets: targets, Fields: fields})
	if err != nil {
		return "", err
	}
	return res.Model, db.WithActor(ctx, w.Pool, "", "system", "", func(tx pgx.Tx) error {
		for _, l := range targets {
			tr := res.Texts[l]
			review := "machine"
			// Hasil sama persis dengan teks asli (dan bukan sekadar merek/kode) → minta dicek manusia.
			if len([]rune(name)) > 3 && strings.TrimSpace(tr["name"]) == strings.TrimSpace(name) {
				review = "needs_review"
			}
			if _, err := tx.Exec(ctx, `
				INSERT INTO product_translations (product_id, lang, name, description, slug, source, review_status, source_lang, model)
				VALUES ($1::uuid, $2::lang_code, $3, $4, $5, 'machine', $6, $7::lang_code, $8)
				ON CONFLICT (product_id, lang) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description,
					slug = EXCLUDED.slug, source = 'machine', review_status = EXCLUDED.review_status,
					source_lang = EXCLUDED.source_lang, model = EXCLUDED.model, reviewed_by = NULL, reviewed_at = NULL
				WHERE NOT product_translations.seller_edited`,
				productID, l, truncate(tr["name"], 160), tr["description"], seller.ProductSlug(tr["name"], name, productID), review, source, res.Model); err != nil {
				return err
			}
			for key, rf := range refs {
				col := "name_i18n"
				if rf.table == "product_option_values" {
					col = "value_i18n"
				}
				if v := strings.TrimSpace(tr[key]); v != "" {
					if _, err := tx.Exec(ctx, `UPDATE `+rf.table+` SET `+col+` = `+col+` || jsonb_build_object($2::text, $3::text) WHERE id = $1::uuid`,
						rf.id, l, truncate(v, 80)); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}
