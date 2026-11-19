package storefront

// Halaman legal (/km/legal/terms …) dan Kontak (/km/contact) — isinya diatur super admin di admin → Legal documents.

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"phsar/internal/legal"
)

// LegalDocs: sumber dokumen legal (legal.Store).
type LegalDocs interface {
	Latest(ctx context.Context, typ, lang string) (legal.Doc, error)
	Settings(ctx context.Context) (legal.Settings, error)
}

type legalLink struct {
	Slug, Title string
	Active      bool
}

type legalData struct {
	Doc       template.HTML
	Version   string
	Effective time.Time
	Links     []legalLink
	S         legal.Settings
	Title     string
	Contact   bool
	FromLang  string // dokumen belum ada di bahasa ini → tampil bahasa Inggris
}

func (h *Handler) legalLinks(lang, active string) []legalLink {
	out := []legalLink{}
	for _, t := range legal.Types {
		out = append(out, legalLink{Slug: "legal/" + t.Slug, Title: T(lang, "legal_"+t.Type), Active: active == t.Slug})
	}
	return append(out, legalLink{Slug: "contact", Title: T(lang, "contact"), Active: active == "contact"})
}

func (h *Handler) legalPage(lang string) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
		if h.Legal == nil {
			return "", 0, errNotFound
		}
		slug := chi.URLParam(r, "slug")
		typ, ok := legal.TypeBySlug(slug)
		if !ok {
			return "", 0, errNotFound
		}
		d, err := h.Legal.Latest(r.Context(), typ, lang)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", 0, errNotFound
		}
		if err != nil {
			return "", 0, err
		}
		st, err := h.Legal.Settings(r.Context())
		if err != nil {
			return "", 0, err
		}
		host := h.BaseURL
		if u, err := url.Parse(h.BaseURL); err == nil && u.Host != "" {
			host = u.Host
		}
		p.Title, p.Active = T(lang, "legal_"+typ)+" · Phsar", "legal"
		p.Description = T(lang, "legal_"+typ) + " — Phsar"
		data := legalData{Title: T(lang, "legal_"+typ), Doc: legal.Render(d.Body, st, "/"+lang, host), Version: d.Version, Effective: d.EffectiveAt, Links: h.legalLinks(lang, slug), S: st}
		if d.Lang != lang {
			data.FromLang = d.Lang
		}
		p.Data = data
		return "legal", http.StatusOK, nil
	}
}

func (h *Handler) contactPage(lang string) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
		if h.Legal == nil {
			return "", 0, errNotFound
		}
		st, err := h.Legal.Settings(r.Context())
		if err != nil {
			return "", 0, err
		}
		p.Title, p.Active, p.Description = T(lang, "contact")+" · Phsar", "legal", T(lang, "contactSub")
		p.Data = legalData{Links: h.legalLinks(lang, "contact"), S: st, Contact: true}
		return "legal", http.StatusOK, nil
	}
}
