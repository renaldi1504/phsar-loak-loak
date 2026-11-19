package admin

// Invoices & tax: pengaturan pajak Phsar + faktur komisi bulanan per toko (izin commission.manage).

import (
	"github.com/go-chi/chi/v5"

	"phsar/internal/invoices"
)

func (h *Handler) mountInvoices(r chi.Router) {
	ih := &invoices.Handler{S: &invoices.Service{Pool: h.pool}, Log: h.log}
	r.Route("/invoices", func(p chi.Router) {
		p.Use(h.can("commission.manage"))
		p.Get("/", ih.AdminList)
		p.Put("/settings", ih.AdminSaveSettings)
		p.Post("/generate", ih.AdminGenerate)
		p.Get("/{id}", ih.AdminGet)
		p.Get("/{id}/html", ih.AdminHTML)
		p.Post("/{id}/void", ih.AdminVoid)
	})
}
