package notify

// courierreturn: pengantaran gagal → paket dikembalikan ke toko.
const (
	CourierReturnStarted   = "order.courier_return_started" // ke seller: paket sedang dibawa kembali
	CourierReturnArrived   = "order.courier_return_arrived" // ke seller: paket sudah di toko → konfirmasi
	CourierReturnClosed    = "order.courier_return_closed"  // ke seller: Phsar / sistem menandai diterima
	BuyerDeliveryReturning = "buyer.delivery_returning"     // ke pembeli: gagal antar, paket kembali ke toko, uang kembali
	CourierReturnTask      = "courier.return_task"          // ke kurir: antar paket kembali ke toko
)

var courierReturnTexts = map[string]map[string]tmpl{
	CourierReturnStarted: {
		"en": {"Package coming back 📦↩️", "The courier could not deliver order {order}. The package is being brought back to your pickup point. You will be asked to confirm when it arrives."},
		"km": {"កញ្ចប់កំពុងត្រឡប់មកវិញ 📦↩️", "អ្នកដឹកមិនអាចដឹកជូនការបញ្ជាទិញ {order} បានទេ។ កញ្ចប់កំពុងត្រូវបាននាំត្រឡប់មកទីតាំងទទួលរបស់អ្នក។ អ្នកនឹងត្រូវបញ្ជាក់នៅពេលវាមកដល់។"},
		"zh": {"包裹正在退回 📦↩️", "快递员无法送达订单 {order}，包裹正在送回您的取件点。到达后需要您确认。"},
	},
	CourierReturnArrived: {
		"en": {"Confirm the returned package", "The courier brought order {order} back to your pickup point. Check the package and confirm in Seller Center within 24 hours, otherwise Phsar treats it as received."},
		"km": {"បញ្ជាក់កញ្ចប់ដែលបានត្រឡប់មកវិញ", "អ្នកដឹកបាននាំការបញ្ជាទិញ {order} ត្រឡប់មកទីតាំងទទួលរបស់អ្នកវិញ។ សូមពិនិត្យកញ្ចប់ ហើយបញ្ជាក់ក្នុងមជ្ឈមណ្ឌលអ្នកលក់ក្នុង 24 ម៉ោង បើមិនដូច្នោះទេ Phsar នឹងចាត់ទុកថាបានទទួល។"},
		"zh": {"请确认退回的包裹", "快递员已将订单 {order} 送回您的取件点。请检查包裹并在 24 小时内在卖家中心确认，否则 Phsar 将视为已收到。"},
	},
	CourierReturnClosed: {
		"en": {"Returned package recorded", "Order {order} was recorded as returned to your shop and cancelled. The buyer gets a refund. Please check your stock."},
		"km": {"បានកត់ត្រាកញ្ចប់ដែលបានត្រឡប់", "ការបញ្ជាទិញ {order} ត្រូវបានកត់ត្រាថាបានត្រឡប់មកហាងរបស់អ្នក ហើយត្រូវបានបោះបង់។ អ្នកទិញនឹងទទួលប្រាក់វិញ។ សូមពិនិត្យស្តុករបស់អ្នក។"},
		"zh": {"退回包裹已记录", "订单 {order} 已记录为退回店铺并取消，买家将获得退款。请检查您的库存。"},
	},
	BuyerDeliveryReturning: {
		"en": {"We could not deliver your order", "The courier could not deliver order {order}. The package goes back to {store} and we will refund you once the store receives it."},
		"km": {"យើងមិនអាចដឹកជូនការបញ្ជាទិញរបស់អ្នកបានទេ", "អ្នកដឹកមិនអាចដឹកជូនការបញ្ជាទិញ {order} បានទេ។ កញ្ចប់ត្រឡប់ទៅ {store} វិញ ហើយយើងនឹងសងប្រាក់ជូនអ្នកនៅពេលហាងទទួលបាន។"},
		"zh": {"您的订单未能送达", "快递员无法送达订单 {order}。包裹将退回 {store}，店铺收到后我们会为您退款。"},
	},
	CourierReturnTask: {
		"en": {"Return a package ↩️", "Bring order {order} back to {store}. Open the courier app for the pickup point."},
		"km": {"ប្រគល់កញ្ចប់ត្រឡប់ ↩️", "នាំការបញ្ជាទិញ {order} ត្រឡប់ទៅ {store} វិញ។ បើកកម្មវិធីអ្នកដឹកដើម្បីមើលទីតាំង។"},
		"zh": {"退回包裹 ↩️", "请将订单 {order} 送回 {store}。打开快递员页面查看取件点。"},
	},
}

func init() {
	for k, v := range courierReturnTexts {
		texts[k] = v
	}
}
