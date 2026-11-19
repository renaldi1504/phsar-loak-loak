package reports

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
)

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// BuyerReport: POST /v1/reports (login, bukan tamu).
func (s *Service) BuyerReport(w http.ResponseWriter, r *http.Request) {
	var in ReportInput
	if !Decode(w, r, &in) {
		return
	}
	id, err := s.Create(r.Context(), uid(r), in)
	if err != nil {
		s.Fail(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]string{"id": id, "status": "open"})
}

// SellerStrikes: GET /v1/seller/stores/{storeID}/strikes (anggota toko).
func (s *Service) SellerStrikes(w http.ResponseWriter, r *http.Request) {
	raw, err := s.Strikes(r.Context(), chi.URLParam(r, "storeID"), true)
	if err != nil {
		s.Fail(w, err)
		return
	}
	WriteRaw(w, raw)
}

// Fail memetakan error paket ini ke jawaban HTTP.
func (s *Service) Fail(w http.ResponseWriter, err error) {
	var fe FieldError
	switch {
	case errors.As(err, &fe):
		WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, ErrLimit):
		writeErr(w, http.StatusTooManyRequests, "report_limit", err.Error())
	case errors.Is(err, ErrOpen):
		writeErr(w, http.StatusConflict, "already_reported", err.Error())
	case errors.Is(err, ErrNothing):
		writeErr(w, http.StatusConflict, "no_open_reports", err.Error())
	case errors.Is(err, ErrRevoked):
		writeErr(w, http.StatusConflict, "already_revoked", err.Error())
	default:
		if s.Log != nil {
			s.Log.Error("laporan", "err", err)
		}
		writeErr(w, http.StatusInternalServerError, "internal", "something went wrong")
	}
}

func Decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_json", "invalid JSON body")
		return false
	}
	return true
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteRaw(w http.ResponseWriter, raw []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(raw)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	WriteJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}
