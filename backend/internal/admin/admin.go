// Package admin: API untuk tim Phsar. Tahap ini: verifikasi toko, KYC, dan titik penjemputan.
// Akses dicek per izin (role_permissions); super_admin punya "*".
package admin

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/auth"
	"phsar/internal/hours"
	"phsar/internal/notify"
	"phsar/internal/platform/db"
	"phsar/internal/seller"
	"phsar/internal/storage"
)

type Handler struct {
	pool     *pgxpool.Pool
	sellers  *seller.Service
	files    storage.Storage
	log      *slog.Logger
	Refunds  *RefundRoutes   // refund manual (paket shop); nil = tidak dipasang
	Delivery *DeliveryRoutes // kurir & pengantaran
	Payouts  *PayoutRoutes   // payout seller
	Payments *PaymentRoutes  // payment gateway (super admin)
	// Integrations: provider SMS & Login Facebook (super admin)
	Integrations *IntegrationRoutes
	Legal        *LegalRoutes // dokumen legal & kontak (super admin)
	Returns      *ReturnRoutes
	Content      *ContentRoutes // content: banner, pengumuman, koleksi, SEO (super admin)
	Ads          *AdsRoutes     // ads: iklan bersponsor
	Askbot       http.Handler   // askbot: aktif/nonaktif bot jawaban + pemakaian (super admin)
	// courierreturn: antrean "Dikembalikan ke toko" (paket courierreturn); izin delivery.dispatch
	CourierReturns http.Handler
}

// DeliveryRoutes: kurir & pengantaran (paket delivery).
type DeliveryRoutes struct{ Couriers, SaveCourier, List, Assign, Proof http.HandlerFunc }

// PayoutRoutes: saldo & payout seller (paket payout).
type PayoutRoutes struct{ GetSettings, SaveSettings, List, Generate, Paid, Cancel http.HandlerFunc }

// ReturnRoutes: retur & sengketa (paket returns).
type ReturnRoutes struct{ List, Get, Decide, Refunded, Photo http.HandlerFunc }

// RefundRoutes: handler refund dari paket shop (supaya admin tidak bergantung langsung pada shop).
type RefundRoutes struct{ List, Done http.HandlerFunc }

func NewHandler(pool *pgxpool.Pool, sellers *seller.Service, files storage.Storage, log *slog.Logger) *Handler {
	return &Handler{pool: pool, sellers: sellers, files: files, log: log}
}

