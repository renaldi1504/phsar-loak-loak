package storefront

// Teks halaman error (errors.go).
var errorStrs = map[string]map[string]string{
	"en": {
		"errCode": "Error {c}", "goBack": "Go back", "tryAgain": "Try again", "errHelp": "Still stuck? Contact Phsar",
		"errSearchPh": "Search for a product instead", "errBrowse": "Or browse a category",
		"err404Title": "Page not found", "err404Text": "The link may be wrong, or the page has moved.",
		"err404ProductText":  "This product may have been sold or removed by the seller.",
		"err404StoreText":    "This store is no longer open on Phsar, or the link is wrong.",
		"err404CategoryText": "This category no longer exists.",
		"err404OrderText":    "We couldn't find this order. Open it from Account → Orders on the device you ordered with.",
		"err410Title":        "No longer available", "err410Text": "This page was removed and won't come back.",
		"err400Title": "Something's wrong with this link", "err400Text": "Check the address and try again.",
		"err403Title": "You can't open this page", "err403Text": "This page is only for the account that owns it.",
		"err429Title": "Too many requests", "err429Text": "Please wait a minute, then try again.",
		"err500Title": "Something went wrong", "err500Text": "We couldn't load this page. It's not your fault. Please try again in a moment.",
		"err503Title": "Phsar is being updated", "err503Text": "We'll be back in a few minutes. Your cart and orders are safe.",
	},
	"km": {
		"errCode": "កំហុស {c}", "goBack": "ត្រឡប់ក្រោយ", "tryAgain": "សាកល្បងម្តងទៀត", "errHelp": "នៅតែមានបញ្ហា? ទាក់ទង Phsar",
		"errSearchPh": "ស្វែងរកផលិតផលជំនួសវិញ", "errBrowse": "ឬមើលតាមប្រភេទ",
		"err404Title": "រកទំព័រមិនឃើញ", "err404Text": "តំណអាចខុស ឬទំព័រនេះត្រូវបានផ្លាស់ទី។",
		"err404ProductText":  "ផលិតផលនេះប្រហែលជាលក់អស់ ឬត្រូវបានអ្នកលក់ដកចេញ។",
		"err404StoreText":    "ហាងនេះលែងបើកលើ Phsar ទៀតហើយ ឬតំណខុស។",
		"err404CategoryText": "ប្រភេទនេះលែងមានទៀតហើយ។",
		"err404OrderText":    "រកការបញ្ជាទិញនេះមិនឃើញ។ សូមបើកពី គណនី → ការបញ្ជាទិញ នៅលើឧបករណ៍ដែលអ្នកបានបញ្ជាទិញ។",
		"err410Title":        "លែងមានទៀតហើយ", "err410Text": "ទំព័រនេះត្រូវបានដកចេញ ហើយនឹងមិនត្រឡប់មកវិញទេ។",
		"err400Title": "តំណនេះមានបញ្ហា", "err400Text": "សូមពិនិត្យអាសយដ្ឋាន ហើយសាកល្បងម្តងទៀត។",
		"err403Title": "អ្នកមិនអាចបើកទំព័រនេះបានទេ", "err403Text": "ទំព័រនេះសម្រាប់តែគណនីដែលជាម្ចាស់ប៉ុណ្ណោះ។",
		"err429Title": "សំណើច្រើនពេក", "err429Text": "សូមរង់ចាំមួយនាទី រួចសាកល្បងម្តងទៀត។",
		"err500Title": "មានបញ្ហាបន្តិចបន្តួច", "err500Text": "យើងមិនអាចបើកទំព័រនេះបានទេ។ មិនមែនជាកំហុសរបស់អ្នកទេ។ សូមសាកល្បងម្តងទៀតបន្តិចទៀត។",
		"err503Title": "Phsar កំពុងធ្វើបច្ចុប្បន្នភាព", "err503Text": "យើងនឹងត្រឡប់មកវិញក្នុងរយៈពេលប៉ុន្មាននាទីទៀត។ កន្ត្រក និងការបញ្ជាទិញរបស់អ្នកនៅដដែល។",
	},
	"zh": {
		"errCode": "错误 {c}", "goBack": "返回上一页", "tryAgain": "重试", "errHelp": "仍有问题？联系 Phsar",
		"errSearchPh": "改为搜索商品", "errBrowse": "或按分类浏览",
		"err404Title": "页面不存在", "err404Text": "链接可能有误，或页面已移动。",
		"err404ProductText":  "该商品可能已售出或已被卖家下架。",
		"err404StoreText":    "该店铺已不在 Phsar 营业，或链接有误。",
		"err404CategoryText": "该分类已不存在。",
		"err404OrderText":    "找不到该订单。请在下单的设备上，从「我的 → 订单」打开。",
		"err410Title":        "已不再提供", "err410Text": "该页面已被移除，不会再恢复。",
		"err400Title": "链接有问题", "err400Text": "请检查网址后重试。",
		"err403Title": "无法打开此页面", "err403Text": "此页面仅限所属账户查看。",
		"err429Title": "请求过于频繁", "err429Text": "请稍等一分钟后再试。",
		"err500Title": "出了点问题", "err500Text": "页面暂时无法加载，这不是您的问题。请稍后再试。",
		"err503Title": "Phsar 正在更新", "err503Text": "几分钟后恢复。您的购物车和订单都不会受影响。",
	},
}

func init() {
	for l, m := range errorStrs {
		for k, v := range m {
			strs[l][k] = v
		}
	}
}
