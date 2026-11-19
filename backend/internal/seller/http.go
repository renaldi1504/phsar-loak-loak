package seller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
	"phsar/internal/hours"
	"phsar/internal/media"
	"phsar/internal/paging"
	"phsar/internal/storage"
)

type Handler struct {
	Orders   http.Handler     // pesanan toko (paket shop); dipasang di /stores/{storeID}/orders
	Payouts  http.Handler     // saldo & payout (paket payout); khusus pemilik
	Returns  http.Handler     // retur & sengketa (paket returns)
	Strikes  http.HandlerFunc // reports: strike toko (hanya baca)
	Tickets  http.Handler     // support: tiket bantuan toko (paket tickets)
	Ads      http.Handler     // ads: iklan bersponsor (paket ads); khusus pemilik
	Flash    http.Handler     // flash: sesi flash sale (paket flash); khusus pemilik
	Askbot   http.Handler     // askbot: draf jawaban pembeli + FAQ toko (paket askbot); pemilik & staf
	Chat     http.Handler     // chat: chat pembeli ↔ toko (paket chat); pemilik & staf
	svc      *Service
	uploader storage.Uploader
	log      *slog.Logger

	// courierreturn: paket gagal antar yang dikembalikan ke toko (pemilik & staf)
	CourierReturns http.Handler
}

func NewHandler(svc *Service, up storage.Uploader, log *slog.Logger) *Handler {
	return &Handler{svc: svc, uploader: up, log: log}
}

// Routes dipasang di /v1/seller (sudah melewati auth.RequireUser).
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/uploads", h.upload)
	r.Get("/stores", h.myStores)
	r.Post("/stores", h.createStore)
	r.Get("/invites/{token}", h.inviteDetails)
	r.Post("/invites/{token}/accept", h.acceptInvite)
	r.Route("/stores/{storeID}", func(r chi.Router) {
		r.Use(h.member)
		r.Get("/", h.getStore)
		r.Get("/today", h.today) // Seller Center v2: layar "Hari ini"
		r.Get("/hours", h.getHours)
		r.With(ownerOnly).Put("/hours", h.putHours)
		r.With(ownerOnly).Patch("/", h.updateStore)
		r.With(ownerOnly).Post("/submit", h.submit)
		r.With(ownerOnly).Get("/kyc", h.getKYC)
		r.With(ownerOnly).Put("/kyc", h.putKYC)
		r.Get("/pickup-points", h.listPickup)
		r.With(ownerOnly).Post("/pickup-points", h.createPickup)
		r.With(ownerOnly).Put("/pickup-points/{pointID}", h.updatePickup)
		r.With(ownerOnly).Delete("/pickup-points/{pointID}", h.deletePickup)
		// Produk: pemilik dan staf toko boleh mengelola.
		r.Get("/products", h.listProducts)
		r.Post("/products", h.createProduct)
		r.Get("/products/{productID}", h.getProduct)
		r.Put("/products/{productID}", h.updateProduct)
		r.Delete("/products/{productID}", h.deleteProduct)
		r.Patch("/products/{productID}/variants", h.updateStockPrice)
		r.Post("/products/{productID}/status", h.setProductStatus)
		r.Put("/products/{productID}/translations/{lang}", h.editTranslation)
		r.Delete("/products/{productID}/translations/{lang}", h.resetTranslation)
		r.Post("/products/{productID}/translations/retry", h.retryTranslation)
		// Tim
		r.With(ownerOnly).Get("/team", h.getTeam)
		r.With(ownerOnly).Post("/invites", h.createInvite)
		r.With(ownerOnly).Delete("/invites/{inviteID}", h.revokeInvite)
		r.Delete("/members/{userID}", h.removeMember) // owner menghapus staf, atau staf keluar sendiri
		// Pesanan: pemilik dan staf.
		if h.Orders != nil {
			r.Mount("/orders", h.Orders)
		}
		if h.Strikes != nil { // reports:
			r.Get("/strikes", h.Strikes)
		}
		if h.Tickets != nil { // support
			r.Mount("/tickets", h.Tickets)
		}
		if h.Returns != nil {
			r.Mount("/returns", h.Returns)
		}
		if h.Askbot != nil { // askbot:
			r.Mount("/assist", h.Askbot)
		}
		if h.Chat != nil { // chat:
			r.Mount("/chat", h.Chat)
		}
		if h.CourierReturns != nil { // courierreturn
			r.Mount("/courier-returns", h.CourierReturns)
		}
		if h.Flash != nil { // flash
			r.With(ownerOnly).Mount("/flash", h.Flash)
		}
		if h.Ads != nil { // ads
			r.With(ownerOnly).Mount("/ads", h.Ads)
		}
		if h.Payouts != nil {
			r.With(ownerOnly).Mount("/payouts", h.Payouts)
		}
	})
	return r
}