// Routes dipasang di /admin/api (sudah melewati auth.RequireUser).
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(h.requireAdmin)
	r.Get("/me", h.me)
	r.Get("/dashboard", h.dashboard)
	r.With(h.can("audit.view")).Get("/audit", h.audit)
	r.Route("/commission", func(p chi.Router) {
		p.Use(h.can("commission.manage"))
		p.Get("/", h.commissionView)
		p.Get("/stores", h.commissionStores)
		p.Post("/rules", h.commissionSet)
		p.Delete("/rules/{id}", h.commissionEnd)
	})
	// Menu admin tambahan (satu file per fitur)
	h.mountReports(r)
	h.mountSupport(r)
	h.mountChat(r) // chat: percakapan pembeli ↔ toko untuk sengketa
	h.mountReconciliation(r)
	h.mountSweep(r)
	h.mountInvoices(r)
	h.mountAds(r)
	h.mountPromos(r)
	h.mountCatalog(r)
	h.mountContent(r)
	h.mountShipping(r)
	h.mountJobs(r)
	h.mountExtras(r)     // extras: strike otomatis & tarik saldo manual (super admin)
	if h.Askbot != nil { // askbot: hanya super admin (izin "*")
		r.With(h.can("askbot.configure")).Mount("/askbot", h.Askbot)
	}
	// CS: pelanggan, pesanan, log pembayaran
	r.With(h.can("customer.view")).Get("/customers", h.listCustomers)
	r.With(h.can("customer.view")).Get("/customers/{id}", h.customerDetail)
	r.With(h.can("support.handle")).Post("/customers/{id}/notes", h.addCustomerNote)
	r.With(h.can("support.handle")).Post("/customers/{id}/status", h.setCustomerStatus)
	r.With(h.can("order.view")).Get("/orders", h.listOrders)
	r.With(h.can("order.view")).Get("/orders/{no}", h.orderDetail)
	r.With(h.can("order.view")).Get("/payment-events", h.listPaymentEvents)
	if h.Returns != nil {
		// CS memutuskan (return.decide); keuangan mengembalikan uang (refund.decide). Keduanya boleh melihat (order.view).
		r.Route("/returns", func(p chi.Router) {
			p.With(h.can("order.view")).Get("/", h.Returns.List)
			p.With(h.can("order.view")).Get("/{id}", h.Returns.Get)
			p.With(h.can("order.view")).Get("/{id}/photos/{i}", h.Returns.Photo)
			p.With(h.can("return.decide")).Post("/{id}/decide", h.Returns.Decide)
			p.With(h.can("refund.decide")).Post("/{id}/refunded", h.Returns.Refunded)
		})
	}
	if h.Legal != nil {
		r.Route("/legal", func(p chi.Router) {
			p.Use(h.can("legal.manage")) // hanya super admin (izin "*")
			p.Get("/", h.legalView)
			p.Put("/settings", h.legalSaveSettings)
			p.Post("/preview", h.legalPreview)
			p.Post("/docs/{type}", h.legalPublish)
		})
	}
	if h.Integrations != nil {
		r.Route("/integrations", func(p chi.Router) {
			p.Use(h.can("integration.configure")) // hanya super admin (izin "*")
			p.Get("/", h.integrationsView)
			p.Put("/sms", h.integrationsSaveSMS)
			p.Post("/sms/test", h.integrationsTestSMS)
			p.Put("/facebook", h.integrationsSaveFacebook)
			p.Post("/facebook/test", h.integrationsTestFacebook)
			p.Put("/ai", h.integrationsSaveAI) // AI: Anthropic utama, Gemini cadangan
			p.Post("/ai/test", h.integrationsTestAI)
			p.Put("/maps", h.integrationsSaveMaps)
			p.Get("/usage", h.integrationsUsage) // pemakaian & pengingat saldo
			p.Put("/usage/budget", h.integrationsSaveBudget)
		})
	}
	if h.Payments != nil {
		r.Route("/payments", func(p chi.Router) {
			p.Use(h.can("payment.configure")) // hanya super admin (izin "*")
			p.Get("/", h.paymentsView)
			p.Put("/settings", h.paymentsSaveGeneral)
			p.Put("/payway", h.paymentsSavePayWay)
			p.Post("/payway/test", h.paymentsTestPayWay)
		})
	}
	r.With(h.can("admin.manage")).Get("/roles", h.roles)
	r.With(h.can("admin.manage")).Get("/admins", h.listAdmins)
	r.With(h.can("admin.manage")).Post("/admins", h.addAdmin)
	r.With(h.can("admin.manage")).Patch("/admins/{userID}", h.updateAdmin)
	r.With(h.can("admin.manage")).Post("/admins/{userID}/mfa-reset", h.resetAdminMFA)
	r.With(h.can("seller.verify")).Get("/stores", h.listStores)
	r.With(h.can("seller.verify")).Get("/stores/{storeID}", h.storeDetail)
	r.With(h.can("seller.verify")).Post("/stores/{storeID}/kyc", h.decideKYC)
	r.With(h.can("seller.verify")).Post("/stores/{storeID}/decision", h.decideStore)
	r.With(h.can("pickup_point.verify")).Post("/pickup-points/{pointID}/decision", h.decidePickup)
	r.With(h.can("seller.verify")).Get("/files/*", h.privateFile)
	r.With(h.can("product.moderate")).Get("/products", h.listProducts)
	r.With(h.can("product.moderate")).Get("/products/{productID}", h.productDetail)
	r.With(h.can("product.moderate")).Post("/products/{productID}/decision", h.decideProduct)
	r.With(h.can("translation.review")).Put("/products/{productID}/translations/{lang}", h.reviewTranslation)
	if h.Refunds != nil {
		r.With(h.can("refund.decide")).Get("/refunds", h.Refunds.List)
		r.With(h.can("refund.decide")).Post("/refunds/{orderID}/done", h.Refunds.Done)
	}
	r.With(h.can("order.rules")).Get("/order-timing", h.getOrderTiming) // hanya super admin (izin "*")
	r.With(h.can("order.rules")).Put("/order-timing", h.saveOrderTiming)
	// pricehistory: aturan harga flash (harga terendah 30 hari) — hanya super admin (izin "*")
	r.With(h.can("flash.rules")).Get("/flash-rules", h.phGetRules)
	r.With(h.can("flash.rules")).Put("/flash-rules", h.phSaveRules)
	if h.Payouts != nil {
		p := r.With(h.can("payout.approve"))
		p.Get("/payout-settings", h.Payouts.GetSettings)
		p.Put("/payout-settings", h.Payouts.SaveSettings)
		p.Get("/payouts", h.Payouts.List)
		p.Post("/payouts/generate", h.Payouts.Generate)
		p.Post("/payouts/{id}/paid", h.Payouts.Paid)
		p.Post("/payouts/{id}/cancel", h.Payouts.Cancel)
	}
	if h.Delivery != nil {
		d := r.With(h.can("delivery.dispatch"))
		d.Get("/couriers", h.Delivery.Couriers)
		d.Post("/couriers", h.Delivery.SaveCourier)
		d.Get("/deliveries", h.Delivery.List)
		d.Post("/deliveries/{id}/assign", h.Delivery.Assign)
		d.Get("/deliveries/{id}/proof", h.Delivery.Proof)
	}
	if h.CourierReturns != nil { // courierreturn
		r.With(h.can("delivery.dispatch")).Mount("/courier-returns", h.CourierReturns)
	}
	return r
}

