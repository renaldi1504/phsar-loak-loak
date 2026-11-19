package admin

// Konten & SEO website pembeli (izin content.manage — hanya super admin "*"):
// banner beranda, pita pengumuman, koleksi pilihan, meta SEO per halaman, gambar OG bawaan.

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"phsar/internal/content"
	"phsar/internal/storage"
)

type ContentRoutes struct{ Store *content.Store }

func (h *Handler) mountContent(r chi.Router) {
	if h.Content == nil {
		return
	}
	r.Route("/content", func(p chi.Router) {
		p.Use(h.can("content.manage"))
		p.Get("/", h.contentView)
		p.Post("/upload", h.contentUpload)
		p.Get("/products", h.contentProducts)
		p.Post("/banners", h.contentSaveBanner)
		p.Put("/banners/{id}", h.contentSaveBanner)
		p.Delete("/banners/{id}", h.contentDeleteBanner)
		p.Post("/collections", h.contentSaveCollection)
		p.Put("/collections/{id}", h.contentSaveCollection)
		p.Delete("/collections/{id}", h.contentDeleteCollection)
		p.Post("/seasons", h.contentSaveSeason) // kategori musiman
		p.Put("/seasons/{id}", h.contentSaveSeason)
		p.Delete("/seasons/{id}", h.contentDeleteSeason)
		p.Put("/categories/{id}/cover", h.contentCategoryCover) // sampul kotak kategori
		p.Put("/announcement", h.contentSaveAnnouncement)
		p.Put("/seo/{page}", h.contentSaveSEO)
		p.Put("/og", h.contentSaveOG)
	})
}

func (h *Handler) contentView(w http.ResponseWriter, r *http.Request) {
	snap, err := h.Content.Store.Fresh(r.Context())
	if err != nil {
		h.log.Error("konten", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not load content")
		return
	}
	type cat struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		CoverKey string `json:"cover_key"`
		CoverURL string `json:"cover_url"`
	}
	cats := []cat{}
	rows, err := h.pool.Query(r.Context(), `SELECT c.id, coalesce(t.name, c.id), coalesce(c.cover_key, '') FROM categories c
		LEFT JOIN category_translations t ON t.category_id = c.id AND t.lang = 'en'
		WHERE c.deleted_at IS NULL AND c.parent_id IS NULL ORDER BY c.sort_order, c.id`)
	if err == nil {
		for rows.Next() {
			var c cat
			if rows.Scan(&c.ID, &c.Name, &c.CoverKey) == nil {
				if c.CoverKey != "" {
					c.CoverURL = h.files.PublicURL(c.CoverKey)
				}
				cats = append(cats, c)
			}
		}
		rows.Close()
	}
	// nama produk untuk koleksi pilihan
	ids := []string{}
	for _, c := range snap.Collections {
		ids = append(ids, c.ProductIDs...)
	}
	for _, x := range snap.Seasons {
		ids = append(ids, x.ProductIDs...)
	}
	products, _ := h.contentProductList(r, `p.id = ANY($1::uuid[])`, ids)
	writeJSON(w, http.StatusOK, map[string]any{"banners": snap.Banners, "collections": snap.Collections, "seasons": snap.Seasons, "announcement": snap.Announcement,
		"seo": snap.SEOList(), "og_image_key": snap.OGImageKey, "og_image_url": snap.OGImageURL, "categories": cats, "products": products})
}

type contentProduct struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Store    string `json:"store"`
	Status   string `json:"status"`
	ImageURL string `json:"image_url"`
}