type roleKey struct{}

func (h *Handler) member(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := auth.UserFrom(r.Context())
		role, err := h.svc.Role(r.Context(), chi.URLParam(r, "storeID"), c.Subject)
		if err != nil {
			h.fail(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), roleKey{}, role)))
	})
}

func ownerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if role, _ := r.Context().Value(roleKey{}).(string); role != "owner" {
			writeErr(w, http.StatusForbidden, "owner_only", ErrOwnerOnly.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// POST /v1/seller/uploads  (multipart: kind, file)
func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, storage.MaxVideoBytes+(1<<20))
	mr, err := r.MultipartReader()
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad_upload", "kirim multipart/form-data dengan field kind lalu file")
		return
	}
	kind := ""
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad_upload", err.Error())
			return
		}
		switch part.FormName() {
		case "kind":
			b, _ := io.ReadAll(io.LimitReader(part, 64))
			kind = string(b)
		case "file":
			if kind == "" {
				writeErr(w, http.StatusBadRequest, "bad_upload", "field kind harus dikirim sebelum file")
				return
			}
			up, err := h.uploader.Save(r.Context(), uid(r), kind, part)
			if err != nil {
				h.fail(w, err)
				return
			}
			writeJSON(w, http.StatusCreated, up)
			return
		}
	}
	writeErr(w, http.StatusBadRequest, "bad_upload", "field file tidak ada")
}

func (h *Handler) myStores(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.MyStores(r.Context(), uid(r))
	h.respond(w, http.StatusOK, map[string]any{"stores": list}, err)
}

func (h *Handler) createStore(w http.ResponseWriter, r *http.Request) {
	var in CreateStoreInput
	if !decode(w, r, &in) {
		return
	}
	st, err := h.svc.CreateStore(r.Context(), uid(r), in)
	h.respond(w, http.StatusCreated, st, err)
}

func (h *Handler) getStore(w http.ResponseWriter, r *http.Request) {
	st, err := h.svc.GetStore(r.Context(), chi.URLParam(r, "storeID"))
	st.Role, _ = r.Context().Value(roleKey{}).(string)
	if err == nil {
		if c, cerr := h.svc.Commission(r.Context(), st.ID); cerr == nil {
			st.Commission = &c
		}
	}
	h.respond(w, http.StatusOK, st, err)
}

func (h *Handler) updateStore(w http.ResponseWriter, r *http.Request) {
	var in UpdateStoreInput
	if !decode(w, r, &in) {
		return
	}
	st, err := h.svc.UpdateStore(r.Context(), chi.URLParam(r, "storeID"), uid(r), in)
	h.respond(w, http.StatusOK, st, err)
}

func (h *Handler) submit(w http.ResponseWriter, r *http.Request) {
	st, err := h.svc.Submit(r.Context(), chi.URLParam(r, "storeID"), uid(r))
	if errors.Is(err, ErrIncomplete) {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": map[string]string{"code": "store_incomplete", "message": err.Error()}, "checklist": st.Checklist})
		return
	}
	h.respond(w, http.StatusOK, st, err)
}

