package notify

// flash: notifikasi sesi flash sale untuk seller.
const (
	FlashInvited  = "flash.invited"  // Phsar memilih produk toko untuk sesi (seller harus menerima)
	FlashApproved = "flash.approved" // pendaftaran seller disetujui
	FlashRejected = "flash.rejected" // pendaftaran ditolak / barang dikeluarkan
)

var flashTexts = map[string]map[string]tmpl{
	FlashInvited: {
		"en": {"Invited to a flash sale ⚡", "Phsar chose {product} of your products for “{note}”. Open Promotions to accept or decline."},
		"km": {"ត្រូវបានអញ្ជើញចូល Flash sale ⚡", "Phsar បានជ្រើសផលិតផលរបស់អ្នក {product} សម្រាប់ «{note}»។ បើក ការផ្សព្វផ្សាយ ដើម្បីទទួល ឬបដិសេធ។"},
		"zh": {"受邀参加限时抢购 ⚡", "Phsar 为“{note}”选择了您的 {product} 个商品。请打开“促销”接受或拒绝。"},
	},
	FlashApproved: {
		"en": {"Flash sale product approved ✅", "{product} is in the flash sale: {note}"},
		"km": {"ផលិតផល Flash sale ត្រូវបានអនុម័ត ✅", "{product} ចូលរួម Flash sale៖ {note}"},
		"zh": {"限时抢购商品已通过 ✅", "{product} 已加入限时抢购：{note}"},
	},
	FlashRejected: {
		"en": {"Flash sale product not accepted", "{product}: {note}"},
		"km": {"ផលិតផល Flash sale មិនត្រូវបានទទួល", "{product}៖ {note}"},
		"zh": {"限时抢购商品未通过", "{product}：{note}"},
	},
}

func init() {
	for k, v := range flashTexts {
		texts[k] = v
	}
}
