// Package webui menyajikan halaman web Seller Center (/seller) dan admin (/admin).
// Halaman ini memanggil JSON API yang sama dengan app Flutter.
package webui

import (
	"context"
	"embed"
	"encoding/json"
	"net/http"
)

//go:embed *.html *.css *.js *.png
var files embed.FS

func page(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := files.ReadFile(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		h := w.Header()
		h.Set("Content-Type", "text/html; charset=utf-8")
		h.Set("Cache-Control", "no-cache")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "same-origin")
		_, _ = w.Write(b)
	}
}

// ServiceWorker: GET /seller/sw.js, /admin/sw.js, /sw.js — service worker Web Push (cakupan = folder URL-nya).
func ServiceWorker(w http.ResponseWriter, r *http.Request) {
	b, err := files.ReadFile("sw.js")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(b)
}

var (
	SellerPage  = page("seller.html")
	AdminPage   = page("admin.html")
	CourierPage = page("courier.html")
)

// Assets: GET /ui/{common.css,common.js,seller.js,admin.js}
func Assets() http.Handler {
	fs := http.FileServer(http.FS(files))
	return http.StripPrefix("/ui/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		fs.ServeHTTP(w, r)
	}))
}

// ConfigJS: GET /ui/config.js — pengaturan publik untuk halaman web (kunci Google Maps browser).
// Kunci ini memang terlihat di browser; amankan dengan pembatasan HTTP referrer + API di Google Cloud.
// login: cara login yang tersedia saat ini (SMS & Facebook diatur super admin, bisa berubah tanpa restart).
// maps: kunci browser Google Maps + Map ID dari Admin → Integrasi (cadangan .env).
func ConfigJS(maps func(ctx context.Context) (key, mapID string), login func(ctx context.Context) map[string]bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mapsKey, mapID := maps(r.Context())
		cfg := map[string]any{"googleMapsKey": mapsKey, "googleMapId": mapID}
		for k, v := range login(r.Context()) {
			cfg[k] = v
		}
		b, _ := json.Marshal(cfg)
		body := []byte("window.PHSAR_CONFIG = " + string(b) + ";\n")
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write(body)
	}
}
