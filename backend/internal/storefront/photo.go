package storefront

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"phsar/internal/photosearch"
)

// PhotoSearcher: diisi httpx (photosearch.Service). nil = fitur mati.
type PhotoSearcher interface {
	Enabled() bool
	Search(ctx context.Context, img []byte, limit int) ([]photosearch.Hit, error)
}

const photoMaxBytes = 10 << 20

// photoLimiter: maks 12 pencarian foto per menit per IP (embedder memakan CPU).
type photoLimiter struct {
	mu   sync.Mutex
	hits map[string][]time.Time
}

func (l *photoLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.hits == nil {
		l.hits = map[string][]time.Time{}
	}
	now := time.Now()
	keep := l.hits[ip][:0]
	for _, t := range l.hits[ip] {
		if now.Sub(t) < time.Minute {
			keep = append(keep, t)
		}
	}
	if len(keep) >= 12 {
		l.hits[ip] = keep
		return false
	}
	l.hits[ip] = append(keep, now)
	if len(l.hits) > 5000 { // bersihkan sesekali
		for k, v := range l.hits {
			if len(v) == 0 || now.Sub(v[len(v)-1]) > time.Minute {
				delete(l.hits, k)
			}
		}
	}
	return true
}

func clientIP(r *http.Request) string {
	if v := r.Header.Get("CF-Connecting-IP"); v != "" {
		return v
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// photoSearch: POST /{lang}/photo-search (multipart "image") → potongan HTML kartu produk.
// Header X-Results = jumlah hasil. Foto tidak disimpan.
func (h *Handler) photoSearch(lang string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		fail := func(code int, msg string) { http.Error(w, msg, code) }
		if h.Photo == nil || !h.Photo.Enabled() {
			fail(http.StatusServiceUnavailable, "photo search disabled")
			return
		}
		if !h.photoRate.allow(clientIP(r)) {
			fail(http.StatusTooManyRequests, "too many requests")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, photoMaxBytes+64<<10)
		f, _, err := r.FormFile("image")
		if err != nil {
			var mbe *http.MaxBytesError
			if errors.As(err, &mbe) {
				fail(http.StatusRequestEntityTooLarge, "too large")
				return
			}
			fail(http.StatusUnprocessableEntity, "image required")
			return
		}
		defer f.Close()
		img, err := io.ReadAll(io.LimitReader(f, photoMaxBytes+1))
		if err != nil || len(img) == 0 {
			fail(http.StatusUnprocessableEntity, "image required")
			return
		}
		if len(img) > photoMaxBytes {
			fail(http.StatusRequestEntityTooLarge, "too large")
			return
		}
		ct := http.DetectContentType(img)
		if !strings.HasPrefix(ct, "image/") && !isHEIC(img) {
			fail(http.StatusUnsupportedMediaType, "not an image")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
		defer cancel()
		hits, err := h.Photo.Search(ctx, img, 60)
		switch {
		case errors.Is(err, photosearch.ErrBadImage):
			fail(http.StatusUnprocessableEntity, "bad image")
			return
		case errors.Is(err, photosearch.ErrDisabled), errors.Is(err, photosearch.ErrEmbedder):
			h.Log.Warn("photo search", "err", err)
			fail(http.StatusServiceUnavailable, "photo search unavailable")
			return
		case err != nil:
			h.Log.Error("photo search", "err", err)
			fail(http.StatusInternalServerError, "error")
			return
		}
		ids := make([]string, len(hits))
		for i, x := range hits {
			ids[i] = x.ProductID
		}
		cards, err := h.Cat.CardsByIDs(ctx, lang, ids, 24)
		if err != nil {
			h.Log.Error("photo search cards", "err", err)
			fail(http.StatusInternalServerError, "error")
			return
		}
		p := &Page{Lang: lang, Base: h.BaseURL, Version: h.version, Riel: h.Cat.RielPerUSD(ctx)}
		var buf bytes.Buffer
		if len(cards) > 0 {
			if err := h.tmpl.ExecuteTemplate(&buf, "grid", gridCtx{p, cards}); err != nil {
				h.Log.Error("photo search template", "err", err)
				fail(http.StatusInternalServerError, "error")
				return
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Results", strconv.Itoa(len(cards)))
		_, _ = w.Write(buf.Bytes())
	}
}

// Foto iPhone (HEIC) tidak dikenali DetectContentType.
func isHEIC(b []byte) bool {
	return len(b) > 12 && string(b[4:8]) == "ftyp" && strings.Contains("heic heix mif1 msf1 hevc", string(b[8:12]))
}