type adminCtx struct {
	userID string
	role   string
	perms  map[string]bool
}

type ctxKey struct{}

func (h *Handler) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := auth.UserFrom(r.Context())
		a := adminCtx{userID: c.Subject, perms: map[string]bool{}}
		err := h.pool.QueryRow(r.Context(), `SELECT role_code FROM admin_users WHERE user_id = $1 AND active AND deleted_at IS NULL`, c.Subject).Scan(&a.role)
		if errors.Is(err, pgx.ErrNoRows) {
			writeErr(w, http.StatusForbidden, "not_admin", "akun ini bukan admin Phsar")
			return
		}
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "internal", err.Error())
			return
		}
		rows, err := h.pool.Query(r.Context(), `SELECT permission FROM role_permissions WHERE role_code = $1`, a.role)
		if err == nil {
			for rows.Next() {
				var p string
				if rows.Scan(&p) == nil {
					a.perms[p] = true
				}
			}
			rows.Close()
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, a)))
	})
}

func (h *Handler) can(perm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			a := r.Context().Value(ctxKey{}).(adminCtx)
			if !a.perms["*"] && !a.perms[perm] {
				writeErr(w, http.StatusForbidden, "forbidden", "peran "+a.role+" tidak punya izin "+perm)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func actor(r *http.Request) string { return r.Context().Value(ctxKey{}).(adminCtx).userID }

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	a := r.Context().Value(ctxKey{}).(adminCtx)
	perms := []string{}
	for p := range a.perms {
		perms = append(perms, p)
	}
	var name, phone string
	var backup int
	_ = h.pool.QueryRow(r.Context(), `SELECT coalesce(u.name, ''), coalesce(u.phone, ''),
		(SELECT count(*) FROM admin_mfa_backup_codes b WHERE b.user_id = u.id AND b.used_at IS NULL)
		FROM users u WHERE u.id = $1`, a.userID).Scan(&name, &phone, &backup)
	writeJSON(w, http.StatusOK, map[string]any{"user_id": a.userID, "role": a.role, "permissions": perms,
		"name": name, "phone": phone, "backup_codes_left": backup})
}

type storeRow struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Status      string     `json:"status"`
	OwnerName   string     `json:"owner_name"`
	OwnerPhone  string     `json:"owner_phone"`
	KYC         string     `json:"kyc"`
	Pickups     int        `json:"pickup_points"`
	SubmittedAt *time.Time `json:"submitted_at"`
}

