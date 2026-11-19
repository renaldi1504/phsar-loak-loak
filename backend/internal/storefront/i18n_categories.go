package storefront

// Tampilan kategori v2: kotak kategori, musim ini, subkategori, Tonton per kategori, menu kategori desktop.
var categoryStrs = map[string]map[string]string{
	"en": {"seasonTitle": "This season", "seasonUntil": "Until {d}", "seasonAllYear": "All year", "newVideos": "{n} new videos today",
		"tileNew": "{n} new videos", "sortPopular": "Best selling", "productsTab": "Products", "watchTab": "Watch",
		"watchTitle": "{c} videos", "watchEmpty": "No products in this category yet.", "watchPhotos": "Photos", "watchHint": "Swipe up for the next video",
		"watchBuy": "Buy", "watchMuted": "Tap for sound", "allCats": "All categories", "seeAllIn": "See all {c}",
		"watchNew": "Watch {c} videos", "seasonEmpty": "Products for this season are coming soon."},
	"km": {"seasonTitle": "រដូវនេះ", "seasonUntil": "ដល់ {d}", "seasonAllYear": "ពេញមួយឆ្នាំ", "newVideos": "វីដេអូថ្មី {n} ថ្ងៃនេះ",
		"tileNew": "វីដេអូថ្មី {n}", "sortPopular": "លក់ដាច់", "productsTab": "ទំនិញ", "watchTab": "មើលវីដេអូ",
		"watchTitle": "វីដេអូ {c}", "watchEmpty": "មិនទាន់មានផលិតផលក្នុងប្រភេទនេះនៅឡើយទេ។", "watchPhotos": "រូបថត", "watchHint": "អូសឡើងលើ ដើម្បីមើលវីដេអូបន្ទាប់",
		"watchBuy": "ទិញ", "watchMuted": "ចុចដើម្បីបើកសំឡេង", "allCats": "ប្រភេទទាំងអស់", "seeAllIn": "មើល{c}ទាំងអស់",
		"watchNew": "មើលវីដេអូ {c}", "seasonEmpty": "ទំនិញសម្រាប់រដូវនេះនឹងមកដល់ឆាប់ៗ។"},
	"zh": {"seasonTitle": "本季精选", "seasonUntil": "至 {d}", "seasonAllYear": "全年", "newVideos": "今日新视频 {n} 个",
		"tileNew": "{n} 个新视频", "sortPopular": "销量", "productsTab": "商品", "watchTab": "看视频",
		"watchTitle": "{c}视频", "watchEmpty": "此分类暂无商品。", "watchPhotos": "照片", "watchHint": "上滑观看下一个视频",
		"watchBuy": "购买", "watchMuted": "点击开启声音", "allCats": "全部分类", "seeAllIn": "查看全部{c}",
		"watchNew": "观看{c}视频", "seasonEmpty": "本季商品即将上架。"},
}

func init() {
	for l, m := range categoryStrs {
		for k, v := range m {
			strs[l][k] = v
		}
	}
}
