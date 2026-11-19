package admin

// chat: CS membaca percakapan pembeli ↔ toko untuk sengketa (izin support.handle, hanya baca).
// Setiap percakapan yang dibuka dicatat di audit log (aksi chat.view).

import (
	"github.com/go-chi/chi/v5"

	"phsar/internal/chat"
)

func (h *Handler) mountChat(r chi.Router) {
	ch := &chat.Handler{S: &chat.Service{Pool: h.pool, Files: h.files, Log: h.log}}
	r.Group(func(r chi.Router) {
		r.Use(h.can("support.handle"))
		ch.AdminRoutes(r, actor)
	})
}
