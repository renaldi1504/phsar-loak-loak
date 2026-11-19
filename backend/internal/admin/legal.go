package admin

// Dokumen legal (syarat, privasi, retur, perjanjian seller, hapus data) + identitas/kontak perusahaan.
// Setiap terbit = versi baru (append-only) untuk ketiga bahasa; versi lama tetap tersimpan sebagai bukti.

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"phsar/internal/legal"
)

type LegalRoutes struct {
	Store   *legal.Store
	BaseURL string
}

func (h *Handler) legalView(w http.ResponseWriter, r *http.Request) {
	docs, err := h.Legal.Store.Admin(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	st, err := h.Legal.Store.Settings(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"documents": docs, "settings": st, "langs": legal.Langs})
}

func (h *Handler) legalPublish(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Version string            `json:"version"`
		Bodies  map[string]string `json:"bodies"`
	}
	if !decode(w, r, &in) {
		return
	}
	h.legalReply(w, r, h.Legal.Store.Publish(r.Context(), actor(r), chi.URLParam(r, "type"), in.Version, in.Bodies))
}

func (h *Handler) legalSaveSettings(w http.ResponseWriter, r *http.Request) {
	var in legal.Settings
	if !decode(w, r, &in) {
		return
	}
	h.legalReply(w, r, h.Legal.Store.SaveSettings(r.Context(), actor(r), in))
}

// POST /legal/preview {body} → HTML seperti yang akan dilihat pembeli (placeholder diisi pengaturan saat ini).
func (h *Handler) legalPreview(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Body string `json:"body"`
	}
	if !decode(w, r, &in) {
		return
	}
	if len(in.Body) > 60000 {
		writeFieldErr(w, "body")
		return
	}
	st, err := h.Legal.Store.Settings(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"html": string(legal.Render(in.Body, st, "/en", strings.TrimPrefix(strings.TrimPrefix(h.Legal.BaseURL, "https://"), "http://")))})
}

func (h *Handler) legalReply(w http.ResponseWriter, r *http.Request, err error) {
	var fe legal.FieldError
	switch {
	case err == nil:
		h.legalView(w, r)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	default:
		h.log.Error("dokumen legal", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save")
	}
}
