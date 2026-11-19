package qa

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"phsar/internal/auth"
	"phsar/internal/config"
	"phsar/internal/httpx"
)

// ---------------------------------------------------------------------------
// askbot: draf jawaban untuk pertanyaan pembeli + FAQ toko + pengaturan super admin.
// ---------------------------------------------------------------------------

func TestAskbot(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	sa := newAdmin(t, e, "super_admin")
	s := newSellable(t, e, mod, "Askbot linen shirt", []string{"5.00", "7.00"}, []int{3, 0})
	other := newSellable(t, e, mod, "Askbot other shop item", []string{"9.00"}, []int{1})
	stranger, _ := e.login(t)
	base := "/v1/seller/stores/" + s.storeID + "/assist"
	suggest := base + "/suggest"

	// Staf toko (lewat undangan).
	staff, _ := e.login(t)
	inv := e.call(t, "POST", "/v1/seller/stores/"+s.storeID+"/invites", nil, s.token).want(201)
	e.call(t, "POST", "/v1/seller/invites/"+inv.str("token")+"/accept", nil, staff).want(200)

	t.Run("tanpa API key: jawaban berbasis aturan dari data produk", func(t *testing.T) {
		r := e.call(t, "POST", suggest, map[string]any{"question": "How much is it and is it in stock?", "product_id": s.productID}, s.token).want(200)
		reply := r.str("reply")
		check(t, strings.Contains(reply, "$5.00") && strings.Contains(reply, "$7.00"), "harga varian tidak ada: %q", reply)
		check(t, len(r.Body) == 1, "kontrak: hanya field reply, dapat %s", truncate(string(r.Raw)))
		// Bahasa mengikuti pertanyaan.
		km := e.call(t, "POST", suggest, map[string]any{"question": "តម្លៃប៉ុន្មាន?", "product_id": s.productID}, s.token).want(200).str("reply")
		check(t, strings.Contains(km, "សួស្តី"), "jawaban harus berbahasa Khmer: %q", km)
		zh := e.call(t, "POST", suggest, map[string]any{"question": "有保修吗？", "product_id": s.productID}, staff).want(200).str("reply")
		check(t, strings.Contains(zh, "您好"), "jawaban harus berbahasa Mandarin: %q", zh)
		// Tidak ada datanya → tidak mengarang.
		x := e.call(t, "POST", suggest, map[string]any{"question": "Is it made of cotton?", "product_id": s.productID}, s.token).want(200).str("reply")
		check(t, strings.Contains(x, "We will check"), "harus bilang akan dicek: %q", x)
		// Tanpa produk pun boleh (FAQ / jam toko).
		e.call(t, "POST", suggest, map[string]any{"question": "Hello?"}, s.token).want(200)
		check(t, e.scalar(t, `SELECT count(*) FROM askbot_usage WHERE store_id = $1 AND mode = 'rules'`, s.storeID) == "5", "pemakaian tidak tercatat")
	})

	t.Run("izin & produk toko lain ditolak", func(t *testing.T) {
		e.call(t, "POST", suggest, map[string]any{"question": "Price?"}, stranger).want(404)
		e.call(t, "POST", suggest, map[string]any{"question": "Price?"}, "").want(401)
		e.call(t, "POST", suggest, map[string]any{"question": "Price?", "product_id": other.productID}, s.token).wantErr(404, "not_found")
		e.call(t, "GET", base+"/faqs", nil, stranger).want(404)
	})

	t.Run("validasi", func(t *testing.T) {
		e.call(t, "POST", suggest, map[string]any{"question": "   "}, s.token).wantErr(422, "validation", "question")
		e.call(t, "POST", suggest, map[string]any{"question": strings.Repeat("a", 1001)}, s.token).wantErr(422, "validation", "question")
		e.call(t, "POST", suggest, map[string]any{"question": "Price?", "product_id": "abc"}, s.token).wantErr(422, "validation", "product_id")
		e.call(t, "POST", suggest, map[string]any{"question": "Price?", "x": 1}, s.token).wantErr(400, "bad_json")
	})

	var storeFAQ, prodFAQ string
	t.Run("FAQ: CRUD oleh pemilik & staf", func(t *testing.T) {
		f := e.call(t, "POST", base+"/faqs", map[string]any{"question": "Do you deliver to Otres beach?", "answer": "Yes, Phsar couriers deliver to Otres every day."}, s.token).want(201)
		f.eq("lang", "en").eq("product_id", "")
		storeFAQ = f.str("id")
		p := e.call(t, "POST", base+"/faqs", map[string]any{"question": "How do I measure the chest?", "answer": "Measure under the arms, around the widest part.",
			"lang": "en", "product_id": s.productID}, staff).want(201)
		p.eq("product_id", s.productID).eq("product_name", "Askbot linen shirt")
		prodFAQ = p.str("id")
		km := e.call(t, "POST", base+"/faqs", map[string]any{"question": "ដឹកជញ្ជូនប៉ុន្មានថ្ងៃ?", "answer": "ក្នុងថ្ងៃតែមួយ។"}, s.token).want(201).eq("lang", "km")

		r := e.call(t, "GET", base+"/faqs", nil, s.token).want(200)
		check(t, r.len("faqs") == 3, "harus 3 FAQ, dapat %d", r.len("faqs"))
		r.eq("enabled", "true").eq("ai", "false")
		check(t, e.call(t, "GET", base+"/faqs?product_id=store", nil, s.token).want(200).len("faqs") == 2, "filter FAQ umum salah")

		e.call(t, "PUT", base+"/faqs/"+storeFAQ, map[string]any{"question": "Do you deliver to Otres beach?", "answer": "Yes, every day before 6 pm."}, staff).want(200).
			eq("answer", "Yes, every day before 6 pm.")
		e.call(t, "DELETE", base+"/faqs/"+km.str("id"), nil, staff).want(200)
		e.call(t, "DELETE", base+"/faqs/"+km.str("id"), nil, staff).wantErr(404, "not_found")
	})

	t.Run("FAQ: validasi & milik toko lain", func(t *testing.T) {
		e.call(t, "POST", base+"/faqs", map[string]any{"question": "", "answer": "x"}, s.token).wantErr(422, "validation", "question")
		e.call(t, "POST", base+"/faqs", map[string]any{"question": "Q?", "answer": " "}, s.token).wantErr(422, "validation", "answer")
		e.call(t, "POST", base+"/faqs", map[string]any{"question": "Q?", "answer": strings.Repeat("a", 2001)}, s.token).wantErr(422, "validation", "answer")
		e.call(t, "POST", base+"/faqs", map[string]any{"question": "Q?", "answer": "A", "lang": "fr"}, s.token).wantErr(422, "validation", "lang")
		e.call(t, "POST", base+"/faqs", map[string]any{"question": "Q?", "answer": "A", "product_id": other.productID}, s.token).wantErr(404, "not_found")
		otherBase := "/v1/seller/stores/" + other.storeID + "/assist"
		e.call(t, "PUT", otherBase+"/faqs/"+storeFAQ, map[string]any{"question": "Q?", "answer": "A"}, other.token).wantErr(404, "not_found")
		e.call(t, "DELETE", otherBase+"/faqs/"+prodFAQ, nil, other.token).wantErr(404, "not_found")
	})

	t.Run("FAQ dipakai bot", func(t *testing.T) {
		r := e.call(t, "POST", suggest, map[string]any{"question": "Do you deliver to Otres beach?"}, s.token).want(200)
		check(t, strings.Contains(r.str("reply"), "every day before 6 pm"), "jawaban FAQ tidak dipakai: %q", r.str("reply"))
	})

	t.Run("FAQ produk tampil di halaman produk", func(t *testing.T) {
		slug := e.scalar(t, `SELECT slug FROM product_translations WHERE product_id = $1 AND lang = 'en'`, s.productID)
		e.page(t, "/en/p/"+slug).has(t, "Questions &amp; answers", "judul blok").has(t, "How do I measure the chest?", "pertanyaan FAQ produk").
			hasNot(t, "Otres", "FAQ umum toko tidak tampil di halaman produk")
		e.page(t, "/km/p/"+slug).has(t, "សំណួរ និងចម្លើយ", "judul Khmer")
		oslug := e.scalar(t, `SELECT slug FROM product_translations WHERE product_id = $1 AND lang = 'en'`, other.productID)
		e.page(t, "/en/p/"+oslug).hasNot(t, `id="faq"`, "produk tanpa FAQ")
	})

	t.Run("super admin: nonaktif, batas per jam, laporan", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/askbot", nil, mod).want(403)
		r := e.call(t, "GET", "/admin/api/askbot", nil, sa).want(200)
		r.eq("settings.enabled", "true").eq("settings.hourly_limit", "60").eq("ai", "false")
		check(t, r.str("usage.week") != "0" && r.str("usage.faqs") == "2", "laporan salah: %s", truncate(string(r.Raw)))
		e.call(t, "PUT", "/admin/api/askbot", map[string]any{"enabled": true, "hourly_limit": 0}, sa).wantErr(422, "validation", "hourly_limit")

		e.call(t, "PUT", "/admin/api/askbot", map[string]any{"enabled": false, "hourly_limit": 60}, sa).want(200).eq("settings.enabled", "false")
		// 503 adalah kontrak (bukan bug) → lewat server tanpa pencatat 5xx milik TestMain.
		plain := askbotServer(t, e, e.cfg)
		plain.call(t, "POST", suggest, map[string]any{"question": "Price?"}, s.token).wantErr(503, "assist_unavailable")
		e.call(t, "GET", base+"/faqs", nil, s.token).want(200).eq("enabled", "false")

		used := e.scalar(t, `SELECT count(*) FROM askbot_usage WHERE store_id = $1`, s.storeID)
		e.call(t, "PUT", "/admin/api/askbot", map[string]any{"enabled": true, "hourly_limit": mustAtoi(used) + 1}, sa).want(200)
		e.call(t, "POST", suggest, map[string]any{"question": "Price?"}, s.token).want(200)
		e.call(t, "POST", suggest, map[string]any{"question": "Price?"}, s.token).wantErr(429, "rate_limited")
		e.call(t, "POST", "/v1/seller/stores/"+other.storeID+"/assist/suggest", map[string]any{"question": "Price?"}, other.token).want(200) // batas per toko
		e.call(t, "PUT", "/admin/api/askbot", map[string]any{"enabled": true, "hourly_limit": 60}, sa).want(200)
	})

	t.Run("dengan Claude (server palsu)", func(t *testing.T) {
		var calls atomic.Int32
		var fail atomic.Bool
		fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			if fail.Load() {
				w.WriteHeader(529)
				_, _ = w.Write([]byte(`{"type":"error","error":{"type":"overloaded_error","message":"Overloaded"}}`))
				return
			}
			b, _ := io.ReadAll(r.Body)
			var in struct {
				Messages []struct {
					Content string `json:"content"`
				} `json:"messages"`
			}
			_ = json.Unmarshal(b, &in)
			txt := "No product data."
			if len(in.Messages) == 1 && strings.Contains(in.Messages[0].Content, "Askbot linen shirt") && strings.Contains(in.Messages[0].Content, "Khmer") {
				txt = "សួស្តី! អាវនេះតម្លៃ $5.00។"
			}
			out, _ := json.Marshal(map[string]any{"model": "claude-qa", "content": []any{map[string]any{"type": "text", "text": txt}},
				"usage": map[string]int{"input_tokens": 200, "output_tokens": 20}})
			_, _ = w.Write(out)
		}))
		defer fake.Close()
		cfg := e.cfg
		cfg.AnthropicKey, cfg.AnthropicModel, cfg.AnthropicBaseURL = "qa-key", "claude-qa", fake.URL
		e2 := askbotServer(t, e, cfg)

		r := e2.call(t, "POST", suggest, map[string]any{"question": "តម្លៃប៉ុន្មាន?", "product_id": s.productID}, s.token).want(200)
		r.eq("reply", "សួស្តី! អាវនេះតម្លៃ $5.00។")
		check(t, calls.Load() == 1, "Claude harus dipanggil sekali, dapat %d", calls.Load())
		check(t, e.scalar(t, `SELECT mode || '/' || model || '/' || input_tokens FROM askbot_usage WHERE store_id = $1 ORDER BY id DESC LIMIT 1`, s.storeID) == "claude/claude-qa/200",
			"pemakaian Claude tidak tercatat")
		e2.call(t, "GET", base+"/faqs", nil, s.token).want(200).eq("ai", "true")
		e2.call(t, "GET", "/admin/api/askbot", nil, sa).want(200).eq("ai", "true").eq("model", "claude-qa")

		// Claude gagal → tetap dapat draf dari aturan (bukan 5xx), error dicatat.
		fail.Store(true)
		r = e2.call(t, "POST", suggest, map[string]any{"question": "How much is it?", "product_id": s.productID}, s.token).want(200)
		check(t, strings.Contains(r.str("reply"), "$5.00"), "cadangan aturan: %q", r.str("reply"))
		check(t, e.scalar(t, `SELECT mode || '/' || (error IS NOT NULL)::text FROM askbot_usage WHERE store_id = $1 ORDER BY id DESC LIMIT 1`, s.storeID) == "rules/true",
			"error Claude tidak tercatat")
	})
}

func mustAtoi(s string) int { n, _ := strconv.Atoi(s); return n }

// askbotServer: server kedua dengan config lain (mis. API key + Claude palsu), memakai database & token yang sama.
func askbotServer(t *testing.T, e *env, cfg config.Config) *env {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := httptest.NewServer(httpx.NewRouter(cfg, log, e.pool, auth.NewService(e.pool, cfg, log)))
	t.Cleanup(srv.Close)
	e2 := *e
	e2.srv = srv
	return &e2
}