// GET /admin/api/stores?status=pending_review (default) — antrean verifikasi, terlama dulu.
func (h *Handler) listStores(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "pending_review"
	}
	rows, err := h.pool.Query(r.Context(), `
		SELECT st.id, st.name, st.slug, st.status::text, coalesce(u.name, ''), coalesce(u.phone, ''),
		       coalesce(tp.verification_status, 'missing'),
		       (SELECT count(*) FROM store_pickup_points p WHERE p.store_id = st.id AND p.deleted_at IS NULL),
		       st.submitted_at
		FROM stores st
		LEFT JOIN store_members sm ON sm.store_id = st.id AND sm.role = 'owner'
		LEFT JOIN users u ON u.id = sm.user_id
		LEFT JOIN seller_tax_profiles tp ON tp.store_id = st.id
		WHERE st.deleted_at IS NULL AND st.status::text = $1
		ORDER BY st.submitted_at NULLS LAST, st.created_at
		LIMIT 200`, status)
	if err != nil {
		h.fail(w, err)
		return
	}
	defer rows.Close()
	out := []storeRow{}
	for rows.Next() {
		var s storeRow
		if err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.Status, &s.OwnerName, &s.OwnerPhone, &s.KYC, &s.Pickups, &s.SubmittedAt); err != nil {
			h.fail(w, err)
			return
		}
		out = append(out, s)
	}
	writeJSON(w, http.StatusOK, map[string]any{"stores": out})
}

func (h *Handler) storeDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "storeID")
	st, err := h.sellers.GetStore(r.Context(), id)
	if err != nil {
		h.fail(w, err)
		return
	}
	kyc, err := h.sellers.GetKYC(r.Context(), id)
	if err != nil {
		h.fail(w, err)
		return
	}
	var docKey *string
	_ = h.pool.QueryRow(r.Context(), `SELECT id_document_key FROM seller_tax_profiles WHERE store_id = $1`, id).Scan(&docKey)
	docURL := ""
	if docKey != nil {
		docURL = "/admin/api/files/" + *docKey
	}
	points, err := h.sellers.ListPickupPoints(r.Context(), id)
	if err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"store": st, "kyc": kyc, "id_document_url": docURL, "pickup_points": points})
}

type decision struct {
	Decision string `json:"decision"` // approve · reject
	Note     string `json:"note"`     // wajib saat reject; terlihat oleh seller
}

func (d decision) valid() bool {
	return d.Decision == "approve" || (d.Decision == "reject" && strings.TrimSpace(d.Note) != "")
}

// POST /admin/api/stores/{id}/kyc {decision, note}
func (h *Handler) decideKYC(w http.ResponseWriter, r *http.Request) {
	var d decision
	if !decode(w, r, &d) || !d.valid() {
		writeErr(w, http.StatusUnprocessableEntity, "validation", "decision approve/reject; note wajib saat reject")
		return
	}
	id := chi.URLParam(r, "storeID")
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		status := "verified"
		if d.Decision == "reject" {
			status = "rejected"
		}
		tag, err := tx.Exec(r.Context(), `UPDATE seller_tax_profiles SET verification_status = $2::text, rejection_note = $3,
			verified_by = CASE WHEN $2::text = 'verified' THEN $4::uuid END, verified_at = CASE WHEN $2::text = 'verified' THEN now() END
			WHERE store_id = $1`, id, status, nilIfBlank(d.Note), actor(r))
		if err == nil && tag.RowsAffected() == 0 {
			return seller.ErrNotFound
		}
		if err != nil || status != "rejected" {
			return err
		}
		return notifyStore(r.Context(), tx, id, notify.KYCRejected, notify.Params{Note: d.Note}, "#/store/"+id)
	})
	if err != nil {
		h.fail(w, err)
		return
	}
	h.storeDetail(w, r)
}

