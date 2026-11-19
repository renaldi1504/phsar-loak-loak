package admin

// Ads: iklan produk bersponsor — tinjau kampanye, jeda, pengaturan (izin ads.review). Logika di paket ads.

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// AdsRoutes: handler admin dari paket ads (supaya admin tidak bergantung langsung pada paket itu).
type AdsRoutes struct{ Handler http.Handler }

func (h *Handler) mountAds(r chi.Router) {
	if h.Ads != nil {
		r.With(h.can("ads.review")).Mount("/ads", h.Ads.Handler)
	}
}
