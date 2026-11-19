package shop

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
	"phsar/internal/i18n"
	"phsar/internal/paging"
	"phsar/internal/payment"
)

type Handler struct {
	S   *Service
	Log *slog.Logger
	// Extra: rute tambahan untuk pembeli yang sudah login (mis. retur).
	Extra func(chi.Router)
}

// Routes untuk /v1/shop. requireUser = middleware login pembeli.
func (h *Handler) Routes(requireUser func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()
	r.Post("/cart/quote", h.quote)
	r.Group(func(r chi.Router) {
		r.Use(requireUser)
		r.Get("/address", h.address)
		r.Post("/checkouts", h.checkout)
		r.Get("/checkouts/{id}", h.getCheckout)
		r.Post("/checkouts/{id}/cancel", h.cancel)
		r.Get("/orders", h.listOrders)
		r.Get("/orders/{no}", h.getOrder)
		if h.Extra != nil {
			h.Extra(r) // retur (paket returns)
		}
	})
	return r
}

func lang(r *http.Request) string { return string(i18n.From(r.Context())) }

func buyer(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

func (h *Handler) quote(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Items       []CartItem `json:"items"`
		VoucherCode string     `json:"voucher_code"` // promos
	}
	if !decode(w, r, &in) {
		return
	}
	q, err := h.S.QuoteWithVoucher(r.Context(), lang(r), in.Items, in.VoucherCode)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, q)
}

func (h *Handler) address(w http.ResponseWriter, r *http.Request) {
	a, err := h.S.DefaultAddress(r.Context(), buyer(r))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"address": a})
}

func (h *Handler) checkout(w http.ResponseWriter, r *http.Request) {
	var in CheckoutInput
	if !decode(w, r, &in) {
		return
	}
	v, err := h.S.Checkout(r.Context(), buyer(r), lang(r), in)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}

func (h *Handler) getCheckout(w http.ResponseWriter, r *http.Request) {
	v, err := h.S.Poll(r.Context(), buyer(r), chi.URLParam(r, "id"))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	v, err := h.S.CancelByBuyer(r.Context(), buyer(r), chi.URLParam(r, "id"))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) listOrders(w http.ResponseWriter, r *http.Request) {
	list, err := h.S.ListOrders(r.Context(), buyer(r))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"orders": list})
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	o, err := h.S.GetOrder(r.Context(), buyer(r), chi.URLParam(r, "no"))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

// Callback ABA PayWay: POST /webhooks/payway. Isi callback hanya dicatat; status diambil dari API resmi.
func (h *Handler) PayWayCallback(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(io.LimitReader(r.Body, 64<<10))
	var m map[string]any
	if json.Unmarshal(body, &m) != nil {
		_ = r.ParseForm()
		m = map[string]any{}
		for k := range r.PostForm {
			m[k] = r.PostForm.Get(k)
		}
	}
	tranID, _ := m["tran_id"].(string)
	verified := h.S.Pay.VerifyCallback(r.Header, body)
	var cid string
	_ = h.S.Pool.QueryRow(r.Context(), `SELECT id FROM checkouts WHERE tran_id = $1`, tranID).Scan(&cid)
	status, _ := json.Marshal(m["status"])
	h.S.event(r.Context(), cid, truncRunes(tranID, 20), "callback", verified, truncRunes(string(status), 20), m)
	if cid == "" {
		writeJSON(w, http.StatusOK, map[string]string{"result": "unknown transaction"})
		return
	}
	if _, err := h.S.Refresh(r.Context(), tranID, "callback"); err != nil {
		h.Log.Warn("callback: cek transaksi gagal (worker akan mencoba lagi)", "tran_id", tranID, "err", err)
	}
	writeJSON(w, http.StatusOK, map[string]string{"result": "ok"})
}

// DevPay (hanya APP_ENV=development + PAYMENT_PROVIDER=fake): tombol "simulasi bayar" di layar QR.
func (h *Handler) DevPay(w http.ResponseWriter, r *http.Request) {
	f, ok := h.S.Pay.(*payment.Fake)
	if fp, isGW := h.S.Pay.(interface {
		FakeProvider(context.Context) (*payment.Fake, bool)
	}); isGW {
		f, ok = fp.FakeProvider(r.Context())
	}
	if !ok {
		writeErr(w, http.StatusNotFound, "not_found", "only when the payment mode is fake (test payments)")
		return
	}
	var tranID, total string
	if err := h.S.Pool.QueryRow(r.Context(), `SELECT tran_id, total::text FROM checkouts WHERE id = $1 AND buyer_id = $2`,
		chi.URLParam(r, "id"), buyer(r)).Scan(&tranID, &total); err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "checkout not found")
		return
	}
	f.Pay(tranID, total)
	if _, err := h.S.Refresh(r.Context(), tranID, "fake"); err != nil {
		h.fail(w, r, err)
		return
	}
	v, err := h.S.GetCheckout(r.Context(), buyer(r), chi.URLParam(r, "id"))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

var notServedMsg = map[string]string{
	"en": "We currently deliver in Sihanoukville only.",
	"km": "បច្ចុប្បន្នយើងដឹកជញ្ជូនតែនៅក្រុងព្រះសីហនុប៉ុណ្ណោះ។",
	"zh": "目前我们只在西哈努克城配送。",
}

