package notify

// Pengaturan notifikasi per user per audience: grup (Pesanan, Chat, …) bisa dimatikan per jalur
// (suara, push browser, Telegram), plus jam tenang dan volume. Lonceng (daftar di aplikasi) selalu terisi.

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// Audience: untuk siapa sebuah notifikasi.
const (
	AudSeller  = "seller"
	AudBuyer   = "buyer"
	AudCourier = "courier"
	AudAdmin   = "admin"
)

// Groups: grup yang bisa diatur, sesuai urutan tampil.
var Groups = map[string][]string{
	AudSeller:  {"orders", "chat", "returns", "payments", "store", "promo"},
	AudBuyer:   {"orders", "returns", "chat", "support"},
	AudCourier: {"tasks"},
	AudAdmin:   {"stores", "products", "refunds", "reports", "support", "returns", "payouts", "deliveries"},
}

// GroupOf: grup sebuah jenis notifikasi ("" = tidak bisa dimatikan, mis. notifikasi uji).
func GroupOf(kind string) string {
	switch kind {
	case TestKind:
		return ""
	case ChatNew, ChatReply:
		return "chat"
	case PayoutSent:
		return "payments"
	case TicketReply:
		return "support"
	case "admin.store_review", "admin.pickup_review":
		return "stores"
	case "admin.product_review":
		return "products"
	case "admin.refund_pending", "admin.return_refund":
		return "refunds"
	case "admin.report_new":
		return "reports"
	case "admin.ticket_new", "admin.ticket_reply":
		return "support"
	case "admin.return_dispute":
		return "returns"
	case "admin.payout_pending":
		return "payouts"
	case "admin.delivery_new", "admin.delivery_failed":
		return "deliveries"
	}
	switch {
	case strings.HasPrefix(kind, "order."), strings.HasPrefix(kind, "buyer.order_"), kind == "buyer.delivery_returning":
		return "orders"
	case strings.HasPrefix(kind, "return."), strings.HasPrefix(kind, "buyer.return_"):
		return "returns"
	case strings.HasPrefix(kind, "courier."):
		return "tasks"
	case strings.HasPrefix(kind, "ad."), strings.HasPrefix(kind, "flash."):
		return "promo"
	}
	return "store" // toko, produk, KYC, titik jemput, staf, strike
}

// Sound: bunyi untuk sebuah jenis notifikasi dan berapa kali diputar.
//
//	order  — pesanan baru untuk seller: 3× (permintaan Komodo, 18 Nov)
//	alert  — antrean admin, tugas kurir, peringatan toko
//	chat   — pesan chat
//	ding   — lainnya
func Sound(kind string) (string, int) {
	switch {
	case kind == OrderNew:
		return "order", 3
	case kind == ChatNew || kind == ChatReply:
		return "chat", 1
	case strings.HasPrefix(kind, "admin."), strings.HasPrefix(kind, "courier."),
		kind == "store.strike", kind == "store.suspended", kind == OrderSellerTimeout, kind == OrderPackTimeout:
		return "alert", 1
	}
	return "ding", 1
}

// GroupPref: jalur per grup.
type GroupPref struct {
	Sound    bool `json:"sound"`
	Push     bool `json:"push"`
	Telegram bool `json:"telegram"`
}

// Prefs: notification_prefs.prefs.
type Prefs struct {
	Groups    map[string]GroupPref `json:"groups"`
	QuietFrom string               `json:"quiet_from"` // "23:00" — kosong = tanpa jam tenang
	QuietTo   string               `json:"quiet_to"`   // "07:00"
	Volume    int                  `json:"volume"`     // 0–100
}

// DefaultPrefs: semua grup nyala; admin tidak dikirimi Telegram kecuali menyalakannya sendiri.
func DefaultPrefs(aud string) Prefs {
	p := Prefs{Groups: map[string]GroupPref{}, Volume: 80}
	for _, g := range Groups[aud] {
		p.Groups[g] = GroupPref{Sound: true, Push: true, Telegram: aud != AudAdmin}
	}
	return p
}

