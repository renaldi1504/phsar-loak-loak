// Package modrules: aturan moderasi produk 3 jalur.
//
//	auto   (🟢) lolos semua cek → langsung tayang, masuk daftar "Baru tayang" untuk dicek admin belakangan
//	manual (🟡) ada alasan (lihat Reason*) → menunggu admin
//	block  (🔴) kata terlarang beraksi "block" → ditolak saat dikirim (dicek di paket seller)
//
// Angka-angka disimpan di platform_settings['moderation_rules'] dan bisa diubah admin.
package modrules

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

const SettingKey = "moderation_rules"

// Rules: platform_settings['moderation_rules'].
type Rules struct {
	AutoPublish bool `json:"auto_publish"` // false = semua produk ditinjau manual (sakelar darurat)
	Probation   int  `json:"probation"`    // produk pertama toko yang wajib disetujui admin
	DiscountPct int  `json:"discount_pct"` // harga coret lebih dari X% di atas harga jual → ditinjau
	DropPct     int  `json:"drop_pct"`     // harga turun lebih dari X% dari harga saat disetujui → ditinjau
	ReportHold  int  `json:"report_hold"`  // laporan serius dari N pembeli berbeda → produk disembunyikan & ditinjau
}

func Defaults() Rules {
	return Rules{AutoPublish: true, Probation: 3, DiscountPct: 70, DropPct: 50, ReportHold: 3}
}

