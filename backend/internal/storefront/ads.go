package storefront

// ads: slot produk bersponsor di atas hasil cari & halaman kategori (paket ads memilih, menghitung tayangan & klik).

import (
	"context"
	"time"

	"phsar/internal/ads"
)

// Sponsor: sumber iklan (ads.Service).
type Sponsor interface {
	Pick(ctx context.Context, categoryID, query string) ([]ads.Slot, error)
	Impress(ctx context.Context, campaignIDs []string) error
}

// AdCard: kartu produk biasa + tautan klik iklan (/ad/{id}).
type AdCard struct {
	Card
	AdURL string
}

type adGridCtx struct {
	P     *Page
	Cards []AdCard
}

func init() {
	for l, s := range map[string]string{"en": "Sponsored", "km": "ផ្សាយពាណិជ្ជកម្ម", "zh": "赞助"} {
		strs[l]["sponsored"] = s
	}
	funcs["adg"] = func(p *Page, cards []AdCard) adGridCtx { return adGridCtx{p, cards} }
}

// sponsored memilih slot iklan untuk halaman ini, membuang produk yang sama dari hasil biasa, dan mencatat tayangan.
func (h *Handler) sponsored(ctx context.Context, p *Page, categoryID, query, ua string, organic []Card) ([]AdCard, []Card) {
	if h.Ads == nil || p.Prelaunch {
		return nil, organic
	}
	slots, err := h.Ads.Pick(ctx, categoryID, query)
	if err != nil {
		h.Log.Warn("iklan: gagal memilih slot", "err", err)
		return nil, organic
	}
	if len(slots) == 0 {
		return nil, organic
	}
	ids := make([]string, len(slots))
	for i, s := range slots {
		ids[i] = s.ProductID
	}
	cards, err := h.Cat.CardsByIDs(ctx, p.Lang, ids, len(ids))
	if err != nil {
		h.Log.Warn("iklan: gagal memuat kartu", "err", err)
		return nil, organic
	}
	byProduct := map[string]string{}
	for _, s := range slots {
		byProduct[s.ProductID] = s.CampaignID
	}
	out, shown := []AdCard{}, []string{}
	for _, c := range cards {
		cid := byProduct[c.ID]
		out = append(out, AdCard{Card: c, AdURL: "/ad/" + cid + "?l=" + p.Lang})
		shown = append(shown, cid)
	}
	rest := organic[:0:0]
	for _, c := range organic {
		if _, dup := byProduct[c.ID]; !dup {
			rest = append(rest, c)
		}
	}
	if len(shown) > 0 && !ads.IsBot(ua) {
		go func() { // tayangan dicatat di latar: satu pernyataan per halaman
			c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if err := h.Ads.Impress(c, shown); err != nil {
				h.Log.Warn("iklan: gagal mencatat tayangan", "err", err)
			}
		}()
	}
	return out, rest
}
