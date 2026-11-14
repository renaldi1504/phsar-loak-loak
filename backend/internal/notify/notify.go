// Package notify: notifikasi untuk seller — tersimpan (lonceng di Seller Center / app) dan dikirim lewat
// Telegram ke user yang login dengan Telegram.
//
// Pola outbox: Enqueue dipanggil di dalam transaksi yang sama dengan perubahan datanya (mis. keputusan admin),
// jadi notifikasi tidak pernah hilang atau terkirim untuk perubahan yang batal. Worker mengirimnya kemudian.
package notify

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/jackc/pgx/v5"
)

// Jenis notifikasi. Penerima: semua anggota toko, kecuali yang ditandai ownerOnly.
const (
	StoreApproved          = "store.approved"
	StoreRejected          = "store.rejected"
	KYCRejected            = "kyc.rejected"
	PickupApproved         = "pickup.approved"
	PickupRejected         = "pickup.rejected"
	ProductApproved        = "product.approved"
	ProductChangesRequired = "product.changes_requested"
	ProductRejected        = "product.rejected"
	StaffJoined            = "staff.joined"
)

var ownerOnly = map[string]bool{StoreApproved: false, StoreRejected: true, KYCRejected: true, PickupApproved: true, PickupRejected: true, StaffJoined: true}

// Params: nilai yang dimasukkan ke teks. Kosong boleh.
type Params struct {
	Store   string `json:"store,omitempty"`
	Product string `json:"product,omitempty"`
	Point   string `json:"point,omitempty"`
	Person  string `json:"person,omitempty"`
	Note    string `json:"note,omitempty"`
}

// Enqueue membuat notifikasi untuk anggota toko (di dalam tx pemanggil).
func Enqueue(ctx context.Context, tx pgx.Tx, storeID, kind string, p Params, link string) error {
	params, _ := json.Marshal(p)
	roleFilter := ""
	if ownerOnly[kind] {
		roleFilter = "AND sm.role = 'owner'"
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO notifications (user_id, store_id, kind, params, link)
		SELECT sm.user_id, sm.store_id, $2, $3, $4 FROM store_members sm
		JOIN users u ON u.id = sm.user_id AND u.deleted_at IS NULL
		WHERE sm.store_id = $1 `+roleFilter, storeID, kind, params, nilIfEmpty(link))
	return err
}

// ---- Teks (en / km / zh) ----

type tmpl struct{ title, body string }

var texts = map[string]map[string]tmpl{
	StoreApproved: {
		"en": {"Your store is approved 🎉", "{store} is now active. You can send products for review."},
		"km": {"ហាងរបស់អ្នកត្រូវបានអនុម័ត 🎉", "{store} សកម្មហើយ។ អ្នកអាចផ្ញើផលិតផលឱ្យពិនិត្យបាន។"},
		"zh": {"店铺审核通过 🎉", "{store} 已上线，现在可以提交商品审核。"},
	},
	StoreRejected: {
		"en": {"Your store needs changes", "{store}: {note}"},
		"km": {"ហាងរបស់អ្នកត្រូវកែប្រែ", "{store}៖ {note}"},
		"zh": {"店铺需要修改", "{store}：{note}"},
	},
	KYCRejected: {
		"en": {"Identity check needs changes", "{store}: {note}"},
		"km": {"ការផ្ទៀងផ្ទាត់អត្តសញ្ញាណត្រូវកែប្រែ", "{store}៖ {note}"},
		"zh": {"身份验证需要修改", "{store}：{note}"},
	},
	PickupApproved: {
		"en": {"Pickup point verified", "{point} ({store}) is verified. Couriers will pick up here."},
		"km": {"ទីតាំងទទួលកញ្ចប់ត្រូវបានផ្ទៀងផ្ទាត់", "{point} ({store}) បានផ្ទៀងផ្ទាត់។ អ្នកដឹកនឹងមកទទួលនៅទីនេះ។"},
		"zh": {"取件点已核实", "{point}（{store}）已核实，快递员会在这里取件。"},
	},
	PickupRejected: {
		"en": {"Pickup point needs changes", "{point} ({store}): {note}"},
		"km": {"ទីតាំងទទួលកញ្ចប់ត្រូវកែប្រែ", "{point} ({store})៖ {note}"},
		"zh": {"取件点需要修改", "{point}（{store}）：{note}"},
	},
	ProductApproved: {
		"en": {"Product is live ✅", "{product} is now visible to buyers."},
		"km": {"ផលិតផលកំពុងលក់ ✅", "{product} អ្នកទិញអាចមើលឃើញហើយ។"},
		"zh": {"商品已上架 ✅", "{product} 现在买家可以看到了。"},
	},
	ProductChangesRequired: {
		"en": {"Product needs changes", "{product}: {note}"},
		"km": {"ផលិតផលត្រូវកែប្រែ", "{product}៖ {note}"},
		"zh": {"商品需要修改", "{product}：{note}"},
	},
	ProductRejected: {
		"en": {"Product can’t be sold", "{product}: {note}"},
		"km": {"ផលិតផលមិនអាចលក់បាន", "{product}៖ {note}"},
		"zh": {"商品不能销售", "{product}：{note}"},
	},
	StaffJoined: {
		"en": {"New staff member", "{person} joined {store} as staff."},
		"km": {"បុគ្គលិកថ្មី", "{person} បានចូលរួមជាបុគ្គលិករបស់ {store}។"},
		"zh": {"新店员", "{person} 已加入 {store} 成为店员。"},
	},
}

var openLabel = map[string]string{"en": "Open Seller Center", "km": "បើកមជ្ឈមណ្ឌលអ្នកលក់", "zh": "打开卖家中心"}

// Render menghasilkan judul & isi dalam bahasa lang (cadangan: en).
func Render(kind, lang string, p Params) (title, body string) {
	byLang, ok := texts[kind]
	if !ok {
		return kind, ""
	}
	t, ok := byLang[lang]
	if !ok {
		t = byLang["en"]
	}
	r := strings.NewReplacer("{store}", p.Store, "{product}", p.Product, "{point}", p.Point, "{person}", p.Person, "{note}", p.Note)
	body = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(r.Replace(t.body)), ":"))
	return r.Replace(t.title), body
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func truncate(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n]) + "…"
	}
	return s
}
