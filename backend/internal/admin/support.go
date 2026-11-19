package admin

// Alat CS: cari pelanggan, lihat semua pesanannya, catatan CS, tangguhkan akun;
// cari pesanan (nomor pesanan / transaksi PayWay / HP) dengan detail lengkap; log pembayaran PayWay.
// Data disusun langsung di Postgres (json_build_object) supaya satu permintaan = satu jawaban lengkap.

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
)

var (
	nonDigits = regexp.MustCompile(`\D`)
)

// phoneDigits: "012 345 678" / "+855 12 345 678" → "12345678" (cocok dengan bagian akhir nomor tersimpan).
func phoneDigits(q string) string {
	d := nonDigits.ReplaceAllString(q, "")
	d = strings.TrimPrefix(d, "855")
	return strings.TrimPrefix(d, "0")
}

func writeRawJSON(w http.ResponseWriter, raw []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(raw)
}

// ---------------------------------------------------------------------------
// Pelanggan
// ---------------------------------------------------------------------------

// GET /customers?q= — HP, nama, @username Telegram, nomor pesanan / retur, atau ID.
func (h *Handler) listCustomers(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	where, args := "u.deleted_at IS NULL", []any{}
	if q != "" {
		conds := []string{"u.name ILIKE '%' || $1 || '%'", "EXISTS (SELECT 1 FROM auth_identities ai WHERE ai.user_id = u.id AND ai.display ILIKE '%' || $1 || '%')",
			"EXISTS (SELECT 1 FROM orders o WHERE o.buyer_id = u.id AND o.order_no = upper($1))",
			"EXISTS (SELECT 1 FROM returns rt WHERE rt.buyer_id = u.id AND rt.return_no = upper($1))",
			"EXISTS (SELECT 1 FROM checkouts c WHERE c.buyer_id = u.id AND c.tran_id = $1)"}
		args = append(args, strings.TrimPrefix(q, "@"))
		if d := phoneDigits(q); len(d) >= 4 {
			args = append(args, d)
			conds = append(conds, "regexp_replace(coalesce(u.phone, ''), '\\D', '', 'g') LIKE '%' || $2 || '%'",
				"EXISTS (SELECT 1 FROM checkouts c WHERE c.buyer_id = u.id AND regexp_replace(c.phone, '\\D', '', 'g') LIKE '%' || $2 || '%')")
		}
		if uuidRe.MatchString(q) {
			conds = append(conds, "u.id::text = lower($1)")
		}
		where += " AND (" + strings.Join(conds, " OR ") + ")"
	}
	var raw []byte
	err := h.pool.QueryRow(r.Context(), `SELECT coalesce(json_agg(x ORDER BY x.last_activity DESC NULLS LAST), '[]') FROM (
		SELECT u.id, u.name, u.phone, u.is_guest, u.status, u.created_at,
		  (SELECT coalesce(array_agg(DISTINCT ai.provider::text), '{}') FROM auth_identities ai WHERE ai.user_id = u.id) AS providers,
		  (SELECT count(*) FROM orders o WHERE o.buyer_id = u.id AND o.status <> 'pending_payment') AS orders,
		  (SELECT coalesce(sum(o.total), 0)::text FROM orders o WHERE o.buyer_id = u.id AND o.paid_at IS NOT NULL) AS spent,
		  (SELECT count(*) FROM returns rt WHERE rt.buyer_id = u.id AND rt.status IN ('requested','disputed','approved')) AS open_returns,
		  EXISTS (SELECT 1 FROM store_members sm WHERE sm.user_id = u.id) AS is_seller,
		  EXISTS (SELECT 1 FROM couriers c WHERE c.user_id = u.id) AS is_courier,
		  EXISTS (SELECT 1 FROM admin_users a WHERE a.user_id = u.id) AS is_admin,
		  greatest(u.created_at, (SELECT max(o.created_at) FROM orders o WHERE o.buyer_id = u.id)) AS last_activity
		FROM users u WHERE `+where+` ORDER BY u.created_at DESC LIMIT 60) x`, args...).Scan(&raw)
	if err != nil {
		h.log.Error("cari pelanggan", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "search failed")
		return
	}
	writeRawJSON(w, []byte(`{"customers":`+string(raw)+`}`))
}

