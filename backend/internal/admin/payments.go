package admin

// Payment gateway: super admin memasukkan kredensial (PayWay) dan memilih mode (tutup / PayWay / uji) dari admin.
// Menyimpan perubahan wajib kode authenticator lagi (step-up), karena kredensial ini mengatur ke mana uang pembeli masuk.

import (
	"errors"
	"net/http"

	"phsar/internal/auth"
	"phsar/internal/paygw"
)

type PaymentRoutes struct {
	Store *paygw.Store
	// ConfirmMFA: cek kode authenticator admin yang sedang login.
	ConfirmMFA func(r *http.Request, userID, code string) error
}

func (h *Handler) paymentsView(w http.ResponseWriter, r *http.Request) {
	v, err := h.Payments.Store.View(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) paymentsSaveGeneral(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Mode        string `json:"mode"`
		LifetimeMin int    `json:"lifetime_min"`
		MFACode     string `json:"mfa_code"`
	}
	if !decode(w, r, &in) || !stepUp(w, r, in.MFACode, h.Payments.ConfirmMFA) {
		return
	}
	h.paymentsReply(w, r, h.Payments.Store.SaveGeneral(r.Context(), actor(r), in.Mode, in.LifetimeMin))
}

func (h *Handler) paymentsSavePayWay(w http.ResponseWriter, r *http.Request) {
	var in struct {
		paygw.PayWayInput
		MFACode string `json:"mfa_code"`
	}
	if !decode(w, r, &in) || !stepUp(w, r, in.MFACode, h.Payments.ConfirmMFA) {
		return
	}
	h.paymentsReply(w, r, h.Payments.Store.SavePayWay(r.Context(), actor(r), in.PayWayInput))
}

func (h *Handler) paymentsTestPayWay(w http.ResponseWriter, r *http.Request) {
	res, err := h.Payments.Store.TestPayWay(r.Context(), actor(r))
	if err != nil {
		h.paymentsReply(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// stepUp: minta kode authenticator lagi sebelum menyimpan kredensial (pembayaran, SMS, Facebook).
func stepUp(w http.ResponseWriter, r *http.Request, code string, confirm func(r *http.Request, userID, code string) error) bool {
	err := confirm(r, actor(r), code)
	switch {
	case err == nil:
		return true
	case errors.Is(err, auth.ErrMFAInvalid):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "mfa_invalid", "field": "mfa_code", "message": err.Error()}})
	case errors.Is(err, auth.ErrMFALocked):
		writeErr(w, http.StatusTooManyRequests, "mfa_locked", err.Error())
	default:
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
	}
	return false
}

func (h *Handler) paymentsReply(w http.ResponseWriter, r *http.Request, err error) {
	var fe paygw.FieldError
	switch {
	case err == nil:
		h.paymentsView(w, r)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	default:
		h.log.Error("pengaturan pembayaran", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save payment settings")
	}
}