func (h *Handler) getKYC(w http.ResponseWriter, r *http.Request) {
	k, err := h.svc.GetKYC(r.Context(), chi.URLParam(r, "storeID"))
	h.respond(w, http.StatusOK, map[string]any{"kyc": k}, err)
}

func (h *Handler) putKYC(w http.ResponseWriter, r *http.Request) {
	var in KYCInput
	if !decode(w, r, &in) {
		return
	}
	k, err := h.svc.PutKYC(r.Context(), chi.URLParam(r, "storeID"), uid(r), in)
	h.respond(w, http.StatusOK, map[string]any{"kyc": k}, err)
}

func (h *Handler) listPickup(w http.ResponseWriter, r *http.Request) {
	l, err := h.svc.ListPickupPoints(r.Context(), chi.URLParam(r, "storeID"))
	h.respond(w, http.StatusOK, map[string]any{"pickup_points": l}, err)
}

func (h *Handler) createPickup(w http.ResponseWriter, r *http.Request) {
	var in PickupInput
	if !decode(w, r, &in) {
		return
	}
	l, err := h.svc.CreatePickupPoint(r.Context(), chi.URLParam(r, "storeID"), uid(r), in)
	h.respond(w, http.StatusCreated, map[string]any{"pickup_points": l}, err)
}

func (h *Handler) updatePickup(w http.ResponseWriter, r *http.Request) {
	var in PickupInput
	if !decode(w, r, &in) {
		return
	}
	l, err := h.svc.UpdatePickupPoint(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "pointID"), uid(r), in)
	h.respond(w, http.StatusOK, map[string]any{"pickup_points": l}, err)
}

func (h *Handler) deletePickup(w http.ResponseWriter, r *http.Request) {
	l, err := h.svc.DeletePickupPoint(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "pointID"), uid(r))
	h.respond(w, http.StatusOK, map[string]any{"pickup_points": l}, err)
}

// listProducts: ?filter=&q=&sort=&page=&size= (tanpa page: semua, maks. 500 — aplikasi lama).
func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	l, err := h.svc.ListProducts(r.Context(), chi.URLParam(r, "storeID"), ProductQuery{
		Filter: q.Get("filter"), Sort: q.Get("sort"), Page: paging.FromRequest(r, 24, 100)})
	h.respond(w, http.StatusOK, l, err)
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var in ProductInput
	if !decode(w, r, &in) {
		return
	}
	p, err := h.svc.CreateProduct(r.Context(), chi.URLParam(r, "storeID"), uid(r), in)
	h.respond(w, http.StatusCreated, p, err)
}

func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.GetProduct(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "productID"))
	h.respond(w, http.StatusOK, p, err)
}

func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	var in ProductInput
	if !decode(w, r, &in) {
		return
	}
	p, err := h.svc.UpdateProduct(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "productID"), uid(r), in)
	h.respond(w, http.StatusOK, p, err)
}

func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.DeleteProduct(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "productID"), uid(r)); err != nil {
		h.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateStockPrice(w http.ResponseWriter, r *http.Request) {
	var in StockPriceInput
	if !decode(w, r, &in) {
		return
	}
	p, err := h.svc.UpdateStockPrice(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "productID"), uid(r), in)
	h.respond(w, http.StatusOK, p, err)
}

func (h *Handler) setProductStatus(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Status string `json:"status"`
	}
	if !decode(w, r, &in) {
		return
	}
	p, err := h.svc.SetProductStatus(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "productID"), uid(r), in.Status)
	h.respond(w, http.StatusOK, p, err)
}

func (h *Handler) editTranslation(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if !decode(w, r, &in) {
		return
	}
	p, err := h.svc.EditTranslation(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "productID"), uid(r),
		chi.URLParam(r, "lang"), in.Name, in.Description)
	h.respond(w, http.StatusOK, p, err)
}

func (h *Handler) resetTranslation(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.ResetTranslation(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "productID"), uid(r), chi.URLParam(r, "lang"))
	h.respond(w, http.StatusOK, p, err)
}

func (h *Handler) retryTranslation(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.RetryTranslation(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "productID"), uid(r))
	h.respond(w, http.StatusOK, p, err)
}