// Validate: batas wajar supaya admin tidak salah ketik.
func (r Rules) Validate() (field string, ok bool) {
	switch {
	case r.Probation < 0 || r.Probation > 50:
		return "probation", false
	case r.DiscountPct < 10 || r.DiscountPct > 95:
		return "discount_pct", false
	case r.DropPct < 10 || r.DropPct > 95:
		return "drop_pct", false
	case r.ReportHold < 1 || r.ReportHold > 50:
		return "report_hold", false
	}
	return "", true
}

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Load: tanpa baris / JSON rusak = Defaults(); field yang tidak ada ikut default.
func Load(ctx context.Context, q querier) (Rules, error) {
	r := Defaults()
	var raw []byte
	err := q.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = $1`, SettingKey).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return r, nil
	}
	if err != nil {
		return r, err
	}
	if json.Unmarshal(raw, &r) != nil {
		return Defaults(), nil
	}
	if _, ok := r.Validate(); !ok {
		return Defaults(), nil
	}
	return r, nil
}

// Alasan produk masuk jalur manual (kode ini juga dipakai UI seller & admin).
const (
	ReasonProbation   = "new_store"           // toko belum melewati masa percobaan
	ReasonRestricted  = "restricted_category" // kategori terbatas
	ReasonKeyword     = "keyword"             // kata dengan aksi "review"
	ReasonPriceFloor  = "price_floor"         // jauh di bawah harga wajar merek/model
	ReasonBigDiscount = "big_discount"        // harga coret tidak masuk akal
	ReasonPriceDrop   = "price_drop"          // harga anjlok setelah disetujui
	ReasonStrikes     = "store_strikes"       // toko punya strike aktif
	ReasonChanges     = "changes_requested"   // admin minta perbaikan → admin yang cek lagi
	ReasonWasRejected = "was_rejected"        // pernah ditolak
	ReasonReported    = "reported"            // dilaporkan beberapa pembeli
	ReasonManualOnly  = "auto_publish_off"    // tayang otomatis dimatikan admin
)

// Floor: aturan harga minimum per kata kunci di nama produk.
type Floor struct {
	Keyword  string          `json:"keyword"`
	MinPrice decimal.Decimal `json:"min_price"`
}

// PriceHit: detail untuk moderator.
type PriceHit struct {
	Reason  string `json:"reason"`
	Keyword string `json:"keyword,omitempty"`
	Limit   string `json:"limit,omitempty"` // harga minimum / harga saat disetujui
	Price   string `json:"price"`
}

// Kata yang menandakan aksesori ("case iphone 15" boleh $3) — cek harga minimum dilewati.
var accessoryWords = []string{
	"case", "casing", "cover", "charger", "charging", "cable", "adapter", "strap", "band", "screen protector", "tempered glass",
	"protector", "film", "skin", "sticker", "stand", "holder", "pouch", "sleeve", "dust bag", "box only", "empty box", "replacement",
	"compatible", "keychain", "mount", "dock", "earpad", "ear tips", "lanyard", "dummy", "toy", "poster",
	"手机壳", "保护套", "保护壳", "充电器", "数据线", "贴膜", "钢化膜", "表带", "支架", "包装盒",
	"ស្រោម", "ខ្សែសាក", "ឆ្នាំងសាក", "កញ្ចក់ការពារ",
}

// IsAccessory: nama produk menyebut aksesori.
func IsAccessory(name string) bool {
	n := strings.ToLower(name)
	for _, w := range accessoryWords {
		if ContainsWord(n, w) {
			return true
		}
	}
	return false
}

// PriceCheck: input cek harga untuk satu produk.
type PriceCheck struct {
	Name        string
	Prices      []decimal.Decimal  // harga varian aktif
	Compares    []*decimal.Decimal // harga coret per varian (nil = tidak ada), urutan sama dengan Prices
	ApprovedMin *decimal.Decimal   // harga termurah saat terakhir disetujui; nil = belum pernah
}

// CheckPrice: daftar kejanggalan harga (kosong = wajar). Harga terlalu mahal tidak dicek (itu salah ketik, bukan penipuan).
func CheckPrice(in PriceCheck, floors []Floor, r Rules) []PriceHit {
	if len(in.Prices) == 0 {
		return nil
	}
	min := in.Prices[0]
	for _, p := range in.Prices[1:] {
		if p.LessThan(min) {
			min = p
		}
	}
	var hits []PriceHit
	if !IsAccessory(in.Name) {
		name := strings.ToLower(in.Name)
		var best *Floor // kata kunci terpanjang menang ("ipad pro" lebih spesifik dari "ipad")
		for i := range floors {
			f := &floors[i]
			if ContainsWord(name, strings.ToLower(f.Keyword)) && (best == nil || len(f.Keyword) > len(best.Keyword)) {
				best = f
			}
		}
		if best != nil && min.LessThan(best.MinPrice) {
			hits = append(hits, PriceHit{Reason: ReasonPriceFloor, Keyword: best.Keyword, Limit: best.MinPrice.StringFixed(2), Price: min.StringFixed(2)})
		}
	}
	limit := decimal.NewFromInt(int64(r.DiscountPct)).Div(decimal.NewFromInt(100))
	for i, c := range in.Compares {
		if c == nil || c.Sign() <= 0 || i >= len(in.Prices) {
			continue
		}
		off := c.Sub(in.Prices[i]).Div(*c)
		if off.GreaterThan(limit) {
			hits = append(hits, PriceHit{Reason: ReasonBigDiscount, Limit: c.StringFixed(2), Price: in.Prices[i].StringFixed(2)})
			break
		}
	}
	if in.ApprovedMin != nil && in.ApprovedMin.Sign() > 0 {
		keep := decimal.NewFromInt(int64(100 - r.DropPct)).Div(decimal.NewFromInt(100))
		if min.LessThan(in.ApprovedMin.Mul(keep)) {
			hits = append(hits, PriceHit{Reason: ReasonPriceDrop, Limit: in.ApprovedMin.StringFixed(2), Price: min.StringFixed(2)})
		}
	}
	return hits
}

// Situation: semua fakta untuk memilih jalur.
type Situation struct {
	PrevModeration string // "" = produk baru / belum pernah dikirim
	MajorEdit      bool   // produk yang sudah tayang diubah besar (nama, kategori, kondisi, foto sampul)
	ContentChanged bool   // ada perubahan selain harga/stok
	Restricted     bool
	ProbationLeft  bool // toko belum punya cukup produk yang disetujui admin
	StrikePoints   int
	KeywordReview  bool
	Reported       bool // sedang ditahan karena laporan pembeli (edit seller tidak melepasnya)
	PriceHits      []PriceHit
}

// Decide: alasan jalur manual (kosong = tayang otomatis).
func Decide(s Situation, r Rules) []string {
	var out []string
	add := func(ok bool, reason string) {
		if ok {
			out = append(out, reason)
		}
	}
	for _, h := range s.PriceHits {
		add(!contains(out, h.Reason), h.Reason)
	}
	if s.PrevModeration == "approved" && !s.MajorEdit {
		// Sudah tayang, perubahan kecil: hanya kata baru & harga yang dicek ulang.
		add(s.ContentChanged && s.KeywordReview, ReasonKeyword)
		return out
	}
	add(!r.AutoPublish, ReasonManualOnly)
	add(s.PrevModeration == "changes_requested", ReasonChanges)
	add(s.PrevModeration == "rejected", ReasonWasRejected)
	add(s.Reported, ReasonReported)
	add(s.ProbationLeft, ReasonProbation)
	add(s.Restricted, ReasonRestricted)
	add(s.StrikePoints > 0, ReasonStrikes)
	add(s.KeywordReview, ReasonKeyword)
	return out
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// ContainsWord: untuk kata latin dicek batas kata ("gun" tidak cocok dengan "burgundy"); aksara lain = substring.
func ContainsWord(text, kw string) bool {
	if kw == "" {
		return false
	}
	latin := true
	for _, r := range kw {
		if r > 127 {
			latin = false
			break
		}
	}
	if !latin {
		return strings.Contains(text, kw)
	}
	re, err := regexp.Compile(`(^|[^a-z0-9])` + regexp.QuoteMeta(kw) + `($|[^a-z0-9])`)
	return err == nil && re.MatchString(text)
}