// POST /admin/api/pickup-points/{id}/decision {decision, note}
func (h *Handler) decidePickup(w http.ResponseWriter, r *http.Request) {
	var d decision
	if !decode(w, r, &d) || !d.valid() {
		writeErr(w, http.StatusUnprocessableEntity, "validation", "decision approve/reject; note wajib saat reject")
		return
	}
	id := chi.URLParam(r, "pointID")
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		var tag interface{ RowsAffected() int64 }
		var err error
		if d.Decision == "approve" {
			tag, err = tx.Exec(r.Context(), `UPDATE store_pickup_points SET verified_at = now(), verified_by = $2, rejection_note = NULL
				WHERE id = $1 AND deleted_at IS NULL`, id, actor(r))
		} else {
			tag, err = tx.Exec(r.Context(), `UPDATE store_pickup_points SET verified_at = NULL, verified_by = NULL, rejection_note = $2
				WHERE id = $1 AND deleted_at IS NULL`, id, d.Note)
		}
		if err == nil && tag.RowsAffected() == 0 {
			return seller.ErrNotFound
		}
		if err != nil {
			return err
		}
		var storeID, label string
		if err := tx.QueryRow(r.Context(), `SELECT store_id, label FROM store_pickup_points WHERE id = $1`, id).Scan(&storeID, &label); err != nil {
			return err
		}
		kind := notify.PickupApproved
		if d.Decision != "approve" {
			kind = notify.PickupRejected
		}
		return notifyStore(r.Context(), tx, storeID, kind, notify.Params{Point: label, Note: d.Note}, "#/store/"+storeID)
	})
	if err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /admin/api/stores/{id}/decision {decision, note}
// approve: butuh KYC verified + minimal 1 titik penjemputan terverifikasi → status active.
// reject: kembali ke onboarding dengan catatan untuk seller.
func (h *Handler) decideStore(w http.ResponseWriter, r *http.Request) {
	var d decision
	if !decode(w, r, &d) || !d.valid() {
		writeErr(w, http.StatusUnprocessableEntity, "validation", "decision approve/reject; note wajib saat reject")
		return
	}
	id := chi.URLParam(r, "storeID")
	if d.Decision == "approve" {
		st, err := h.sellers.GetStore(r.Context(), id)
		if err != nil {
			h.fail(w, err)
			return
		}
		if st.Checklist.KYC != "verified" || st.Checklist.PickupVerified == 0 {
			writeErr(w, http.StatusUnprocessableEntity, "not_ready", "verifikasi KYC dan minimal satu titik penjemputan dulu")
			return
		}
	}
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		var tag interface{ RowsAffected() int64 }
		var err error
		if d.Decision == "approve" {
			tag, err = tx.Exec(r.Context(), `UPDATE stores SET status = 'active', verified_at = now(), verified_by = $2, review_note = NULL
				WHERE id = $1 AND status = 'pending_review'`, id, actor(r))
		} else {
			tag, err = tx.Exec(r.Context(), `UPDATE stores SET status = 'onboarding', review_note = $2
				WHERE id = $1 AND status = 'pending_review'`, id, d.Note)
		}
		if err == nil && tag.RowsAffected() == 0 {
			return errNotPending
		}
		if err != nil {
			return err
		}
		kind := notify.StoreApproved
		if d.Decision != "approve" {
			kind = notify.StoreRejected
		}
		return notifyStore(r.Context(), tx, id, kind, notify.Params{Note: d.Note}, "#/store/"+id)
	})
	if err != nil {
		h.fail(w, err)
		return
	}
	h.storeDetail(w, r)
}

