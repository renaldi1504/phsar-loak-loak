// Package payment: pembayaran KHQR lewat ABA PayWay (dan versi palsu untuk pengembangan).
//
// Prinsip keamanan: callback dari luar TIDAK PERNAH dipercaya begitu saja. Callback hanya memicu
// pengecekan ke API check-transaction PayWay (ditandatangani dengan API key kita); hanya hasil
// pengecekan itu — status APPROVED dan nominal yang sama — yang membuat pesanan lunas.
package payment

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    string `json:"price"`
}

type QRRequest struct {
	TranID      string
	Amount      string // "12.50" (USD)
	FirstName   string
	Phone       string
	Items       []Item
	CallbackURL string
	LifetimeMin int
}

type QR struct {
	String   string // isi KHQR (EMV)
	Image    string // data:image/png;base64,… (bisa kosong untuk provider palsu)
	Deeplink string // buka aplikasi ABA Mobile
}

// Status hasil pengecekan resmi.
type Status struct {
	State       string // APPROVED · PENDING · DECLINED · CANCELLED · REFUNDED · NOT_FOUND
	Amount      string // nominal yang dibayar
	APV         string
	PaymentType string
	Raw         json.RawMessage
}

func (s Status) Paid() bool { return s.State == "APPROVED" }

type Provider interface {
	Name() string
	CreateQR(ctx context.Context, r QRRequest) (QR, error)
	Check(ctx context.Context, tranID string) (Status, error)
	// VerifyCallback: tanda tangan header cocok (hanya dicatat; keputusan tetap dari Check).
	VerifyCallback(header http.Header, body []byte) bool
}

// ---------------------------------------------------------------------------
// ABA PayWay
// ---------------------------------------------------------------------------

type PayWay struct {
	BaseURL    string // https://checkout-sandbox.payway.com.kh (sandbox)
	MerchantID string
	APIKey     string
	HTTP       *http.Client
}

func (p *PayWay) Name() string { return "payway" }

func (p *PayWay) client() *http.Client {
	if p.HTTP != nil {
		return p.HTTP
	}
	return &http.Client{Timeout: 20 * time.Second}
}

func (p *PayWay) sign(parts ...string) string {
	m := hmac.New(sha512.New, []byte(p.APIKey))
	m.Write([]byte(strings.Join(parts, "")))
	return base64.StdEncoding.EncodeToString(m.Sum(nil))
}

func reqTime() string { return time.Now().UTC().Format("20060102150405") }

