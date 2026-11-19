package storefront

// courierreturn: teks riwayat pesanan pembeli saat pengantaran gagal & paket dikembalikan ke toko (digabung ke shopStrs).
var courierReturnStrs = map[string]map[string]string{
	"en": {
		"ev_delivery_failed": "Delivery attempt failed", "ev_delivery_retry": "Courier is trying again",
		"ev_return_to_store": "Going back to the store — you will get a refund", "ev_returned_to_store": "Back at the store",
		"ev_return_received": "Store received the package",
	},
	"km": {
		"ev_delivery_failed": "ការដឹកជញ្ជូនមិនបានសម្រេច", "ev_delivery_retry": "អ្នកដឹកកំពុងព្យាយាមម្តងទៀត",
		"ev_return_to_store": "កំពុងត្រឡប់ទៅហាងវិញ — អ្នកនឹងទទួលប្រាក់វិញ", "ev_returned_to_store": "បានត្រឡប់ដល់ហាងវិញ",
		"ev_return_received": "ហាងបានទទួលកញ្ចប់វិញ",
	},
	"zh": {
		"ev_delivery_failed": "配送未成功", "ev_delivery_retry": "快递员正在再次配送",
		"ev_return_to_store": "正在退回店铺——您将获得退款", "ev_returned_to_store": "已退回店铺",
		"ev_return_received": "店铺已收到包裹",
	},
}

func init() {
	for l, m := range courierReturnStrs {
		for k, v := range m {
			shopStrs[l][k] = v
		}
	}
}