var errNotPending = errors.New("toko tidak sedang menunggu verifikasi")

// notifyStore mengisi nama toko lalu membuat notifikasi untuk anggota toko (dalam tx keputusan admin).
func notifyStore(ctx context.Context, tx pgx.Tx, storeID, kind string, p notify.Params, link string) error {
	if err := tx.QueryRow(ctx, `SELECT name FROM stores WHERE id = $1`, storeID).Scan(&p.Store); err != nil {
		return err
	}
	return notify.Enqueue(ctx, tx, storeID, kind, p, link)
}

// ---- Moderasi produk ----

type productRow struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	StoreName   string          `json:"store_name"`
	StoreID     string          `json:"store_id"`
	CategoryID  string          `json:"category_id"`
	Condition   string          `json:"condition"`
	Moderation  string          `json:"moderation_status"`
	CoverURL    string          `json:"cover_url"`
	MinPrice    *string         `json:"min_price"`
	Flags       int             `json:"flags"`
	Reasons     json.RawMessage `json:"reasons"` // alasan jalur manual (modrules.Reason*)
	Translating bool            `json:"translating"`
	SubmittedAt *time.Time      `json:"submitted_at"`
}

// GET /admin/api/products?status=pending_review — antrean moderasi, terlama dulu.
func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	where, order := "p.moderation_status::text = $1", "p.submitted_at NULLS LAST, p.updated_at"
	switch status {
	case "pending_review", "approved", "changes_requested", "rejected":
	case "auto": // tayang otomatis, belum dilihat admin — terbaru dulu
		where, order = "p.moderation_status = 'approved' AND p.post_review_pending AND $1 <> ''", "p.published_at DESC NULLS LAST"
	default:
		status = "pending_review"
	}
	rows, err := h.pool.Query(r.Context(), `
		SELECT p.id, coalesce(t.name, ''), st.name, st.id, p.category_id, p.condition::text, p.moderation_status::text,
		       coalesce((SELECT coalesce(m.image_url, m.storage_key) FROM product_media m WHERE m.product_id = p.id AND m.deleted_at IS NULL AND m.kind = 'photo'
		                 ORDER BY m.sort_order LIMIT 1), ''),
		       p.min_price::text, coalesce(jsonb_array_length(p.auto_flags->'keywords'), 0), coalesce(p.auto_flags->'reasons', '[]'),
		       EXISTS (SELECT 1 FROM translation_jobs j WHERE j.entity_table = 'product_translations' AND j.entity_id = p.id::text
		               AND j.status IN ('queued','running')),
		       p.submitted_at
		FROM products p JOIN stores st ON st.id = p.store_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE p.deleted_at IS NULL AND p.status <> 'draft' AND `+where+`
		ORDER BY `+order+`
		LIMIT 200`, status)
	if err != nil {
		h.fail(w, err)
		return
	}
	defer rows.Close()
	out := []productRow{}
	for rows.Next() {
		var p productRow
		var cover string
		if err := rows.Scan(&p.ID, &p.Name, &p.StoreName, &p.StoreID, &p.CategoryID, &p.Condition, &p.Moderation, &cover,
			&p.MinPrice, &p.Flags, &p.Reasons, &p.Translating, &p.SubmittedAt); err != nil {
			h.fail(w, err)
			return
		}
		p.CoverURL = cover
		if !strings.HasPrefix(cover, "http") {
			p.CoverURL = h.files.PublicURL(cover)
		}
		out = append(out, p)
	}
	writeJSON(w, http.StatusOK, map[string]any{"products": out})
}

type reviewRow struct {
	Decision   string    `json:"decision"`
	ReasonCode *string   `json:"reason_code"`
	Note       *string   `json:"note"`
	Reviewer   string    `json:"reviewer"`
	CreatedAt  time.Time `json:"created_at"`
}