func (h *Handler) getTeam(w http.ResponseWriter, r *http.Request) {
	t, err := h.svc.GetTeam(r.Context(), chi.URLParam(r, "storeID"))
	h.respond(w, http.StatusOK, t, err)
}

func (h *Handler) createInvite(w http.ResponseWriter, r *http.Request) {
	inv, err := h.svc.CreateInvite(r.Context(), chi.URLParam(r, "storeID"), uid(r))
	h.respond(w, http.StatusCreated, inv, err)
}

func (h *Handler) revokeInvite(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.RevokeInvite(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "inviteID"), uid(r)); err != nil {
		h.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) removeMember(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(roleKey{}).(string)
	if err := h.svc.RemoveMember(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "userID"), uid(r), role); err != nil {
		h.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) inviteDetails(w http.ResponseWriter, r *http.Request) {
	info, err := h.svc.InviteDetails(r.Context(), chi.URLParam(r, "token"), uid(r))
	h.respond(w, http.StatusOK, info, err)
}

func (h *Handler) acceptInvite(w http.ResponseWriter, r *http.Request) {
	storeID, err := h.svc.AcceptInvite(r.Context(), chi.URLParam(r, "token"), uid(r))
	h.respond(w, http.StatusOK, map[string]string{"store_id": storeID}, err)
}

func (h *Handler) respond(w http.ResponseWriter, status int, v any, err error) {
	if err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, status, v)
}

func (h *Handler) fail(w http.ResponseWriter, err error) {
	var ve validationError
	var he hours.FieldError
	switch {
	case errors.As(err, &ve):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": ve.field, "message": ve.msg}})
	case errors.As(err, &he):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": he.Field, "message": he.Msg}})
	case errors.Is(err, ErrNotMember):
		writeErr(w, http.StatusNotFound, "store_not_found", err.Error()) // jangan bocorkan toko orang lain ada
	case errors.Is(err, ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, ErrOutsideArea):
		writeErr(w, http.StatusUnprocessableEntity, "outside_service_area", err.Error())
	case errors.Is(err, ErrNotEditable):
		writeErr(w, http.StatusConflict, "store_under_review", err.Error())
	case errors.Is(err, ErrSlugTaken):
		writeErr(w, http.StatusConflict, "slug_taken", err.Error())
	case errors.Is(err, ErrStoreNotActive):
		writeErr(w, http.StatusConflict, "store_not_active", err.Error())
	case errors.Is(err, ErrProductBanned):
		writeErr(w, http.StatusConflict, "product_banned", err.Error())
	case errors.Is(err, media.ErrTooLong):
		writeErr(w, http.StatusUnprocessableEntity, "video_too_long", err.Error())
	case errors.Is(err, ErrInviteInvalid):
		writeErr(w, http.StatusGone, "invite_invalid", err.Error())
	case errors.Is(err, ErrTooManyStaff):
		writeErr(w, http.StatusConflict, "too_many_staff", err.Error())
	case errors.Is(err, ErrCannotRemove):
		writeErr(w, http.StatusConflict, "cannot_remove_owner", err.Error())
	case errors.Is(err, ErrOwnerOnly):
		writeErr(w, http.StatusForbidden, "owner_only", err.Error())
	case errors.Is(err, ErrTooManyStores):
		writeErr(w, http.StatusConflict, "too_many_stores", err.Error())
	case errors.Is(err, storage.ErrNotOwner):
		writeErr(w, http.StatusUnprocessableEntity, "upload_not_found", err.Error())
	case errors.Is(err, storage.ErrKind), errors.Is(err, storage.ErrType):
		writeErr(w, http.StatusUnsupportedMediaType, "upload_type", err.Error())
	case errors.Is(err, storage.ErrTooLarge):
		writeErr(w, http.StatusRequestEntityTooLarge, "upload_too_large", err.Error())
	default:
		h.log.Error("seller error", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "internal error")
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(io.LimitReader(r.Body, 256<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_json", err.Error())
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}
