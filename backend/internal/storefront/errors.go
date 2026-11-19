package storefront

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// Halaman error untuk pembeli (404, 410, 429, 500, 503, …) dengan tampilan & bahasa Phsar.
//
//   - Handler halaman cukup mengembalikan errNotFound atau httpError(kode); wrap yang merender.
//   - Rute di luar storefront (URL acak, method salah, panic) memakai Handler.ErrorPage dari router.
//   - Error 5xx dirender dengan template "errorstatic" yang tidak menyentuh database,
//     jadi tetap tampil walau penyebab error-nya database.

type statusError struct{ code int }

func (e statusError) Error() string { return "http " + strconv.Itoa(e.code) }

// httpError: kembalikan dari handler halaman, mis. return "", 0, httpError(http.StatusGone).
func httpError(code int) error { return statusError{code} }

// errorStatus: kode HTTP dari error handler; ok=false berarti error biasa (dianggap 500).
func errorStatus(err error) (int, bool) {
	var se statusError
	switch {
	case errors.Is(err, errNotFound):
		return http.StatusNotFound, true
	case errors.As(err, &se):
		return se.code, true
	}
	return http.StatusInternalServerError, false
}

// ErrorView: isi halaman error (Page.Data untuk template "errorpage" / "errorstatic").
type ErrorView struct {
	Code               int
	Icon, Title, Text  string
	Search, Retry, Cat bool // tampilkan kotak cari / tombol coba lagi / chip kategori
}

func errorView(lang string, code int, path string) ErrorView {
	v := ErrorView{Code: code, Icon: "search"}
	key := "err" + strconv.Itoa(code)
	switch code {
	case http.StatusNotFound:
		v.Search, v.Cat = true, true
		// Teks sesuai jenis halaman yang dicari.
		rest := path
		if len(rest) > 3 && validLang(strings.Trim(rest[:3], "/")) {
			rest = rest[3:]
		}
		switch {
		case strings.HasPrefix(rest, "/p/"):
			key = "err404Product"
		case strings.HasPrefix(rest, "/s/"):
			key, v.Icon = "err404Store", "store"
		case strings.HasPrefix(rest, "/c/"):
			key = "err404Category"
		case strings.HasPrefix(rest, "/orders/"):
			key, v.Icon, v.Search, v.Cat = "err404Order", "bag", false, false
		}
	case http.StatusGone:
		v.Icon, v.Search, v.Cat = "box", true, true
	case http.StatusForbidden, http.StatusUnauthorized:
		key, v.Icon = "err403", "shield"
	case http.StatusTooManyRequests:
		v.Icon, v.Retry = "shield", true
	case http.StatusServiceUnavailable, http.StatusBadGateway, http.StatusGatewayTimeout:
		key, v.Icon, v.Retry = "err503", "store", true
	default:
		if code >= 500 {
			key, v.Icon, v.Retry = "err500", "box", true
		} else {
			key, v.Icon = "err400", "search"
		}
	}
	title := map[string]string{"err404Product": "err404", "err404Store": "err404", "err404Category": "err404", "err404Order": "err404"}[key]
	if title == "" {
		title = key
	}
	v.Title, v.Text = T(lang, title+"Title"), T(lang, key+"Text")
	return v
}

// fillError: ubah Page menjadi halaman error.
func (h *Handler) fillError(p *Page, r *http.Request, code int) {
	v := errorView(p.Lang, code, r.URL.Path)
	p.Data, p.Title, p.Robots, p.Description = v, v.Title+" · Phsar", "noindex", v.Text
	p.Canonical = ""
}

// langOf: bahasa dari awalan path (/km/…), lalu cookie, lalu default.
func (h *Handler) langOf(r *http.Request) string {
	if seg := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)[0]; validLang(seg) {
		return seg
	}
	if c, err := r.Cookie("phsar_lang"); err == nil && validLang(c.Value) {
		return c.Value
	}
	return h.Default
}

// ErrorPage: halaman error HTML untuk rute di luar storefront (dipanggil router: 404, 405, panic).
func (h *Handler) ErrorPage(w http.ResponseWriter, r *http.Request, code int) {
	lang := h.langOf(r)
	if code >= 500 { // jangan sentuh database: bisa jadi itu penyebabnya
		h.writeStaticError(w, r, lang, code)
		return
	}
	p := h.newPage(r, lang)
	h.fillError(p, r, code)
	h.writeError(w, r, p, code)
}

// writeError: render "errorpage" (layout lengkap); kalau gagal turun ke versi statis.
func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, p *Page, code int) {
	var buf bytes.Buffer
	if err := h.tmpl.ExecuteTemplate(&buf, "errorpage", p); err != nil {
		h.Log.Error("storefront error page", "code", code, "err", err)
		h.writeStaticError(w, r, p.Lang, code)
		return
	}
	writeErrorHeaders(w, p.Lang, code)
	_, _ = w.Write(buf.Bytes())
}

// writeStaticError: halaman mandiri tanpa data database (CSS inline). Pegangan terakhir untuk 5xx.
func (h *Handler) writeStaticError(w http.ResponseWriter, r *http.Request, lang string, code int) {
	p := &Page{Lang: lang, Base: h.BaseURL, Version: h.version}
	h.fillError(p, r, code)
	var buf bytes.Buffer
	if err := h.tmpl.ExecuteTemplate(&buf, "errorstatic", p); err != nil {
		h.Log.Error("storefront static error page", "code", code, "err", err)
		http.Error(w, http.StatusText(code), code)
		return
	}
	writeErrorHeaders(w, lang, code)
	_, _ = w.Write(buf.Bytes())
}

func writeErrorHeaders(w http.ResponseWriter, lang string, code int) {
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Content-Language", htmlLang[lang])
	h.Set("Cache-Control", "no-store")
	h.Set("X-Robots-Tag", "noindex")
	if code == http.StatusServiceUnavailable || code == http.StatusTooManyRequests {
		h.Set("Retry-After", "60")
	}
	w.WriteHeader(code)
}
