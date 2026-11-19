package admin

// Commission sweep (izin payout.approve): saldo rekening penampung PayWay, komisi yang boleh dipindahkan
// ke rekening perusahaan, catatan sweep (tidak pernah dihapus, hanya dibatalkan), ringkasan bulanan & ekspor CSV.

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"phsar/internal/sweep"
)

func (h *Handler) mountSweep(r chi.Router) {
	s := &sweep.Service{Pool: h.pool}
	r.Route("/sweep", func(p chi.Router) {
		p.Use(h.can("payout.approve"))
		p.Get("/", func(w http.ResponseWriter, r *http.Request) { h.sweepView(w, r, s) })
		p.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var in sweep.Input
			if !decode(w, r, &in) {
				return
			}
			_, err := s.Record(r.Context(), actor(r), in)
			h.sweepReply(w, r, s, err, http.StatusCreated)
		})
		p.Post("/{id}/void", func(w http.ResponseWriter, r *http.Request) {
			var in struct {
				Note string `json:"note"`
			}
			if !decode(w, r, &in) {
				return
			}
			_, err := s.Void(r.Context(), actor(r), chi.URLParam(r, "id"), in.Note)
			h.sweepReply(w, r, s, err, http.StatusOK)
		})
		p.Get("/export.csv", func(w http.ResponseWriter, r *http.Request) { h.sweepCSV(w, r, s) })
	})
}

func (h *Handler) sweepView(w http.ResponseWriter, r *http.Request, s *sweep.Service) {
	h.sweepWrite(w, r, s, http.StatusOK)
}

func (h *Handler) sweepWrite(w http.ResponseWriter, r *http.Request, s *sweep.Service, status int) {
	sum, err := s.Summary(r.Context())
	if err != nil {
		h.fail(w, err)
		return
	}
	list, err := s.List(r.Context(), 100)
	if err != nil {
		h.fail(w, err)
		return
	}
	months, err := s.Monthly(r.Context(), 12)
	if err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, status, map[string]any{"summary": sum, "sweeps": list, "monthly": months})
}

func (h *Handler) sweepReply(w http.ResponseWriter, r *http.Request, s *sweep.Service, err error, status int) {
	var fe sweep.FieldError
	switch {
	case err == nil:
		h.sweepWrite(w, r, s, status)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, sweep.ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, sweep.ErrVoided):
		writeErr(w, http.StatusConflict, "already_voided", err.Error())
	default:
		h.fail(w, err)
	}
}

// csvSafe mencegah formula spreadsheet (=, +, -, @) dari teks yang diketik admin.
func csvSafe(s string) string {
	if s != "" && strings.ContainsRune("=+-@\t\r", rune(s[0])) {
		return "'" + s
	}
	return s
}

// GET /sweep/export.csv?what=monthly|sweeps
func (h *Handler) sweepCSV(w http.ResponseWriter, r *http.Request, s *sweep.Service) {
	what := r.URL.Query().Get("what")
	var rows [][]string
	switch what {
	case "sweeps":
		list, err := s.List(r.Context(), 10000)
		if err != nil {
			h.fail(w, err)
			return
		}
		rows = append(rows, []string{"date", "amount_usd", "to_account", "reference", "note", "recorded_by", "recorded_at", "voided_at", "void_note"})
		for _, x := range list {
			voided, vnote, note := "", "", ""
			if x.VoidedAt != nil {
				voided = x.VoidedAt.UTC().Format(time.RFC3339)
			}
			if x.VoidNote != nil {
				vnote = *x.VoidNote
			}
			if x.Note != nil {
				note = *x.Note
			}
			rows = append(rows, []string{x.SweptOn, x.Amount, csvSafe(x.ToAccount), csvSafe(x.Reference), csvSafe(note), csvSafe(x.By),
				x.CreatedAt.UTC().Format(time.RFC3339), voided, csvSafe(vnote)})
		}
	default:
		what = "monthly"
		months, err := s.Monthly(r.Context(), 24)
		if err != nil {
			h.fail(w, err)
			return
		}
		rows = append(rows, []string{"month", "collected_usd", "refunds_usd", "seller_payouts_usd", "commission_net_usd", "swept_usd"})
		for _, m := range months {
			rows = append(rows, []string{m.Month, m.Collected, m.Refunds, m.Payouts, m.Commission, m.Swept})
		}
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="phsar-sweep-`+what+`-`+time.Now().Format("20060102")+`.csv"`)
	w.Header().Set("Cache-Control", "no-store")
	cw := csv.NewWriter(w)
	_ = cw.WriteAll(rows)
}