// GET /customers/{id}
func (h *Handler) customerDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !uuidRe.MatchString(id) {
		writeErr(w, http.StatusNotFound, "not_found", "customer not found")
		return
	}
	var raw []byte
	err := h.pool.QueryRow(r.Context(), `SELECT json_build_object(
		'id', u.id, 'name', u.name, 'phone', u.phone, 'email', u.email, 'language', u.language, 'is_guest', u.is_guest,
		'status', u.status, 'created_at', u.created_at,
		'is_admin', EXISTS (SELECT 1 FROM admin_users a WHERE a.user_id = u.id),
		'identities', (SELECT coalesce(json_agg(json_build_object('provider', ai.provider, 'display', ai.display, 'last_used_at', ai.last_used_at,
			'created_at', ai.created_at) ORDER BY ai.created_at), '[]') FROM auth_identities ai WHERE ai.user_id = u.id),
		'sessions', (SELECT count(*) FROM sessions s WHERE s.user_id = u.id AND s.revoked_at IS NULL AND s.expires_at > now()),
		'last_seen_at', (SELECT max(s.last_seen_at) FROM sessions s WHERE s.user_id = u.id),
		'stores', (SELECT coalesce(json_agg(json_build_object('id', st.id, 'name', st.name, 'status', st.status, 'role', sm.role)), '[]')
			FROM store_members sm JOIN stores st ON st.id = sm.store_id WHERE sm.user_id = u.id),
		'courier', (SELECT json_build_object('name', c.display_name, 'active', c.active) FROM couriers c WHERE c.user_id = u.id),
		'stats', json_build_object(
			'orders', (SELECT count(*) FROM orders o WHERE o.buyer_id = u.id AND o.status <> 'pending_payment'),
			'spent', (SELECT coalesce(sum(o.total), 0)::text FROM orders o WHERE o.buyer_id = u.id AND o.paid_at IS NOT NULL),
			'cancelled', (SELECT count(*) FROM orders o WHERE o.buyer_id = u.id AND o.status IN ('cancelled','refunded') AND o.paid_at IS NOT NULL),
			'returns', (SELECT count(*) FROM returns rt WHERE rt.buyer_id = u.id),
			'unpaid_checkouts', (SELECT count(*) FROM checkouts c WHERE c.buyer_id = u.id AND c.status <> 'paid')),
		'orders', (SELECT coalesce(json_agg(json_build_object('order_no', o.order_no, 'status', o.status, 'store_name', st.name, 'total', o.total::text,
			'item_count', o.item_count, 'refund_status', o.refund_status, 'created_at', o.created_at) ORDER BY o.created_at DESC), '[]')
			FROM (SELECT * FROM orders WHERE buyer_id = u.id AND status <> 'pending_payment' ORDER BY created_at DESC LIMIT 50) o JOIN stores st ON st.id = o.store_id),
		'returns', (SELECT coalesce(json_agg(json_build_object('id', rt.id, 'return_no', rt.return_no, 'order_no', o.order_no, 'status', rt.status,
			'reason', rt.reason, 'amount', coalesce(rt.refund_amount, rt.requested_amount)::text, 'created_at', rt.created_at) ORDER BY rt.created_at DESC), '[]')
			FROM returns rt JOIN orders o ON o.id = rt.order_id WHERE rt.buyer_id = u.id),
		'addresses', (SELECT coalesce(json_agg(DISTINCT jsonb_build_object('recipient', c.recipient, 'phone', c.phone, 'address', c.address_line)), '[]')
			FROM (SELECT recipient, phone, address_line FROM checkouts WHERE buyer_id = u.id ORDER BY created_at DESC LIMIT 20) c),
		'notes', (SELECT coalesce(json_agg(json_build_object('id', n.id, 'kind', n.kind, 'note', n.note, 'created_at', n.created_at,
			'author', coalesce(a.name, a.phone)) ORDER BY n.id DESC), '[]') FROM customer_notes n JOIN users a ON a.id = n.author_id WHERE n.user_id = u.id)
		) FROM users u WHERE u.id = $1 AND u.deleted_at IS NULL`, id).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		writeErr(w, http.StatusNotFound, "not_found", "customer not found")
		return
	}
	if err != nil {
		h.log.Error("detail pelanggan", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not load customer")
		return
	}
	_, _ = h.pool.Exec(r.Context(), `INSERT INTO audit_logs (actor_type, actor_user_id, action, entity_table, entity_id)
		VALUES ('admin', $1, 'customer.view', 'users', $2)`, actor(r), id)
	writeRawJSON(w, raw)
}

func cleanNote(n string) (string, bool) {
	n = strings.TrimSpace(n)
	c := utf8.RuneCountInString(n)
	return n, c >= 2 && c <= 2000
}

// POST /customers/{id}/notes {note}
func (h *Handler) addCustomerNote(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Note string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	note, ok := cleanNote(in.Note)
	if !ok {
		writeFieldErr(w, "note")
		return
	}
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(r.Context(), `INSERT INTO customer_notes (user_id, author_id, note) SELECT id, $2, $3 FROM users WHERE id = $1 AND deleted_at IS NULL`,
			chi.URLParam(r, "id"), actor(r), note)
		if err == nil && tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
		return err
	})
	h.supportReply(w, r, err)
}

