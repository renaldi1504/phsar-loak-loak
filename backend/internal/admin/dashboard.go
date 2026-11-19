package admin

// Dashboard admin, audit log, dan matriks peran — mengikuti "Mock-up Admin Phsar".

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type queueItem struct {
	Key    string     `json:"key"`
	Count  int        `json:"count"`
	Oldest *time.Time `json:"oldest"`
}

// GET /admin/api/dashboard — angka ringkas + antrean yang perlu tindakan (sesuai izin) + aktivitas admin terbaru.
func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	a := r.Context().Value(ctxKey{}).(adminCtx)
	can := func(p string) bool { return a.perms["*"] || a.perms[p] }
	ctx := r.Context()
	out := map[string]any{}

	var gmv, comm string
	var orders, active, registered int
	if err := h.pool.QueryRow(ctx, `SELECT coalesce(sum(total), 0)::text, coalesce(sum(commission_amount), 0)::text, count(*)
		FROM orders WHERE paid_at > now() - interval '7 days' AND status NOT IN ('pending_payment', 'cancelled')`).Scan(&gmv, &comm, &orders); err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	_ = h.pool.QueryRow(ctx, `SELECT count(*) FILTER (WHERE status = 'active'), count(*) FROM stores WHERE deleted_at IS NULL`).Scan(&active, &registered)
	var toTransfer string
	var toTransferN int
	_ = h.pool.QueryRow(ctx, `SELECT coalesce(sum(amount), 0)::text, count(*) FROM payouts WHERE status = 'pending'`).Scan(&toTransfer, &toTransferN)
	out["kpis"] = map[string]any{"gmv_7d": gmv, "orders_7d": orders, "commission_7d": comm, "active_sellers": active,
		"registered_sellers": registered, "to_transfer": toTransfer, "to_transfer_count": toTransferN}

	queues := []queueItem{}
	add := func(perm, key, q string) {
		if !can(perm) {
			return
		}
		var it queueItem
		it.Key = key
		if err := h.pool.QueryRow(ctx, q).Scan(&it.Count, &it.Oldest); err == nil {
			queues = append(queues, it)
		}
	}
	add("seller.verify", "stores", `SELECT count(*), min(submitted_at) FROM stores WHERE status = 'pending_review' AND deleted_at IS NULL`)
	add("product.moderate", "products", `SELECT count(*), min(updated_at) FROM products WHERE moderation_status = 'pending_review' AND status <> 'draft' AND deleted_at IS NULL`)
	add("refund.decide", "refunds", `SELECT count(*), min(cancelled_at) FROM orders WHERE refund_status = 'pending'`)
	add("report.handle", "reports", `SELECT count(DISTINCT (target_type, coalesce(product_id, store_id))), min(created_at) FROM reports WHERE status = 'open'`) // reports:
	add("support.handle", "tickets", `SELECT count(*), min(last_user_message_at) FROM tickets WHERE status = 'open'
		AND (last_admin_reply_at IS NULL OR last_user_message_at > last_admin_reply_at)`) // support
	add("return.decide", "disputes", `SELECT count(*), min(created_at) FROM returns WHERE status = 'disputed'`)
	add("refund.decide", "return_refunds", `SELECT count(*), min(decided_at) FROM returns WHERE status = 'approved'`)
	add("payout.approve", "payouts", `SELECT count(*), min(created_at) FROM payouts WHERE status = 'pending'`)
	add("delivery.dispatch", "deliveries", `SELECT count(*), min(created_at) FROM deliveries WHERE status = 'unassigned'`)
	add("delivery.dispatch", "failed", `SELECT count(*), min(updated_at) FROM deliveries WHERE status = 'failed'`)
	add("delivery.dispatch", "courier_returns", `SELECT count(*), min(started_at) FROM courier_returns WHERE status IN ('returning','awaiting_seller')`) // courierreturn
	add("seller.verify", "pickups", `SELECT count(*), min(p.created_at) FROM store_pickup_points p JOIN stores s ON s.id = p.store_id
		WHERE p.verified_at IS NULL AND p.deleted_at IS NULL AND s.deleted_at IS NULL AND s.status = 'active'`)
	// Pengingat saldo / error AI & SMS (hanya yang boleh mengatur integrasi).
	if h.Integrations != nil && can("integration.configure") {
		alerts := h.Integrations.Store.Alerts(ctx)
		if len(alerts) > 0 {
			queues = append(queues, queueItem{Key: "integrations", Count: len(alerts)})
		}
		out["integration_alerts"] = alerts
	}
	out["queues"] = queues

	act := []map[string]any{}
	if rows, err := h.pool.Query(ctx, `SELECT l.occurred_at, coalesce(nullif(u.name, ''), u.phone, ''), coalesce(au.role_code, ''), l.action, l.entity_table, coalesce(l.entity_id, '')
		FROM audit_logs l LEFT JOIN users u ON u.id = l.actor_user_id LEFT JOIN admin_users au ON au.user_id = l.actor_user_id
		WHERE l.actor_type = 'admin' AND l.occurred_at > now() - interval '30 days'
		ORDER BY l.occurred_at DESC LIMIT 8`); err == nil {
		for rows.Next() {
			var at time.Time
			var actor, role, action, table, entity string
			if rows.Scan(&at, &actor, &role, &action, &table, &entity) == nil {
				act = append(act, map[string]any{"at": at, "actor": actor, "role": role, "action": action, "table": table, "entity": entity})
			}
		}
		rows.Close()
	}
	out["activity"] = act
	out["can_audit"] = can("audit.view")
	writeJSON(w, http.StatusOK, out)
}

