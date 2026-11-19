package storefront

import (
	"fmt"
	"html/template"
	"sort"
	"strconv"
	"strings"
)

// glyphDefs: gradien + mask semua ikon, ditulis sekali di awal <body>.
func glyphDefs() template.HTML {
	names := make([]string, 0, len(glyphShapes))
	for n := range glyphShapes {
		names = append(names, n)
	}
	sort.Strings(names)
	var b strings.Builder
	b.WriteString(`<svg class="glyph-defs" width="0" height="0" aria-hidden="true" focusable="false"><defs>`)
	b.WriteString(`<linearGradient id="pg" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="#ff7ba8"/><stop offset=".56" stop-color="#f50750"/><stop offset="1" stop-color="#c30963"/></linearGradient>`)
	for _, n := range names {
		b.WriteString(`<mask id="gm-` + n + `" x="0" y="0" width="24" height="24" maskUnits="userSpaceOnUse" style="mask-type:luminance"><g fill="white" stroke="none">`)
		b.WriteString(glyphShapes[n])
		b.WriteString(`</g></mask>`)
	}
	b.WriteString(`</defs></svg>`)
	return template.HTML(b.String())
}

// glyph "cart" → ikon gradien pink; glyph "cart" "mono" → warna mengikuti teks (currentColor).
func glyph(name string, variant ...string) template.HTML {
	if _, ok := glyphShapes[name]; !ok {
		name = "bag"
	}
	fill := "url(#pg)"
	cls := "g"
	if len(variant) > 0 && variant[0] == "mono" {
		fill, cls = "currentColor", "g mono"
	}
	return template.HTML(`<svg class="` + cls + `" viewBox="0 0 24 24" aria-hidden="true" focusable="false"><rect width="24" height="24" fill="` + fill + `" mask="url(#gm-` + name + `)"/></svg>`)
}

// Ikon kategori dari kolom categories.icon.
var categoryGlyph = map[string]string{"headphones": "headphones", "beauty": "beauty", "shirt": "shirt", "star": "star", "cooker": "cooker",
	// catalog: ikon tambahan yang bisa dipilih admin di Categories & brands
	"groceries": "groceries", "sofa": "sofa", "bag": "bag", "box": "box", "heart": "heart", "home": "home", "video": "video",
	"tshirt": "tshirt", "watch": "watch", "phone": "phone", "laptop": "laptop", "health": "health", "baby": "baby", "ball": "ball", "book": "book"}

func catGlyph(icon string) template.HTML {
	if g, ok := categoryGlyph[icon]; ok {
		return glyph(g)
	}
	return glyph("bag")
}

// rielOnly: "756,000 ៛" — harga riel di bawah harga dolar (gaya website Phsar pertama).
func rielOnly(p *Page, usd string) string {
	s := p.RielOf(usd)
	if s == "" {
		return ""
	}
	return strings.TrimSuffix(strings.TrimPrefix(s, "≈ "), "៛") + " ៛"
}

// durs: 16 → "16s", 75 → "1:15".
func durs(sec int) string {
	if sec < 60 {
		return strconv.Itoa(sec) + "s"
	}
	return strconv.Itoa(sec/60) + ":" + fmt.Sprintf("%02d", sec%60)
}

// hd membuat map untuk template "head": hd "ID" "h-new" "Title" "…" "Sub" "…".
func hd(kv ...any) map[string]any {
	m := map[string]any{}
	for i := 0; i+1 < len(kv); i += 2 {
		m[fmt.Sprint(kv[i])] = kv[i+1]
	}
	return m
}

// tileCols: kolom kotak kategori — semua dalam satu baris kalau muat (≤ max), selain itu dibagi rata ke baris
// sesedikit mungkin (12 kategori, max 6 → 2 baris × 6; 13 kategori, max 12 → 2 baris × 7).
func tileCols(n, max int) int {
	if n <= max {
		if n < 1 {
			return 1
		}
		return n
	}
	rows := (n + max - 1) / max
	return (n + rows - 1) / rows
}