func (h *Handler) contentProductList(r *http.Request, where string, arg any) ([]contentProduct, error) {
	rows, err := h.pool.Query(r.Context(), `SELECT p.id, coalesce(t.name, ''), st.name, p.status::text,
		coalesce((SELECT coalesce(m.image_url, m.storage_key) FROM product_media m WHERE m.product_id = p.id AND m.deleted_at IS NULL AND m.kind = 'photo' ORDER BY m.sort_order LIMIT 1), '')
		FROM products p JOIN stores st ON st.id = p.store_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE p.deleted_at IS NULL AND `+where+` ORDER BY p.published_at DESC NULLS LAST LIMIT 20`, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []contentProduct{}
	for rows.Next() {
		var p contentProduct
		if err := rows.Scan(&p.ID, &p.Name, &p.Store, &p.Status, &p.ImageURL); err != nil {
			return nil, err
		}
		if p.ImageURL != "" && !strings.HasPrefix(p.ImageURL, "http") && !strings.HasPrefix(p.ImageURL, "/") {
			p.ImageURL = h.files.PublicURL(p.ImageURL)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GET /content/products?q= — cari produk tampil untuk koleksi pilihan.
func (h *Handler) contentProducts(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if len([]rune(q)) > 80 {
		q = string([]rune(q)[:80])
	}
	like := "%" + strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(q) + "%"
	l, err := h.contentProductList(r, `p.status IN ('active','sold_out') AND p.moderation_status = 'approved'
		AND (EXISTS (SELECT 1 FROM product_translations x WHERE x.product_id = p.id AND x.name ILIKE $1) OR st.name ILIKE $1)`, like)
	if err != nil {
		h.log.Error("cari produk konten", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "search failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"products": l})
}

// POST /content/upload (multipart: file) → {key, url}
func (h *Handler) contentUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 6<<20)
	f, _, err := r.FormFile("file")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad_upload", "send multipart with a file field")
		return
	}
	defer f.Close()
	up, err := storage.Uploader{Store: h.files, Pool: h.pool}.Save(r.Context(), actor(r), "content_image", io.Reader(f))
	switch {
	case errors.Is(err, storage.ErrTooLarge):
		writeErr(w, http.StatusRequestEntityTooLarge, "too_large", "image is larger than 5 MB")
	case errors.Is(err, storage.ErrType):
		writeErr(w, http.StatusUnsupportedMediaType, "bad_type", "use a JPG, PNG or WebP image")
	case err != nil:
		h.log.Error("upload konten", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "upload failed")
	default:
		writeJSON(w, http.StatusCreated, up)
	}
}

func (h *Handler) contentSaveBanner(w http.ResponseWriter, r *http.Request) {
	var in content.BannerInput
	if !decode(w, r, &in) {
		return
	}
	_, err := h.Content.Store.SaveBanner(r.Context(), actor(r), chi.URLParam(r, "id"), in)
	h.contentReply(w, r, err)
}

func (h *Handler) contentDeleteBanner(w http.ResponseWriter, r *http.Request) {
	h.contentReply(w, r, h.Content.Store.DeleteBanner(r.Context(), actor(r), chi.URLParam(r, "id")))
}

func (h *Handler) contentSaveCollection(w http.ResponseWriter, r *http.Request) {
	var in content.CollectionInput
	if !decode(w, r, &in) {
		return
	}
	_, err := h.Content.Store.SaveCollection(r.Context(), actor(r), chi.URLParam(r, "id"), in)
	h.contentReply(w, r, err)
}

func (h *Handler) contentDeleteCollection(w http.ResponseWriter, r *http.Request) {
	h.contentReply(w, r, h.Content.Store.DeleteCollection(r.Context(), actor(r), chi.URLParam(r, "id")))
}

func (h *Handler) contentSaveAnnouncement(w http.ResponseWriter, r *http.Request) {
	var in content.Announcement
	if !decode(w, r, &in) {
		return
	}
	h.contentReply(w, r, h.Content.Store.SaveAnnouncement(r.Context(), actor(r), in))
}

func (h *Handler) contentSaveSEO(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title       content.Text `json:"title"`
		Description content.Text `json:"description"`
	}
	if !decode(w, r, &in) {
		return
	}
	h.contentReply(w, r, h.Content.Store.SaveSEO(r.Context(), actor(r), chi.URLParam(r, "page"), in.Title, in.Description))
}

func (h *Handler) contentSaveOG(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ImageKey string `json:"image_key"`
	}
	if !decode(w, r, &in) {
		return
	}
	h.contentReply(w, r, h.Content.Store.SaveOGImage(r.Context(), actor(r), in.ImageKey))
}

func (h *Handler) contentReply(w http.ResponseWriter, r *http.Request, err error) {
	var fe content.FieldError
	switch {
	case err == nil:
		h.contentView(w, r)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, content.ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", "not found")
	default:
		h.log.Error("simpan konten", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save")
	}
}

func (h *Handler) contentSaveSeason(w http.ResponseWriter, r *http.Request) {
	var in content.SeasonInput
	if !decode(w, r, &in) {
		return
	}
	_, err := h.Content.Store.SaveSeason(r.Context(), actor(r), chi.URLParam(r, "id"), in)
	h.contentReply(w, r, err)
}

func (h *Handler) contentDeleteSeason(w http.ResponseWriter, r *http.Request) {
	h.contentReply(w, r, h.Content.Store.DeleteSeason(r.Context(), actor(r), chi.URLParam(r, "id")))
}

func (h *Handler) contentCategoryCover(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ImageKey string `json:"image_key"`
	}
	if !decode(w, r, &in) {
		return
	}
	h.contentReply(w, r, h.Content.Store.SetCategoryCover(r.Context(), actor(r), chi.URLParam(r, "id"), strings.TrimSpace(in.ImageKey)))
}
