package qa

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

// chatUpload: upload foto chat (multipart) ke jalur pembeli / toko.
func (e *env) chatUpload(t *testing.T, token, path string, data []byte) resp {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("kind", "chat_photo")
	fw, _ := mw.CreateFormFile("file", "photo.jpg")
	_, _ = fw.Write(data)
	_ = mw.Close()
	req, _ := http.NewRequest(http.MethodPost, e.srv.URL+path, &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	return e.send(t, req)
}

// ---------------------------------------------------------------------------
// Chat pembeli ↔ toko: buka dari produk/pesanan, kirim, tanda dibaca, staf, foto, notifikasi berbatas, admin, batas kirim
// ---------------------------------------------------------------------------

func TestChat(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs")
	s := newSellable(t, e, mod, "Chat sneaker", []string{"9.00"}, []int{20})
	o := newSellable(t, e, mod, "Other shop item", []string{"4.00"}, []int{5})
	buyer, _ := e.login(t)
	other, _ := e.login(t)
	staff, staffPhone := e.login(t)
	e.sql(t, `INSERT INTO store_members (store_id, user_id, role) SELECT $1, id, 'staff' FROM users WHERE phone = '+855' || substr($2, 2)`, s.storeID, staffPhone)
	c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-chat-1-xyz", "9.00", item(s.variants[0], 1)), buyer).want(201)
	e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
	no := c.str("orders.0.order_no")
	sb := "/v1/seller/stores/" + s.storeID + "/chat"

	var cid string
	t.Run("pembeli membuka chat dari produk / pesanan", func(t *testing.T) {
		e.call(t, "POST", "/v1/shop/chat", map[string]any{}, buyer).wantErr(422, "validation", "store_id")
		e.call(t, "POST", "/v1/shop/chat", map[string]any{"product_id": "00000000-0000-0000-0000-000000000000"}, buyer).wantErr(422, "validation", "product_id")
		e.call(t, "POST", "/v1/shop/chat", map[string]any{"store_id": s.storeID, "product_id": o.productID}, buyer).wantErr(422, "validation", "product_id")
		e.call(t, "POST", "/v1/shop/chat", map[string]any{"order_no": no}, other).wantErr(422, "validation", "order_no")
		e.call(t, "POST", "/v1/shop/chat", map[string]any{"store_id": s.storeID}, s.token).wantErr(422, "validation", "store_id")
		e.call(t, "POST", "/v1/shop/chat", map[string]any{"store_id": s.storeID, "x": 1}, buyer).wantErr(400, "bad_json")
		e.call(t, "GET", "/v1/shop/chat", nil, "").want(401)
		r := e.call(t, "POST", "/v1/shop/chat", map[string]any{"product_id": s.productID}, buyer).want(201)
		r.eq("store_id", s.storeID).eq("card.type", "product").eq("card.product_id", s.productID)
		cid = r.str("id")
		check(t, r.str("buyer_phone") == "" && r.len("messages") == 0, "data admin tidak boleh bocor: %s", truncate(string(r.Raw)))
		e.call(t, "POST", "/v1/shop/chat", map[string]any{"order_no": no}, buyer).want(200).eq("id", cid).eq("card.type", "order").eq("card.order_no", no)
		check(t, e.call(t, "GET", "/v1/shop/chat", nil, buyer).want(200).len("conversations") == 0, "percakapan kosong tidak tampil di daftar")
	})

	t.Run("kirim pesan & validasi", func(t *testing.T) {
		base := "/v1/shop/chat/" + cid + "/messages"
		e.call(t, "POST", base, map[string]any{"body": "   "}, buyer).wantErr(422, "validation", "body")
		e.call(t, "POST", base, map[string]any{"body": strings.Repeat("a", 2001)}, buyer).wantErr(422, "validation", "body")
		e.call(t, "POST", base, map[string]any{"body": "hi", "product_id": o.productID}, buyer).wantErr(422, "validation", "product_id")
		e.call(t, "POST", base, map[string]any{"body": "hi", "photo_key": "private/chat/nope.jpg"}, buyer).wantErr(422, "validation", "photo_key")
		e.call(t, "POST", base, map[string]any{"body": "hi"}, other).want(404)
		e.call(t, "GET", "/v1/shop/chat/"+cid, nil, other).want(404)
		e.call(t, "GET", "/v1/shop/chat/not-a-uuid", nil, buyer).want(404)
		m := e.call(t, "POST", base, map[string]any{"body": "<b>Hi</b> is size 42 available?\u0007", "product_id": s.productID}, buyer).want(201)
		m.eq("sender", "buyer").eq("body", "<b>Hi</b> is size 42 available?").eq("card.type", "product")
		l := e.call(t, "GET", "/v1/shop/chat", nil, buyer).want(200)
		check(t, l.len("conversations") == 1 && l.str("unread") == "0", "daftar pembeli: %s", truncate(string(l.Raw)))
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'chat.new' AND store_id = $1`, s.storeID) == "2", "pemilik & staf harus dapat notifikasi")
		check(t, strings.Contains(e.scalar(t, `SELECT link FROM notifications WHERE kind = 'chat.new' ORDER BY id DESC LIMIT 1`), "/chat/"+cid), "link notifikasi ke percakapan")
		// pesan kedua dalam 10 menit → tanpa notifikasi baru
		e.call(t, "POST", base, map[string]any{"body": "Also, do you deliver today?"}, buyer).want(201)
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'chat.new'`) == "2", "notifikasi harus dibatasi per 10 menit")
	})

	t.Run("toko: daftar, belum dibaca, staf membalas", func(t *testing.T) {
		e.call(t, "GET", sb, nil, buyer).wantErr(404, "store_not_found")
		e.call(t, "GET", "/v1/seller/stores/"+o.storeID+"/chat/"+cid, nil, o.token).want(404)
		l := e.call(t, "GET", sb, nil, s.token).want(200)
		l.eq("unread", "2").eq("waiting", "1").eq("conversations.0.id", cid)
		check(t, l.str("conversations.0.buyer_phone") == "", "toko tidak boleh melihat HP pembeli")
		e.call(t, "GET", "/v1/seller/stores/"+s.storeID+"/today", nil, s.token).want(200).eq("chat_unread", "1")
		d := e.call(t, "GET", sb+"/"+cid, nil, s.token).want(200)
		d.eq("messages.0.card.type", "product").eq("product_id", s.productID)
		check(t, d.len("messages") == 2, "dua pesan: %s", truncate(string(d.Raw)))
		e.call(t, "GET", sb, nil, s.token).want(200).eq("unread", "0").eq("waiting", "1")

		r := e.call(t, "POST", sb+"/"+cid+"/messages", map[string]any{"body": "Yes, size 42 is in stock!"}, staff).want(201)
		r.eq("sender", "seller")
		e.call(t, "GET", "/v1/seller/stores/"+s.storeID+"/today", nil, s.token).want(200).eq("chat_unread", "0")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'chat.reply'`) == "1", "pembeli harus dapat notifikasi")
		check(t, strings.Contains(e.scalar(t, `SELECT link FROM notifications WHERE kind = 'chat.reply'`), "/chat#"+cid), "link notifikasi pembeli")
		e.call(t, "POST", sb+"/"+cid+"/messages", map[string]any{"body": "We deliver today."}, s.token).want(201)
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'chat.reply'`) == "1", "notifikasi pembeli dibatasi")
		e.call(t, "POST", sb+"/"+cid+"/messages", map[string]any{"body": "x", "product_id": o.productID}, s.token).wantErr(422, "validation", "product_id")
	})

	t.Run("pembeli: belum dibaca, polling, kartu pesanan", func(t *testing.T) {
		e.call(t, "GET", "/v1/shop/chat", nil, buyer).want(200).eq("unread", "2").eq("conversations.0.last_sender", "seller")
		full := e.call(t, "GET", "/v1/shop/chat/"+cid, nil, buyer).want(200)
		check(t, full.len("messages") == 4 && full.str("messages.2.sender_name") == "", "nama staf tidak dikirim ke pembeli: %s", truncate(string(full.Raw)))
		last := full.str("messages.3.id")
		e.call(t, "GET", "/v1/shop/chat", nil, buyer).want(200).eq("unread", "0")
		check(t, e.call(t, "GET", "/v1/shop/chat/"+cid+"?after="+last, nil, buyer).want(200).len("messages") == 0, "polling tanpa pesan baru")
		e.call(t, "POST", "/v1/shop/chat/"+cid+"/messages", map[string]any{"body": "Where is my order?", "order_no": no}, buyer).want(201).eq("card.type", "order").eq("card.order_no", no)
		p := e.call(t, "GET", "/v1/shop/chat/"+cid+"?after="+last, nil, buyer).want(200)
		check(t, p.len("messages") == 1, "polling: hanya pesan baru: %s", truncate(string(p.Raw)))
	})

	var mid string
	t.Run("foto chat privat", func(t *testing.T) {
		up := e.chatUpload(t, buyer, "/v1/shop/chat/uploads", jpegBytes).want(201)
		key := up.str("key")
		check(t, strings.HasPrefix(key, "private/chat/") && up.str("url") == "", "foto chat harus privat: %s", truncate(string(up.Raw)))
		e.chatUpload(t, buyer, "/v1/shop/chat/uploads", []byte("not an image at all")).want(415)
		e.call(t, "POST", sb+"/"+cid+"/messages", map[string]any{"photo_key": key}, s.token).wantErr(422, "validation", "photo_key")
		m := e.call(t, "POST", "/v1/shop/chat/"+cid+"/messages", map[string]any{"photo_key": key}, buyer).want(201)
		m.eq("has_photo", "true")
		mid = m.str("id")
		check(t, e.call(t, "GET", "/v1/shop/chat/"+cid+"/messages/"+mid+"/photo", nil, buyer).Status == 200, "pembeli melihat fotonya")
		check(t, e.call(t, "GET", sb+"/"+cid+"/messages/"+mid+"/photo", nil, staff).Status == 200, "staf toko melihat foto")
		e.call(t, "GET", "/v1/shop/chat/"+cid+"/messages/"+mid+"/photo", nil, other).want(404)
		e.call(t, "GET", "/v1/seller/stores/"+o.storeID+"/chat/"+cid+"/messages/"+mid+"/photo", nil, o.token).want(404)
		check(t, e.scalar(t, `SELECT last_preview FROM chat_conversations WHERE id = $1`, cid) == "📷", "cuplikan foto")
		su := e.chatUpload(t, s.token, sb+"/uploads", jpegBytes).want(201)
		e.call(t, "POST", sb+"/"+cid+"/messages", map[string]any{"photo_key": su.str("key"), "body": "Here it is"}, s.token).want(201).eq("has_photo", "true")
	})

	t.Run("admin: hanya support.handle, akses dicatat", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/chats", nil, mod).want(403)
		l := e.call(t, "GET", "/admin/api/chats?q="+no, nil, cs).want(200)
		check(t, l.len("conversations") == 1 && l.str("conversations.0.id") == cid, "cari nomor pesanan: %s", truncate(string(l.Raw)))
		check(t, e.call(t, "GET", "/admin/api/chats?q=nothing-matches-this", nil, cs).want(200).len("conversations") == 0, "pencarian kosong")
		before := e.scalar(t, `SELECT seller_unread FROM chat_conversations WHERE id = $1`, cid)
		d := e.call(t, "GET", "/admin/api/chats/"+cid, nil, cs).want(200)
		check(t, d.str("buyer_phone") != "" && d.len("messages") >= 6, "admin melihat isi lengkap: %s", truncate(string(d.Raw)))
		check(t, e.scalar(t, `SELECT count(*) FROM audit_logs WHERE action = 'chat.view' AND entity_id = $1`, cid) == "1", "akses admin harus dicatat di audit")
		check(t, e.scalar(t, `SELECT seller_unread FROM chat_conversations WHERE id = $1`, cid) == before, "admin tidak mengubah tanda dibaca")
		check(t, e.call(t, "GET", "/admin/api/chats/"+cid+"/messages/"+mid+"/photo", nil, cs).Status == 200, "admin melihat foto")
		e.call(t, "GET", "/admin/api/chats/00000000-0000-0000-0000-000000000000", nil, cs).want(404)
	})

	t.Run("batas kirim & akun tamu", func(t *testing.T) {
		uid := e.scalar(t, `SELECT buyer_id::text FROM chat_conversations WHERE id = $1`, cid)
		e.sql(t, `INSERT INTO chat_messages (conversation_id, sender_type, sender_id, body) SELECT $1, 'buyer', $2, 'flood' FROM generate_series(1, 20)`, cid, uid)
		e.call(t, "POST", "/v1/shop/chat/"+cid+"/messages", map[string]any{"body": "one more"}, buyer).wantErr(429, "chat_rate_limited")
		g := e.call(t, "POST", "/v1/app/auth/guest", nil, "").want(200)
		e.call(t, "GET", "/v1/shop/chat", nil, g.str("access_token")).wantErr(403, "guest_not_allowed")
	})

	t.Run("halaman website", func(t *testing.T) {
		r := e.page(t, "/en/chat")
		check(t, r.Status == 200 && strings.Contains(string(r.Raw), "/assets/chat.js"), "halaman /en/chat harus ada")
		slug := e.scalar(t, `SELECT slug FROM product_translations WHERE product_id = $1 AND lang = 'en'`, s.productID)
		p := e.page(t, "/en/p/"+slug)
		check(t, strings.Contains(string(p.Raw), "/en/chat?store="+s.storeID+"&amp;product="+s.productID), "tombol chat di halaman produk")
		st := e.page(t, "/km/s/"+e.scalar(t, `SELECT slug FROM stores WHERE id = $1`, s.storeID))
		check(t, strings.Contains(string(st.Raw), "/km/chat?store="+s.storeID), "tombol chat di halaman toko")
	})
}
