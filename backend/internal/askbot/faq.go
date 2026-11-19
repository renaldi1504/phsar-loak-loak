package askbot

// FAQ toko: tanya-jawab umum (pengiriman, jam buka, cara ukur, …) per toko, opsional per produk.
// Dipakai bot sebagai fakta, dan FAQ per produk tampil di halaman produk storefront.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
)

type FAQ struct {
	ID          string    `json:"id"`
	ProductID   *string   `json:"product_id"`
	ProductName string    `json:"product_name,omitempty"`
	Lang        string    `json:"lang"`
	Question    string    `json:"question"`
	Answer      string    `json:"answer"`
	Position    int       `json:"position"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FAQInput struct {
	Question  string  `json:"question"`
	Answer    string  `json:"answer"`
	Lang      string  `json:"lang,omitempty"` // kosong = dideteksi dari pertanyaan
	ProductID *string `json:"product_id"`     // null / kosong = seluruh toko
	Position  int     `json:"position,omitempty"`
}

func (in *FAQInput) normalize() error {
	in.Question, in.Answer = strings.TrimSpace(in.Question), strings.TrimSpace(in.Answer)
	if in.ProductID != nil && strings.TrimSpace(*in.ProductID) == "" {
		in.ProductID = nil
	}
	switch {
	case in.Question == "":
		return FieldError{"question", "required"}
	case utf8.RuneCountInString(in.Question) > MaxFAQQuestion:
		return FieldError{"question", fmt.Sprintf("at most %d characters", MaxFAQQuestion)}
	case in.Answer == "":
		return FieldError{"answer", "required"}
	case utf8.RuneCountInString(in.Answer) > MaxFAQAnswer:
		return FieldError{"answer", fmt.Sprintf("at most %d characters", MaxFAQAnswer)}
	case in.ProductID != nil && !isUUID(*in.ProductID):
		return FieldError{"product_id", "invalid id"}
	case in.Position < 0 || in.Position > 999:
		return FieldError{"position", "must be between 0 and 999"}
	}
	switch in.Lang {
	case "":
		in.Lang = DetectLang(in.Question + " " + in.Answer)
	case "en", "km", "zh":
	default:
		return FieldError{"lang", "must be en, km or zh"}
	}
	return nil
}

// ListFAQs: semua FAQ toko; productID = "" → semua, "store" → hanya umum, uuid → produk itu + umum.
func (s *Service) ListFAQs(ctx context.Context, storeID, lang, productID string) ([]FAQ, error) {
	where, args := "", []any{storeID, lang}
	switch {
	case productID == "store":
		where = " AND f.product_id IS NULL"
	case productID != "":
		if !isUUID(productID) {
			return nil, FieldError{"product_id", "invalid id"}
		}
		where = " AND (f.product_id IS NULL OR f.product_id = $3)"
		args = append(args, productID)
	}
	rows, err := s.Pool.Query(ctx, `SELECT f.id, f.product_id, coalesce(t.name, o.name, ''), f.lang::text, f.question, f.answer, f.position, f.updated_at
		FROM store_faqs f
		LEFT JOIN products p ON p.id = f.product_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = $2::lang_code
		LEFT JOIN product_translations o ON o.product_id = p.id AND o.lang = p.source_lang
		WHERE f.store_id = $1`+where+` ORDER BY f.product_id IS NOT NULL, f.position, f.created_at`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []FAQ{}
	for rows.Next() {
		var f FAQ
		if err := rows.Scan(&f.ID, &f.ProductID, &f.ProductName, &f.Lang, &f.Question, &f.Answer, &f.Position, &f.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (s *Service) checkProduct(ctx context.Context, tx pgx.Tx, storeID string, pid *string) error {
	if pid == nil {
		return nil
	}
	var ok bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM products WHERE id = $1 AND store_id = $2 AND deleted_at IS NULL)`, *pid, storeID).Scan(&ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func (s *Service) CreateFAQ(ctx context.Context, storeID, userID string, in FAQInput) (FAQ, error) {
	if err := in.normalize(); err != nil {
		return FAQ{}, err
	}
	var id string
	err := db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		if err := s.checkProduct(ctx, tx, storeID, in.ProductID); err != nil {
			return err
		}
		var n int
		// Kunci baris toko supaya batas jumlah tidak terlewati oleh dua permintaan bersamaan.
		if _, err := tx.Exec(ctx, `SELECT 1 FROM stores WHERE id = $1 FOR UPDATE`, storeID); err != nil {
			return err
		}
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM store_faqs WHERE store_id = $1`, storeID).Scan(&n); err != nil {
			return err
		}
		if n >= MaxFAQs {
			return FieldError{"question", fmt.Sprintf("at most %d questions per store", MaxFAQs)}
		}
		return tx.QueryRow(ctx, `INSERT INTO store_faqs (store_id, product_id, lang, question, answer, position, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, storeID, in.ProductID, in.Lang, in.Question, in.Answer, in.Position, nullable(userID)).Scan(&id)
	})
	if err != nil {
		return FAQ{}, err
	}
	return s.getFAQ(ctx, storeID, id)
}

