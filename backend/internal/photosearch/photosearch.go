// Package photosearch: cari produk dengan foto.
//
// Alur: worker mengirim foto produk ke layanan embedder (embedder/server.py, model CLIP) → vektor 512
// angka disimpan di product_image_vectors. Saat pembeli mengunggah foto, foto itu juga dijadikan vektor
// lalu dibandingkan (cosine) dengan semua vektor di memori. Foto pembeli tidak disimpan.
package photosearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/jobs"
	"phsar/internal/storage"
)

var (
	ErrDisabled = errors.New("photo search disabled") // EMBEDDER_URL kosong
	ErrBadImage = errors.New("bad image")
	ErrEmbedder = errors.New("embedder unavailable")
)

// PhotosPerProduct: berapa foto pertama tiap produk yang diindeks.
const PhotosPerProduct = 3

// Embedder mengubah gambar jadi vektor ternormalisasi.
type Embedder interface {
	Embed(ctx context.Context, img []byte) (vec []float32, model string, err error)
}

// HTTPEmbedder memanggil embedder/server.py (POST /embed).
type HTTPEmbedder struct {
	URL    string
	Client *http.Client
}

func (e *HTTPEmbedder) Embed(ctx context.Context, img []byte) ([]float32, string, error) {
	if e == nil || e.URL == "" {
		return nil, "", ErrDisabled
	}
	c := e.Client
	if c == nil {
		c = &http.Client{Timeout: 20 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(e.URL, "/")+"/embed", bytes.NewReader(img))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	res, err := c.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrEmbedder, err)
	}
	defer res.Body.Close()
	var out struct {
		Model  string    `json:"model"`
		Vector []float32 `json:"vector"`
		Error  string    `json:"error"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&out); err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrEmbedder, err)
	}
	switch {
	case res.StatusCode == http.StatusUnprocessableEntity || res.StatusCode == http.StatusRequestEntityTooLarge:
		return nil, "", ErrBadImage
	case res.StatusCode != http.StatusOK || len(out.Vector) == 0:
		return nil, "", fmt.Errorf("%w: status %d %s", ErrEmbedder, res.StatusCode, out.Error)
	}
	return normalize(out.Vector), out.Model, nil
}

func normalize(v []float32) []float32 {
	var s float64
	for _, x := range v {
		s += float64(x) * float64(x)
	}
	n := float32(math.Sqrt(s))
	if n == 0 {
		return v
	}
	out := make([]float32, len(v))
	for i, x := range v {
		out[i] = x / n
	}
	return out
}

// Hit: satu produk hasil pencarian (skor = kemiripan foto terbaiknya, 0..1).
type Hit struct {
	ProductID string
	Score     float32
}

type entry struct {
	product string
	vec     []float32
}

// Service: indeks di memori + worker pengindeks.
type Service struct {
	Pool  *pgxpool.Pool
	Files storage.Storage
	Emb   Embedder
	Log   *slog.Logger

	mu       sync.RWMutex
	entries  []entry
	model    string
	loadedAt time.Time
	dirty    bool
}

func (s *Service) Enabled() bool {
	if h, ok := s.Emb.(*HTTPEmbedder); ok {
		return h != nil && h.URL != ""
	}
	return s.Emb != nil
}

// Search: foto pembeli → produk paling mirip (maks limit). Produk yang tidak tampil difilter pemanggil.
func (s *Service) Search(ctx context.Context, img []byte, limit int) ([]Hit, error) {
	if !s.Enabled() {
		return nil, ErrDisabled
	}
	q, model, err := s.Emb.Embed(ctx, img)
	if err != nil {
		return nil, err
	}
	if err := s.ensureLoaded(ctx, model); err != nil {
		return nil, err
	}
	s.mu.RLock()
	best := map[string]float32{}
	for _, e := range s.entries {
		if len(e.vec) != len(q) {
			continue
		}
		var d float32
		for i, x := range e.vec {
			d += x * q[i]
		}
		if d > best[e.product] {
			best[e.product] = d
		}
	}
	s.mu.RUnlock()
	hits := make([]Hit, 0, len(best))
	for p, sc := range best {
		hits = append(hits, Hit{p, sc})
	}
	sort.Slice(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })
	return cut(hits, limit), nil
}

// cut: buang hasil yang jelas tidak mirip. Skor CLIP foto-ke-foto: barang sejenis ±0.75+, tidak terkait ±0.6.
func cut(h []Hit, limit int) []Hit {
	const minScore, window = 0.62, 0.10
	out := h[:0]
	for _, x := range h {
		if x.Score < minScore || (len(out) > 0 && x.Score < h[0].Score-window) || len(out) >= limit {
			break
		}
		out = append(out, x)
	}
	return out
}

func (s *Service) ensureLoaded(ctx context.Context, model string) error {
	s.mu.RLock()
	fresh := s.model == model && !s.dirty && time.Since(s.loadedAt) < 2*time.Minute
	s.mu.RUnlock()
	if fresh {
		return nil
	}
	rows, err := s.Pool.Query(ctx, `SELECT v.product_id::text, v.vec FROM product_image_vectors v
		JOIN products p ON p.id = v.product_id
		JOIN product_media m ON m.id = v.media_id
		WHERE v.model = $1 AND v.vec IS NOT NULL AND p.deleted_at IS NULL AND m.deleted_at IS NULL`, model)
	if err != nil {
		return err
	}
	defer rows.Close()
	var es []entry
	for rows.Next() {
		var e entry
		if err := rows.Scan(&e.product, &e.vec); err != nil {
			return err
		}
		es = append(es, e)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	s.entries, s.model, s.loadedAt, s.dirty = es, model, time.Now(), false
	s.mu.Unlock()
	return nil
}

// IndexPending memproses foto produk yang belum punya vektor. Mengembalikan jumlah foto yang diproses.
func (s *Service) IndexPending(ctx context.Context, batch int) (int, error) {
	if !s.Enabled() {
		return 0, ErrDisabled
	}
	rows, err := s.Pool.Query(ctx, `SELECT m.id::text, m.product_id::text, m.storage_key FROM (
			SELECT m.*, row_number() OVER (PARTITION BY m.product_id ORDER BY m.sort_order, m.id) AS rn
			FROM product_media m WHERE m.kind = 'photo' AND m.deleted_at IS NULL) m
		JOIN products p ON p.id = m.product_id AND p.deleted_at IS NULL
		WHERE m.rn <= $1 AND NOT EXISTS (SELECT 1 FROM product_image_vectors v WHERE v.media_id = m.id)
		ORDER BY m.product_id, m.rn LIMIT $2`, PhotosPerProduct, batch)
	if err != nil {
		return 0, err
	}
	type job struct{ media, product, key string }
	var jobs []job
	for rows.Next() {
		var j job
		if err := rows.Scan(&j.media, &j.product, &j.key); err != nil {
			rows.Close()
			return 0, err
		}
		jobs = append(jobs, j)
	}
	rows.Close()
	done := 0
	for _, j := range jobs {
		img, err := s.read(ctx, j.key)
		var vec []float32
		model := ""
		if err == nil {
			vec, model, err = s.Emb.Embed(ctx, img)
		}
		if errors.Is(err, ErrEmbedder) || errors.Is(err, ErrDisabled) || ctx.Err() != nil {
			return done, err // embedder mati: coba lagi nanti, jangan tandai gagal
		}
		var errText *string
		if err != nil {
			t := err.Error()
			errText, vec, model = &t, nil, "error"
		}
		if _, err := s.Pool.Exec(ctx, `INSERT INTO product_image_vectors (media_id, product_id, model, vec, error)
			VALUES ($1, $2, $3, $4, $5) ON CONFLICT (media_id) DO NOTHING`, j.media, j.product, model, vec, errText); err != nil {
			return done, err
		}
		done++
	}
	if done > 0 {
		s.mu.Lock()
		s.dirty = true
		s.mu.Unlock()
	}
	return done, nil
}

func (s *Service) read(ctx context.Context, key string) ([]byte, error) {
	r, err := s.Files.Open(ctx, key)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(io.LimitReader(r, 15<<20))
}

func (s *Service) indexBatch(ctx context.Context) (int, error) { return s.IndexPending(ctx, 40) } // jobs:

// Worker: indeks foto baru setiap 20 detik (diam saja kalau EMBEDDER_URL kosong).
func (s *Service) Worker(ctx context.Context) {
	if !s.Enabled() {
		return
	}
	s.Log.Info("cari-dengan-foto aktif")
	t := time.NewTicker(20 * time.Second)
	defer t.Stop()
	warned := false
	for {
		n, err := jobs.Track(ctx, s.Pool, "photo.index", s.indexBatch) // jobs:
		switch {
		case err != nil && !warned:
			s.Log.Warn("indeks foto gagal (dicoba lagi nanti)", "err", err)
			warned = true
		case err == nil:
			warned = false
			if n > 0 {
				s.Log.Info("indeks foto", "foto", n)
				continue // masih ada antrean: lanjut tanpa menunggu
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}