var hhmm = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)

// ErrBadPrefs: isian pengaturan tidak valid.
var ErrBadPrefs = errors.New("pengaturan notifikasi tidak valid")

// Clean: buang grup yang tidak dikenal, isi grup yang hilang dengan default, cek jam & volume.
func (p Prefs) Clean(aud string) (Prefs, error) {
	out := DefaultPrefs(aud)
	for g, v := range p.Groups {
		if _, ok := out.Groups[g]; ok {
			out.Groups[g] = v
		}
	}
	if (p.QuietFrom == "") != (p.QuietTo == "") ||
		(p.QuietFrom != "" && (!hhmm.MatchString(p.QuietFrom) || !hhmm.MatchString(p.QuietTo) || p.QuietFrom == p.QuietTo)) {
		return out, ErrBadPrefs
	}
	if p.Volume < 0 || p.Volume > 100 {
		return out, ErrBadPrefs
	}
	out.QuietFrom, out.QuietTo, out.Volume = p.QuietFrom, p.QuietTo, p.Volume
	return out, nil
}

func (p Prefs) group(kind string) GroupPref {
	g := GroupOf(kind)
	if g == "" {
		return GroupPref{Sound: true, Push: true, Telegram: true}
	}
	if v, ok := p.Groups[g]; ok {
		return v
	}
	return GroupPref{Sound: true, Push: true, Telegram: true}
}

var phnomPenh = func() *time.Location {
	if l, err := time.LoadLocation("Asia/Phnom_Penh"); err == nil {
		return l
	}
	return time.FixedZone("ICT", 7*3600)
}()

// Quiet: sekarang di dalam jam tenang (waktu Kamboja; boleh melewati tengah malam, mis. 23:00–07:00).
func (p Prefs) Quiet(now time.Time) bool {
	if p.QuietFrom == "" {
		return false
	}
	cur := now.In(phnomPenh).Format("15:04")
	if p.QuietFrom < p.QuietTo {
		return cur >= p.QuietFrom && cur < p.QuietTo
	}
	return cur >= p.QuietFrom || cur < p.QuietTo
}

// SoundFor: bunyi yang diputar di browser ("none" kalau suara grup mati atau jam tenang).
func (p Prefs) SoundFor(kind string, now time.Time) (string, int) {
	if kind != TestKind && (!p.group(kind).Sound || p.Quiet(now)) {
		return "none", 0
	}
	return Sound(kind)
}

// PushAllowed / TelegramAllowed: jalur grup ini nyala.
func (p Prefs) PushAllowed(kind string) bool     { return p.group(kind).Push }
func (p Prefs) TelegramAllowed(kind string) bool { return p.group(kind).Telegram }

type rowQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// LoadPrefs: pengaturan tersimpan digabung dengan default (tanpa baris = default).
func LoadPrefs(ctx context.Context, q rowQuerier, userID, aud string) (Prefs, error) {
	var raw []byte
	err := q.QueryRow(ctx, `SELECT prefs FROM notification_prefs WHERE user_id = $1 AND audience = $2`, userID, aud).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return DefaultPrefs(aud), nil
	}
	if err != nil {
		return DefaultPrefs(aud), err
	}
	var p Prefs
	if json.Unmarshal(raw, &p) != nil {
		return DefaultPrefs(aud), nil
	}
	if p.Groups == nil && p.Volume == 0 && p.QuietFrom == "" {
		p.Volume = 80
	}
	clean, err := p.Clean(aud)
	if err != nil { // data lama tidak valid → abaikan bagian yang rusak
		return DefaultPrefs(aud), nil
	}
	return clean, nil
}

// AudienceOfKind: audience untuk EnqueueUser (pembeli atau kurir).
func AudienceOfKind(kind string) string {
	if strings.HasPrefix(kind, "courier.") {
		return AudCourier
	}
	return AudBuyer
}