func (s *Service) UpdateFAQ(ctx context.Context, storeID, userID, id string, in FAQInput) (FAQ, error) {
	if !isUUID(id) {
		return FAQ{}, ErrNotFound
	}
	if err := in.normalize(); err != nil {
		return FAQ{}, err
	}
	err := db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		if err := s.checkProduct(ctx, tx, storeID, in.ProductID); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, `UPDATE store_faqs SET product_id = $3, lang = $4, question = $5, answer = $6, position = $7
			WHERE id = $1 AND store_id = $2`, id, storeID, in.ProductID, in.Lang, in.Question, in.Answer, in.Position)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	if err != nil {
		return FAQ{}, err
	}
	return s.getFAQ(ctx, storeID, id)
}

func (s *Service) DeleteFAQ(ctx context.Context, storeID, userID, id string) error {
	if !isUUID(id) {
		return ErrNotFound
	}
	return db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `DELETE FROM store_faqs WHERE id = $1 AND store_id = $2`, id, storeID)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
}

func (s *Service) getFAQ(ctx context.Context, storeID, id string) (FAQ, error) {
	var f FAQ
	err := s.Pool.QueryRow(ctx, `SELECT f.id, f.product_id, coalesce(t.name, ''), f.lang::text, f.question, f.answer, f.position, f.updated_at
		FROM store_faqs f
		LEFT JOIN products p ON p.id = f.product_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE f.id = $1 AND f.store_id = $2`, id, storeID).
		Scan(&f.ID, &f.ProductID, &f.ProductName, &f.Lang, &f.Question, &f.Answer, &f.Position, &f.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return f, ErrNotFound
	}
	return f, err
}

// QA: satu tanya-jawab untuk storefront.
type QA struct{ Question, Answer, Lang string }

// ProductFAQ: tanya-jawab khusus produk untuk halaman produk (bahasa halaman didahulukan).
func (s *Service) ProductFAQ(ctx context.Context, lang, productID string) ([]QA, error) {
	rows, err := s.Pool.Query(ctx, `SELECT question, answer, lang::text FROM store_faqs
		WHERE product_id = $1 ORDER BY (lang = $2::lang_code) DESC, position, created_at LIMIT 20`, productID, lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []QA{}
	for rows.Next() {
		var q QA
		if err := rows.Scan(&q.Question, &q.Answer, &q.Lang); err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

// ---- Laporan admin ----

type StoreUsage struct {
	StoreID string `json:"store_id"`
	Name    string `json:"name"`
	Count   int    `json:"count"`
}

type Usage struct {
	Day       int          `json:"day"`
	Week      int          `json:"week"`
	ClaudeWk  int          `json:"claude_week"`
	RulesWk   int          `json:"rules_week"`
	ErrorsWk  int          `json:"errors_week"`
	TokensIn  int          `json:"tokens_in_week"`
	TokensOut int          `json:"tokens_out_week"`
	FAQs      int          `json:"faqs"`
	Stores    int          `json:"stores_with_faqs"`
	Top       []StoreUsage `json:"top_stores"`
}

func (s *Service) Usage(ctx context.Context) (Usage, error) {
	u := Usage{Top: []StoreUsage{}}
	err := s.Pool.QueryRow(ctx, `SELECT
		count(*) FILTER (WHERE created_at > now() - interval '1 day'), count(*),
		count(*) FILTER (WHERE mode = 'claude'), count(*) FILTER (WHERE mode = 'rules'), count(*) FILTER (WHERE error IS NOT NULL),
		coalesce(sum(input_tokens), 0), coalesce(sum(output_tokens), 0)
		FROM askbot_usage WHERE created_at > now() - interval '7 days'`).
		Scan(&u.Day, &u.Week, &u.ClaudeWk, &u.RulesWk, &u.ErrorsWk, &u.TokensIn, &u.TokensOut)
	if err != nil {
		return u, err
	}
	if err := s.Pool.QueryRow(ctx, `SELECT count(*), count(DISTINCT store_id) FROM store_faqs`).Scan(&u.FAQs, &u.Stores); err != nil {
		return u, err
	}
	rows, err := s.Pool.Query(ctx, `SELECT u.store_id, s.name, count(*) FROM askbot_usage u JOIN stores s ON s.id = u.store_id
		WHERE u.created_at > now() - interval '7 days' GROUP BY u.store_id, s.name ORDER BY count(*) DESC, s.name LIMIT 5`)
	if err != nil {
		return u, err
	}
	defer rows.Close()
	for rows.Next() {
		var x StoreUsage
		if err := rows.Scan(&x.StoreID, &x.Name, &x.Count); err != nil {
			return u, err
		}
		u.Top = append(u.Top, x)
	}
	return u, rows.Err()
}
