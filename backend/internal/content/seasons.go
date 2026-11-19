package content

// Kategori musiman ("Musim ini" di beranda): Tahun Baru, Choul Chnam, Pchum Ben, Pantai, dll.
// Tampil paling depan selama tanggalnya berlaku. Isi produk: satu kategori, produk pilihan, atau kata pencarian.

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
)

type Season struct {
	ID         string     `json:"id"`
	Title      Text       `json:"title"`
	Color      string     `json:"color"`  // pink · blue · sand · green
	Source     string     `json:"source"` // category · products · search
	CategoryID string     `json:"category_id"`
	ProductIDs []string   `json:"product_ids"`
	Query      string     `json:"query"`
	StartsAt   *time.Time `json:"starts_at"`
	EndsAt     *time.Time `json:"ends_at"`
	SortOrder  int        `json:"sort_order"`
	Active     bool       `json:"active"`
}

func (x Season) Live(now time.Time) bool { return live(x.Active, x.StartsAt, x.EndsAt, now) }

// LiveSeasons: kategori musiman yang tampil sekarang (urut sort_order).
func (s *Snapshot) LiveSeasons(now time.Time) []Season {
	out := []Season{}
	for _, x := range s.Seasons {
		if x.Live(now) {
			out = append(out, x)
		}
	}
	return out
}

// SeasonByID: kategori musiman yang tampil sekarang, cocok dengan id penuh atau 8 karakter awal (URL pendek).
func (s *Snapshot) SeasonByID(id string, now time.Time) *Season {
	for _, x := range s.Seasons {
		if (x.ID == id || (len(id) == 8 && strings.HasPrefix(x.ID, id))) && x.Live(now) {
			x := x
			return &x
		}
	}
	return nil
}

func (s *Store) loadSeasons(ctx context.Context, snap *Snapshot) error {
	rows, err := s.Pool.Query(ctx, `SELECT id, title, color, source, coalesce(category_id, ''), product_ids::text[], query, starts_at, ends_at, sort_order, active
		FROM season_picks ORDER BY sort_order, created_at`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var x Season
		if err := rows.Scan(&x.ID, &x.Title, &x.Color, &x.Source, &x.CategoryID, &x.ProductIDs, &x.Query, &x.StartsAt, &x.EndsAt, &x.SortOrder, &x.Active); err != nil {
			return err
		}
		if x.ProductIDs == nil {
			x.ProductIDs = []string{}
		}
		snap.Seasons = append(snap.Seasons, x)
	}
	return rows.Err()
}

type SeasonInput struct {
	Title      Text       `json:"title"`
	Color      string     `json:"color"`
	Source     string     `json:"source"`
	CategoryID string     `json:"category_id"`
	ProductIDs []string   `json:"product_ids"`
	Query      string     `json:"query"`
	StartsAt   *time.Time `json:"starts_at"`
	EndsAt     *time.Time `json:"ends_at"`
	SortOrder  int        `json:"sort_order"`
	Active     bool       `json:"active"`
}

