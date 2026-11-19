package admin

// extras: pengaturan super admin untuk strike otomatis (pesanan telat) dan tarik saldo manual seller.
//   GET/PUT /admin/api/auto-strikes   (izin order.rules — hanya super admin)
//   GET/PUT /admin/api/manual-payout  (izin payout.rules — hanya super admin)

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"phsar/internal/extras"
	"phsar/internal/platform/db"
)

func (h *Handler) mountExtras(r chi.Router) {
	r.With(h.can("order.rules")).Get("/auto-strikes", func(w http.ResponseWriter, r *http.Request) {
		st, err := extras.LoadAutoStrikes(r.Context(), h.pool)
		h.exReply(w, st, err)
	})
	r.With(h.can("order.rules")).Put("/auto-strikes", func(w http.ResponseWriter, r *http.Request) {
		var in extras.AutoStrikes
		if !exDecode(w, r, &in) {
			return
		}
		if err := in.Check(); err != nil {
			h.exReply(w, nil, err)
			return
		}
		h.exReply(w, in, h.exSave(r, extras.KeyAutoStrikes, in))
	})
	r.With(h.can("payout.rules")).Get("/manual-payout", func(w http.ResponseWriter, r *http.Request) {
		st, err := extras.LoadManualPayout(r.Context(), h.pool)
		h.exReply(w, st, err)
	})
	r.With(h.can("payout.rules")).Put("/manual-payout", func(w http.ResponseWriter, r *http.Request) {
		var in extras.ManualPayout
		if !exDecode(w, r, &in) {
			return
		}
		if err := in.Check(); err != nil {
			h.exReply(w, nil, err)
			return
		}
		h.exReply(w, in, h.exSave(r, extras.KeyManualPayout, in))
	})
}

func (h *Handler) exSave(r *http.Request, key string, v any) error {
	return db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		return extras.Save(r.Context(), tx, key, actor(r), v)
	})
}

func (h *Handler) exReply(w http.ResponseWriter, v any, err error) {
	var fe extras.FieldError
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, v)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	default:
		h.fail(w, err)
	}
}

// exDecode: JSON ketat (field tak dikenal ditolak).
func exDecode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(io.LimitReader(r.Body, 16<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_json", "invalid JSON body")
		return false
	}
	return true
}