// POST /customers/{id}/status {status: suspended|active, note} — tangguhkan: semua sesi dicabut, tidak bisa login.
func (h *Handler) setCustomerStatus(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	if in.Status != "suspended" && in.Status != "active" {
		writeFieldErr(w, "status")
		return
	}
	note, ok := cleanNote(in.Note)
	if !ok || utf8.RuneCountInString(note) < 5 {
		writeFieldErr(w, "note")
		return
	}
	id := chi.URLParam(r, "id")
	if id == actor(r) {
		writeErr(w, http.StatusConflict, "self", "you cannot change your own account")
		return
	}
	var isAdmin bool
	_ = h.pool.QueryRow(r.Context(), `SELECT EXISTS (SELECT 1 FROM admin_users WHERE user_id = $1)`, id).Scan(&isAdmin)
	if isAdmin {
		writeErr(w, http.StatusConflict, "is_admin", "admin accounts are managed in Admins & roles")
		return
	}
	kind := map[string]string{"suspended": "suspended", "active": "reactivated"}[in.Status]
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(r.Context(), `UPDATE users SET status = $2 WHERE id = $1 AND deleted_at IS NULL AND status <> $2 AND status <> 'deleted'`, id, in.Status)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return errNotPending
		}
		if in.Status == "suspended" {
			if _, err := tx.Exec(r.Context(), `UPDATE sessions SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL`, id); err != nil {
				return err
			}
		}
		_, err = tx.Exec(r.Context(), `INSERT INTO customer_notes (user_id, author_id, kind, note) VALUES ($1, $2, $3, $4)`, id, actor(r), kind, note)
		return err
	})
	h.supportReply(w, r, err)
}

func (h *Handler) supportReply(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case err == nil:
		h.customerDetail(w, r)
	case errors.Is(err, pgx.ErrNoRows):
		writeErr(w, http.StatusNotFound, "not_found", "customer not found")
	case errors.Is(err, errNotPending):
		writeErr(w, http.StatusConflict, "no_change", "account already has this status")
	default:
		h.log.Error("pelanggan", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save")
	}
}

// ---------------------------------------------------------------------------
// Pesanan
// ---------------------------------------------------------------------------

var orderStatuses = map[string]bool{"pending_payment": true, "paid": true, "accepted": true, "ready_for_pickup": true, "picked_up": true,
	"delivered": true, "completed": true, "cancelled": true, "refunded": true}

// GET /orders?q=&status= — nomor pesanan, ID transaksi PayWay, nomor retur, HP pembeli/penerima, nama toko.
func (h *Handler) listOrders(w http.ResponseWriter, r *http.Request) {
	q, status := strings.TrimSpace(r.URL.Query().Get("q")), r.URL.Query().Get("status")
	where, args := []string{"o.status <> 'pending_payment'"}, []any{}
	if status != "" {
		if status == "pending_payment" {
			where = nil
		}
		if !orderStatuses[status] {
			writeFieldErr(w, "status")
			return
		}
		args = append(args, status)
		where = append(where, "o.status = $1::order_status")
	}
	if q != "" {
		args = append(args, q)
		n := len(args)
		p := "$" + strconv.Itoa(n)
		conds := []string{"o.order_no LIKE upper(" + p + ") || '%'", "c.tran_id = " + p, "st.name ILIKE '%' || " + p + " || '%'",
			"EXISTS (SELECT 1 FROM returns rt WHERE rt.order_id = o.id AND rt.return_no = upper(" + p + "))"}
		if d := phoneDigits(q); len(d) >= 4 {
			args = append(args, d)
			p2 := "$" + strconv.Itoa(len(args))
			conds = append(conds, "regexp_replace(c.phone, '\\D', '', 'g') LIKE '%' || "+p2+" || '%'",
				"regexp_replace(coalesce(u.phone, ''), '\\D', '', 'g') LIKE '%' || "+p2+" || '%'")
		}
		where = append(where, "("+strings.Join(conds, " OR ")+")")
	}
	w2 := "TRUE"
	if len(where) > 0 {
		w2 = strings.Join(where, " AND ")
	}
	var raw []byte
	err := h.pool.QueryRow(r.Context(), `SELECT coalesce(json_agg(x), '[]') FROM (
		SELECT o.order_no, o.status, o.total::text AS total, o.item_count, o.refund_status, o.created_at, o.paid_at,
		  st.name AS store_name, coalesce(nullif(u.name, ''), c.recipient) AS buyer_name, c.phone AS buyer_phone, c.tran_id, d.status AS delivery_status,
		  (SELECT rt.status FROM returns rt WHERE rt.order_id = o.id ORDER BY rt.created_at DESC LIMIT 1) AS return_status
		FROM orders o JOIN checkouts c ON c.id = o.checkout_id JOIN stores st ON st.id = o.store_id JOIN users u ON u.id = o.buyer_id
		LEFT JOIN deliveries d ON d.order_id = o.id
		WHERE `+w2+` ORDER BY o.created_at DESC LIMIT 100) x`, args...).Scan(&raw)
	if err != nil {
		h.log.Error("cari pesanan", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "search failed")
		return
	}
	writeRawJSON(w, []byte(`{"orders":`+string(raw)+`}`))
}