func (h *Handler) fail(w http.ResponseWriter, r *http.Request, err error) {
	var ve validationError
	var ce *CartChangedError
	var vce *VoucherError // promos
	switch {
	case errors.As(err, &vce):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": vce.Code, "field": "voucher_code",
			"message": voucherMsg(vce, lang(r))}})
	case errors.As(err, &ve):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": ve.field, "message": ve.msg}})
	case errors.As(err, &ce):
		writeJSON(w, http.StatusConflict, map[string]any{"error": map[string]string{"code": "cart_changed", "message": "Some items changed. Please review your cart."}, "quote": ce.Quote})
	case errors.Is(err, ErrNotServed):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "not_served", "field": "address.address_line", "message": notServedMsg[lang(r)]}})
	case errors.Is(err, ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", "not found")
	case errors.Is(err, ErrClosed):
		writeErr(w, http.StatusConflict, "checkout_closed", closedMsg[lang(r)])
	case errors.Is(err, ErrPayment):
		h.Log.Error("pembayaran", "err", err)
		writeErr(w, http.StatusBadGateway, "payment_unavailable", "Payment is temporarily unavailable. Please try again in a few minutes.")
	default:
		h.Log.Error("shop", "path", r.URL.Path, "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "something went wrong")
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(io.LimitReader(r.Body, 128<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_json", err.Error())
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}

// ---------------------------------------------------------------------------
// Seller Center: /v1/seller/stores/{storeID}/orders (dipasang di dalam route toko yang sudah cek anggota)
// ---------------------------------------------------------------------------

func (h *Handler) SellerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// ?filter=&q=&from=YYYY-MM-DD&to=YYYY-MM-DD&page=&size= (tanpa page: maks. 200 — aplikasi lama)
		q := SellerOrderQuery{Filter: r.URL.Query().Get("filter"), Page: paging.FromRequest(r, 20, 100)}
		q.From, q.To = dayRange(r.URL.Query().Get("from"), r.URL.Query().Get("to"))
		l, err := h.S.SellerListOrders(r.Context(), chi.URLParam(r, "storeID"), q)
		h.reply(w, r, l, err)
	})
	r.Get("/{no}", func(w http.ResponseWriter, r *http.Request) {
		o, err := h.S.SellerGetOrder(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "no"))
		h.reply(w, r, o, err)
	})
	r.Post("/{no}/accept", func(w http.ResponseWriter, r *http.Request) {
		o, err := h.S.Accept(r.Context(), chi.URLParam(r, "storeID"), buyer(r), chi.URLParam(r, "no"))
		h.reply(w, r, o, err)
	})
	r.Post("/{no}/reject", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Reason string `json:"reason"`
			Note   string `json:"note"`
		}
		if !decode(w, r, &in) {
			return
		}
		o, err := h.S.Reject(r.Context(), chi.URLParam(r, "storeID"), buyer(r), chi.URLParam(r, "no"), in.Reason, in.Note)
		h.reply(w, r, o, err)
	})
	r.Post("/{no}/ready", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			PickupPointID string `json:"pickup_point_id"`
		}
		if r.ContentLength > 0 && !decode(w, r, &in) {
			return
		}
		o, err := h.S.Ready(r.Context(), chi.URLParam(r, "storeID"), buyer(r), chi.URLParam(r, "no"), in.PickupPointID)
		h.reply(w, r, o, err)
	})
	return r
}

// Admin: refund manual. adminID diambil dari token (izin dicek oleh router admin).
func (h *Handler) AdminListRefunds(w http.ResponseWriter, r *http.Request) {
	l, err := h.S.ListRefunds(r.Context(), r.URL.Query().Get("status"))
	h.reply(w, r, map[string]any{"refunds": l}, err)
}

func (h *Handler) AdminMarkRefunded(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Ref  string `json:"refund_ref"`
		Note string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	err := h.S.MarkRefunded(r.Context(), buyer(r), chi.URLParam(r, "orderID"), in.Ref, in.Note)
	h.reply(w, r, map[string]string{"result": "ok"}, err)
}

func (h *Handler) reply(w http.ResponseWriter, r *http.Request, v any, err error) {
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, v)
	case errors.Is(err, ErrBadTransition):
		writeErr(w, http.StatusConflict, "bad_status", "This order can no longer be changed this way. Refresh the page.")
	case errors.Is(err, ErrNoPickupPoint):
		writeErr(w, http.StatusConflict, "no_pickup_point", "Add a pickup point first.")
	default:
		h.fail(w, r, err)
	}
}

var closedMsg = map[string]string{
	"en": "Phsar is opening soon. Shopping is not open yet.",
	"km": "ផ្សារនឹងបើកឆាប់ៗនេះ។ មិនទាន់អាចទិញបានទេ។",
	"zh": "Phsar 即将开业，暂时还不能下单。",
}

// dayRange: tanggal "YYYY-MM-DD" (zona Phnom Penh) → [awal hari from, awal hari setelah to). Kosong/salah = tanpa batas.
func dayRange(from, to string) (*time.Time, *time.Time) {
	loc, err := time.LoadLocation("Asia/Phnom_Penh")
	if err != nil {
		loc = time.FixedZone("ICT", 7*3600)
	}
	var f, t *time.Time
	if d, err := time.ParseInLocation("2006-01-02", from, loc); err == nil {
		f = &d
	}
	if d, err := time.ParseInLocation("2006-01-02", to, loc); err == nil {
		d = d.AddDate(0, 0, 1)
		t = &d
	}
	return f, t
}
