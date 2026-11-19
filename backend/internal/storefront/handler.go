// Package storefront: website untuk pembeli (server-side rendering, ramah Google).
//
// URL: /km/…, /en/…, /zh/… (bahasa di URL supaya tiap bahasa terindeks Google, dengan hreflang).
//
//	/{lang}/              beranda
//	/{lang}/c/{kategori}  kategori
//	/{lang}/search?q=     pencarian
//	/{lang}/p/{slug}      detail produk (slug per bahasa)
//	/{lang}/s/{toko}      halaman toko
//	/sitemap.xml /robots.txt
package storefront

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"math"
	"net/http"
	"phsar/internal/shipping"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed templates/*.html assets/*
var files embed.FS

// LoginState: login SMS & Facebook tersedia atau tidak (auth.Service).
type LoginState interface {
	SMSEnabled(ctx context.Context) bool
	FacebookEnabled(ctx context.Context) bool
}

// PaymentState: status payment gateway saat ini (paygw.Store).
type PaymentState interface {
	Open(ctx context.Context) bool
	Test(ctx context.Context) bool
}

type Handler struct {
	Cat     *Catalog
	BaseURL string // PUBLIC_BASE_URL tanpa garis miring akhir
	Default string // bahasa default (DEFAULT_LANG)
	Log     *slog.Logger
	// Halaman belanja
	TelegramLogin bool          // tombol "Lanjut dengan Telegram" di login checkout
	Logins        LoginState    // cara login yang tersedia (SMS / Facebook diatur super admin)
	Legal         LegalDocs     // halaman legal & kontak
	Ads           Sponsor       // ads: produk bersponsor (nil = tanpa iklan)
	Payments      PaymentState  // diatur super admin: off = pita "segera dibuka" + checkout ditutup; fake = tombol "simulasi bayar"
	Photo         PhotoSearcher // cari dengan foto (nil = mati)
	Content       ContentSource // content: banner, pengumuman, koleksi, SEO (nil = tidak ada)
	FAQ           FAQSource     // askbot: tanya jawab produk dari seller (nil = tidak ada)
	photoRate     photoLimiter
	tmpl          *template.Template
	version       string
}

func New(cat *Catalog, baseURL, defaultLang string, log *slog.Logger) *Handler {
	h := &Handler{Cat: cat, BaseURL: strings.TrimRight(baseURL, "/"), Default: defaultLang, Log: log,
		version: strconv.FormatInt(time.Now().Unix(), 36)}
	if !validLang(h.Default) {
		h.Default = "km"
	}
	h.tmpl = template.Must(template.New("").Funcs(funcs).ParseFS(files, "templates/*.html"))
	return h
}

func (h *Handler) prelaunch(ctx context.Context) bool {
	return h.Payments != nil && !h.Payments.Open(ctx)
}

// Routes dipasang di root router.
func (h *Handler) Routes(r chi.Router) {
	r.Get("/", h.root)
	r.Get("/robots.txt", h.robots)
	r.Get("/sitemap.xml", h.sitemap)
	assets, _ := fs.Sub(files, "assets")
	r.Handle("/assets/*", http.StripPrefix("/assets/", cacheForever(http.FileServer(http.FS(assets)))))
	for _, l := range langs {
		lang := l
		r.Route("/"+lang, func(r chi.Router) {
			r.Use(middleware.GetHead) // HEAD dijawab seperti GET (crawler, monitor uptime)
			r.Get("/", h.wrap(lang, h.home))
			r.Get("/c/{slug}", h.wrap(lang, h.category))
			r.Get("/c/{slug}/watch", h.wrap(lang, h.watch)) // Tonton: video per kategori / sub-kategori
			r.Get("/season/{id}", h.wrap(lang, h.season))   // kategori musiman
			r.Get("/search", h.wrap(lang, h.search))
			r.Get("/new", h.wrap(lang, h.newest)) // daftar "Produk terbaru" (lanjutan beranda)
			r.Get("/p/{slug}", h.wrap(lang, h.product))
			r.Get("/s/{slug}", h.wrap(lang, h.store))
			r.Get("/legal/{slug}", h.wrap(lang, h.legalPage(lang)))
			r.Get("/contact", h.wrap(lang, h.contactPage(lang)))
			r.Get("/cart", h.wrap(lang, h.shopPage("cart")))
			r.Get("/checkout", h.wrap(lang, h.shopPage("checkout")))
			r.Get("/orders", h.wrap(lang, h.shopPage("orders")))
			r.Get("/orders/{no}", h.wrap(lang, h.shopPage("order")))
			r.Get("/support", h.wrap(lang, h.shopPage("support"))) // support: tiket bantuan
			r.Get("/chat", h.wrap(lang, h.shopPage("chat")))       // chat: chat pembeli ↔ toko
			r.Get("/flash-sale", h.wrap(lang, h.flashPage))        // flash: sesi flash sale
			r.Post("/photo-search", h.photoSearch(lang))
		})
	}
}

// ---- Data untuk template ----

type Page struct {
	Lang, Title, Description, Canonical, Image, Robots, Active string
	Alternates                                                 map[string]string // hreflang → URL absolut
	LangLinks                                                  map[string]string // ganti bahasa → path relatif
	JSONLD                                                     []template.JS
	Data                                                       any
	Riel                                                       float64
	Base, Version                                              string
	Year                                                       int
	Query                                                      string
	Categories                                                 []Category
	Seasons                                                    []SeasonView // kategori musiman yang tampil sekarang
	Prelaunch                                                  bool
	Ship                                                       shipping.Info // shipping: ongkir zona default (tampilan)
	Announce                                                   *Announce     // content: pita pengumuman
}

func (p *Page) T(key string, args ...any) string { return T(p.Lang, key, args...) }

func (p *Page) Langs() []string { return langs }

func (p *Page) LangName(l string) string { return langShort[l] }

func (p *Page) HTMLLang() string { return htmlLang[p.Lang] }

func (p *Page) HTMLLangOf() map[string]string { return htmlLang }

func (p *Page) LangFull() map[string]string { return langNames }

func (p *Page) OGLocale() string { return ogLocale[p.Lang] }

// Riel: "≈ 98.000៛" (dibulatkan ke 100 riel).
func (p *Page) RielOf(usd string) string {
	v, err := strconv.ParseFloat(usd, 64)
	if err != nil || v <= 0 {
		return ""
	}
	r := math.Round(v*p.Riel/100) * 100
	return "≈ " + thousands(int64(r)) + "៛"
}

func (p *Page) U(path string) string { return "/" + p.Lang + path }

type handlerFunc func(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error)

// newPage: Page dasar (header, kategori, ongkir, pengumuman, link bahasa) untuk path permintaan ini.
func (h *Handler) newPage(r *http.Request, lang string) *Page {
	p := &Page{Lang: lang, Base: h.BaseURL, Version: h.version, Prelaunch: h.prelaunch(r.Context()), Year: time.Now().Year(), Riel: h.Cat.RielPerUSD(r.Context()),
		Alternates: map[string]string{}, LangLinks: map[string]string{}, Image: h.BaseURL + "/assets/og.png"}
	p.Categories, _ = h.Cat.Categories(r.Context(), lang)
	p.Ship = shipping.GetInfo(r.Context(), h.Cat.Pool) // shipping
	h.applyContent(r, p)                               // content
	p.Seasons = h.seasonViews(p, h.contentSnap(r.Context()))
	// Default: halaman yang sama di bahasa lain (path sama).
	rest := strings.TrimPrefix(r.URL.Path, "/"+lang)
	for _, l := range langs {
		p.LangLinks[l] = "/" + l + rest
		if r.URL.RawQuery != "" {
			p.LangLinks[l] += "?" + r.URL.RawQuery
		}
		p.Alternates[l] = h.BaseURL + "/" + l + rest
	}
	p.Canonical = h.BaseURL + r.URL.Path
	return p
}

// wrap: siapkan Page, panggil handler halaman, render template (atau halaman error, lihat errors.go).
func (h *Handler) wrap(lang string, fn handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := h.newPage(r, lang)
		http.SetCookie(w, &http.Cookie{Name: "phsar_lang", Value: lang, Path: "/", MaxAge: 365 * 24 * 3600, SameSite: http.SameSiteLaxMode})

		name, status, err := fn(w, r, p)
		if err != nil {
			code, known := errorStatus(err)
			if !known {
				h.Log.Error("storefront", "path", r.URL.Path, "err", err)
			}
			if code >= 500 {
				h.writeStaticError(w, r, lang, code)
				return
			}
			h.fillError(p, r, code)
			h.writeError(w, r, p, code)
			return
		}
		if name == "" { // handler sudah menulis respons (redirect)
			return
		}
		var buf bytes.Buffer
		if err := h.tmpl.ExecuteTemplate(&buf, name, p); err != nil {
			h.Log.Error("storefront template", "page", name, "err", err)
			h.writeStaticError(w, r, lang, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Language", htmlLang[lang])
		w.Header().Set("Vary", "Accept-Language, Cookie")
		if status == http.StatusOK && w.Header().Get("Cache-Control") == "" {
			w.Header().Set("Cache-Control", "public, max-age=60")
		}
		w.WriteHeader(status)
		_, _ = w.Write(buf.Bytes())
	}
}

// root: / → bahasa dari cookie, lalu Accept-Language, lalu default.
func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	lang := ""
	if c, err := r.Cookie("phsar_lang"); err == nil && validLang(c.Value) {
		lang = c.Value
	}
	if lang == "" {
		for _, part := range strings.Split(r.Header.Get("Accept-Language"), ",") {
			code := strings.ToLower(strings.TrimSpace(strings.SplitN(part, ";", 2)[0]))
			if len(code) >= 2 && validLang(code[:2]) {
				lang = code[:2]
				break
			}
		}
	}
	if lang == "" {
		lang = h.Default
	}
	w.Header().Set("Vary", "Accept-Language, Cookie")
	http.Redirect(w, r, "/"+lang+"/", http.StatusFound)
}

// ---- Halaman ----

type homeData struct {
	Videos, Newest []Card
	Tiles          []Tile           // kotak kategori (ikon / sampul admin)
	Banners        []homeBanner     // content
	Collections    []homeCollection // content
	Flash          *FlashSession    // flash: sesi berjalan (atau berikutnya)
	NewestMore     string           // link halaman ke-2 "Produk terbaru" (infinite scroll)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
	videos, err := h.Cat.VideoFeed(r.Context(), p.Lang, 12)
	if err != nil {
		return "", 0, err
	}
	newest, err := h.Cat.Newest(r.Context(), p.Lang, pageSize+1, 0)
	if err != nil {
		return "", 0, err
	}
	p.Title = "Phsar — " + p.T("tagline")
	p.Description = p.T("metaDesc")
	banners, cols := h.homeContent(r, p) // content
	var fl *FlashSession                 // flash
	if list, err := h.Cat.FlashSessions(r.Context(), p.Lang, 1, 12); err == nil && len(list) > 0 {
		fl = &list[0]
	}
	more := ""
	if len(newest) > pageSize {
		newest, more = newest[:pageSize], "/"+p.Lang+"/new?page=2"
	}
	tiles, err := h.Cat.CategoryTiles(r.Context(), p.Lang, p.Categories)
	if err != nil {
		h.Log.Warn("kotak kategori", "err", err)
	}
	p.Active, p.Data = "home", homeData{Videos: videos, Newest: newest, Tiles: tiles, Banners: banners, Collections: cols, Flash: fl, NewestMore: more}
	h.applySEO(r, p, "home") // content
	p.JSONLD = append(p.JSONLD, ld(map[string]any{
		"@context": "https://schema.org", "@type": "WebSite", "name": "Phsar", "url": h.BaseURL + "/" + p.Lang + "/",
		"inLanguage": htmlLang[p.Lang],
		"potentialAction": map[string]any{"@type": "SearchAction", "query-input": "required name=q",
			"target": h.BaseURL + "/" + p.Lang + "/search?q={q}"},
	}))
	return "home", http.StatusOK, nil
}

type listData struct {
	Heading, Sub, Sort, Query, Empty string
	Category                         *Category
	Parent                           *Category  // kategori utama yang sedang dibuka (atau induk sub-kategori)
	Subs                             []Category // sub-kategori untuk baris chip kedua
	InSub                            bool       // sedang di halaman sub-kategori
	Cards                            []Card
	Sponsored                        []AdCard // ads
	NextPage                         string
	ShowSort                         bool
}

const pageSize = 24

// ---- Halaman bertingkat (infinite scroll) ----
//
// site.js memuat halaman berikutnya dengan ?page=N&fragment=1 dan hanya menerima kartu + link berikutnya
// (template "more"). Tanpa JS, link "Muat lagi" membuka halaman penuh ?page=N. Parameter (bukan header)
// supaya cache/CDN tidak pernah tertukar antara halaman penuh dan potongan.

type moreData struct {
	Cards    []Card
	NextPage string
}

func pageParam(r *http.Request) int {
	n, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if n < 1 {
		return 1
	}
	if n > 200 { // batas wajar; mencegah OFFSET raksasa
		return 200
	}
	return n
}

func isFragment(r *http.Request) bool { return r.URL.Query().Get("fragment") == "1" }

// nextPage: URL halaman berikutnya (path lengkap; parameter lain seperti sort/q tetap), tanpa fragment.
func nextPage(r *http.Request, page int) string {
	q := r.URL.Query()
	q.Del("fragment")
	q.Set("page", strconv.Itoa(page+1))
	return r.URL.Path + "?" + q.Encode()
}

// pageOf: potong hasil (limit pageSize+1) menjadi satu halaman + link berikutnya.
func pageOf(r *http.Request, cards []Card, page int) ([]Card, string) {
	if len(cards) > pageSize {
		return cards[:pageSize], nextPage(r, page)
	}
	return cards, ""
}

// fragment: balas potongan kartu untuk infinite scroll.
func (h *Handler) fragment(w http.ResponseWriter, p *Page, cards []Card, next string) (string, int, error) {
	w.Header().Set("X-Robots-Tag", "noindex")
	p.Data = moreData{cards, next}
	return "more", http.StatusOK, nil
}

// newest: /{lang}/new — semua produk terbaru, bertingkat.
func (h *Handler) newest(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
	page := pageParam(r)
	cards, err := h.Cat.Newest(r.Context(), p.Lang, pageSize+1, (page-1)*pageSize)
	if err != nil {
		return "", 0, err
	}
	cards, next := pageOf(r, cards, page)
	if isFragment(r) {
		return h.fragment(w, p, cards, next)
	}
	p.Title, p.Description = p.T("newTitle")+" · Phsar", p.T("newSub")
	if page > 1 {
		p.Robots = "noindex, follow"
	}
	p.Data = listData{Heading: p.T("newTitle"), Sub: p.T("newSub"), Cards: cards, NextPage: next, Empty: p.T("empty")}
	return "list", http.StatusOK, nil
}

func (h *Handler) category(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
	slug := chi.URLParam(r, "slug")
	cat, parent := findCategoryParent(p, slug)
	if cat == nil {
		return "", 0, errNotFound
	}
	// Link lama ?sub=slug (tampilan kategori v2) → halaman sub-kategori sendiri.
	if sub := r.URL.Query().Get("sub"); sub != "" {
		for _, c := range cat.Children {
			if c.Slug == sub {
				http.Redirect(w, r, "/"+p.Lang+"/c/"+c.Slug, http.StatusMovedPermanently)
				return "", 0, nil
			}
		}
	}
	if slug != cat.Slug { // slug kategori bahasa lain → slug bahasa ini
		http.Redirect(w, r, "/"+p.Lang+"/c/"+cat.Slug, http.StatusMovedPermanently)
		return "", 0, nil
	}
	page := pageParam(r)
	sort := r.URL.Query().Get("sort")
	if _, ok := sorts[sort]; !ok {
		sort = "new"
	}
	cards, err := h.Cat.ByCategory(r.Context(), p.Lang, cat.ID, sort, pageSize+1, (page-1)*pageSize)
	if err != nil {
		return "", 0, err
	}
	d := listData{Heading: cat.Name, Sub: p.T("products", "n", cat.Count), Sort: sort, Category: cat, ShowSort: true, Empty: p.T("empty")}
	// Baris sub-kategori: di kategori utama → anak-anaknya; di sub-kategori → saudara-saudaranya.
	d.Parent = cat
	if parent != nil {
		d.Parent = parent
	}
	d.Subs = d.Parent.Children
	d.InSub = parent != nil
	d.Cards, d.NextPage = pageOf(r, cards, page)
	if isFragment(r) {
		return h.fragment(w, p, d.Cards, d.NextPage)
	}
	if page == 1 { // ads: slot bersponsor hanya di halaman pertama
		d.Sponsored, d.Cards = h.sponsored(r.Context(), p, cat.ID, "", r.UserAgent(), d.Cards)
	}
	// Slug kategori berbeda per bahasa → link ganti bahasa & hreflang ke slug masing-masing.
	for _, l := range langs {
		cs, _ := h.Cat.Categories(r.Context(), l)
		for _, c := range flatCategories(cs) {
			if c.ID == cat.ID {
				p.LangLinks[l] = "/" + l + "/c/" + c.Slug
				p.Alternates[l] = h.BaseURL + "/" + l + "/c/" + c.Slug
			}
		}
	}
	p.Title = cat.Name + " · Phsar"
	p.Description = cat.Name + " — " + p.T("metaDesc")
	h.applySEO(r, p, "c:"+cat.ID) // content
	if page > 1 || sort != "new" {
		p.Robots = "noindex, follow"
	}
	p.Data = d
	return "list", http.StatusOK, nil
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if len([]rune(q)) > 80 {
		q = string([]rune(q)[:80])
	}
	p.Query, p.Active, p.Robots = q, "search", "noindex, follow"
	d := listData{Heading: p.T("searchTitle"), Query: q}
	if q != "" {
		page := pageParam(r)
		cards, err := h.Cat.Search(r.Context(), p.Lang, q, pageSize+1, (page-1)*pageSize)
		if err != nil {
			return "", 0, err
		}
		d.Cards, d.NextPage = pageOf(r, cards, page)
		if isFragment(r) {
			return h.fragment(w, p, d.Cards, d.NextPage)
		}
		total, err := h.Cat.SearchCount(r.Context(), q)
		if err != nil {
			return "", 0, err
		}
		if page == 1 {
			d.Sponsored, d.Cards = h.sponsored(r.Context(), p, "", q, r.UserAgent(), d.Cards) // ads
		}
		d.Sub = p.T("results", "n", total, "q", q)
		d.Empty = p.T("noResults", "q", q)
		p.Title = q + " · " + p.T("search") + " · Phsar"
	} else {
		p.Title = p.T("search") + " · Phsar"
	}
	p.Data = d
	return "list", http.StatusOK, nil
}

type productData struct {
	P       *Product
	Related []Card
	Langs   map[string]string
	FAQ     []QA // askbot:
}

func (h *Handler) product(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
	slug := chi.URLParam(r, "slug")
	prod, err := h.Cat.ProductBySlug(r.Context(), p.Lang, slug)
	if err != nil {
		return "", 0, err
	}
	if slug != prod.Slug { // slug lama / bahasa lain → URL resmi untuk bahasa ini
		http.Redirect(w, r, "/"+p.Lang+"/p/"+prod.Slug, http.StatusMovedPermanently)
		return "", 0, nil
	}
	related, err := h.Cat.Related(r.Context(), p.Lang, prod.CategoryID, prod.ID, 8)
	if err != nil {
		return "", 0, err
	}
	for _, l := range langs {
		p.LangLinks[l] = "/" + l + "/p/" + prod.Slugs[l]
		p.Alternates[l] = h.BaseURL + "/" + l + "/p/" + prod.Slugs[l]
	}
	p.Canonical = h.BaseURL + "/" + p.Lang + "/p/" + prod.Slug
	price := "$" + money(prod.MinPrice)
	p.Title = prod.Name + " – " + price + " · " + prod.Store.Name + " · Phsar"
	p.Description = metaDesc(prod.Description, prod.Name+" · "+price+" · "+p.T("trust1"))
	for _, m := range prod.Media {
		if m.Kind == "photo" {
			p.Image = m.URL
			break
		}
		if m.PosterURL != "" && p.Image == h.BaseURL+"/assets/og.png" {
			p.Image = m.PosterURL
		}
	}
	p.Active, p.Data = "product", productData{P: prod, Related: related, FAQ: h.productFAQ(r.Context(), p.Lang, prod.ID)} // askbot: FAQ
	crumbs := []any{map[string]any{"@type": "ListItem", "position": 1, "name": "Phsar", "item": h.BaseURL + "/" + p.Lang + "/"}}
	if prod.ParentSlug != "" {
		crumbs = append(crumbs, map[string]any{"@type": "ListItem", "position": len(crumbs) + 1, "name": prod.ParentName, "item": h.BaseURL + "/" + p.Lang + "/c/" + prod.ParentSlug})
	}
	crumbs = append(crumbs,
		map[string]any{"@type": "ListItem", "position": len(crumbs) + 1, "name": prod.CategoryName, "item": h.BaseURL + "/" + p.Lang + "/c/" + prod.CategorySlug},
		map[string]any{"@type": "ListItem", "position": len(crumbs) + 2, "name": prod.Name})
	p.JSONLD = append(p.JSONLD, productLD(h.BaseURL, p, prod), ld(map[string]any{
		"@context": "https://schema.org", "@type": "BreadcrumbList", "itemListElement": crumbs}))
	return "product", http.StatusOK, nil
}

type storeData struct {
	S        *StoreCard
	Cards    []Card
	NextPage string
}

func (h *Handler) store(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
	st, err := h.Cat.StoreBySlug(r.Context(), p.Lang, chi.URLParam(r, "slug"))
	if err != nil {
		return "", 0, err
	}
	page := pageParam(r)
	cards, err := h.Cat.ByStore(r.Context(), p.Lang, st.ID, pageSize+1, (page-1)*pageSize)
	if err != nil {
		return "", 0, err
	}
	cards, next := pageOf(r, cards, page)
	if isFragment(r) {
		return h.fragment(w, p, cards, next)
	}
	if page > 1 {
		p.Robots = "noindex, follow"
	}
	p.Title = st.Name + " · Phsar"
	p.Description = metaDesc(st.About, st.Name+" — "+p.T("metaDesc"))
	if st.LogoURL != "" {
		p.Image = st.LogoURL
	}
	p.Data = storeData{st, cards, next}
	p.JSONLD = append(p.JSONLD, ld(map[string]any{"@context": "https://schema.org", "@type": "Store", "name": st.Name,
		"url": p.Canonical, "image": nonEmpty(st.LogoURL), "description": nonEmpty(st.About),
		"address": map[string]any{"@type": "PostalAddress", "addressLocality": "Sihanoukville", "addressCountry": "KH"}}))
	return "store", http.StatusOK, nil
}

type shopData struct {
	Mode, OrderNo       string
	TelegramLogin, Test bool
	SMS, Closed, FB     bool
}

// shopPage: keranjang, checkout, pesanan — kerangka HTML; isinya dirender shop.js (data milik pembeli, tidak di-cache).
func (h *Handler) shopPage(mode string) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
		title := map[string]string{"cart": "cartTitle", "checkout": "checkout", "orders": "myOrders", "order": "myOrders", "support": "supportTitle", "chat": "chatTitle"}[mode]
		p.Title, p.Robots, p.Active = p.T(title)+" · Phsar", "noindex, nofollow", "cart"
		if mode == "orders" || mode == "order" {
			p.Active = "orders"
		}
		w.Header().Set("Cache-Control", "no-store")
		p.Data = shopData{Mode: mode, OrderNo: chi.URLParam(r, "no"), TelegramLogin: h.TelegramLogin, Test: h.Payments != nil && h.Payments.Test(r.Context()), SMS: h.Logins != nil && h.Logins.SMSEnabled(r.Context()), FB: h.Logins != nil && h.Logins.FacebookEnabled(r.Context()), Closed: h.prelaunch(r.Context())}
		return "shop", http.StatusOK, nil
	}
}

// ---- robots & sitemap ----

func (h *Handler) robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "User-agent: *\nAllow: /\nDisallow: /seller/\nDisallow: /admin/\nDisallow: /v1/\nDisallow: /dev/\nDisallow: /ad/\nDisallow: /*/search\nDisallow: /*/cart\n\nSitemap: %s/sitemap.xml\n", h.BaseURL)
}

func (h *Handler) sitemap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">` + "\n")
	write := func(paths map[string]string, mod time.Time) {
		for _, l := range langs {
			b.WriteString("<url><loc>" + xmlEsc(h.BaseURL+paths[l]) + "</loc>")
			for _, a := range langs {
				b.WriteString(`<xhtml:link rel="alternate" hreflang="` + htmlLang[a] + `" href="` + xmlEsc(h.BaseURL+paths[a]) + `"/>`)
			}
			if !mod.IsZero() {
				b.WriteString("<lastmod>" + mod.UTC().Format("2006-01-02") + "</lastmod>")
			}
			b.WriteString("</url>\n")
		}
	}
	home := map[string]string{}
	for _, l := range langs {
		home[l] = "/" + l + "/"
	}
	write(home, time.Time{})
	if h.Legal != nil {
		for _, slug := range []string{"legal/terms", "legal/privacy", "legal/returns", "legal/seller-terms", "contact"} {
			ps := map[string]string{}
			for _, l := range langs {
				ps[l] = "/" + l + "/" + slug
			}
			write(ps, time.Time{})
		}
	}
	catSlugs := map[string]map[string]string{}
	for _, l := range langs {
		cs, err := h.Cat.Categories(ctx, l)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		for _, c := range flatCategories(cs) {
			if catSlugs[c.ID] == nil {
				catSlugs[c.ID] = map[string]string{}
			}
			catSlugs[c.ID][l] = "/" + l + "/c/" + c.Slug
		}
	}
	for _, paths := range catSlugs {
		write(paths, time.Time{})
	}
	stores, _ := h.Cat.SitemapStores(ctx)
	for _, s := range stores {
		paths := map[string]string{}
		for _, l := range langs {
			paths[l] = "/" + l + "/s/" + s
		}
		write(paths, time.Time{})
	}
	prods, err := h.Cat.SitemapProducts(ctx, 15000)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	for _, p := range prods {
		paths := map[string]string{}
		fallback := ""
		for _, s := range p.Slugs {
			fallback = s
			break
		}
		for _, l := range langs {
			s := p.Slugs[l]
			if s == "" {
				s = fallback
			}
			paths[l] = "/" + l + "/p/" + s
		}
		write(paths, p.Updated)
	}
	b.WriteString("</urlset>\n")
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=600")
	_, _ = w.Write([]byte(b.String()))
}

// ---- JSON-LD ----

func ld(v any) template.JS {
	b, _ := json.Marshal(v) // json.Marshal meng-escape <, >, & → aman di dalam <script>
	return template.JS(b)
}

func productLD(base string, p *Page, prod *Product) template.JS {
	cond := map[string]string{"new": "NewCondition", "refurbished": "RefurbishedCondition"}[prod.Condition]
	if cond == "" {
		cond = "UsedCondition"
	}
	avail := "https://schema.org/InStock"
	if prod.Stock <= 0 {
		avail = "https://schema.org/OutOfStock"
	}
	images := []string{}
	var video map[string]any
	for _, m := range prod.Media {
		switch m.Kind {
		case "photo":
			images = append(images, m.URL)
		case "video":
			if video == nil {
				video = map[string]any{"@type": "VideoObject", "name": prod.Name, "description": metaDesc(prod.Description, prod.Name),
					"contentUrl": m.URL, "thumbnailUrl": nonEmpty(m.PosterURL), "uploadDate": publishedISO(prod)}
			}
		}
	}
	offer := map[string]any{"@type": "Offer", "priceCurrency": "USD", "price": money2(prod.MinPrice), "availability": avail,
		"itemCondition": "https://schema.org/" + cond, "url": p.Canonical,
		"seller": map[string]any{"@type": "Organization", "name": prod.Store.Name},
		"shippingDetails": map[string]any{"@type": "OfferShippingDetails",
			"shippingRate":        map[string]any{"@type": "MonetaryAmount", "value": 0, "currency": "USD"},
			"shippingDestination": map[string]any{"@type": "DefinedRegion", "addressCountry": "KH"}}}
	if prod.MinPrice != prod.MaxPrice {
		offer = map[string]any{"@type": "AggregateOffer", "priceCurrency": "USD", "lowPrice": money2(prod.MinPrice), "highPrice": money2(prod.MaxPrice),
			"offerCount": len(prod.Variants), "availability": avail, "itemCondition": "https://schema.org/" + cond, "url": p.Canonical}
	}
	m := map[string]any{"@context": "https://schema.org", "@type": "Product", "name": prod.Name, "description": metaDesc(prod.Description, prod.Name),
		"sku": prod.ID, "image": images, "category": prod.CategoryName, "offers": offer}
	if len(prod.Variants) > 0 && prod.Variants[0].SKU != "" {
		m["sku"] = prod.Variants[0].SKU
	}
	if video != nil {
		m["subjectOf"] = video
	}
	return ld(m)
}

func publishedISO(p *Product) string {
	if p.PublishedAt != nil {
		return p.PublishedAt.UTC().Format(time.RFC3339)
	}
	return p.UpdatedAt.UTC().Format(time.RFC3339)
}

// ---- helper ----

var funcs = template.FuncMap{
	"money": money,
	"cond":  func(p *Page, c string) string { return p.T(c) },
	"pct": func(price, compare string) int {
		a, b := num(price), num(compare)
		if a <= 0 || b <= a {
			return 0
		}
		return int(math.Round((b - a) / b * 100))
	},
	"dur":    func(s int) string { return fmt.Sprintf("%d:%02d", s/60, s%60) },
	"add":    func(a, b int) int { return a + b },
	"trimAt": func(s string) string { return strings.TrimPrefix(strings.TrimSpace(s), "@") },
	"eq2":    func(a, b string) bool { return a == b },
	"has": func(m []Media, kind string) bool {
		for _, x := range m {
			if x.Kind == kind {
				return true
			}
		}
		return false
	},
	"js":    func(s string) template.JS { return template.JS(s) },
	"month": func(t time.Time) string { return t.Format("Jan 2006") },
	"join":  strings.Join,
	"pc":    func(p *Page, c Card) cardCtx { return cardCtx{p, c} },
	"pg":    func(p *Page, cards []Card) gridCtx { return gridCtx{p, cards} },
	"fc":    func(p *Page, f FlashCard) map[string]any { return map[string]any{"P": p, "F": f} },
	"initial": func(s string) string {
		for _, r := range strings.TrimSpace(s) {
			return strings.ToUpper(string(r))
		}
		return "?"
	},
	"langname":  func(l string) string { return langNames[l] },
	"glyph":     glyph,
	"catGlyph":  catGlyph,
	"tileCols":  tileCols,
	"glyphDefs": glyphDefs,
	"rielOnly":  rielOnly,
	"durs":      durs,
	"hd":        hd,
	"money2":    money2,
	"take": func(c []Card, n int) []Card {
		if len(c) > n {
			return c[:n]
		}
		return c
	},
}

type cardCtx struct {
	P *Page
	C Card
}

type gridCtx struct {
	P     *Page
	Cards []Card
}

// money: "24" atau "24.50".
func money(s string) string {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	if v == math.Trunc(v) {
		return strconv.FormatFloat(v, 'f', 0, 64)
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func money2(s string) string {
	v, _ := strconv.ParseFloat(s, 64)
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func thousands(n int64) string {
	s := strconv.FormatInt(n, 10)
	var out []byte
	for i, c := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, c)
	}
	return string(out)
}

func metaDesc(text, fallback string) string {
	t := strings.Join(strings.Fields(text), " ")
	if t == "" {
		t = fallback
	}
	if r := []rune(t); len(r) > 155 {
		t = strings.TrimSpace(string(r[:152])) + "…"
	}
	return t
}

func nonEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func xmlEsc(s string) string {
	var b bytes.Buffer
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func cacheForever(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("v") != "" {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		next.ServeHTTP(w, r)
	})
}
