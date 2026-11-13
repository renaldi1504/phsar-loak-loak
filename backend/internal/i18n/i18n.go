// Package i18n menentukan bahasa request: en (default), km, zh.
package i18n

import (
	"context"
	"net/http"
	"strings"
)

type Lang string

const (
	EN Lang = "en"
	KM Lang = "km"
	ZH Lang = "zh"
)

var Supported = []Lang{EN, KM, ZH}

func Parse(s string) (Lang, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch {
	case strings.HasPrefix(s, "km"):
		return KM, true
	case strings.HasPrefix(s, "zh"):
		return ZH, true
	case strings.HasPrefix(s, "en"):
		return EN, true
	}
	return "", false
}

type ctxKey struct{}

// Middleware: ?lang= → header X-Lang (app) → Accept-Language → default.
func Middleware(def Lang) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := def
			if l, ok := Parse(r.URL.Query().Get("lang")); ok {
				lang = l
			} else if l, ok := Parse(r.Header.Get("X-Lang")); ok {
				lang = l
			} else {
				for _, part := range strings.Split(r.Header.Get("Accept-Language"), ",") {
					if l, ok := Parse(strings.SplitN(part, ";", 2)[0]); ok {
						lang = l
						break
					}
				}
			}
			w.Header().Set("Content-Language", string(lang))
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, lang)))
		})
	}
}

func From(ctx context.Context) Lang {
	if l, ok := ctx.Value(ctxKey{}).(Lang); ok {
		return l
	}
	return EN
}

// Pick memilih teks dari map {"en":…,"km":…,"zh":…} dengan cadangan ke en.
func Pick(m map[string]string, l Lang) string {
	if v := m[string(l)]; v != "" {
		return v
	}
	return m[string(EN)]
}
