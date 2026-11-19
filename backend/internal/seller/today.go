package seller

// today: ringkasan untuk layar "Hari ini" di Seller Center v2 — daftar tugas yang perlu dikerjakan
// sekarang (dengan tenggat terdekat) + angka singkat hari ini. Satu panggilan supaya beranda cepat.

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Today struct {
	NewOrders     int        `json:"new_orders"`
	NewDue        *time.Time `json:"new_due,omitempty"` // batas terima pesanan paling dekat
	ToPack        int        `json:"to_pack"`
	PackDue       *time.Time `json:"pack_due,omitempty"` // batas kemas paling dekat
	Ready         int        `json:"ready"`
	Returns       int        `json:"returns"`
	ReturnsDue    *time.Time `json:"returns_due,omitempty"`
	FlashInvites  int        `json:"flash_invites"`
	FlashStarts   *time.Time `json:"flash_starts,omitempty"`
	ProductFixes  int        `json:"product_fixes"`  // produk yang diminta diperbaiki Phsar
	TicketReplies int        `json:"ticket_replies"` // tiket bantuan yang menunggu balasan toko
	ChatUnread    int        `json:"chat_unread"`    // chat: percakapan yang pesan terakhirnya dari pembeli (belum dibalas)
	OrdersToday   int        `json:"orders_today"`
	SalesToday    string     `json:"sales_today"`
	Timing        *hoursView `json:"timing,omitempty"` // aturan batas waktu + toko buka sekarang?

	// courierreturn: paket gagal antar yang sudah dikembalikan kurir dan menunggu konfirmasi toko
	CourierReturns    int        `json:"courier_returns"`
	CourierReturnsDue *time.Time `json:"courier_returns_due,omitempty"`
}

// Commission: tarif komisi yang berlaku sekarang (untuk kalkulator "Anda terima" di editor produk).
// Aturan toko menang; tanpa itu dipakai tarif kategori, lalu global.
type Commission struct {
	Store      string            `json:"store,omitempty"`
	Global     string            `json:"global"`
	Categories map[string]string `json:"categories"`
}

const commissionActive = `valid_from <= now() AND (valid_to IS NULL OR valid_to > now())`

func (s *Service) Commission(ctx context.Context, storeID string) (Commission, error) {
	c := Commission{Global: "8", Categories: map[string]string{}}
	rows, err := s.pool.Query(ctx, `SELECT DISTINCT ON (scope, coalesce(category_id::text, ''))
		  scope::text, coalesce(category_id::text, ''), coalesce(store_id::text, ''), rate_percent::text
		FROM commission_rules WHERE `+commissionActive+` AND (scope <> 'store' OR store_id = $1)
		ORDER BY scope, coalesce(category_id::text, ''), valid_from DESC`, storeID)
	if err != nil {
		return c, err
	}
	defer rows.Close()
	for rows.Next() {
		var scope, cat, store, rate string
		if err := rows.Scan(&scope, &cat, &store, &rate); err != nil {
			return c, err
		}
		switch scope {
		case "global":
			c.Global = rate
		case "store":
			c.Store = rate
		case "category":
			c.Categories[cat] = rate
		}
	}
	if err := rows.Err(); err != nil {
		return c, err
	}
	// Tarif kategori utama berlaku juga untuk sub-kategorinya (kecuali sub-kategori punya tarif sendiri).
	kids, err := s.pool.Query(ctx, `SELECT id, parent_id FROM categories WHERE parent_id IS NOT NULL`)
	if err != nil {
		return c, err
	}
	defer kids.Close()
	for kids.Next() {
		var id, parent string
		if err := kids.Scan(&id, &parent); err != nil {
			return c, err
		}
		if r, ok := c.Categories[parent]; ok {
			if _, own := c.Categories[id]; !own {
				c.Categories[id] = r
			}
		}
	}
	return c, kids.Err()
}

func (s *Service) Today(ctx context.Context, storeID string) (Today, error) {
	var d Today
	err := s.pool.QueryRow(ctx, `SELECT
		  count(*) FILTER (WHERE status = 'paid'), min(respond_by) FILTER (WHERE status = 'paid'),
		  count(*) FILTER (WHERE status = 'accepted'), min(ready_by) FILTER (WHERE status = 'accepted'), count(*) FILTER (WHERE status = 'ready_for_pickup'),
		  count(*) FILTER (WHERE paid_at >= (date_trunc('day', now() AT TIME ZONE 'Asia/Phnom_Penh') AT TIME ZONE 'Asia/Phnom_Penh')
		                     AND status NOT IN ('cancelled','refunded')),
		  coalesce(sum(subtotal) FILTER (WHERE paid_at >= (date_trunc('day', now() AT TIME ZONE 'Asia/Phnom_Penh') AT TIME ZONE 'Asia/Phnom_Penh')
		                     AND status NOT IN ('cancelled','refunded')), 0)::numeric(12,2)::text
		FROM orders WHERE store_id = $1 AND paid_at IS NOT NULL`, storeID).
		Scan(&d.NewOrders, &d.NewDue, &d.ToPack, &d.PackDue, &d.Ready, &d.OrdersToday, &d.SalesToday)
	if err != nil {
		return d, err
	}
	if err := s.pool.QueryRow(ctx, `SELECT count(*), min(seller_deadline) FROM returns
		WHERE store_id = $1 AND status = 'requested' AND seller_responded_at IS NULL`, storeID).Scan(&d.Returns, &d.ReturnsDue); err != nil {
		return d, err
	}
	if err := s.pool.QueryRow(ctx, `SELECT count(*), min(fs.starts_at) FROM flash_items i JOIN flash_sessions fs ON fs.id = i.session_id
		WHERE i.store_id = $1 AND i.status = 'invited' AND fs.cancelled_at IS NULL AND fs.ends_at > now()`, storeID).Scan(&d.FlashInvites, &d.FlashStarts); err != nil {
		return d, err
	}
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM products WHERE store_id = $1 AND deleted_at IS NULL
		AND moderation_status = 'changes_requested'`, storeID).Scan(&d.ProductFixes); err != nil {
		return d, err
	}
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM tickets WHERE store_id = $1 AND status = 'pending'`, storeID).Scan(&d.TicketReplies); err != nil {
		return d, err
	}
	// chat: belum dibalas
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM chat_conversations WHERE store_id = $1 AND last_sender = 'buyer'`, storeID).Scan(&d.ChatUnread); err != nil {
		return d, err
	}
	// courierreturn
	if err := s.pool.QueryRow(ctx, `SELECT count(*), min(confirm_by) FROM courier_returns WHERE store_id = $1 AND status = 'awaiting_seller'`, storeID).
		Scan(&d.CourierReturns, &d.CourierReturnsDue); err != nil {
		return d, err
	}
	return d, nil
}

// GET /v1/seller/stores/{storeID}/today
func (h *Handler) today(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "storeID")
	d, err := h.svc.Today(r.Context(), id)
	if err == nil {
		if v, verr := h.svc.hoursView(r, id); verr == nil {
			d.Timing = &v
		}
	}
	h.respond(w, http.StatusOK, d, err)
}