// Kunci yang tidak pernah ditampilkan di audit log (nilai rahasia / terenkripsi).
func sensitiveKey(k string) bool {
	k = strings.ToLower(k)
	return strings.HasSuffix(k, "_enc") || strings.HasSuffix(k, "_hash") || strings.Contains(k, "token") || strings.Contains(k, "secret") ||
		k == "code" || k == "password"
}

func maskJSON(raw []byte) map[string]any {
	if len(raw) == 0 {
		return nil
	}
	var m map[string]any
	if json.Unmarshal(raw, &m) != nil {
		return nil
	}
	for k := range m {
		if sensitiveKey(k) && m[k] != nil {
			m[k] = "•••"
		}
	}
	return m
}

// GET /admin/api/audit?actor=&table=&action=&before= — 50 baris per halaman, terbaru dulu.
func (h *Handler) audit(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	where := []string{"l.occurred_at > now() - interval '180 days'"}
	args := []any{}
	arg := func(v any) string { args = append(args, v); return "$" + strconv.Itoa(len(args)) }
	switch v := q.Get("actor"); v {
	case "":
	case "admin", "seller", "buyer", "courier", "system":
		where = append(where, "l.actor_type = "+arg(v)+"::actor_type")
	default:
		if uuidRe.MatchString(v) {
			where = append(where, "l.actor_user_id = "+arg(v)+"::uuid")
		}
	}
	if v := q.Get("table"); v != "" {
		where = append(where, "l.entity_table = "+arg(v))
	}
	if v := q.Get("action"); v != "" {
		where = append(where, "l.action = "+arg(v))
	}
	if v, err := strconv.ParseInt(q.Get("before"), 10, 64); err == nil && v > 0 {
		where = append(where, "l.id < "+arg(v))
	}
	rows, err := h.pool.Query(r.Context(), `SELECT l.id, l.occurred_at, l.actor_type::text, coalesce(nullif(u.name, ''), u.phone, ''), coalesce(au.role_code, ''),
			l.action, l.entity_table, coalesce(l.entity_id, ''), l.old_data, l.new_data, coalesce(l.changed_fields, '{}'),
			coalesce(l.request_id, ''), coalesce(host(l.ip), ''), coalesce(l.user_agent, '')
		FROM audit_logs l LEFT JOIN users u ON u.id = l.actor_user_id LEFT JOIN admin_users au ON au.user_id = l.actor_user_id
		WHERE `+strings.Join(where, " AND ")+` ORDER BY l.id DESC LIMIT 50`, args...)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	defer rows.Close()
	type rec struct {
		ID        int64          `json:"id"`
		At        time.Time      `json:"at"`
		ActorType string         `json:"actor_type"`
		Actor     string         `json:"actor"`
		Role      string         `json:"role"`
		Action    string         `json:"action"`
		Table     string         `json:"table"`
		Entity    string         `json:"entity"`
		Old       map[string]any `json:"old"`
		New       map[string]any `json:"new"`
		Changed   []string       `json:"changed"`
		RequestID string         `json:"request_id"`
		IP        string         `json:"ip"`
		UA        string         `json:"user_agent"`
	}
	out := []rec{}
	for rows.Next() {
		var x rec
		var oldRaw, newRaw []byte
		if err := rows.Scan(&x.ID, &x.At, &x.ActorType, &x.Actor, &x.Role, &x.Action, &x.Table, &x.Entity, &oldRaw, &newRaw, &x.Changed,
			&x.RequestID, &x.IP, &x.UA); err != nil {
			writeErr(w, http.StatusInternalServerError, "internal", err.Error())
			return
		}
		x.Old, x.New = maskJSON(oldRaw), maskJSON(newRaw)
		out = append(out, x)
	}
	resp := map[string]any{"records": out}
	if q.Get("before") == "" { // pilihan filter, hanya di halaman pertama
		lists := map[string][]string{}
		for key, sql := range map[string]string{
			"tables":  `SELECT DISTINCT entity_table FROM audit_logs WHERE occurred_at > now() - interval '180 days' ORDER BY 1 LIMIT 100`,
			"actions": `SELECT DISTINCT action FROM audit_logs WHERE occurred_at > now() - interval '180 days' ORDER BY 1 LIMIT 100`,
		} {
			lists[key] = []string{}
			if rs, err := h.pool.Query(r.Context(), sql); err == nil {
				for rs.Next() {
					var s string
					if rs.Scan(&s) == nil {
						lists[key] = append(lists[key], s)
					}
				}
				rs.Close()
			}
		}
		resp["tables"], resp["actions"] = lists["tables"], lists["actions"]
	}
	writeJSON(w, http.StatusOK, resp)
}

// GET /admin/api/roles — matriks izin per peran.
func (h *Handler) roles(w http.ResponseWriter, r *http.Request) {
	rows, err := h.pool.Query(r.Context(), `SELECT role_code, permission FROM role_permissions ORDER BY role_code, permission`)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	defer rows.Close()
	m := map[string][]string{}
	for rows.Next() {
		var role, perm string
		if rows.Scan(&role, &perm) == nil {
			m[role] = append(m[role], perm)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"roles": m})
}
