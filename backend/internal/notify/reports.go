package notify

// reports: notifikasi strike & penangguhan toko (paket reports).
const (
	StoreStrike     = "store.strike"     // toko mendapat strike
	StoreSuspended  = "store.suspended"  // 3 strike aktif → ditangguhkan
	StoreReinstated = "store.reinstated" // strike dicabut → toko aktif lagi
)

var reportTexts = map[string]map[string]tmpl{
	StoreStrike: {
		"en": {"Your store received a strike ⚠️", "{store}: {note} ({point} active of 3). At 3 active strikes the store is suspended."},
		"km": {"ហាងរបស់អ្នកទទួលបានការព្រមាន ⚠️", "{store}៖ {note} (សកម្ម {point} ក្នុងចំណោម 3)។ ពេលមានការព្រមានសកម្ម 3 ហាងនឹងត្រូវផ្អាក។"},
		"zh": {"您的店铺收到一次警告 ⚠️", "{store}：{note}（有效警告 {point}/3）。累计 3 次有效警告，店铺将被暂停。"},
	},
	StoreSuspended: {
		"en": {"Your store is suspended", "{store} has 3 active strikes and is hidden from buyers. Contact Phsar support to appeal."},
		"km": {"ហាងរបស់អ្នកត្រូវបានផ្អាក", "{store} មានការព្រមានសកម្ម 3 ហើយត្រូវបានលាក់ពីអ្នកទិញ។ សូមទាក់ទងផ្នែកជំនួយ Phsar ដើម្បីប្តឹងឧទ្ធរណ៍។"},
		"zh": {"您的店铺已被暂停", "{store} 已有 3 次有效警告，买家暂时看不到店铺。如需申诉请联系 Phsar 客服。"},
	},
	StoreReinstated: {
		"en": {"Your store is active again ✅", "A strike on {store} was removed: {note}"},
		"km": {"ហាងរបស់អ្នកសកម្មឡើងវិញ ✅", "ការព្រមានលើ {store} ត្រូវបានដកចេញ៖ {note}"},
		"zh": {"您的店铺已恢复 ✅", "{store} 的一次警告已撤销：{note}"},
	},
}

func init() {
	for k, v := range reportTexts {
		texts[k] = v
	}
}
