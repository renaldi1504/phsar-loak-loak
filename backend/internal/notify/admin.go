package notify

// Notifikasi untuk admin, dibuat oleh trigger database (migrasi 00050, fungsi notify_admins) ke setiap admin
// aktif yang perannya punya izin antrean itu. Params: store, product, order, note, point.
// Juga: notifikasi uji ("Kirim notifikasi uji" di pengaturan) untuk semua audience.

// TestKind: notifikasi uji dari halaman pengaturan — selalu berbunyi dan dikirim push.
const TestKind = "system.test"

func init() {
	for k, v := range map[string]map[string]tmpl{
		TestKind: {
			"en": {"Test notification ✅", "Notifications and sound are working on this device."},
			"km": {"ការជូនដំណឹងសាកល្បង ✅", "ការជូនដំណឹង និងសំឡេងដំណើរការលើឧបករណ៍នេះ។"},
			"zh": {"测试通知 ✅", "此设备的通知和提示音工作正常。"},
		},
		"admin.store_review": {
			"en": {"New store to verify", "{store} is waiting for verification."},
			"km": {"ហាងថ្មីត្រូវផ្ទៀងផ្ទាត់", "{store} កំពុងរង់ចាំការផ្ទៀងផ្ទាត់។"},
			"zh": {"新店铺待审核", "{store} 正在等待审核。"},
		},
		"admin.pickup_review": {
			"en": {"Pickup point to verify", "{point} ({store})"},
			"km": {"ទីតាំងទទួលកញ្ចប់ត្រូវផ្ទៀងផ្ទាត់", "{point} ({store})"},
			"zh": {"取件点待核实", "{point}（{store}）"},
		},
		"admin.product_review": {
			"en": {"Product waiting for review", "{product} · {store}"},
			"km": {"ផលិតផលរង់ចាំការពិនិត្យ", "{product} · {store}"},
			"zh": {"商品待审核", "{product} · {store}"},
		},
		"admin.refund_pending": {
			"en": {"Refund to process", "Order {order} was cancelled — refund the buyer."},
			"km": {"ការសងប្រាក់ត្រូវដំណើរការ", "ការបញ្ជាទិញ {order} ត្រូវបានលុបចោល — សងប្រាក់ឲ្យអ្នកទិញ។"},
			"zh": {"待处理退款", "订单 {order} 已取消——请给买家退款。"},
		},
		"admin.return_refund": {
			"en": {"Return approved — refund due", "Order {order}"},
			"km": {"ការប្រគល់វិញត្រូវបានអនុម័ត — ត្រូវសងប្រាក់", "ការបញ្ជាទិញ {order}"},
			"zh": {"退货已批准——待退款", "订单 {order}"},
		},
		"admin.report_new": {
			"en": {"New buyer report", "{store} · {note}"},
			"km": {"របាយការណ៍ថ្មីពីអ្នកទិញ", "{store} · {note}"},
			"zh": {"新的买家举报", "{store} · {note}"},
		},
		"admin.ticket_new": {
			"en": {"New support ticket", "{order}: {note}"},
			"km": {"សំបុត្រជំនួយថ្មី", "{order}៖ {note}"},
			"zh": {"新的客服工单", "{order}：{note}"},
		},
		"admin.ticket_reply": {
			"en": {"Customer replied to a ticket", "{order}: {note}"},
			"km": {"អតិថិជនបានឆ្លើយសំបុត្រជំនួយ", "{order}៖ {note}"},
			"zh": {"客户回复了工单", "{order}：{note}"},
		},
		"admin.return_dispute": {
			"en": {"Return dispute to decide", "Order {order} — the seller rejected or didn't answer."},
			"km": {"វិវាទការប្រគល់វិញត្រូវសម្រេច", "ការបញ្ជាទិញ {order} — អ្នកលក់បដិសេធ ឬមិនបានឆ្លើយ។"},
			"zh": {"待裁决的退货纠纷", "订单 {order}——卖家拒绝或未回复。"},
		},
		"admin.payout_pending": {
			"en": {"Payout to approve", "{store} · {product}"},
			"km": {"ការទូទាត់ត្រូវអនុម័ត", "{store} · {product}"},
			"zh": {"待批准的结算", "{store} · {product}"},
		},
		"admin.delivery_new": {
			"en": {"Delivery needs a courier", "Order {order}"},
			"km": {"ការដឹកជញ្ជូនត្រូវការអ្នកដឹក", "ការបញ្ជាទិញ {order}"},
			"zh": {"配送待分配快递员", "订单 {order}"},
		},
		"admin.delivery_failed": {
			"en": {"Delivery failed", "Order {order}: {note}"},
			"km": {"ការដឹកជញ្ជូនបរាជ័យ", "ការបញ្ជាទិញ {order}៖ {note}"},
			"zh": {"配送失败", "订单 {order}：{note}"},
		},
	} {
		texts[k] = v
	}
}
