package storefront

import (
	"encoding/json"
	"html/template"
)

// tour: teks tur "cara pakai" untuk pengunjung baru website (assets/tour-site.js, mesin tur /ui/tour.js).
var tourStrs = map[string]map[string]string{
	"en": {
		"next": "Next", "back": "Back", "skip": "Skip", "done": "Start shopping", "of": "{a}/{b}",
		"videoT": "Watch, then shop", "videoB": "Short videos from shops in Sihanoukville. Tap a video to see the product and buy it.",
		"searchT": "Search", "searchB": "Type what you need — in Khmer, English or Chinese.",
		"photoT": "Search with a photo", "photoB": "Have a picture of something you want? Tap the camera and we find similar products.",
		"cartT": "Your cart", "cartB": "Products you add wait here. Pay with KHQR / ABA when you are ready.",
		"langT": "Choose your language", "langB": "ភាសាខ្មែរ, English or 中文 — change it anytime.",
	},
	"km": {
		"next": "បន្ទាប់", "back": "ថយក្រោយ", "skip": "រំលង", "done": "ចាប់ផ្តើមទិញ", "of": "{a}/{b}",
		"videoT": "មើល រួចទិញ", "videoB": "វីដេអូខ្លីៗពីហាងនានានៅក្រុងព្រះសីហនុ។ ចុចលើវីដេអូ ដើម្បីមើលផលិតផល ហើយទិញ។",
		"searchT": "ស្វែងរក", "searchB": "វាយបញ្ចូលអ្វីដែលអ្នកត្រូវការ — ជាភាសាខ្មែរ អង់គ្លេស ឬចិន។",
		"photoT": "ស្វែងរកដោយរូបថត", "photoB": "មានរូបភាពរបស់ដែលអ្នកចង់បាន? ចុចកាមេរ៉ា យើងនឹងរកផលិតផលស្រដៀងគ្នាឱ្យ។",
		"cartT": "កន្ត្រករបស់អ្នក", "cartB": "ផលិតផលដែលអ្នកបន្ថែម រង់ចាំនៅទីនេះ។ បង់ប្រាក់តាម KHQR / ABA ពេលអ្នករួចរាល់។",
		"langT": "ជ្រើសរើសភាសា", "langB": "ភាសាខ្មែរ, English ឬ 中文 — ប្តូរបានគ្រប់ពេល។",
	},
	"zh": {
		"next": "下一步", "back": "上一步", "skip": "跳过", "done": "开始购物", "of": "{a}/{b}",
		"videoT": "先看视频，再下单", "videoB": "来自西哈努克港本地店铺的短视频。点击视频即可查看并购买商品。",
		"searchT": "搜索", "searchB": "输入你想要的东西——支持高棉语、英语和中文。",
		"photoT": "拍照搜索", "photoB": "有想买的东西的照片？点击相机图标，我们帮你找到相似商品。",
		"cartT": "购物车", "cartB": "加入的商品都在这里。准备好后用 KHQR / ABA 付款。",
		"langT": "选择语言", "langB": "ភាសាខ្មែរ、English 或中文——随时可以切换。",
	},
}

// TourJS: teks tur untuk bahasa halaman (kosong → Inggris), aman disisipkan di <script>.
func (p *Page) TourJS() template.JS {
	m := map[string]string{}
	for k, v := range tourStrs["en"] {
		m[k] = v
	}
	for k, v := range tourStrs[p.Lang] {
		m[k] = v
	}
	b, _ := json.Marshal(m)
	return template.JS(b)
}
