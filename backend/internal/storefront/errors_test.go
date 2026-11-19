package storefront

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testHandler() *Handler {
	return New(&Catalog{}, "https://phsar.test", "km", slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestErrorViewContext(t *testing.T) {
	cases := map[string]string{
		"/en/p/abc":    "err404ProductText",
		"/en/s/x":      "err404StoreText",
		"/km/c/shoes":  "err404CategoryText",
		"/zh/orders/1": "err404OrderText",
		"/en/whatever": "err404Text",
	}
	for path, key := range cases {
		lang := path[1:3]
		if got, want := errorView(lang, 404, path).Text, T(lang, key); got != want {
			t.Errorf("%s: text %q, want %q", path, got, want)
		}
	}
	if v := errorView("en", 500, "/en/"); !v.Retry || v.Search {
		t.Errorf("500 view: %+v", v)
	}
}

func TestErrorStatusMapping(t *testing.T) {
	if c, ok := errorStatus(errNotFound); c != 404 || !ok {
		t.Fatal("errNotFound → 404")
	}
	if c, ok := errorStatus(httpError(410)); c != 410 || !ok {
		t.Fatal("httpError(410) → 410")
	}
	if c, ok := errorStatus(io.EOF); c != 500 || ok {
		t.Fatal("other → 500")
	}
}

// Halaman 404 lengkap (layout) dan 500 statis harus bisa dirender di ketiga bahasa.
func TestErrorPagesRender(t *testing.T) {
	h := testHandler()
	for _, lang := range langs {
		p := &Page{Lang: lang, Base: h.BaseURL, Version: "t", Alternates: map[string]string{}, LangLinks: map[string]string{},
			Categories: []Category{{ID: "1", Name: "Fashion", Slug: "fashion", Icon: "shirt"}}}
		r := httptest.NewRequest("GET", "/"+lang+"/p/nope", nil)
		h.fillError(p, r, 404)
		var buf bytes.Buffer
		if err := h.tmpl.ExecuteTemplate(&buf, "errorpage", p); err != nil {
			t.Fatalf("%s errorpage: %v", lang, err)
		}
		out := buf.String()
		for _, want := range []string{T(lang, "err404ProductText"), "err-search", "/" + lang + "/c/fashion", `noindex`} {
			if !strings.Contains(out, want) {
				t.Errorf("%s 404 missing %q", lang, want)
			}
		}

		w := httptest.NewRecorder()
		h.writeStaticError(w, httptest.NewRequest("GET", "/"+lang+"/", nil), lang, 500)
		if w.Code != 500 || !strings.Contains(w.Body.String(), T(lang, "err500Title")) || w.Header().Get("Cache-Control") != "no-store" {
			t.Errorf("%s static 500: code %d", lang, w.Code)
		}
	}
	// 5xx lewat ErrorPage tidak menyentuh database (Catalog kosong → akan panic kalau disentuh).
	w := httptest.NewRecorder()
	h.ErrorPage(w, httptest.NewRequest("GET", "/en/x", nil), http.StatusServiceUnavailable)
	if w.Code != 503 || w.Header().Get("Retry-After") == "" {
		t.Errorf("503: code %d", w.Code)
	}
}

func TestTileCols(t *testing.T) {
	for _, c := range []struct{ n, max, want int }{{5, 12, 5}, {12, 12, 12}, {12, 6, 6}, {13, 12, 7}, {7, 6, 4}, {0, 6, 1}} {
		if got := tileCols(c.n, c.max); got != c.want {
			t.Errorf("tileCols(%d, %d) = %d, harusnya %d", c.n, c.max, got, c.want)
		}
	}
}