func b64(s string) string {
	if s == "" {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// CreateQR: POST /api/payment-gateway/v1/payments/generate-qr
// Urutan hash sesuai dokumentasi PayWay: req_time, merchant_id, tran_id, amount, items, first_name, last_name,
// email, phone, purchase_type, payment_option, callback_url, return_deeplink, currency, custom_fields,
// return_params, payout, lifetime, qr_image_template.
func (p *PayWay) CreateQR(ctx context.Context, r QRRequest) (QR, error) {
	// Batas PayWay: items (base64) maks. 500 karakter, first_name maks. 20, qr_image_template wajib.
	items := fitItems(r.Items, 500)
	f := map[string]string{
		"req_time": reqTime(), "merchant_id": p.MerchantID, "tran_id": r.TranID, "amount": r.Amount, "items": items,
		"first_name": cut(r.FirstName, 20), "last_name": "", "email": "", "phone": r.Phone, "purchase_type": "purchase",
		"payment_option": "abapay_khqr", "callback_url": b64(r.CallbackURL), "return_deeplink": "", "currency": "USD",
		"custom_fields": "", "return_params": "", "payout": "", "lifetime": fmt.Sprint(r.LifetimeMin), "qr_image_template": QRTemplate,
	}
	order := []string{"req_time", "merchant_id", "tran_id", "amount", "items", "first_name", "last_name", "email", "phone",
		"purchase_type", "payment_option", "callback_url", "return_deeplink", "currency", "custom_fields", "return_params",
		"payout", "lifetime", "qr_image_template"}
	parts := make([]string, len(order))
	body := map[string]any{}
	for i, k := range order {
		parts[i] = f[k]
		if f[k] != "" {
			body[k] = f[k]
		}
	}
	body["lifetime"] = r.LifetimeMin
	body["hash"] = p.sign(parts...)

	var out struct {
		QRString string `json:"qrString"`
		QRImage  string `json:"qrImage"`
		Deeplink string `json:"abapay_deeplink"`
		Status   struct {
			Code    any    `json:"code"`
			Message string `json:"message"`
		} `json:"status"`
	}
	raw, err := p.post(ctx, "/api/payment-gateway/v1/payments/generate-qr", body, &out)
	if err != nil {
		return QR{}, err
	}
	code := fmt.Sprint(out.Status.Code)
	if (code != "0" && code != "00") || out.QRString == "" {
		return QR{}, fmt.Errorf("payway generate-qr: kode %s: %s (%s)", code, out.Status.Message, trunc(string(raw)))
	}
	return QR{String: out.QRString, Image: out.QRImage, Deeplink: out.Deeplink}, nil
}

// QRTemplate: desain gambar QR dari PayWay (wajib diisi di generate-qr).
const QRTemplate = "template3_color"

// fitItems: daftar barang (hanya informasi di PayWay) dipotong sampai base64-nya muat; kosong kalau satu pun tidak muat.
func fitItems(items []Item, max int) string {
	for n := len(items); n > 0; n-- {
		list := make([]Item, n)
		for i := range list {
			list[i] = items[i]
			list[i].Name = cut(list[i].Name, 40)
		}
		b, _ := json.Marshal(list)
		if enc := b64(string(b)); len(enc) <= max {
			return enc
		}
	}
	return ""
}

// cut memotong teks ke n huruf (bukan byte), tanpa spasi di ujung.
func cut(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) > n {
		r = r[:n]
	}
	return strings.TrimSpace(string(r))
}

// Check: POST /api/payment-gateway/v1/payments/check-transaction-2 (hash = req_time + merchant_id + tran_id).
func (p *PayWay) Check(ctx context.Context, tranID string) (Status, error) {
	rt := reqTime()
	body := map[string]any{"req_time": rt, "merchant_id": p.MerchantID, "tran_id": tranID, "hash": p.sign(rt, p.MerchantID, tranID)}
	var out struct {
		Data struct {
			PaymentStatus string `json:"payment_status"`
			PaymentAmount any    `json:"payment_amount"`
			TotalAmount   any    `json:"total_amount"`
			APV           string `json:"apv"`
		} `json:"data"`
		Status struct {
			Code    any    `json:"code"`
			Message string `json:"message"`
		} `json:"status"`
	}
	raw, err := p.post(ctx, "/api/payment-gateway/v1/payments/check-transaction-2", body, &out)
	if err != nil {
		return Status{}, err
	}
	switch code := fmt.Sprint(out.Status.Code); code {
	case "00", "0":
	case "6":
		return Status{State: "NOT_FOUND", Raw: raw}, nil
	default:
		return Status{}, fmt.Errorf("payway check-transaction: kode %s: %s", code, out.Status.Message)
	}
	amt := fmt.Sprint(out.Data.PaymentAmount)
	if out.Data.PaymentAmount == nil {
		amt = fmt.Sprint(out.Data.TotalAmount)
	}
	return Status{State: strings.ToUpper(out.Data.PaymentStatus), Amount: amt, APV: out.Data.APV, Raw: raw}, nil
}

// VerifyCallback: header X-PAYWAY-HMAC-SHA512 = base64(HMAC-SHA512(api key, nilai payload diurutkan per key)).
func (p *PayWay) VerifyCallback(h http.Header, body []byte) bool {
	got := h.Get("X-PAYWAY-HMAC-SHA512")
	if got == "" {
		return false
	}
	var m map[string]any
	if json.Unmarshal(body, &m) != nil {
		return false
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		switch v := m[k].(type) {
		case string:
			b.WriteString(v)
		case nil:
		case map[string]any, []any:
			j, _ := json.Marshal(v)
			b.Write(j)
		default:
			j, _ := json.Marshal(v)
			b.Write(j)
		}
	}
	want := p.sign(b.String())
	return hmac.Equal([]byte(got), []byte(want))
}

func (p *PayWay) post(ctx context.Context, path string, body any, out any) (json.RawMessage, error) {
	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.BaseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := p.client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("payway: tidak bisa menghubungi server: %w", err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	if res.StatusCode >= 300 && res.StatusCode != 400 {
		return raw, fmt.Errorf("payway: HTTP %d: %s", res.StatusCode, trunc(string(raw)))
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return raw, fmt.Errorf("payway: balasan tidak bisa dibaca: %s", trunc(string(raw)))
	}
	return raw, nil
}

func trunc(s string) string {
	if len(s) > 300 {
		return s[:300]
	}
	return s
}

// ---------------------------------------------------------------------------
// Palsu (pengembangan & QA): tidak ada uang. "Bayar" lewat Fake.Pay / tombol simulasi.
// ---------------------------------------------------------------------------

type Fake struct {
	mu     sync.Mutex
	paid   map[string]string // tran_id → nominal
	Secret string            // untuk menandatangani callback palsu di tes
}

func (f *Fake) Name() string { return "fake" }

func (f *Fake) CreateQR(_ context.Context, r QRRequest) (QR, error) {
	if r.TranID == "" || r.Amount == "" {
		return QR{}, errors.New("fake: tran_id/amount kosong")
	}
	return QR{String: "00020101021229180014fake.phsar.test" + r.TranID + "5303840540" + r.Amount + "6304TEST"}, nil
}

// Pay mensimulasikan pembeli membayar (dipakai tombol "simulasi bayar" di mode pengembangan).
func (f *Fake) Pay(tranID, amount string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.paid == nil {
		f.paid = map[string]string{}
	}
	f.paid[tranID] = amount
}

func (f *Fake) Check(_ context.Context, tranID string) (Status, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if a, ok := f.paid[tranID]; ok {
		return Status{State: "APPROVED", Amount: a, APV: "FAKE" + tranID[len(tranID)-4:], PaymentType: "KHQR"}, nil
	}
	return Status{State: "PENDING"}, nil
}

func (f *Fake) VerifyCallback(h http.Header, body []byte) bool {
	if f.Secret == "" {
		return false
	}
	return (&PayWay{APIKey: f.Secret}).VerifyCallback(h, body)
}
