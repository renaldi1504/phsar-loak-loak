package qa

import (
	"context"
	"io"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"phsar/internal/httpx"
)

// ---------------------------------------------------------------------------
// Batas waktu pesanan diatur super admin + hitungan hanya saat toko buka + batas kemas.
// Juga: 1 akun hanya boleh memiliki 1 toko.
// ---------------------------------------------------------------------------

func TestOrderTiming(t *testing.T) {
	e := getEnv(t)
	ctx := context.Background()
	sa := newAdmin(t, e, "super_admin")
	mod := newAdmin(t, e, "moderator")
	s := newSellable(t, e, mod, "Timing kettle", []string{"10.00"}, []int{50})
	buyer, _ := e.login(t)
	shop := httpx.NewShop(e.cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), e.pool)
	base := "/v1/seller/stores/" + s.storeID
	everyDay := []map[string]string{}
	for i := 0; i < 7; i++ {
		everyDay = append(everyDay, map[string]string{"open": "08:00", "close": "22:00"})
	}
	rules := func(accept, pack int, biz bool) map[string]any {
		return map[string]any{"accept_minutes": accept, "pack_enabled": true, "pack_minutes": pack, "business_hours": biz, "default_hours": everyDay}
	}
	t.Cleanup(func() { e.call(t, "PUT", "/admin/api/order-timing", rules(120, 720, true), sa) })
	n := 0
	pay := func() string {
		n++
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-timing-"+time.Now().Format("150405.000")+string(rune('a'+n)), "10.00", item(s.variants[0], 1)), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		return e.scalar(t, `SELECT order_no FROM orders WHERE checkout_id = $1`, c.str("id"))
	}
	minutesLeft := func(col, no string) float64 {
		v := e.scalar(t, `SELECT round(extract(epoch FROM `+col+` - now()) / 60)::text FROM orders WHERE order_no = $1`, no)
		f, _ := strconv.ParseFloat(v, 64)
		return f
	}

	t.Run("hanya super admin yang bisa mengatur; isian divalidasi", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/order-timing", nil, mod).want(403)
		e.call(t, "GET", "/admin/api/order-timing", nil, sa).want(200).eq("accept_minutes", "120")
		e.call(t, "PUT", "/admin/api/order-timing", rules(5, 720, true), sa).wantErr(422, "validation")
		closed := rules(120, 720, true)
		closed["default_hours"] = []any{nil, nil, nil, nil, nil, nil, nil}
		e.call(t, "PUT", "/admin/api/order-timing", closed, sa).wantErr(422, "validation")
	})

	t.Run("hitungan biasa (24 jam): batas = sekarang + menit dari admin", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/order-timing", rules(45, 120, false), sa).want(200)
		no := pay()
		m := minutesLeft("respond_by", no)
		check(t, m >= 43 && m <= 45, "respond_by harus ±45 menit, dapat %v", m)
		e.call(t, "POST", base+"/orders/"+no+"/accept", nil, s.token).want(200)
		p := minutesLeft("ready_by", no)
		check(t, p >= 118 && p <= 120, "ready_by harus ±120 menit, dapat %v", p)
	})

	t.Run("hanya saat toko buka: toko tutup sampai lusa → hitungan menunggu", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/order-timing", rules(60, 120, true), sa).want(200)
		ict := time.FixedZone("ICT", 7*3600)
		openDay := int(time.Now().In(ict).AddDate(0, 0, 2).Weekday()) // hanya buka 2 hari lagi
		week := make([]any, 7)
		week[openDay] = map[string]string{"open": "08:00", "close": "22:00"}
		e.call(t, "PUT", base+"/hours", map[string]any{"hours": week}, s.token).want(200).eq("open_now", "false")
		no := pay()
		m := minutesLeft("respond_by", no)
		check(t, m > 24*60, "batas harus menunggu toko buka (>24 jam), dapat %v menit", m)
		at := e.scalar(t, `SELECT to_char(respond_by AT TIME ZONE 'Asia/Phnom_Penh', 'HH24:MI') FROM orders WHERE order_no = $1`, no)
		check(t, at == "09:00", "toko buka 08:00 + 60 menit → 09:00, dapat %s", at)
		e.call(t, "GET", base+"/today", nil, s.token).want(200).eq("timing.business_hours", "true").eq("timing.open_now", "false")
		// null = kembali ke jam standar Phsar
		e.call(t, "PUT", base+"/hours", map[string]any{"hours": nil}, s.token).want(200).eq("hours", "")
	})

	t.Run("lewat batas kemas → batal otomatis + refund + seller diberi tahu", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/order-timing", rules(60, 120, false), sa).want(200)
		no := pay()
		e.call(t, "POST", base+"/orders/"+no+"/accept", nil, s.token).want(200)
		e.sql(t, `UPDATE orders SET ready_by = now() - interval '1 minute' WHERE order_no = $1`, no)
		if n, err := shop.ExpireUnpacked(ctx); err != nil || n < 1 {
			t.Fatalf("ExpireUnpacked: %d %v", n, err)
		}
		e.call(t, "GET", base+"/orders/"+no, nil, s.token).want(200).eq("status", "cancelled").eq("cancel_reason", "pack_timeout")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'order.pack_timeout' AND store_id = $1`, s.storeID) >= "1", "seller tidak diberi tahu")
		check(t, e.scalar(t, `SELECT refund_status FROM orders WHERE order_no = $1`, no) == "pending", "refund harus diantrekan")
	})

	t.Run("1 akun hanya boleh memiliki 1 toko", func(t *testing.T) {
		e.call(t, "POST", "/v1/seller/stores", map[string]string{"name": "Second shop"}, s.token).wantErr(409, "too_many_stores")
	})
}
