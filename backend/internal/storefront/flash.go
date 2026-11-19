package storefront

// flash: bagian "Flash Sale" di beranda + halaman /{lang}/flash-sale (sesi berjalan & berikutnya).
// Harga yang dibayar tetap dihitung server di keranjang/checkout (paket shop, view flash_live).

import (
	"context"
	"net/http"
	"time"
)

// FlashCard: produk di sesi flash (harga flash termurah di antara variannya + progres terjual).
type FlashCard struct {
	Card
	Price, Normal string
	Off, SoldPct  int
	Sold, Cap     int
	SoldOut       bool
}

type FlashSession struct {
	ID, Name         string
	StartsAt, EndsAt time.Time
	Live             bool
	Cards            []FlashCard
}

func flashName(lang, en, km, zh string) string {
	switch {
	case lang == "km" && km != "":
		return km
	case lang == "zh" && zh != "":
		return zh
	}
	return en
}

// FlashSessions: sesi terbit yang sedang berjalan + maks. `upcoming` sesi berikutnya (dalam 7 hari), masing-masing dengan produknya.
func (c *Catalog) FlashSessions(ctx context.Context, lang string, upcoming, perSession int) ([]FlashSession, error) {
	rows, err := c.Pool.Query(ctx, `SELECT id, name_en, name_km, name_zh, starts_at, ends_at, starts_at <= now() FROM flash_sessions
		WHERE published AND cancelled_at IS NULL AND ends_at > now() AND starts_at < now() + interval '7 days'
		  AND EXISTS (SELECT 1 FROM flash_items i WHERE i.session_id = flash_sessions.id AND i.status = 'approved')
		ORDER BY starts_at LIMIT 10`)
	if err != nil {
		return nil, err
	}
	var list []FlashSession
	future := 0
	for rows.Next() {
		var s FlashSession
		var en, km, zh string
		if err := rows.Scan(&s.ID, &en, &km, &zh, &s.StartsAt, &s.EndsAt, &s.Live); err != nil {
			rows.Close()
			return nil, err
		}
		if !s.Live {
			if future >= upcoming {
				continue
			}
			future++
		}
		s.Name = flashName(lang, en, km, zh)
		list = append(list, s)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	out := list[:0]
	for _, s := range list {
		cards, err := c.sessionCards(ctx, lang, s.ID, perSession)
		if err != nil {
			return nil, err
		}
		if len(cards) > 0 {
			s.Cards = cards
			out = append(out, s)
		}
	}
	return out, nil
}

func (c *Catalog) sessionCards(ctx context.Context, lang, sessionID string, limit int) ([]FlashCard, error) {
	rows, err := c.Pool.Query(ctx, `SELECT x.product_id, x.sale, x.normal, x.cap, x.sold FROM (
		  SELECT i.product_id::text, min(i.sale_price) AS sale,
		    (array_agg(i.normal_price ORDER BY i.sale_price))[1] AS normal, sum(i.qty_cap)::int AS cap,
		    sum(coalesce((SELECT sum(oi.qty) FROM order_items oi JOIN orders o ON o.id = oi.order_id WHERE oi.flash_item_id = i.id
		      AND o.status NOT IN ('cancelled','refunded')), 0))::int AS sold
		  FROM flash_items i JOIN product_variants v ON v.id = i.variant_id
		  WHERE i.session_id::text = $1 AND i.status = 'approved' AND v.deleted_at IS NULL AND v.status = 'active'
		  GROUP BY i.product_id) x
		ORDER BY (x.sold >= x.cap), (1 - x.sale / x.normal) DESC, x.product_id LIMIT $2`, sessionID, limit)
	if err != nil {
		return nil, err
	}
	type agg struct {
		price, normal float64
		cap, sold     int
	}
	ids, by := []string{}, map[string]agg{}
	for rows.Next() {
		var id string
		var a agg
		if err := rows.Scan(&id, &a.price, &a.normal, &a.cap, &a.sold); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
		by[id] = a
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	cards, err := c.CardsByIDs(ctx, lang, ids, len(ids))
	if err != nil {
		return nil, err
	}
	out := make([]FlashCard, 0, len(cards))
	for _, k := range cards {
		a := by[k.ID]
		f := FlashCard{Card: k, Price: fmt2(a.price), Normal: fmt2(a.normal), Cap: a.cap, Sold: a.sold}
		if a.normal > 0 {
			f.Off = int((1 - a.price/a.normal) * 100)
		}
		if a.cap > 0 {
			f.SoldPct = a.sold * 100 / a.cap
		}
		f.SoldOut = a.sold >= a.cap || k.Stock <= 0
		if f.SoldOut {
			f.SoldPct = 100
		}
		out = append(out, f)
	}
	return out, nil
}

type flashPageData struct {
	Sessions []FlashSession
	Active   string
}

// flashPage: GET /{lang}/flash-sale — tab per sesi (yang berjalan dulu).
func (h *Handler) flashPage(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
	list, err := h.Cat.FlashSessions(r.Context(), p.Lang, 3, 60)
	if err != nil {
		return "", 0, err
	}
	active := r.URL.Query().Get("s")
	found := false
	for _, s := range list {
		found = found || s.ID == active
	}
	if !found && len(list) > 0 {
		active = list[0].ID
	}
	p.Title, p.Active, p.Description = p.T("flashSale")+" · Phsar", "flash", p.T("flashMeta")
	p.Data = flashPageData{list, active}
	return "flashsale", http.StatusOK, nil
}
