package admin

// Support tickets: kotak masuk CS untuk tiket pembeli & seller (izin support.handle).

import (
	"github.com/go-chi/chi/v5"

	"phsar/internal/tickets"
)

func (h *Handler) mountSupport(r chi.Router) {
	th := &tickets.Handler{S: &tickets.Service{Pool: h.pool, Log: h.log}}
	r.Group(func(r chi.Router) {
		r.Use(h.can("support.handle"))
		th.AdminRoutes(r, actor)
	})
}
