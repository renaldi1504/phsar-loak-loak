package admin

// Sesi flash sale (izin commission.manage, di bawah /promos): sesi, barang (pilih banyak sekaligus), setujui pendaftaran seller.

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"phsar/internal/flash"
)

func (h *Handler) flashSvc() *flash.Service {
	return &flash.Service{Pool: h.pool, Files: h.files, Log: h.log}
}

func (h *Handler) mountFlash(p chi.Router) {
	p.Get("/flash", h.flashSessions)
	p.Post("/flash", h.flashSaveSession)
	p.Get("/flash/{id}", h.flashSession)
	p.Put("/flash/{id}", h.flashSaveSession)
	p.Post("/flash/{id}/cancel", h.flashCancel)
	p.Post("/flash/{id}/items", h.flashAddItems)
	p.Put("/flash/items/{item}", h.flashUpdateItem)
	p.Post("/flash/items/{item}/decide", h.flashDecide)
}

func (h *Handler) flashReply(w http.ResponseWriter, err error, ok func()) {
	var fe flash.FieldError
	switch {
	case err == nil:
		ok()
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, flash.ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", "not found")
	case errors.Is(err, flash.ErrState):
		writeErr(w, http.StatusConflict, "flash_state", "this session or item can no longer be changed")
	default:
		h.log.Error("flash sale", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save")
	}
}

func (h *Handler) flashSessions(w http.ResponseWriter, r *http.Request) {
	list, err := h.flashSvc().Sessions(r.Context())
	h.flashReply(w, err, func() { writeJSON(w, http.StatusOK, map[string]any{"sessions": list}) })
}

func (h *Handler) flashSessionJSON(w http.ResponseWriter, r *http.Request, id string) {
	svc := h.flashSvc()
	ses, err := svc.Session(r.Context(), id)
	if err != nil {
		h.flashReply(w, err, nil)
		return
	}
	items, err := svc.SessionItems(r.Context(), id)
	h.flashReply(w, err, func() { writeJSON(w, http.StatusOK, map[string]any{"session": ses, "items": items}) })
}

func (h *Handler) flashSession(w http.ResponseWriter, r *http.Request) {
	h.flashSessionJSON(w, r, chi.URLParam(r, "id"))
}

func (h *Handler) flashSaveSession(w http.ResponseWriter, r *http.Request) {
	var in flash.SessionInput
	if !decode(w, r, &in) {
		return
	}
	id, err := h.flashSvc().SaveSession(r.Context(), actor(r), chi.URLParam(r, "id"), in)
	h.flashReply(w, err, func() { h.flashSessionJSON(w, r, id) })
}

func (h *Handler) flashCancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.flashSvc().CancelSession(r.Context(), actor(r), id)
	h.flashReply(w, err, func() { h.flashSessionJSON(w, r, id) })
}

func (h *Handler) flashAddItems(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Items []flash.ItemInput `json:"items"`
	}
	if !decode(w, r, &in) {
		return
	}
	id := chi.URLParam(r, "id")
	_, err := h.flashSvc().AddItems(r.Context(), actor(r), id, in.Items)
	h.flashReply(w, err, func() { h.flashSessionJSON(w, r, id) })
}

func (h *Handler) flashUpdateItem(w http.ResponseWriter, r *http.Request) {
	var in flash.ItemInput
	if !decode(w, r, &in) {
		return
	}
	err := h.flashSvc().UpdateItem(r.Context(), actor(r), chi.URLParam(r, "item"), in)
	h.flashReply(w, err, func() { writeJSON(w, http.StatusOK, map[string]bool{"ok": true}) })
}

func (h *Handler) flashDecide(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Action string `json:"action"`
		Note   string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	err := h.flashSvc().Decide(r.Context(), actor(r), chi.URLParam(r, "item"), in.Action, in.Note)
	h.flashReply(w, err, func() { writeJSON(w, http.StatusOK, map[string]bool{"ok": true}) })
}
