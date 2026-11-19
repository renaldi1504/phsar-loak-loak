package notify

// extras: notifikasi strike otomatis (pesanan batal karena telat) — teks 3 bahasa, nomor pesanan di {order}.
const (
	StoreStrikeLateAccept = "store.strike_late_accept" // strike otomatis: pesanan tidak diterima tepat waktu
	StoreStrikeLatePack   = "store.strike_late_pack"   // strike otomatis: pesanan tidak dikemas tepat waktu
)

var extrasTexts = map[string]map[string]tmpl{
	StoreStrikeLateAccept: {
		"en": {"Strike: order not accepted in time ⚠️", "{store}: order {order} was cancelled because it was not accepted in time ({point} active of 3). At 3 active strikes the store is suspended."},
		"km": {"ការព្រមាន៖ មិនបានទទួលការបញ្ជាទិញទាន់ពេល ⚠️", "{store}៖ ការបញ្ជាទិញ {order} ត្រូវបានបោះបង់ ព្រោះមិនបានទទួលទាន់ពេល (សកម្ម {point} ក្នុងចំណោម 3)។ ពេលមានការព្រមានសកម្ម 3 ហាងនឹងត្រូវផ្អាក។"},
		"zh": {"警告：未及时接单 ⚠️", "{store}：订单 {order} 因未及时接单已取消（有效警告 {point}/3）。累计 3 次有效警告，店铺将被暂停。"},
	},
	StoreStrikeLatePack: {
		"en": {"Strike: order not packed in time ⚠️", "{store}: order {order} was cancelled because it was not packed in time ({point} active of 3). At 3 active strikes the store is suspended."},
		"km": {"ការព្រមាន៖ មិនបានវេចខ្ចប់ការបញ្ជាទិញទាន់ពេល ⚠️", "{store}៖ ការបញ្ជាទិញ {order} ត្រូវបានបោះបង់ ព្រោះមិនបានវេចខ្ចប់ទាន់ពេល (សកម្ម {point} ក្នុងចំណោម 3)។ ពេលមានការព្រមានសកម្ម 3 ហាងនឹងត្រូវផ្អាក។"},
		"zh": {"警告：未及时打包 ⚠️", "{store}：订单 {order} 因未及时打包已取消（有效警告 {point}/3）。累计 3 次有效警告，店铺将被暂停。"},
	},
}

func init() {
	for k, v := range extrasTexts {
		texts[k] = v
	}
}
