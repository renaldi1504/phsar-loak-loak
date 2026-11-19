package storefront

// content: banner beranda, pita pengumuman, koleksi pilihan, meta SEO, gambar OG bawaan — diatur super admin (paket content).

import (
	"context"
	"net/http"
	"strings"
	"time"

	"phsar/internal/content"
)

// ContentSource: snapshot konten (content.Store, cache 30 detik).
type ContentSource interface {
	Snapshot(ctx context.Context) *content.Snapshot
}

type Announce struct {
	Text, Link, Color string
	External          bool
}

type homeBanner struct {
	Image, Link, Title, Sub string
	External                bool
}

type homeCollection struct {
	ID, Title, More string
	Cards           []Card
}

func external(link string) bool { return strings.HasPrefix(link, "https://") }

func (h *Handler) contentSnap(ctx context.Context) *content.Snapshot {
	if h.Content == nil {
		return nil
	}
	return h.Content.Snapshot(ctx)
}

// applyContent: dipanggil di wrap untuk setiap halaman — pita pengumuman & gambar OG bawaan.
func (h *Handler) applyContent(r *http.Request, p *Page) {
	snap := h.contentSnap(r.Context())
	if snap == nil {
		return
	}
	if a := snap.LiveAnnouncement(time.Now()); a != nil {
		p.Announce = &Announce{Text: a.Text.In(p.Lang), Link: a.Link, Color: a.Color, External: external(a.Link)}
	}
	if snap.OGImageURL != "" {
		p.Image = snap.OGImageURL
	}
}

// applySEO: judul/deskripsi pengganti dari admin (page = home · c:<id>).
func (h *Handler) applySEO(r *http.Request, p *Page, page string) {
	snap := h.contentSnap(r.Context())
	if snap == nil {
		return
	}
	t, d := snap.Meta(page, p.Lang)
	if t != "" {
		p.Title = t
	}
	if d != "" {
		p.Description = d
	}
}

// homeContent: banner & koleksi yang tampil di beranda.
func (h *Handler) homeContent(r *http.Request, p *Page) ([]homeBanner, []homeCollection) {
	snap := h.contentSnap(r.Context())
	if snap == nil {
		return nil, nil
	}
	var banners []homeBanner
	for _, b := range snap.LiveBanners(time.Now()) {
		if b.ImageURL == "" {
			continue
		}
		banners = append(banners, homeBanner{Image: b.ImageURL, Link: b.Link, Title: b.Title.In(p.Lang), Sub: b.Subtitle.In(p.Lang), External: external(b.Link)})
	}
	var cols []homeCollection
	for _, c := range snap.LiveCollections() {
		hc := homeCollection{ID: "col-" + c.ID[:8], Title: c.Title.In(p.Lang)}
		var cards []Card
		var err error
		if c.Source == "category" {
			cards, err = h.Cat.ByCategory(r.Context(), p.Lang, c.CategoryID, "new", c.MaxItems, 0)
			for _, k := range flatCategories(p.Categories) {
				if k.ID == c.CategoryID {
					hc.More = "/" + p.Lang + "/c/" + k.Slug
				}
			}
		} else {
			cards, err = h.Cat.CardsByIDs(r.Context(), p.Lang, c.ProductIDs, c.MaxItems)
		}
		if err != nil {
			h.Log.Warn("koleksi beranda", "id", c.ID, "err", err)
			continue
		}
		if len(cards) == 0 {
			continue
		}
		hc.Cards = cards
		cols = append(cols, hc)
	}
	return banners, cols
}
