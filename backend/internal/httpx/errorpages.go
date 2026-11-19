package httpx

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

// wantsHTML: permintaan navigasi browser (bukan API/fetch) yang layak dijawab halaman HTML.
func wantsHTML(r *http.Request) bool {
	p := r.URL.Path
	if strings.HasPrefix(p, "/v1/") || strings.HasPrefix(p, "/webhooks/") || p == "/healthz" || p == "/readyz" {
		return false
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	return strings.Contains(r.Header.Get("Accept"), "text/html")
}

// recoverer: pengganti middleware.Recoverer. Panic dicatat beserta stack trace; browser mendapat
// halaman error 500 Phsar, API mendapat JSON.
func recoverer(log *slog.Logger, page func(http.ResponseWriter, *http.Request, int)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rec := recover()
				if rec == nil {
					return
				}
				if rec == http.ErrAbortHandler { // koneksi diputus klien: biarkan net/http menanganinya
					panic(rec)
				}
				log.Error("panic", "method", r.Method, "path", r.URL.Path, "err", rec, "stack", string(debug.Stack()))
				if r.Header.Get("Connection") == "Upgrade" {
					return
				}
				if wantsHTML(r) && page != nil {
					page(w, r, http.StatusInternalServerError)
					return
				}
				Error(w, http.StatusInternalServerError, "internal", "internal error")
			}()
			next.ServeHTTP(w, r)
		})
	}
}
