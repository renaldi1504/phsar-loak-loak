package storefront

import (
	"fmt"
	"strings"
)

var langs = []string{"km", "en", "zh"}

var langNames = map[string]string{"en": "English", "km": "ខ្មែរ", "zh": "中文"}
var langShort = map[string]string{"en": "EN", "km": "ខ្មែរ", "zh": "中文"}
var htmlLang = map[string]string{"en": "en", "km": "km", "zh": "zh-Hans"}
var ogLocale = map[string]string{"en": "en_US", "km": "km_KH", "zh": "zh_CN"}

var strs = map[string]map[string]string{
	"en": {
		"weAccept": "We accept", "soonTitle": "Phsar is opening soon.", "soonText": "Sellers in Sihanoukville can register their shop now — shopping opens soon.", "soonSell": "Open your shop", "tagline": "Marketplace Cambodia", "metaDesc": "Watch real product videos from local sellers in Sihanoukville and buy with free delivery.",
		"home": "Home", "search": "Search", "videos": "Videos", "cart": "Cart", "me": "Account", "language": "Language", "mainMenu": "Main menu",
		"searchPh": "Search products, stores…", "searchBtn": "Search",
		"hero1": "See it on video.", "hero2": "Buy with confidence.", "heroSub": "Watch • Discover • Shop",
		"trust1": "Free delivery in Sihanoukville", "trust2": "Real product videos", "trust3": "Checked by Phsar",
		"categories": "Categories", "videosTitle": "Watch & shop", "videosSub": "Real videos from local sellers",
		"newTitle": "New arrivals", "newSub": "Just added by sellers", "seeAll": "See all", "all": "All",
		"from": "from", "soldOut": "Sold out", "new": "New", "like_new": "Like new", "used_good": "Used – good", "used_fair": "Used – fair", "refurbished": "Refurbished",
		"back": "Back", "productVideo": "Product video", "photo": "Photo", "visitStore": "Visit store", "products": "{n} products", "since": "On Phsar since {d}",
		"returns": "{n}-day returns", "noReturns": "No returns", "warranty_store": "{n}-month store warranty", "warranty_official": "{n}-month official warranty",
		"delivery": "Delivery", "deliveryVal": "Free", "deliveryNote": "Delivered by a Phsar courier within Sihanoukville.",
		"description": "Description", "translatedFrom": "Translated automatically from {lang}.", "showOriginal": "Show original", "showTranslation": "Show translation",
		"conditionTitle": "Condition", "usedFor": "Used for about {n} months", "conditionPhotos": "Condition photos",
		"choose": "Choose {opt}", "left": "Only {n} left", "inStock": "In stock", "outOfStock": "Out of stock in this option",
		"addToCart": "Add to cart", "buyNow": "Buy now", "comingSoon": "Ordering opens soon. Save this page and come back!",
		"related": "More like this", "storeProducts": "Products from this store", "aboutStore": "About the store",
		"sortNew": "Newest", "sortPriceAsc": "Price: low to high", "sortPriceDesc": "Price: high to low", "sort": "Sort",
		"results": "{n} results for “{q}”", "noResults": "Nothing found for “{q}”. Try another word or browse a category.", "empty": "No products here yet.",
		"searchTitle": "Search", "loadMore": "Load more", "allProducts": "All products",
		"notFound": "Page not found", "notFoundText": "This product may have been sold or removed.", "goHome": "Go to home",
		"footerAbout": "Phsar is a video marketplace for Sihanoukville, Cambodia.", "sell": "Sell on Phsar", "sellTab": "Sell", "copyright": "© {y} Phsar",
		"legal_buyer_terms": "Terms of use", "legal_privacy_policy": "Privacy policy", "legal_return_policy": "Returns & refunds", "legal_seller_agreement": "Seller agreement", "legal_data_deletion": "Delete your data", "contact": "Contact us", "contactSub": "Questions about an order, a return or your account? We are happy to help.", "company": "Company", "address": "Address", "phoneLabel": "Phone", "hours": "Opening hours", "contactSeller": "For questions about a product, chat with the seller from the product page.", "legalNav": "Legal", "legalVersion": "Version {v} · effective {d}", "legalOnlyEn": "This document is not yet available in your language. Showing the English version.",
		"share": "Share", "copied": "Link copied",
	},
	"km": {
		"weAccept": "យើងទទួល", "soonTitle": "ផ្សារនឹងបើកឆាប់ៗនេះ។", "soonText": "អ្នកលក់នៅក្រុងព្រះសីហនុអាចចុះឈ្មោះហាងបានឥឡូវនេះ — ការទិញនឹងបើកឆាប់ៗ។", "soonSell": "បើកហាងរបស់អ្នក", "tagline": "ផ្សារអនឡាញកម្ពុជា", "metaDesc": "មើលវីដេអូផលិតផលពិតពីអ្នកលក់ក្នុងក្រុងព្រះសីហនុ ហើយទិញដោយដឹកជញ្ជូនឥតគិតថ្លៃ។",
		"home": "ទំព័រដើម", "search": "ស្វែងរក", "videos": "វីដេអូ", "cart": "កន្ត្រក", "me": "គណនី", "language": "ភាសា", "mainMenu": "ម៉ឺនុយ",
		"searchPh": "ស្វែងរកផលិតផល ហាង…", "searchBtn": "ស្វែងរក",
		"hero1": "មើលតាមវីដេអូ។", "hero2": "ទិញដោយទំនុកចិត្ត។", "heroSub": "មើល • ស្វែងរក • ទិញ",
		"trust1": "ដឹកឥតគិតថ្លៃក្នុងក្រុងព្រះសីហនុ", "trust2": "វីដេអូផលិតផលពិត", "trust3": "ពិនិត្យដោយ Phsar",
		"categories": "ប្រភេទ", "videosTitle": "មើល ហើយទិញ", "videosSub": "វីដេអូពិតពីអ្នកលក់ក្នុងស្រុក",
		"newTitle": "ទំនិញថ្មី", "newSub": "ទើបតែដាក់លក់", "seeAll": "មើលទាំងអស់", "all": "ទាំងអស់",
		"from": "ចាប់ពី", "soldOut": "អស់ស្តុក", "new": "ថ្មី", "like_new": "ដូចថ្មី", "used_good": "ប្រើរួច – ល្អ", "used_fair": "ប្រើរួច – មធ្យម", "refurbished": "ជួសជុលឡើងវិញ",
		"back": "ត្រឡប់", "productVideo": "វីដេអូផលិតផល", "photo": "រូបថត", "visitStore": "ចូលហាង", "products": "ផលិតផល {n}", "since": "នៅលើ Phsar តាំងពី {d}",
		"returns": "ប្តូរវិញក្នុង {n} ថ្ងៃ", "noReturns": "មិនទទួលប្តូរ", "warranty_store": "ធានាដោយហាង {n} ខែ", "warranty_official": "ធានាផ្លូវការ {n} ខែ",
		"delivery": "ការដឹកជញ្ជូន", "deliveryVal": "ឥតគិតថ្លៃ", "deliveryNote": "ដឹកដោយអ្នកដឹករបស់ Phsar ក្នុងក្រុងព្រះសីហនុ។",
		"description": "ការពិពណ៌នា", "translatedFrom": "បកប្រែដោយស្វ័យប្រវត្តិពី{lang}។", "showOriginal": "បង្ហាញអត្ថបទដើម", "showTranslation": "បង្ហាញការបកប្រែ",
		"conditionTitle": "ស្ថានភាព", "usedFor": "បានប្រើប្រហែល {n} ខែ", "conditionPhotos": "រូបថតស្ថានភាព",
		"choose": "ជ្រើស{opt}", "left": "នៅសល់តែ {n}", "inStock": "មានស្តុក", "outOfStock": "អស់ស្តុកក្នុងជម្រើសនេះ",
		"addToCart": "ដាក់ក្នុងកន្ត្រក", "buyNow": "ទិញឥឡូវ", "comingSoon": "ការបញ្ជាទិញនឹងបើកឆាប់ៗនេះ។ សូមរក្សាទំព័រនេះទុក!",
		"related": "ផលិតផលស្រដៀងគ្នា", "storeProducts": "ផលិតផលពីហាងនេះ", "aboutStore": "អំពីហាង",
		"sortNew": "ថ្មីបំផុត", "sortPriceAsc": "តម្លៃ៖ ទាបទៅខ្ពស់", "sortPriceDesc": "តម្លៃ៖ ខ្ពស់ទៅទាប", "sort": "តម្រៀប",
		"results": "លទ្ធផល {n} សម្រាប់ “{q}”", "noResults": "រកមិនឃើញ “{q}”។ សាកពាក្យផ្សេង ឬមើលតាមប្រភេទ។", "empty": "មិនទាន់មានផលិតផលនៅទីនេះទេ។",
		"searchTitle": "ស្វែងរក", "loadMore": "មើលបន្ថែម", "allProducts": "ផលិតផលទាំងអស់",
		"notFound": "រកទំព័រមិនឃើញ", "notFoundText": "ផលិតផលនេះប្រហែលជាលក់អស់ ឬត្រូវបានដកចេញ។", "goHome": "ទៅទំព័រដើម",
		"footerAbout": "Phsar ជាផ្សារវីដេអូសម្រាប់ក្រុងព្រះសីហនុ កម្ពុជា។", "sell": "លក់នៅលើ Phsar", "sellTab": "លក់", "copyright": "© {y} Phsar",
		"legal_buyer_terms": "លក្ខខណ្ឌប្រើប្រាស់", "legal_privacy_policy": "គោលការណ៍ឯកជនភាព", "legal_return_policy": "ការប្ដូរ និងសងប្រាក់វិញ", "legal_seller_agreement": "កិច្ចព្រមព្រៀងអ្នកលក់", "legal_data_deletion": "លុបទិន្នន័យរបស់អ្នក", "contact": "ទំនាក់ទំនងយើង", "contactSub": "មានសំណួរអំពីការបញ្ជាទិញ ការប្ដូរ ឬគណនីរបស់អ្នក? យើងរីករាយជួយ។", "company": "ក្រុមហ៊ុន", "address": "អាសយដ្ឋាន", "phoneLabel": "ទូរស័ព្ទ", "hours": "ម៉ោងធ្វើការ", "contactSeller": "សម្រាប់សំណួរអំពីផលិតផល សូមជជែកជាមួយអ្នកលក់ពីទំព័រផលិតផល។", "legalNav": "ច្បាប់", "legalVersion": "កំណែ {v} · មានសុពលភាព {d}", "legalOnlyEn": "ឯកសារនេះមិនទាន់មានជាភាសារបស់អ្នកទេ។ កំពុងបង្ហាញកំណែភាសាអង់គ្លេស។",
		"share": "ចែករំលែក", "copied": "បានចម្លងតំណ",
	},
	"zh": {
		"weAccept": "我们接受", "soonTitle": "Phsar 即将开业。", "soonText": "西哈努克的卖家现在就可以注册店铺——购物功能即将开放。", "soonSell": "开设店铺", "tagline": "柬埔寨视频商城", "metaDesc": "观看西哈努克本地卖家的真实商品视频，免费配送，放心购买。",
		"home": "首页", "search": "搜索", "videos": "视频", "cart": "购物车", "me": "我的", "language": "语言", "mainMenu": "主菜单",
		"searchPh": "搜索商品、店铺…", "searchBtn": "搜索",
		"hero1": "先看视频，", "hero2": "再放心购买。", "heroSub": "观看 • 发现 • 购物",
		"trust1": "西港免费配送", "trust2": "真实商品视频", "trust3": "Phsar 审核",
		"categories": "分类", "videosTitle": "看视频购物", "videosSub": "本地卖家的真实视频",
		"newTitle": "新品上架", "newSub": "卖家刚刚上架", "seeAll": "查看全部", "all": "全部",
		"from": "起", "soldOut": "已售罄", "new": "全新", "like_new": "几乎全新", "used_good": "二手 – 良好", "used_fair": "二手 – 一般", "refurbished": "翻新",
		"back": "返回", "productVideo": "商品视频", "photo": "照片", "visitStore": "进店逛逛", "products": "{n} 件商品", "since": "{d} 加入 Phsar",
		"returns": "{n} 天可退换", "noReturns": "不支持退换", "warranty_store": "店铺保修 {n} 个月", "warranty_official": "官方保修 {n} 个月",
		"delivery": "配送", "deliveryVal": "免费", "deliveryNote": "由 Phsar 快递员在西哈努克市内配送。",
		"description": "商品描述", "translatedFrom": "由{lang}自动翻译。", "showOriginal": "查看原文", "showTranslation": "查看翻译",
		"conditionTitle": "成色", "usedFor": "已使用约 {n} 个月", "conditionPhotos": "成色照片",
		"choose": "选择{opt}", "left": "仅剩 {n} 件", "inStock": "有货", "outOfStock": "该规格已售罄",
		"addToCart": "加入购物车", "buyNow": "立即购买", "comingSoon": "下单功能即将开放，请收藏此页面！",
		"related": "相似商品", "storeProducts": "本店商品", "aboutStore": "关于店铺",
		"sortNew": "最新", "sortPriceAsc": "价格从低到高", "sortPriceDesc": "价格从高到低", "sort": "排序",
		"results": "“{q}” 共 {n} 个结果", "noResults": "没有找到 “{q}”。换个关键词，或按分类浏览。", "empty": "这里还没有商品。",
		"searchTitle": "搜索", "loadMore": "加载更多", "allProducts": "全部商品",
		"notFound": "页面不存在", "notFoundText": "该商品可能已售出或下架。", "goHome": "返回首页",
		"footerAbout": "Phsar 是柬埔寨西哈努克的视频商城。", "sell": "在 Phsar 开店", "sellTab": "开店", "copyright": "© {y} Phsar",
		"legal_buyer_terms": "使用条款", "legal_privacy_policy": "隐私政策", "legal_return_policy": "退货与退款", "legal_seller_agreement": "卖家协议", "legal_data_deletion": "删除您的数据", "contact": "联系我们", "contactSub": "关于订单、退货或账户有疑问？我们很乐意帮助您。", "company": "公司", "address": "地址", "phoneLabel": "电话", "hours": "营业时间", "contactSeller": "关于商品的问题，请在商品页面与卖家聊天。", "legalNav": "法律", "legalVersion": "版本 {v} · 生效日期 {d}", "legalOnlyEn": "此文件尚无您的语言版本，正在显示英文版。",
		"share": "分享", "copied": "链接已复制",
	},
}

// T menerjemahkan key; args berpasangan: "n", 3, "q", "kaos".
func T(lang, key string, args ...any) string {
	s, ok := strs[lang][key]
	if !ok {
		s, ok = strs["en"][key]
		if !ok {
			return key
		}
	}
	for i := 0; i+1 < len(args); i += 2 {
		s = strings.ReplaceAll(s, "{"+fmt.Sprint(args[i])+"}", fmt.Sprint(args[i+1]))
	}
	return s
}

func validLang(l string) bool { return l == "en" || l == "km" || l == "zh" }
