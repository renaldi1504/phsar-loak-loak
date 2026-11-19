package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mac(key, s string) string {
	m := hmac.New(sha512.New, []byte(key))
	m.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(m.Sum(nil))
}

func TestPayWayCreateQRHashAndCheck(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		switch r.URL.Path {
		case "/api/payment-gateway/v1/payments/generate-qr":
			_, _ = w.Write([]byte(`{"status":{"code":"0","message":"Success"},"qrString":"000201KHQR","qrImage":"data:image/png;base64,AAA","abapay_deeplink":"abamobilebank://x"}`))
		case "/api/payment-gateway/v1/payments/check-transaction-2":
			_, _ = w.Write([]byte(`{"data":{"payment_status_code":0,"payment_status":"APPROVED","payment_amount":12.5,"apv":"123456"},"status":{"code":"00","message":"Success"}}`))
		}
	}))
	defer srv.Close()
	p := &PayWay{BaseURL: srv.URL, MerchantID: "m1", APIKey: "secret"}

	qr, err := p.CreateQR(context.Background(), QRRequest{TranID: "PABC", Amount: "12.50", FirstName: "Dara", Phone: "012345678",
		Items: []Item{{Name: "Tee", Quantity: 1, Price: "12.50"}}, CallbackURL: "https://x/webhooks/payway", LifetimeMin: 15})
	if err != nil {
		t.Fatal(err)
	}
	if qr.String != "000201KHQR" || qr.Deeplink == "" || qr.Image == "" {
		t.Fatalf("qr: %+v", qr)
	}
	s := func(k string) string { v, _ := got[k].(string); return v }
	items := base64.StdEncoding.EncodeToString([]byte(`[{"name":"Tee","quantity":1,"price":"12.50"}]`))
	if s("items") != items {
		t.Fatalf("items harus base64 JSON: %s", s("items"))
	}
	if s("callback_url") != base64.StdEncoding.EncodeToString([]byte("https://x/webhooks/payway")) {
		t.Fatal("callback_url harus base64")
	}
	want := mac("secret", s("req_time")+"m1"+"PABC"+"12.50"+items+"Dara"+""+""+"012345678"+"purchase"+"abapay_khqr"+s("callback_url")+""+"USD"+""+""+""+"15"+QRTemplate)
	if s("qr_image_template") != QRTemplate {
		t.Fatal("qr_image_template wajib dikirim")
	}
	if s("hash") != want {
		t.Fatal("hash generate-qr tidak sesuai urutan dokumentasi PayWay")
	}
	if len(s("req_time")) != 14 {
		t.Fatalf("req_time harus YYYYMMDDHHmmss: %s", s("req_time"))
	}

	st, err := p.Check(context.Background(), "PABC")
	if err != nil {
		t.Fatal(err)
	}
	if !st.Paid() || st.Amount != "12.5" || st.APV != "123456" {
		t.Fatalf("check: %+v", st)
	}
	if s("hash") != mac("secret", s("req_time")+"m1"+"PABC") {
		t.Fatal("hash check-transaction salah")
	}
}

func TestPayWayErrorsAndCallbackSignature(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "check") {
			_, _ = w.Write([]byte(`{"status":{"code":"6","message":"Transaction not found"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":{"code":"1","message":"Wrong hash"}}`))
	}))
	defer srv.Close()
	p := &PayWay{BaseURL: srv.URL, MerchantID: "m1", APIKey: "secret"}
	if _, err := p.CreateQR(context.Background(), QRRequest{TranID: "P1", Amount: "1.00", LifetimeMin: 5}); err == nil {
		t.Fatal("kode error dari PayWay harus jadi error")
	}
	if st, err := p.Check(context.Background(), "P1"); err != nil || st.State != "NOT_FOUND" || st.Paid() {
		t.Fatalf("not found: %+v %v", st, err)
	}

	body := []byte(`{"tran_id":"P1","status":0,"apv":"9","payment_amount":"1.00"}`)
	// diurutkan per key: apv, payment_amount, status, tran_id
	h := http.Header{}
	h.Set("X-PAYWAY-HMAC-SHA512", mac("secret", "9"+"1.00"+"0"+"P1"))
	if !p.VerifyCallback(h, body) {
		t.Fatal("tanda tangan callback yang benar ditolak")
	}
	h.Set("X-PAYWAY-HMAC-SHA512", mac("other", "91.000P1"))
	if p.VerifyCallback(h, body) {
		t.Fatal("tanda tangan palsu diterima")
	}
}

func TestPayWayFieldLimits(t *testing.T) {
	if got := cut("  Sreyneang Sophea Chanthou  ", 20); len([]rune(got)) > 20 || got != "Sreyneang Sophea Cha" {
		t.Fatalf("cut: %q", got)
	}
	if got := cut("សុភា ចាន់ធូ សុខា ស្រីនាង", 20); len([]rune(got)) > 20 {
		t.Fatalf("cut khmer: %q", got)
	}
	many := make([]Item, 30)
	for i := range many {
		many[i] = Item{Name: strings.Repeat("Product name ", 6), Quantity: 1, Price: "1.00"}
	}
	enc := fitItems(many, 500)
	if enc == "" || len(enc) > 500 {
		t.Fatalf("items harus dipotong sampai muat 500 karakter: %d", len(enc))
	}
	if fitItems(nil, 500) != "" {
		t.Fatal("tanpa barang = kosong")
	}
}
