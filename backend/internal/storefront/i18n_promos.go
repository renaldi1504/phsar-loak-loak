package storefront

// promos: teks flash sale & voucher (digabung ke strs dan shopStrs).
var promoStrs = map[string]map[string]string{
	"en": {"flashSale": "Flash sale", "flashEndsIn": "ends in {t}", "flashStartsIn": "starts in {t}", "flashLive": "Live now", "flashNext": "Upcoming", "flashSeeAll": "See all", "flashSold": "{n} sold", "flashSoldOut": "Sold out", "flashAlmost": "Almost gone", "flashNone": "No flash sale right now. Check back soon!", "flashMeta": "Limited-time deals from sellers in Sihanoukville. Quantities are limited."},
	"km": {"flashSale": "Flash sale", "flashEndsIn": "បញ្ចប់ក្នុង {t}", "flashStartsIn": "ចាប់ផ្តើមក្នុង {t}", "flashLive": "កំពុងដំណើរការ", "flashNext": "នឹងមកដល់", "flashSeeAll": "មើលទាំងអស់", "flashSold": "លក់បាន {n}", "flashSoldOut": "អស់ហើយ", "flashAlmost": "ជិតអស់ហើយ", "flashNone": "មិនមាន Flash sale ពេលនេះទេ។ សូមត្រឡប់មកវិញឆាប់ៗ!", "flashMeta": "ការផ្តល់ជូនរយៈពេលកំណត់ពីអ្នកលក់ក្នុងក្រុងព្រះសីហនុ។ ចំនួនមានកំណត់។"},
	"zh": {"flashSale": "限时抢购", "flashEndsIn": "{t} 后结束", "flashStartsIn": "{t} 后开始", "flashLive": "进行中", "flashNext": "即将开始", "flashSeeAll": "查看全部", "flashSold": "已售 {n}", "flashSoldOut": "已抢光", "flashAlmost": "即将售罄", "flashNone": "目前没有限时抢购，请稍后再来！", "flashMeta": "西港本地卖家的限时优惠，数量有限。"},
}

var promoShopStrs = map[string]map[string]string{
	"en": {"voucher": "Voucher code", "voucherPh": "Enter code", "voucherApply": "Apply", "voucherRemove": "Remove", "voucherApplied": "Voucher {c} applied",
		"discount": "Voucher discount", "flashBadge": "Flash sale",
		"err_voucher_flash": "Flash sale items in this session cannot be combined with vouchers.", "err_voucher_not_found": "This voucher code is not valid.", "err_voucher_expired": "This voucher has expired.",
		"err_voucher_not_started": "This voucher is not active yet.", "err_voucher_min_spend": "Spend at least ${a} on items to use this voucher.",
		"err_voucher_used_up": "This voucher has been fully used.", "err_voucher_user_limit": "You have already used this voucher.",
		"err_voucher_first_order": "This voucher is for your first order only.", "err_voucher_login_required": "Sign in with Telegram, phone or Facebook to use this voucher."},
	"km": {"voucher": "លេខកូដប័ណ្ណ", "voucherPh": "បញ្ចូលលេខកូដ", "voucherApply": "ប្រើ", "voucherRemove": "ដកចេញ", "voucherApplied": "បានប្រើប័ណ្ណ {c}",
		"discount": "បញ្ចុះតម្លៃប័ណ្ណ", "flashBadge": "Flash sale",
		"err_voucher_flash": "ទំនិញ Flash sale ក្នុងវគ្គនេះមិនអាចប្រើជាមួយប័ណ្ណបានទេ។", "err_voucher_not_found": "លេខកូដប័ណ្ណនេះមិនត្រឹមត្រូវទេ។", "err_voucher_expired": "ប័ណ្ណនេះផុតកំណត់ហើយ។",
		"err_voucher_not_started": "ប័ណ្ណនេះមិនទាន់អាចប្រើបានទេ។", "err_voucher_min_spend": "ទិញយ៉ាងហោចណាស់ ${a} ដើម្បីប្រើប័ណ្ណនេះ។",
		"err_voucher_used_up": "ប័ណ្ណនេះត្រូវបានប្រើអស់ហើយ។", "err_voucher_user_limit": "អ្នកបានប្រើប័ណ្ណនេះរួចហើយ។",
		"err_voucher_first_order": "ប័ណ្ណនេះសម្រាប់ការបញ្ជាទិញលើកដំបូងប៉ុណ្ណោះ។", "err_voucher_login_required": "សូមចូលដោយ Telegram លេខទូរស័ព្ទ ឬ Facebook ដើម្បីប្រើប័ណ្ណនេះ។"},
	"zh": {"voucher": "优惠码", "voucherPh": "输入优惠码", "voucherApply": "使用", "voucherRemove": "移除", "voucherApplied": "已使用优惠券 {c}",
		"discount": "优惠券折扣", "flashBadge": "限时抢购",
		"err_voucher_flash": "本场限时抢购商品不能与优惠券同时使用。", "err_voucher_not_found": "此优惠码无效。", "err_voucher_expired": "此优惠券已过期。",
		"err_voucher_not_started": "此优惠券尚未开始。", "err_voucher_min_spend": "商品满 ${a} 才能使用此优惠券。",
		"err_voucher_used_up": "此优惠券已被领完。", "err_voucher_user_limit": "您已使用过此优惠券。",
		"err_voucher_first_order": "此优惠券仅限首单使用。", "err_voucher_login_required": "请使用 Telegram、手机号或 Facebook 登录后使用此优惠券。"},
}

func init() {
	for l, m := range promoStrs {
		for k, v := range m {
			strs[l][k] = v
		}
	}
	for l, m := range promoShopStrs {
		for k, v := range m {
			shopStrs[l][k] = v
		}
	}
}