func (h *Handler) productDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "productID")
	p, err := h.sellers.GetProduct(r.Context(), "", id)
	if err != nil {
		h.fail(w, err)
		return
	}
	st, err := h.sellers.GetStore(r.Context(), p.StoreID)
	if err != nil {
		h.fail(w, err)
		return
	}
	st.Checklist = nil
	rows, err := h.pool.Query(r.Context(), `
		SELECT r.decision::text, r.reason_code, r.note, coalesce(u.name, u.phone, ''), r.created_at
		FROM product_moderation_reviews r LEFT JOIN users u ON u.id = r.reviewer_id
		WHERE r.product_id = $1 ORDER BY r.created_at DESC LIMIT 50`, id)
	if err != nil {
		h.fail(w, err)
		return
	}
	defer rows.Close()
	history := []reviewRow{}
	for rows.Next() {
		var rv reviewRow
		if err := rows.Scan(&rv.Decision, &rv.ReasonCode, &rv.Note, &rv.Reviewer, &rv.CreatedAt); err != nil {
			h.fail(w, err)
			return
		}
		history = append(history, rv)
	}
	writeJSON(w, http.StatusOK, map[string]any{"product": p, "store": st, "history": history})
}

var reasonCodes = map[string]bool{"defect_photo_missing": true, "unclear_condition": true, "prohibited": true,
	"counterfeit": true, "wrong_category": true, "price": true, "other": true}

// POST /admin/api/products/{id}/decision {decision: approve|changes_requested|reject, reason_code, note}
// changes_requested: seller memperbaiki lalu kirim ulang. reject: dilarang dijual (barang terlarang, palsu).
func (h *Handler) decideProduct(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Decision   string `json:"decision"`
		ReasonCode string `json:"reason_code"`
		Note       string `json:"note"`
	}
	if !decode(w, r, &d) {
		writeErr(w, http.StatusBadRequest, "bad_json", "JSON tidak valid")
		return
	}
	d.Note = strings.TrimSpace(d.Note)
	switch {
	case d.Decision == "approve":
		d.ReasonCode = ""
	case d.Decision == "changes_requested" || d.Decision == "reject":
		if d.Note == "" || !reasonCodes[d.ReasonCode] {
			writeErr(w, http.StatusUnprocessableEntity, "validation", "pilih alasan dan tulis catatan untuk seller")
			return
		}
	default:
		writeErr(w, http.StatusUnprocessableEntity, "validation", "decision: approve, changes_requested, atau reject")
		return
	}
	id := chi.URLParam(r, "productID")
	decision := map[string]string{"approve": "approved", "changes_requested": "changes_requested", "reject": "rejected"}[d.Decision]
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		var storeStatus, status string
		var flags []byte
		err := tx.QueryRow(r.Context(), `SELECT st.status::text, p.status::text, p.auto_flags FROM products p JOIN stores st ON st.id = p.store_id
			WHERE p.id = $1 AND p.deleted_at IS NULL FOR UPDATE OF p`, id).Scan(&storeStatus, &status, &flags)
		if errors.Is(err, pgx.ErrNoRows) {
			return seller.ErrNotFound
		}
		if err != nil {
			return err
		}
		if status == "draft" {
			return errProductDraft
		}
		if d.Decision == "approve" && storeStatus != "active" {
			return errStoreInactive
		}
		if _, err := tx.Exec(r.Context(), `
			UPDATE products SET moderation_status = $2::moderation_status, moderation_note = $3,
				published_at = CASE WHEN $2 = 'approved' THEN coalesce(published_at, now()) ELSE published_at END,
				robots_index = ($2 = 'approved'),
				approved_min_price = CASE WHEN $2 = 'approved' THEN min_price ELSE approved_min_price END,
				post_review_pending = false
			WHERE id = $1`, id, decision, nilIfBlank(d.Note)); err != nil {
			return err
		}
		if _, err = tx.Exec(r.Context(), `INSERT INTO product_moderation_reviews (product_id, decision, reason_code, note, auto_flags, reviewer_id)
			VALUES ($1, $2::moderation_status, $3, $4, $5, $6)`, id, decision, nilIfBlank(d.ReasonCode), nilIfBlank(d.Note), flags, actor(r)); err != nil {
			return err
		}
		var storeID, name string
		if err := tx.QueryRow(r.Context(), `SELECT p.store_id, coalesce(t.name, '') FROM products p
			LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang WHERE p.id = $1`, id).Scan(&storeID, &name); err != nil {
			return err
		}
		kind := map[string]string{"approved": notify.ProductApproved, "changes_requested": notify.ProductChangesRequired, "rejected": notify.ProductRejected}[decision]
		return notifyStore(r.Context(), tx, storeID, kind, notify.Params{Product: name, Note: d.Note}, "#/store/"+storeID+"/product/"+id)
	})
	if err != nil {
		h.fail(w, err)
		return
	}
	h.productDetail(w, r)
}

