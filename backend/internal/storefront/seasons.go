package storefront

// Kategori musiman ("Musim ini"), halaman musim, dan feed Tonton per kategori.

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"phsar/internal/content"
)

// SeasonView: kartu musim di beranda & menu desktop.
type SeasonView struct {
	ID, Title, Color, Tag, URL string
}

var kmMonths = []string{"មករា", "កុម្ភៈ", "មីនា", "មេសា", "ឧសភា", "មិថុនា", "កក្កដា", "សីហា", "កញ្ញា", "តុលា", "វិច្ឆិកា", "ធ្នូ"}

var phnomPenh = func() *time.Location {
	if l, err := time.LoadLocation("Asia/Phnom_Penh"); err == nil {
		return l
	}
	return time.FixedZone("ICT", 7*3600)
}()

// dayMonth: "1 Jan" · "1 មករា" · "1月1日" (waktu Kamboja).
func dayMonth(lang string, t time.Time) string {
	t = t.In(phnomPenh)
	switch lang {
	case "km":
		return strconv.Itoa(t.Day()) + " " + kmMonths[t.Month()-1]
	case "zh":
		return strconv.Itoa(int(t.Month())) + "月" + strconv.Itoa(t.Day()) + "日"
	}
	return t.Format("2 Jan")
}

func (h *Handler) seasonViews(p *Page, snap *content.Snapshot) []SeasonView {
	if snap == nil {
		return nil
	}
	out := []SeasonView{}
	for _, x := range snap.LiveSeasons(time.Now()) {
		v := SeasonView{ID: x.ID[:8], Title: x.Title.In(p.Lang), Color: x.Color, URL: p.U("/season/" + x.ID[:8])}
		if x.EndsAt != nil {
			// tanggal terakhir tampil = sehari sebelum ends_at kalau berakhir tengah malam
			last := x.EndsAt.In(phnomPenh).Add(-time.Second)
			v.Tag = p.T("seasonUntil", "d", dayMonth(p.Lang, last))
		} else {
			v.Tag = p.T("seasonAllYear")
		}
		out = append(out, v)
	}
	return out
}

type seasonData struct {
	Season SeasonView
	Cards  []Card
}

// season: /{lang}/season/{id8} — produk kategori musiman.
func (h *Handler) season(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
	snap := h.contentSnap(r.Context())
	if snap == nil {
		return "", 0, errNotFound
	}
	x := snap.SeasonByID(chi.URLParam(r, "id"), time.Now())
	if x == nil {
		return "", 0, errNotFound
	}
	var cards []Card
	var err error
	switch x.Source {
	case "category":
		cards, err = h.Cat.ByCategory(r.Context(), p.Lang, x.CategoryID, "popular", 60, 0)
	case "products":
		cards, err = h.Cat.CardsByIDs(r.Context(), p.Lang, x.ProductIDs, 60)
	default:
		cards, err = h.Cat.SearchCards(r.Context(), p.Lang, x.Query, 60)
	}
	if err != nil {
		return "", 0, err
	}
	var v SeasonView
	for _, s := range p.Seasons {
		if s.ID == x.ID[:8] {
			v = s
		}
	}
	p.Title = v.Title + " · Phsar"
	p.Description = v.Title + " — " + p.T("metaDesc")
	p.Data = seasonData{Season: v, Cards: cards}
	return "season", http.StatusOK, nil
}

type watchData struct {
	Category *Category // kategori atau sub-kategori yang diputar
	Parent   *Category // kategori utama (baris chip sub-kategori)
	Items    []WatchItem
	Back     string
}

// watch: /{lang}/c/{slug}/watch — video layar penuh kategori itu (geser ke atas untuk video berikutnya).
func (h *Handler) watch(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
	slug := chi.URLParam(r, "slug")
	cat, parent := findCategoryParent(p, slug)
	if cat == nil {
		return "", 0, errNotFound
	}
	// Link lama ?sub=slug → URL sub-kategori sendiri.
	if sub := r.URL.Query().Get("sub"); sub != "" {
		for _, c := range cat.Children {
			if c.Slug == sub {
				http.Redirect(w, r, p.U("/c/"+c.Slug+"/watch"), http.StatusMovedPermanently)
				return "", 0, nil
			}
		}
	}
	if slug != cat.Slug { // slug bahasa lain → slug bahasa ini
		http.Redirect(w, r, p.U("/c/"+cat.Slug+"/watch"), http.StatusMovedPermanently)
		return "", 0, nil
	}
	if parent == nil {
		parent = cat
	}
	items, err := h.Cat.WatchFeed(r.Context(), p.Lang, cat.ID, 40)
	if err != nil {
		return "", 0, err
	}
	p.Title = p.T("watchTitle", "c", cat.Name) + " · Phsar"
	p.Description = cat.Name + " — " + p.T("metaDesc")
	p.Robots = "noindex, follow"
	p.Active = "watch"
	p.Data = watchData{Category: cat, Parent: parent, Items: items, Back: p.U("/c/" + cat.Slug)}
	return "watch", http.StatusOK, nil
}

// findCategoryParent: kategori (utama atau sub) dengan slug/ID itu, plus induknya (nil kalau kategori utama).
func findCategoryParent(p *Page, slug string) (cat, parent *Category) {
	for i := range p.Categories {
		top := &p.Categories[i]
		if top.Slug == slug || top.ID == slug {
			return top, nil
		}
		for j := range top.Children {
			if k := &top.Children[j]; k.Slug == slug || k.ID == slug {
				return k, top
			}
		}
	}
	return nil, nil
}
