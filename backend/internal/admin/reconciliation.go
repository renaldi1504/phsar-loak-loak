package admin

// Rekonsiliasi bank (izin reconciliation.run): unggah ekspor transaksi ABA PayWay (CSV), cocokkan dengan
// pembayaran & refund di Phsar, selesaikan selisih dengan catatan. Logika ada di paket reconciliation.

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"phsar/internal/reconciliation"
)

func (h *Handler) mountReconciliation(r chi.Router) {
	r.Route("/reconciliation", func(p chi.Router) {
		p.Use(h.can("reconciliation.run"))
		p.Get("/", h.reconList)
		p.Post("/batches", h.reconUpload)
		p.Get("/batches/{id}", h.reconGet)
		p.Post("/batches/{id}/rematch", h.reconRematch)
		p.Delete("/batches/{id}", h.reconDelete)
		p.Post("/rows/{id}/resolve", h.reconResolve)
	})
}

func (h *Handler) recon() *reconciliation.Service { return &reconciliation.Service{Pool: h.pool} }

func (h *Handler) reconList(w http.ResponseWriter, r *http.Request) {
	raw, err := h.recon().List(r.Context())
	if err != nil {
		h.reconFail(w, err)
		return
	}
	writeRawJSON(w, []byte(`{"batches":`+string(raw)+`}`))
}

// POST /reconciliation/batches (multipart: file, period_from?, period_to?)
func (h *Handler) reconUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, reconciliation.MaxBytes+(64<<10))
	if err := r.ParseMultipartForm(reconciliation.MaxBytes); err != nil {
		writeJSON(w, http.StatusRequestEntityTooLarge, map[string]any{"error": map[string]string{"code": "validation", "field": "file", "message": "file too large (max 5 MB) or not multipart"}})
		return
	}
	f, hdr, err := r.FormFile("file")
	if err != nil {
		writeFieldErr(w, "file")
		return
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, reconciliation.MaxBytes+1))
	if err != nil {
		writeFieldErr(w, "file")
		return
	}
	id, err := h.recon().Import(r.Context(), actor(r), hdr.Filename, data, r.FormValue("period_from"), r.FormValue("period_to"))
	if err != nil {
		h.reconFail(w, err)
		return
	}
	h.reconWrite(w, r, id, "open")
}

func (h *Handler) reconGet(w http.ResponseWriter, r *http.Request) {
	h.reconWrite(w, r, chi.URLParam(r, "id"), r.URL.Query().Get("filter"))
}

func (h *Handler) reconWrite(w http.ResponseWriter, r *http.Request, id, filter string) {
	if !uuidRe.MatchString(id) {
		writeErr(w, http.StatusNotFound, "not_found", "batch not found")
		return
	}
	raw, err := h.recon().Get(r.Context(), id, filter)
	if err != nil {
		h.reconFail(w, err)
		return
	}
	writeRawJSON(w, raw)
}

func (h *Handler) reconRematch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !uuidRe.MatchString(id) {
		writeErr(w, http.StatusNotFound, "not_found", "batch not found")
		return
	}
	if err := h.recon().Rematch(r.Context(), actor(r), id); err != nil {
		h.reconFail(w, err)
		return
	}
	h.reconWrite(w, r, id, r.URL.Query().Get("filter"))
}

func (h *Handler) reconDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !uuidRe.MatchString(id) {
		writeErr(w, http.StatusNotFound, "not_found", "batch not found")
		return
	}
	if err := h.recon().Delete(r.Context(), actor(r), id); err != nil {
		h.reconFail(w, err)
		return
	}
	h.reconList(w, r)
}

func (h *Handler) reconResolve(w http.ResponseWriter, r *http.Request) {
	rowID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "row not found")
		return
	}
	var in struct {
		Note string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	batchID, err := h.recon().Resolve(r.Context(), actor(r), rowID, in.Note)
	if err != nil {
		h.reconFail(w, err)
		return
	}
	h.reconWrite(w, r, batchID, r.URL.Query().Get("filter"))
}

func (h *Handler) reconFail(w http.ResponseWriter, err error) {
	var fe reconciliation.FieldError
	switch {
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, reconciliation.ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", "not found")
	case errors.Is(err, reconciliation.ErrNotOpen):
		writeErr(w, http.StatusConflict, "not_exception", err.Error())
	default:
		h.log.Error("rekonsiliasi", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "reconciliation failed")
	}
}
