// Package paging: parameter & metadata halaman untuk daftar API (Seller Center, admin).
//
//	GET …?page=2&size=20&q=kata  →  {"…": [...], "page": {"page": 2, "size": 20, "total": 57, "pages": 3}}
//
// Tanpa parameter page, handler boleh tetap memakai perilaku lama (dipakai aplikasi lama) — lihat Params.On.
package paging

import (
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Params struct {
	Page, Size int
	Q          string // kata kunci (sudah di-trim, maks. 80 karakter)
	On         bool   // klien meminta halaman (?page=…)
}

// FromRequest membaca page (≥1), size (1…max, default def) dan q.
func FromRequest(r *http.Request, def, max int) Params {
	v := r.URL.Query()
	p := Params{Page: 1, Size: def, On: v.Has("page")}
	if n, err := strconv.Atoi(v.Get("page")); err == nil && n > 0 {
		p.Page = min(n, 100000)
	}
	if n, err := strconv.Atoi(v.Get("size")); err == nil && n > 0 {
		p.Size = min(n, max)
	}
	q := strings.TrimSpace(v.Get("q"))
	if utf8.RuneCountInString(q) > 80 {
		q = string([]rune(q)[:80])
	}
	p.Q = q
	return p
}

func (p Params) Offset() int { return (p.Page - 1) * p.Size }

type Meta struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

func NewMeta(p Params, total int) Meta {
	pages := (total + p.Size - 1) / p.Size
	if pages < 1 {
		pages = 1
	}
	return Meta{Page: p.Page, Size: p.Size, Total: total, Pages: pages}
}

// Like: pola ILIKE "%q%" dengan karakter khusus di-escape.
func Like(q string) string {
	return "%" + strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(q) + "%"
}
