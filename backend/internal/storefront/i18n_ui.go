package storefront

// Teks untuk tampilan yang mengikuti website Phsar pertama (banner, keunggulan, cari foto).
var uiStrs = map[string]map[string]string{
	"en": {
		"bannerLine1": "Watch videos, shop", "bannerLine2": "on Phsar",
		"trust1": "Verified sellers", "trust2": "Delivery in Sihanoukville", "trust3": "Easy returns",
		"featuredTitle": "Featured products", "featuredSub": "Watch real product videos, then buy",
		"videosTitle": "Short videos for you", "videosSub": "Real videos from local sellers",
		"viewProduct": "View product", "addShort": "Add", "buyShort": "Buy now", "storeLabel": "Store",
		"photoBtn": "Photo", "photoSearch": "Search by photo", "photoTitle": "Search by photo",
		"photoIntro":  "Take or choose a photo of a product. Phsar finds similar products from local sellers.",
		"photoChoose": "Take or choose a photo", "photoAgain": "Try another photo",
		"photoWorking": "Looking for similar products…", "photoResults": "Similar products",
		"photoNone":   "No similar products yet. Try a close-up of one product on a plain background.",
		"photoOff":    "Photo search is not available right now. Please search by name.",
		"photoTooBig": "The photo is too large (max 10 MB).", "photoBad": "This file is not a photo.",
		"photoBusy":    "Too many photo searches. Please wait a minute.",
		"photoPrivacy": "Your photo is only used for this search and is not saved.",
		"chooseLang":   "Choose language", "langLabel": "Language", "prev": "Previous", "next": "Next",
	},
	"km": {
		"bannerLine1": "មើលវីដេអូ ទិញទំនិញ", "bannerLine2": "នៅ Phsar",
		"trust1": "អ្នកលក់បានផ្ទៀងផ្ទាត់", "trust2": "ដឹកជញ្ជូនក្នុងក្រុងព្រះសីហនុ", "trust3": "ងាយស្រួលត្រឡប់",
		"featuredTitle": "ទំនិញពិសេស", "featuredSub": "ផលិតផលឧបត្ថម្ភ មើល និងទិញ",
		"videosTitle": "វីដេអូពេញនិយម", "videosSub": "វីដេអូពិតពីអ្នកលក់ក្នុងស្រុក",
		"viewProduct": "មើលទំនិញ", "addShort": "ដាក់", "buyShort": "ទិញឥឡូវ", "storeLabel": "ហាង",
		"photoBtn": "រូបថត", "photoSearch": "ស្វែងរកតាមរូបថត", "photoTitle": "ស្វែងរកតាមរូបថត",
		"photoIntro":  "ថត ឬជ្រើសរូបថតទំនិញ។ Phsar នឹងរកទំនិញស្រដៀងគ្នាពីអ្នកលក់ក្នុងស្រុក។",
		"photoChoose": "ថត ឬជ្រើសរូបថត", "photoAgain": "សាករូបថតផ្សេង",
		"photoWorking": "កំពុងរកទំនិញស្រដៀងគ្នា…", "photoResults": "ទំនិញស្រដៀងគ្នា",
		"photoNone":   "មិនទាន់មានទំនិញស្រដៀងគ្នាទេ។ សូមថតជិតទំនិញតែមួយលើផ្ទៃសាមញ្ញ។",
		"photoOff":    "ការស្វែងរកតាមរូបថតមិនអាចប្រើបានឥឡូវនេះទេ។ សូមស្វែងរកតាមឈ្មោះ។",
		"photoTooBig": "រូបថតធំពេក (អតិបរមា 10 MB)។", "photoBad": "ឯកសារនេះមិនមែនជារូបថតទេ។",
		"photoBusy":    "ស្វែងរកតាមរូបថតច្រើនពេក។ សូមរង់ចាំមួយនាទី។",
		"photoPrivacy": "រូបថតរបស់អ្នកប្រើសម្រាប់ការស្វែងរកនេះប៉ុណ្ណោះ ហើយមិនត្រូវបានរក្សាទុកទេ។",
		"chooseLang":   "ជ្រើសភាសា", "langLabel": "ភាសា", "prev": "មុន", "next": "បន្ទាប់",
	},
	"zh": {
		"bannerLine1": "看视频，买好物", "bannerLine2": "就在 Phsar",
		"trust1": "认证卖家", "trust2": "西港配送", "trust3": "轻松退换",
		"featuredTitle": "精选商品", "featuredSub": "先看真实视频，再购买",
		"videosTitle": "为你推荐的短视频", "videosSub": "本地卖家的真实视频",
		"viewProduct": "查看商品", "addShort": "加购", "buyShort": "立即购买", "storeLabel": "店铺",
		"photoBtn": "拍照", "photoSearch": "以图搜物", "photoTitle": "以图搜物",
		"photoIntro":  "拍摄或选择一张商品照片，Phsar 会为你找到本地卖家的相似商品。",
		"photoChoose": "拍照或选择照片", "photoAgain": "换一张照片",
		"photoWorking": "正在查找相似商品…", "photoResults": "相似商品",
		"photoNone":   "暂无相似商品。请拍摄单件商品的近照，背景尽量简单。",
		"photoOff":    "以图搜物暂不可用，请按名称搜索。",
		"photoTooBig": "照片太大（最大 10 MB）。", "photoBad": "该文件不是照片。",
		"photoBusy":    "搜索太频繁，请稍等一分钟。",
		"photoPrivacy": "照片仅用于本次搜索，不会被保存。",
		"chooseLang":   "选择语言", "langLabel": "语言", "prev": "上一个", "next": "下一个",
	},
}

func init() {
	for l, m := range uiStrs {
		for k, v := range m {
			strs[l][k] = v
		}
	}
}
