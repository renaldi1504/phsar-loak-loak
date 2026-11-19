package payout

// extras: tarik saldo manual oleh seller (pemilik toko) — "Tarik sekarang".
//
//	syarat: dinyalakan super admin · saldo ≥ minimal penarikan manual · rekening ada & tidak ditahan
//	        · belum ada payout pending · belum melewati N kali per 7 hari
//	hasil : payout "pending" yang sama seperti jadwal (saldo dipotong), requested_by = pemilik.
//	        Admin keuangan mentransfer manual lalu menandai "paid" seperti biasa (antrean admin & dashboard).
//
// Aman dari klik ganda / jadwal bersamaan: kunci advisory per toko yang sama dengan Generate + indeks unik
// payouts_one_pending; semua syarat dihitung ulang di dalam transaksi.

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"phsar/internal/extras"
	"phsar/internal/platform/db"
)

// WithdrawState: apakah tombol "Tarik sekarang" boleh tampil, dan kalau tidak, satu alasan.
type WithdrawState struct {
	Enabled     bool       `json:"enabled"`
	Can         bool       `json:"can"`
	Reason      string     `json:"reason,omitempty"` // off · pending_exists · no_account · account_hold · limit · below_min
	Available   string     `json:"available"`
	MinAmount   string     `json:"min_amount"`
	MaxPer7Days int        `json:"max_per_7_days"`
	Used7Days   int        `json:"used_7_days"`
	HoldUntil   *time.Time `json:"hold_until,omitempty"` // account_hold: sampai kapan
	NextAt      *time.Time `json:"next_at,omitempty"`    // limit: kapan boleh lagi
}

// ErrWithdrawBlocked: syarat tidak terpenuhi saat ditekan (409, dengan alasan).
type withdrawBlocked struct{ reason string }

func (e withdrawBlocked) Error() string { return "withdraw blocked: " + e.reason }

type rowQuerier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func (s *Service) withdrawCheck(ctx context.Context, q rowQuerier, storeID string) (WithdrawState, *Account, decimal.Decimal, error) {
	st, err := extras.LoadManualPayout(ctx, q)
	if err != nil {
		return WithdrawState{}, nil, decimal.Zero, err
	}
	w := WithdrawState{Enabled: st.Enabled, MinAmount: st.MinAmount, MaxPer7Days: st.MaxPer7Days}
	bal, err := s.balance(ctx, q, storeID)
	if err != nil {
		return w, nil, bal, err
	}
	w.Available = bal.StringFixed(2)
	if !st.Enabled {
		w.Reason = "off"
		return w, nil, bal, nil
	}
	var pending bool
	var oldest *time.Time
	if err := q.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM payouts WHERE store_id = $1 AND status = 'pending'),
		  (SELECT count(*) FROM payouts WHERE store_id = $1 AND requested_by IS NOT NULL AND status <> 'cancelled' AND created_at > now() - interval '7 days'),
		  (SELECT min(created_at) FROM payouts WHERE store_id = $1 AND requested_by IS NOT NULL AND status <> 'cancelled' AND created_at > now() - interval '7 days')`,
		storeID).Scan(&pending, &w.Used7Days, &oldest); err != nil {
		return w, nil, bal, err
	}
	acc, _, err := s.account(ctx, q, storeID)
	if err != nil {
		return w, nil, bal, err
	}
	min, _ := decimal.NewFromString(st.MinAmount)
	switch {
	case pending:
		w.Reason = "pending_exists"
	case acc == nil:
		w.Reason = "no_account"
	case acc.HoldUntil != nil && acc.HoldUntil.After(time.Now()):
		w.Reason, w.HoldUntil = "account_hold", acc.HoldUntil
	case w.Used7Days >= st.MaxPer7Days:
		w.Reason = "limit"
		if oldest != nil {
			next := oldest.Add(7 * 24 * time.Hour)
			w.NextAt = &next
		}
	case !bal.IsPositive() || bal.LessThan(min):
		w.Reason = "below_min"
	default:
		w.Can = true
	}
	return w, acc, bal, nil
}

func (s *Service) WithdrawState(ctx context.Context, storeID string) (WithdrawState, error) {
	w, _, _, err := s.withdrawCheck(ctx, s.Pool, storeID)
	return w, err
}

// Withdraw: pemilik toko menarik seluruh saldo sekarang → payout pending. Hasil: id payout & jumlah.
func (s *Service) Withdraw(ctx context.Context, storeID, userID string) (string, string, error) {
	var pid, amount string
	err := db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('payout:' || $1))`, storeID); err != nil {
			return err
		}
		w, acc, bal, err := s.withdrawCheck(ctx, tx, storeID)
		if err != nil {
			return err
		}
		if !w.Can {
			return withdrawBlocked{w.Reason}
		}
		var orders int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM seller_ledger l WHERE l.store_id = $1 AND l.kind = 'order_earning'
			AND l.id > coalesce((SELECT max(id) FROM seller_ledger WHERE store_id = $1 AND kind = 'payout'), 0)`, storeID).Scan(&orders); err != nil {
			return err
		}
		if err := tx.QueryRow(ctx, `INSERT INTO payouts (store_id, amount, account_id, bank, account_name, account_last4, order_count, created_by, requested_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8) RETURNING id`, storeID, bal, acc.ID, acc.Bank, acc.AccountName, acc.Last4, orders, userID).Scan(&pid); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `INSERT INTO seller_ledger (store_id, kind, amount, payout_id, note, created_by) VALUES ($1, 'payout', $2, $3, 'requested by seller', $4)`,
			storeID, bal.Neg(), pid, userID); err != nil {
			return err
		}
		amount = bal.StringFixed(2)
		_, err = tx.Exec(ctx, `INSERT INTO audit_logs (actor_type, actor_user_id, store_id, action, entity_table, entity_id, new_data)
			VALUES ('seller', $1, $2, 'payout.request', 'payouts', $3, $4)`, userID, storeID, pid, map[string]string{"amount": amount})
		return err
	})
	return pid, amount, err
}

func (h *Handler) withdrawState(w http.ResponseWriter, r *http.Request) {
	st, err := h.S.WithdrawState(r.Context(), chi.URLParam(r, "storeID"))
	h.reply(w, r, st, err)
}

func (h *Handler) withdraw(w http.ResponseWriter, r *http.Request) {
	storeID := chi.URLParam(r, "storeID")
	pid, amount, err := h.S.Withdraw(r.Context(), storeID, uid(r))
	var b withdrawBlocked
	if errors.As(err, &b) {
		writeJSON(w, http.StatusConflict, map[string]any{"error": map[string]string{"code": "withdraw_blocked", "reason": b.reason, "message": "You can't withdraw right now."}})
		return
	}
	if err != nil {
		h.reply(w, r, nil, err)
		return
	}
	st, err := h.S.WithdrawState(r.Context(), storeID)
	h.reply(w, r, map[string]any{"payout_id": pid, "amount": amount, "withdraw": st}, err)
}
