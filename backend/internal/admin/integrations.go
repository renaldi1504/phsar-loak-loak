package admin

// Login & SMS: super admin memasukkan kredensial PlasGate (SMS OTP) dan Facebook (Login Facebook).
// Menyimpan wajib kode authenticator; rahasia tidak pernah dikirim balik ke browser.

import (
	"errors"
	"net/http"

	"phsar/internal/auth"
	"phsar/internal/integ"
)

type IntegrationRoutes struct {
	Store      *integ.Store
	ConfirmMFA func(r *http.Request, userID, code string) error
	TestAI     integ.Tester // uji kunci Anthropic / Gemini (terjemahan satu kata)
}

// PUT /integrations/ai {provider, enabled, model, key?, clear?, mfa_code} — Anthropic (utama) / Gemini (cadangan).
func (h *Handler) integrationsSaveAI(w http.ResponseWriter, r *http.Request) {
	var in struct {
		integ.AIInput
		MFACode string `json:"mfa_code"`
	}
	if !decode(w, r, &in) || !stepUp(w, r, in.MFACode, h.Integrations.ConfirmMFA) {
		return
	}
	h.integrationsReply(w, r, h.Integrations.Store.SaveAI(r.Context(), actor(r), in.AIInput))
}

// POST /integrations/ai/test {provider} — terjemahkan satu kata dengan kunci yang berlaku.
func (h *Handler) integrationsTestAI(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Provider string `json:"provider"`
	}
	if !decode(w, r, &in) {
		return
	}
	if h.Integrations.TestAI == nil {
		writeErr(w, http.StatusServiceUnavailable, "unavailable", "AI test is not available")
		return
	}
	res, err := h.Integrations.Store.TestAI(r.Context(), actor(r), in.Provider, h.Integrations.TestAI)
	if err != nil {
		h.integrationsReply(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// PUT /integrations/maps {map_id, key?, clear?, mfa_code} — kunci browser Google Maps.
func (h *Handler) integrationsSaveMaps(w http.ResponseWriter, r *http.Request) {
	var in struct {
		integ.MapsInput
		MFACode string `json:"mfa_code"`
	}
	if !decode(w, r, &in) || !stepUp(w, r, in.MFACode, h.Integrations.ConfirmMFA) {
		return
	}
	h.integrationsReply(w, r, h.Integrations.Store.SaveMaps(r.Context(), actor(r), in.MapsInput))
}

func (h *Handler) integrationsView(w http.ResponseWriter, r *http.Request) {
	v, err := h.Integrations.Store.View(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) integrationsSaveSMS(w http.ResponseWriter, r *http.Request) {
	var in struct {
		integ.SMSInput
		MFACode string `json:"mfa_code"`
	}
	if !decode(w, r, &in) || !stepUp(w, r, in.MFACode, h.Integrations.ConfirmMFA) {
		return
	}
	h.integrationsReply(w, r, h.Integrations.Store.SaveSMS(r.Context(), actor(r), in.SMSInput))
}

func (h *Handler) integrationsSaveFacebook(w http.ResponseWriter, r *http.Request) {
	var in struct {
		integ.FacebookInput
		MFACode string `json:"mfa_code"`
	}
	if !decode(w, r, &in) || !stepUp(w, r, in.MFACode, h.Integrations.ConfirmMFA) {
		return
	}
	h.integrationsReply(w, r, h.Integrations.Store.SaveFacebook(r.Context(), actor(r), in.FacebookInput))
}

// POST /integrations/sms/test {phone} — kirim SMS uji (biaya 1 SMS) ke nomor yang diketik admin.
func (h *Handler) integrationsTestSMS(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Phone string `json:"phone"`
	}
	if !decode(w, r, &in) {
		return
	}
	phone, err := auth.NormalizePhone(in.Phone)
	if err != nil {
		writeFieldErr(w, "test_phone")
		return
	}
	res, err := h.Integrations.Store.TestSMS(r.Context(), actor(r), phone)
	if err != nil {
		h.integrationsReply(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) integrationsTestFacebook(w http.ResponseWriter, r *http.Request) {
	res, err := h.Integrations.Store.TestFacebook(r.Context(), actor(r))
	if err != nil {
		h.integrationsReply(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) integrationsReply(w http.ResponseWriter, r *http.Request, err error) {
	var fe integ.FieldError
	switch {
	case err == nil:
		h.integrationsView(w, r)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	default:
		h.log.Error("pengaturan integrasi", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save settings")
	}
}

// GET /admin/api/integrations/usage — pemakaian AI/SMS (hari ini, 7 & 30 hari, 14 hari terakhir) + perkiraan saldo + peringatan.
func (h *Handler) integrationsUsage(w http.ResponseWriter, r *http.Request) {
	rep, err := h.Integrations.Store.Usage(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rep)
}

// PUT /admin/api/integrations/usage/budget — saldo di console penyedia, batas peringatan, harga (wajib kode MFA).
func (h *Handler) integrationsSaveBudget(w http.ResponseWriter, r *http.Request) {
	var in struct {
		integ.BudgetInput
		MFACode string `json:"mfa_code"`
	}
	if !decode(w, r, &in) || !stepUp(w, r, in.MFACode, h.Integrations.ConfirmMFA) {
		return
	}
	var fe integ.FieldError
	switch err := h.Integrations.Store.SaveBudget(r.Context(), actor(r), in.BudgetInput); {
	case err == nil:
		h.integrationsUsage(w, r)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	default:
		h.log.Error("anggaran integrasi", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save settings")
	}
}
