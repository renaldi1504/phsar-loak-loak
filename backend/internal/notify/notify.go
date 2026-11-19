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
	"github.com/jackc/pgx/v5/pgconn"
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
	OrderNew               = "order.new"
	OrderSellerTimeout     = "order.seller_timeout" // ke seller: pesanan batal karena tidak diterima tepat waktu
	OrderPackTimeout       = "order.pack_timeout"   // ke seller: pesanan batal karena tidak dikemas tepat waktu
	// Untuk pembeli (EnqueueUser)
	BuyerOrderAccepted  = "buyer.order_accepted"
	BuyerOrderCancelled = "buyer.order_cancelled" // seller menolak / tidak merespons → uang dikembalikan
	BuyerOrderRefunded  = "buyer.order_refunded"
	BuyerOrderOnTheWay  = "buyer.order_on_the_way"
	BuyerOrderDelivered = "buyer.order_delivered"
	// Untuk kurir (EnqueueUser)
	CourierNewTask = "courier.new_task"
	// Untuk seller
	OrderPickedUp  = "order.picked_up"
	OrderDelivered = "order.delivered"
	PayoutSent     = "payout.sent"
	// Retur & sengketa
	ReturnRequested     = "return.requested" // ke seller: jawab dalam 48 jam
	ReturnEscalated     = "return.escalated" // ke seller: tidak dijawab → Phsar memutuskan
	ReturnDecided       = "return.decided"   // ke seller: keputusan Phsar
	ReturnRefunded      = "return.refunded"  // ke seller: saldo dipotong
	BuyerReturnApproved = "buyer.return_approved"
	BuyerReturnReview   = "buyer.return_review"
	BuyerReturnRejected = "buyer.return_rejected"
	// ads: iklan bersponsor (ke seller)
	AdApproved = "ad.approved"
	AdRejected = "ad.rejected"
	AdPaused   = "ad.paused"
)

var ownerOnly = map[string]bool{StoreApproved: false, StoreRejected: true, KYCRejected: true, PickupApproved: true, PickupRejected: true, StaffJoined: true}

// Params: nilai yang dimasukkan ke teks. Kosong boleh.
type Params struct {
	Store   string `json:"store,omitempty"`
	Product string `json:"product,omitempty"`
	Point   string `json:"point,omitempty"`
	Person  string `json:"person,omitempty"`
	Note    string `json:"note,omitempty"`
	Order   string `json:"order,omitempty"`
}