var (
	errProductDraft  = errors.New("produk masih draf, belum dikirim seller")
	errStoreInactive = errors.New("toko belum aktif; setujui tokonya dulu")
)

// PUT /admin/api/products/{id}/translations/{lang} — tim Phsar memperbaiki terjemahan AI.
func (h *Handler) reviewTranslation(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if !decode(w, r, &in) || strings.TrimSpace(in.Name) == "" || len([]rune(in.Name)) > 160 {
		writeErr(w, http.StatusUnprocessableEntity, "validation", "nama 1–160 karakter")
		return
	}
	id, lang := chi.URLParam(r, "productID"), chi.URLParam(r, "lang")
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(r.Context(), `
			UPDATE product_translations t SET name = $3, description = $4, source = 'reviewer', review_status = 'reviewed',
				reviewed_by = $5, reviewed_at = now()
			FROM products p WHERE p.id = t.product_id AND t.product_id = $1 AND t.lang::text = $2 AND t.lang <> p.source_lang`,
			id, lang, strings.TrimSpace(in.Name), strings.TrimSpace(in.Description), actor(r))
		if err == nil && tag.RowsAffected() == 0 {
			return seller.ErrNotFound
		}
		return err
	})
	if err != nil {
		h.fail(w, err)
		return
	}
	h.productDetail(w, r)
}

// GET /admin/api/files/private/kyc/… — dokumen KYC hanya bisa dibuka admin berizin.
func (h *Handler) privateFile(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "*")
	var ctype string
	if err := h.pool.QueryRow(r.Context(), `SELECT content_type FROM uploads WHERE key = $1`, key).Scan(&ctype); err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "file tidak ditemukan")
		return
	}
	f, err := h.files.Open(r.Context(), key)
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "file tidak ditemukan")
		return
	}
	defer f.Close()
	_, _ = h.pool.Exec(r.Context(), `INSERT INTO audit_logs (actor_type, actor_user_id, action, entity_table, entity_id)
		VALUES ('admin', $1, 'kyc_document.view', 'uploads', $2)`, actor(r), key)
	w.Header().Set("Content-Type", ctype)
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = io.Copy(w, f)
}

func (h *Handler) fail(w http.ResponseWriter, err error) {
	var he hours.FieldError
	switch {
	case errors.As(err, &he):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": he.Field, "message": he.Msg}})
	case errors.Is(err, seller.ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, errNotPending):
		writeErr(w, http.StatusConflict, "not_pending", err.Error())
	case errors.Is(err, errProductDraft):
		writeErr(w, http.StatusConflict, "product_draft", err.Error())
	case errors.Is(err, errStoreInactive):
		writeErr(w, http.StatusConflict, "store_not_active", err.Error())
	default:
		h.log.Error("admin error", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "internal error")
	}
}

func nilIfBlank(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(io.LimitReader(r.Body, 16<<10)).Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_json", "invalid JSON body")
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