// GET /orders/{no} — semua yang dibutuhkan CS/keuangan untuk satu pesanan.
func (h *Handler) orderDetail(w http.ResponseWriter, r *http.Request) {
	no := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "no")))
	var raw []byte
	err := h.pool.QueryRow(r.Context(), `SELECT json_build_object(
		'order_no', o.order_no, 'status', o.status, 'created_at', o.created_at, 'paid_at', o.paid_at, 'accepted_at', o.accepted_at,
		'ready_at', o.ready_at, 'completed_at', o.completed_at, 'cancelled_at', o.cancelled_at, 'cancel_reason', o.cancel_reason,
		'reject_reason', o.reject_reason, 'seller_note', o.seller_note,
		'subtotal', o.subtotal::text, 'shipping_fee', o.shipping_fee::text, 'total', o.total::text, 'commission_rate', o.commission_rate::text,
		'commission_amount', o.commission_amount::text, 'seller_amount', o.seller_amount::text,
		'refund', json_build_object('status', o.refund_status, 'amount', o.refund_amount::text, 'ref', o.refund_ref, 'at', o.refunded_at, 'note', o.refund_note),
		'buyer', json_build_object('id', u.id, 'name', u.name, 'phone', u.phone, 'is_guest', u.is_guest, 'status', u.status),
		'store', json_build_object('id', st.id, 'name', st.name, 'slug', st.slug, 'phone', (SELECT us.phone FROM store_members sm JOIN users us ON us.id = sm.user_id
			WHERE sm.store_id = st.id AND sm.role = 'owner' LIMIT 1)),
		'address', json_build_object('recipient', c.recipient, 'phone', c.phone, 'line', c.address_line, 'note', c.address_note, 'lat', c.lat, 'lng', c.lng),
		'payment', json_build_object('tran_id', c.tran_id, 'provider', c.provider, 'checkout_status', c.status, 'checkout_total', c.total::text,
			'paid_amount', c.paid_amount::text, 'paid_at', c.paid_at, 'type', c.payment_type, 'ref', c.payment_ref, 'payer', c.payer_account,
			'orders_in_checkout', (SELECT count(*) FROM orders o2 WHERE o2.checkout_id = c.id)),
		'items', (SELECT coalesce(json_agg(json_build_object('name', oi.name, 'options', oi.options, 'sku', oi.sku, 'image_url', oi.image_url,
			'unit_price', oi.unit_price::text, 'qty', oi.qty, 'line_total', oi.line_total::text, 'product_id', oi.product_id) ORDER BY oi.id), '[]')
			FROM order_items oi WHERE oi.order_id = o.id),
		'delivery', (SELECT json_build_object('status', d.status, 'courier', cr.display_name, 'courier_phone', cu.phone, 'pickup', d.pickup_label,
			'attempts', d.attempts, 'fail_reason', d.fail_reason, 'fail_note', d.fail_note, 'has_proof', d.proof_key IS NOT NULL, 'id', d.id,
			'assigned_at', d.assigned_at, 'picked_up_at', d.picked_up_at, 'delivered_at', d.delivered_at, 'failed_at', d.failed_at)
			FROM deliveries d LEFT JOIN couriers cr ON cr.user_id = d.courier_id LEFT JOIN users cu ON cu.id = d.courier_id WHERE d.order_id = o.id),
		'events', (SELECT coalesce(json_agg(json_build_object('status', e.status, 'actor_type', e.actor_type, 'actor', coalesce(nullif(ea.name, ''), ea.phone),
			'note', e.note, 'at', e.created_at) ORDER BY e.id), '[]') FROM order_events e LEFT JOIN users ea ON ea.id = e.actor_id WHERE e.order_id = o.id),
		'payment_events', (SELECT coalesce(json_agg(json_build_object('id', pe.id, 'source', pe.source, 'verified', pe.verified, 'status', pe.status,
			'payload', pe.payload, 'at', pe.created_at) ORDER BY pe.id), '[]') FROM payment_events pe WHERE pe.checkout_id = c.id),
		'returns', (SELECT coalesce(json_agg(json_build_object('id', rt.id, 'return_no', rt.return_no, 'status', rt.status, 'reason', rt.reason,
			'amount', coalesce(rt.refund_amount, rt.requested_amount)::text, 'created_at', rt.created_at) ORDER BY rt.created_at), '[]') FROM returns rt WHERE rt.order_id = o.id),
		'ledger', (SELECT coalesce(json_agg(json_build_object('kind', l.kind, 'amount', l.amount::text, 'at', l.created_at) ORDER BY l.id), '[]')
			FROM seller_ledger l WHERE l.order_id = o.id)
		) FROM orders o JOIN checkouts c ON c.id = o.checkout_id JOIN stores st ON st.id = o.store_id JOIN users u ON u.id = o.buyer_id
		WHERE o.order_no = $1`, no).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		writeErr(w, http.StatusNotFound, "not_found", "order not found")
		return
	}
	if err != nil {
		h.log.Error("detail pesanan", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not load order")
		return
	}
	writeRawJSON(w, raw)
}