// Enqueue membuat notifikasi untuk anggota toko (di dalam tx pemanggil).
func Enqueue(ctx context.Context, tx pgx.Tx, storeID, kind string, p Params, link string) error {
	params, _ := json.Marshal(p)
	roleFilter := ""
	if ownerOnly[kind] {
		roleFilter = "AND sm.role = 'owner'"
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO notifications (user_id, store_id, kind, params, link, audience)
		SELECT sm.user_id, sm.store_id, $2, $3, $4, 'seller' FROM store_members sm
		JOIN users u ON u.id = sm.user_id AND u.deleted_at IS NULL
		WHERE sm.store_id = $1 `+roleFilter, storeID, kind, params, nilIfEmpty(link))
	return err
}

// EnqueueUser: notifikasi untuk satu user (pembeli, atau kurir untuk jenis courier.*) — tidak muncul di lonceng
// Seller Center. Sampai lewat website (realtime + Web Push) dan Telegram.
func EnqueueUser(ctx context.Context, tx pgx.Tx, userID, kind string, p Params, link string) error {
	params, _ := json.Marshal(p)
	_, err := tx.Exec(ctx, `INSERT INTO notifications (user_id, kind, params, link, audience) VALUES ($1, $2, $3, $4, $5)`,
		userID, kind, params, nilIfEmpty(link), AudienceOfKind(kind))
	return err
}

// EnqueueTest: notifikasi uji untuk diri sendiri (tombol di pengaturan).
func EnqueueTest(ctx context.Context, q interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}, userID, aud string) error {
	_, err := q.Exec(ctx, `INSERT INTO notifications (user_id, kind, params, audience, telegram_status) VALUES ($1, $2, '{}', $3, 'skipped')`,
		userID, TestKind, aud)
	return err
}

// LinkURL: alamat lengkap sebuah notifikasi (untuk Telegram & push). Link "#/…" milik aplikasi web audience-nya.
func LinkURL(base, aud, link string) string {
	base = strings.TrimRight(base, "/")
	home := map[string]string{AudSeller: "/seller/", AudAdmin: "/admin/", AudCourier: "/courier/"}[aud]
	switch {
	case link == "" && home != "":
		return base + home
	case link == "":
		return base + "/"
	case strings.HasPrefix(link, "#") && home != "":
		return base + home + link
	case strings.HasPrefix(link, "/"):
		return base + link
	}
	return base + "/" + link
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
	OrderNew: {
		"en": {"New paid order 🛍️", "Order {order} ({note}) is paid. Please prepare the items for courier pickup."},
		"km": {"ការបញ្ជាទិញថ្មីបានបង់ប្រាក់ 🛍️", "ការបញ្ជាទិញ {order} ({note}) បានបង់ប្រាក់ហើយ។ សូមរៀបចំទំនិញសម្រាប់អ្នកដឹកមកទទួល។"},
		"zh": {"新订单已付款 🛍️", "订单 {order}（{note}）已付款，请备好商品等待快递员取件。"},
	},
	OrderSellerTimeout: {
		"en": {"Order cancelled — not accepted in time", "Order {order} was not accepted before the time limit and was cancelled. The buyer gets a refund."},
		"km": {"ការបញ្ជាទិញត្រូវបានបោះបង់ — មិនបានទទួលទាន់ពេល", "ការបញ្ជាទិញ {order} មិនត្រូវបានទទួលមុនពេលកំណត់ ហើយត្រូវបានបោះបង់។ អ្នកទិញនឹងទទួលប្រាក់វិញ។"},
		"zh": {"订单已取消——未及时接单", "订单 {order} 未在时限内接单，已自动取消，买家将获得退款。"},
	},
	OrderPackTimeout: {
		"en": {"Order cancelled — not packed in time", "Order {order} was not marked ready for pickup before the time limit and was cancelled. The buyer gets a refund."},
		"km": {"ការបញ្ជាទិញត្រូវបានបោះបង់ — មិនបានវេចខ្ចប់ទាន់ពេល", "ការបញ្ជាទិញ {order} មិនត្រូវបានសម្គាល់ថារួចរាល់ឱ្យមកយកមុនពេលកំណត់ ហើយត្រូវបានបោះបង់។ អ្នកទិញនឹងទទួលប្រាក់វិញ។"},
		"zh": {"订单已取消——未及时打包", "订单 {order} 未在时限内标记为待取件，已自动取消，买家将获得退款。"},
	},
	BuyerOrderAccepted: {
		"en": {"Your order is being prepared", "{store} accepted order {order}. A Phsar courier will pick it up and deliver it to you."},
		"km": {"ការបញ្ជាទិញរបស់អ្នកកំពុងរៀបចំ", "{store} បានទទួលការបញ្ជាទិញ {order}។ អ្នកដឹក Phsar នឹងមកយក ហើយដឹកជូនអ្នក។"},
		"zh": {"您的订单正在备货", "{store} 已接受订单 {order}，Phsar 快递员将取件并送到您手中。"},
	},
	BuyerOrderCancelled: {
		"en": {"Order cancelled — you will get your money back", "Order {order} from {store} was cancelled ({note}). We will refund {product} to you."},
		"km": {"ការបញ្ជាទិញត្រូវបានបោះបង់ — អ្នកនឹងទទួលប្រាក់វិញ", "ការបញ្ជាទិញ {order} ពី {store} ត្រូវបានបោះបង់ ({note})។ យើងនឹងសងប្រាក់ {product} ជូនអ្នក។"},
		"zh": {"订单已取消——款项将退还", "来自 {store} 的订单 {order} 已取消（{note}），我们将退还 {product}。"},
	},
	BuyerOrderOnTheWay: {
		"en": {"Your order is on the way 🛵", "Courier {person} picked up order {order} from {store} and is heading to you."},
		"km": {"ការបញ្ជាទិញកំពុងដឹកមក 🛵", "អ្នកដឹក {person} បានយកការបញ្ជាទិញ {order} ពី {store} ហើយកំពុងមករកអ្នក។"},
		"zh": {"您的订单正在配送 🛵", "快递员 {person} 已从 {store} 取走订单 {order}，正在前往您处。"},
	},
	BuyerOrderDelivered: {
		"en": {"Delivered ✅", "Order {order} from {store} was delivered. Thank you for shopping on Phsar!"},
		"km": {"បានដឹកដល់ ✅", "ការបញ្ជាទិញ {order} ពី {store} បានដឹកដល់ហើយ។ សូមអរគុណដែលបានទិញនៅ Phsar!"},
		"zh": {"已送达 ✅", "来自 {store} 的订单 {order} 已送达，感谢您在 Phsar 购物！"},
	},
	CourierNewTask: {
		"en": {"New pickup 📦", "Order {order} is ready at {store}. Open the courier app to take it."},
		"km": {"មានកញ្ចប់ត្រូវយក 📦", "ការបញ្ជាទិញ {order} រួចរាល់នៅ {store}។ បើកកម្មវិធីអ្នកដឹកដើម្បីទទួលយក។"},
		"zh": {"新取件任务 📦", "订单 {order} 已在 {store} 备好，打开快递员页面接单。"},
	},
	OrderPickedUp: {
		"en": {"Courier picked up the order", "Order {order} was picked up by courier {person}."},
		"km": {"អ្នកដឹកបានយកការបញ្ជាទិញ", "ការបញ្ជាទិញ {order} ត្រូវបានយកដោយអ្នកដឹក {person}។"},
		"zh": {"快递员已取件", "订单 {order} 已由快递员 {person} 取走。"},
	},
	OrderDelivered: {
		"en": {"Order delivered ✅", "Order {order} was delivered to the buyer."},
		"km": {"បានដឹកដល់អ្នកទិញ ✅", "ការបញ្ជាទិញ {order} បានដឹកដល់អ្នកទិញហើយ។"},
		"zh": {"订单已送达 ✅", "订单 {order} 已送达买家。"},
	},
	PayoutSent: {
		"en": {"Payout sent 💸", "Phsar transferred {note} to your bank account (reference {product})."},
		"km": {"បានផ្ទេរប្រាក់ 💸", "Phsar បានផ្ទេរ {note} ទៅគណនីធនាគាររបស់អ្នក (លេខយោង {product})។"},
		"zh": {"已打款 💸", "Phsar 已向您的银行账户转账 {note}（参考号 {product}）。"},
	},
	BuyerOrderRefunded: {
		"en": {"Refund sent", "We refunded {product} for order {order}. It may take a few minutes to appear in your bank app."},
		"km": {"បានសងប្រាក់វិញ", "យើងបានសងប្រាក់ {product} សម្រាប់ការបញ្ជាទិញ {order}។ អាចចំណាយពេលពីរបីនាទីដើម្បីបង្ហាញក្នុងកម្មវិធីធនាគារ។"},
		"zh": {"退款已发出", "订单 {order} 的 {product} 已退款，可能需要几分钟才会显示在您的银行 App 中。"},
	},
}

var returnTexts = map[string]map[string]tmpl{
	ReturnRequested: {
		"en": {"Return request ↩️", "The buyer asked to return order {order} ({note}). Please accept or reject within 48 hours, otherwise Phsar decides."},
		"km": {"សំណើប្ដូរ/សងទំនិញ ↩️", "អ្នកទិញស្នើប្ដូរការបញ្ជាទិញ {order} ({note})។ សូមទទួល ឬបដិសេធក្នុង 48 ម៉ោង បើមិនដូច្នោះទេ Phsar នឹងសម្រេច។"},
		"zh": {"退货申请 ↩️", "买家申请退回订单 {order}（{note}）。请在 48 小时内同意或拒绝，否则由 Phsar 决定。"},
	},
	ReturnEscalated: {
		"en": {"Return sent to Phsar review", "You did not answer the return for order {order} within 48 hours. Phsar will decide."},
		"km": {"ការប្ដូរទំនិញបានបញ្ជូនទៅ Phsar ពិនិត្យ", "អ្នកមិនបានឆ្លើយសំណើប្ដូរការបញ្ជាទិញ {order} ក្នុង 48 ម៉ោង។ Phsar នឹងសម្រេច។"},
		"zh": {"退货已转交 Phsar 审核", "您未在 48 小时内回复订单 {order} 的退货申请，将由 Phsar 决定。"},
	},
	ReturnDecided: {
		"en": {"Phsar decided a return", "Order {order}: {note}"},
		"km": {"Phsar បានសម្រេចលើការប្ដូរទំនិញ", "ការបញ្ជាទិញ {order}៖ {note}"},
		"zh": {"Phsar 已对退货作出决定", "订单 {order}：{note}"},
	},
	ReturnRefunded: {
		"en": {"Return refunded", "The buyer was refunded for order {order}. {product} was deducted from your balance."},
		"km": {"បានសងប្រាក់សម្រាប់ការប្ដូរទំនិញ", "អ្នកទិញបានទទួលប្រាក់វិញសម្រាប់ការបញ្ជាទិញ {order}។ {product} ត្រូវបានកាត់ពីសមតុល្យរបស់អ្នក។"},
		"zh": {"退货已退款", "订单 {order} 已向买家退款，已从您的余额中扣除 {product}。"},
	},
	BuyerReturnApproved: {
		"en": {"Return approved ✅", "Your return for order {order} was approved. We will refund {product}. {note}"},
		"km": {"ការប្ដូរទំនិញត្រូវបានយល់ព្រម ✅", "ការប្ដូរការបញ្ជាទិញ {order} ត្រូវបានយល់ព្រម។ យើងនឹងសងប្រាក់ {product}។ {note}"},
		"zh": {"退货已同意 ✅", "订单 {order} 的退货已同意，我们将退款 {product}。{note}"},
	},
	BuyerReturnReview: {
		"en": {"Phsar is reviewing your return", "The seller did not accept the return for order {order}. Phsar will review it and decide, usually within 2 working days."},
		"km": {"Phsar កំពុងពិនិត្យការប្ដូរទំនិញរបស់អ្នក", "អ្នកលក់មិនបានទទួលការប្ដូរការបញ្ជាទិញ {order}។ Phsar នឹងពិនិត្យ និងសម្រេច ជាធម្មតាក្នុង 2 ថ្ងៃធ្វើការ។"},
		"zh": {"Phsar 正在审核您的退货", "卖家未同意订单 {order} 的退货，Phsar 将审核并决定，通常在 2 个工作日内。"},
	},
	BuyerReturnRejected: {
		"en": {"Return not approved", "Order {order}: {note}"},
		"km": {"ការប្ដូរទំនិញមិនត្រូវបានយល់ព្រម", "ការបញ្ជាទិញ {order}៖ {note}"},
		"zh": {"退货未获同意", "订单 {order}：{note}"},
	},
}

// ads: teks notifikasi iklan bersponsor.
var adTexts = map[string]map[string]tmpl{
	AdApproved: {
		"en": {"Ad approved 📣", "The ad for {product} is approved and will show in search and category pages."},
		"km": {"ការផ្សាយពាណិជ្ជកម្មត្រូវបានអនុម័ត 📣", "ការផ្សាយសម្រាប់ {product} ត្រូវបានអនុម័ត ហើយនឹងបង្ហាញក្នុងការស្វែងរក និងទំព័រប្រភេទ។"},
		"zh": {"广告已通过 📣", "{product} 的广告已通过，将显示在搜索和分类页面。"},
	},
	AdRejected: {
		"en": {"Ad not approved", "{product}: {note}"},
		"km": {"ការផ្សាយពាណិជ្ជកម្មមិនត្រូវបានអនុម័ត", "{product}៖ {note}"},
		"zh": {"广告未通过", "{product}：{note}"},
	},
	AdPaused: {
		"en": {"Ad paused by Phsar", "{product}: {note}"},
		"km": {"Phsar បានផ្អាកការផ្សាយពាណិជ្ជកម្ម", "{product}៖ {note}"},
		"zh": {"广告已被 Phsar 暂停", "{product}：{note}"},
	},
}

func init() {
	for k, v := range adTexts {
		texts[k] = v
	}
	for k, v := range returnTexts {
		texts[k] = v
	}
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
	r := strings.NewReplacer("{store}", p.Store, "{product}", p.Product, "{point}", p.Point, "{person}", p.Person, "{note}", p.Note, "{order}", p.Order)
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
