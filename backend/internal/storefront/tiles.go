package storefront

// Kotak kategori di beranda (ikon kategori / sampul admin) dan feed "Tonton" per kategori
// (video layar penuh, geser ke atas).

import (
	"context"
	"sync"
	"time"
)

// Tile: kotak kategori di beranda — ikon kategori, atau sampul kalau admin memasangnya (tanpa foto/video produk).
type Tile struct {
	Category
	Cover string // sampul pilihan admin (kosong = ikon kategori)
}

// TileTTL: lama kotak kategori disimpan di memori (variabel supaya tes bisa mematikannya).
var TileTTL = 60 * time.Second

var tileCache struct {
	sync.Mutex
	at   map[string]time.Time
	data map[string][]Tile
}

// CategoryTiles: kotak kategori untuk beranda, urutan sesuai admin; kategori kosong tetap tampil (seller baru akan mengisi).
func (c *Catalog) CategoryTiles(ctx context.Context, lang string, cats []Category) ([]Tile, error) {
	tileCache.Lock()
	if tileCache.data == nil {
		tileCache.at, tileCache.data = map[string]time.Time{}, map[string][]Tile{}
	}
	if t, ok := tileCache.data[lang]; ok && time.Since(tileCache.at[lang]) < TileTTL && len(t) == len(cats) {
		tileCache.Unlock()
		return t, nil
	}
	tileCache.Unlock()

	ids := make([]string, len(cats))
	for i, k := range cats {
		ids[i] = k.ID
	}
	covers := map[string]string{}
	rows, err := c.Pool.Query(ctx, `SELECT id, cover_key FROM categories WHERE id = ANY($1::text[]) AND coalesce(cover_key, '') <> ''`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, key string
		if err := rows.Scan(&id, &key); err != nil {
			return nil, err
		}
		covers[id] = c.url(key)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	out := make([]Tile, 0, len(cats))
	for _, k := range cats {
		out = append(out, Tile{Category: k, Cover: covers[k.ID]})
	}
	tileCache.Lock()
	tileCache.data[lang], tileCache.at[lang] = out, time.Now()
	tileCache.Unlock()
	return out, nil
}

// WatchItem: satu produk di feed Tonton — video kalau ada, kalau tidak "reel" foto (geser otomatis + zoom pelan).
type WatchItem struct {
	Card
	Video  string
	Photos []string // dipakai kalau Video kosong (maks. 5 foto)
}

// WatchFeed: semua produk di kategori (atau subkategori) itu, terbaru dulu. Produk ber-video diputar videonya;
// produk tanpa video tampil sebagai reel foto supaya feed tetap penuh.
func (c *Catalog) WatchFeed(ctx context.Context, lang, categoryID string, limit int) ([]WatchItem, error) {
	cards, err := c.cards(ctx, lang, inCategory("$2"), "p.published_at DESC NULLS LAST, p.id", limit, categoryID)
	if err != nil || len(cards) == 0 {
		return []WatchItem{}, err
	}
	ids := make([]string, len(cards))
	for i, k := range cards {
		ids[i] = k.ID
	}
	// Video: hasil kompres (video_url) kalau sudah ada, kalau belum file asli (storage_key) — sama seperti halaman produk.
	rows, err := c.Pool.Query(ctx, `SELECT DISTINCT ON (product_id) product_id::text, coalesce(nullif(video_url, ''), storage_key) FROM product_media
		WHERE product_id = ANY($1::uuid[]) AND kind = 'video' AND deleted_at IS NULL
		ORDER BY product_id, sort_order`, ids)
	if err != nil {
		return nil, err
	}
	urls := map[string]string{}
	for rows.Next() {
		var id, u string
		if err := rows.Scan(&id, &u); err != nil {
			rows.Close()
			return nil, err
		}
		urls[id] = c.url(u)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var noVideo []string
	for _, k := range cards {
		if urls[k.ID] == "" {
			noVideo = append(noVideo, k.ID)
		}
	}
	photos := map[string][]string{}
	if len(noVideo) > 0 {
		prow, err := c.Pool.Query(ctx, `SELECT product_id::text, coalesce(image_url, storage_key) FROM (
			SELECT product_id, image_url, storage_key, row_number() OVER (PARTITION BY product_id ORDER BY sort_order) AS n
			FROM product_media WHERE product_id = ANY($1::uuid[]) AND kind = 'photo' AND deleted_at IS NULL) m
			WHERE n <= 5 ORDER BY product_id, n`, noVideo)
		if err != nil {
			return nil, err
		}
		for prow.Next() {
			var id, u string
			if err := prow.Scan(&id, &u); err != nil {
				prow.Close()
				return nil, err
			}
			photos[id] = append(photos[id], c.url(u))
		}
		prow.Close()
		if err := prow.Err(); err != nil {
			return nil, err
		}
	}
	out := make([]WatchItem, 0, len(cards))
	for _, k := range cards {
		switch {
		case urls[k.ID] != "":
			out = append(out, WatchItem{Card: k, Video: urls[k.ID]})
		case len(photos[k.ID]) > 0:
			out = append(out, WatchItem{Card: k, Photos: photos[k.ID]})
		}
	}
	return out, nil
}

// SearchCards: produk untuk kategori musiman berbasis kata pencarian.
func (c *Catalog) SearchCards(ctx context.Context, lang, q string, limit int) ([]Card, error) {
	return c.Search(ctx, lang, q, limit, 0)
}