// ---------------------------------------------------------------------------
// Log pembayaran (PayWay): create · callback · check · fake
// ---------------------------------------------------------------------------

// GET /payment-events?filter=all|callback|problems&q=tran_id/nomor pesanan
func (h *Handler) listPaymentEvents(w http.ResponseWriter, r *http.Request) {
	where, args := []string{"TRUE"}, []any{}
	switch r.URL.Query().Get("filter") {
	case "callback":
		where = append(where, "pe.source = 'callback'")
	case "problems":
		// callback tanpa tanda tangan yang cocok, ditolak/gagal, atau dibayar setelah checkout tidak lagi menunggu
		where = append(where, `(NOT pe.verified AND pe.source = 'callback'
			OR upper(coalesce(pe.status, '')) IN ('DECLINED','FAILED','ERROR','CANCELLED','REVERSED')
			OR (upper(coalesce(pe.status, '')) = 'APPROVED' AND c.status IN ('expired','cancelled')))`)
	}
	if q := strings.TrimSpace(r.URL.Query().Get("q")); q != "" {
		args = append(args, q)
		where = append(where, "(pe.tran_id = $1 OR EXISTS (SELECT 1 FROM orders o WHERE o.checkout_id = pe.checkout_id AND o.order_no = upper($1)))")
	}
	var raw []byte
	err := h.pool.QueryRow(r.Context(), `SELECT coalesce(json_agg(x), '[]') FROM (
		SELECT pe.id, pe.source, pe.verified, pe.status, pe.tran_id, pe.created_at AS at, pe.payload,
		  c.status AS checkout_status, c.total::text AS checkout_total, c.paid_amount::text AS paid_amount, c.phone AS buyer_phone,
		  (SELECT coalesce(array_agg(o.order_no ORDER BY o.order_no), '{}') FROM orders o WHERE o.checkout_id = pe.checkout_id) AS orders
		FROM payment_events pe LEFT JOIN checkouts c ON c.id = pe.checkout_id
		WHERE `+strings.Join(where, " AND ")+` ORDER BY pe.id DESC LIMIT 200) x`, args...).Scan(&raw)
	if err != nil {
		h.log.Error("log pembayaran", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not load")
		return
	}
	var counts struct {
		Today    int `json:"today"`
		Problems int `json:"problems"`
		Paid24h  int `json:"paid_24h"`
	}
	_ = h.pool.QueryRow(r.Context(), `SELECT
		(SELECT count(*) FROM payment_events WHERE created_at > now() - interval '24 hours'),
		(SELECT count(*) FROM payment_events pe LEFT JOIN checkouts c ON c.id = pe.checkout_id WHERE pe.created_at > now() - interval '7 days' AND
		  (NOT pe.verified AND pe.source = 'callback' OR upper(coalesce(pe.status, '')) IN ('DECLINED','FAILED','ERROR','CANCELLED','REVERSED')
		   OR (upper(coalesce(pe.status, '')) = 'APPROVED' AND c.status IN ('expired','cancelled')))),
		(SELECT count(*) FROM checkouts WHERE paid_at > now() - interval '24 hours')`).Scan(&counts.Today, &counts.Problems, &counts.Paid24h)
	c, _ := json.Marshal(counts)
	writeRawJSON(w, []byte(`{"events":`+string(raw)+`,"counts":`+string(c)+`}`))
}
