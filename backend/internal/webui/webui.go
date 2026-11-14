// Package webui menyajikan halaman web Seller Center (/seller) dan admin (/admin).
// Halaman ini memanggil JSON API yang sama dengan app Flutter.
package webui

import (
	"embed"
	"encoding/json"
	"net/http"
)

//go:embed *.html *.css *.js
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

var (
	SellerPage = page("seller.html")
	AdminPage  = page("admin.html")
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
func ConfigJS(mapsKey, mapID string) http.HandlerFunc {
	b, _ := json.Marshal(map[string]string{"googleMapsKey": mapsKey, "googleMapId": mapID})
	body := []byte("window.PHSAR_CONFIG = " + string(b) + ";\n")
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write(body)
	}
}