func (s *Store) cleanSeason(ctx context.Context, in SeasonInput) (SeasonInput, error) {
	var err error
	if in.Title, err = cleanText("title", in.Title, 40); err != nil {
		return in, err
	}
	if in.Title["en"] == "" && in.Title["km"] == "" && in.Title["zh"] == "" {
		return in, FieldError{"title_en", "required"}
	}
	switch in.Color {
	case "":
		in.Color = "pink"
	case "pink", "blue", "sand", "green":
	default:
		return in, FieldError{"color", "pink, blue, sand or green"}
	}
	in.Query = strings.Join(strings.Fields(in.Query), " ")
	switch in.Source {
	case "category":
		var ok bool
		if err := s.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1 AND deleted_at IS NULL)`, in.CategoryID).Scan(&ok); err != nil {
			return in, err
		}
		if !ok {
			return in, FieldError{"category_id", "choose a category"}
		}
		in.ProductIDs, in.Query = []string{}, ""
	case "products":
		in.CategoryID, in.Query = "", ""
		seen, ids := map[string]bool{}, []string{}
		for _, p := range in.ProductIDs {
			p = strings.ToLower(strings.TrimSpace(p))
			if !uuidRe.MatchString(p) {
				return in, FieldError{"product_ids", "invalid product"}
			}
			if !seen[p] {
				seen[p] = true
				ids = append(ids, p)
			}
		}
		if len(ids) == 0 || len(ids) > 60 {
			return in, FieldError{"product_ids", "pick 1–60 products"}
		}
		var n int
		if err := s.Pool.QueryRow(ctx, `SELECT count(*) FROM products WHERE id = ANY($1::uuid[]) AND deleted_at IS NULL`, ids).Scan(&n); err != nil {
			return in, err
		}
		if n != len(ids) {
			return in, FieldError{"product_ids", "a product was not found"}
		}
		in.ProductIDs = ids
	case "search":
		in.CategoryID, in.ProductIDs = "", []string{}
		if n := utf8.RuneCountInString(in.Query); n < 2 || n > 60 {
			return in, FieldError{"query", "2–60 characters"}
		}
	default:
		return in, FieldError{"source", "category, products or search"}
	}
	if in.StartsAt != nil && in.EndsAt != nil && !in.EndsAt.After(*in.StartsAt) {
		return in, FieldError{"ends_at", "must be after the start"}
	}
	if in.SortOrder < -1000 || in.SortOrder > 1000 {
		return in, FieldError{"sort_order", "between -1000 and 1000"}
	}
	return in, nil
}

// SaveSeason: id kosong = baru.
func (s *Store) SaveSeason(ctx context.Context, actorID, id string, in SeasonInput) (string, error) {
	in, err := s.cleanSeason(ctx, in)
	if err != nil {
		return "", err
	}
	if id != "" && !uuidRe.MatchString(id) {
		return "", ErrNotFound
	}
	err = db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		if id == "" {
			var count int
			if err := tx.QueryRow(ctx, `SELECT count(*) FROM season_picks`).Scan(&count); err != nil {
				return err
			}
			if count >= 30 {
				return FieldError{"seasons", "at most 30 seasonal picks — delete old ones first"}
			}
			return tx.QueryRow(ctx, `INSERT INTO season_picks (title, color, source, category_id, product_ids, query, starts_at, ends_at, sort_order, active, created_by)
				VALUES ($1, $2, $3, nullif($4, ''), $5::uuid[], $6, $7, $8, $9, $10, $11) RETURNING id`, in.Title, in.Color, in.Source, in.CategoryID,
				in.ProductIDs, in.Query, in.StartsAt, in.EndsAt, in.SortOrder, in.Active, actorID).Scan(&id)
		}
		tag, err := tx.Exec(ctx, `UPDATE season_picks SET title = $2, color = $3, source = $4, category_id = nullif($5, ''), product_ids = $6::uuid[],
			query = $7, starts_at = $8, ends_at = $9, sort_order = $10, active = $11 WHERE id = $1`, id, in.Title, in.Color, in.Source, in.CategoryID,
			in.ProductIDs, in.Query, in.StartsAt, in.EndsAt, in.SortOrder, in.Active)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	s.invalidate()
	return id, err
}

func (s *Store) DeleteSeason(ctx context.Context, actorID, id string) error {
	return s.del(ctx, actorID, "season_picks", id)
}

// SetCategoryCover: sampul kategori pilihan admin (key kosong = otomatis foto produk terlaris).
func (s *Store) SetCategoryCover(ctx context.Context, actorID, categoryID, key string) error {
	if key != "" {
		if err := s.checkImage(ctx, "image_key", key); err != nil {
			return err
		}
	}
	err := db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE categories SET cover_key = nullif($2, '') WHERE id = $1 AND deleted_at IS NULL`, categoryID, key)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	return err
}
