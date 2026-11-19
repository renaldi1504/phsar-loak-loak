package storefront

// Halaman "Jual di Phsar" (/{lang}/sell): ajakan untuk calon seller — keuntungan, langkah daftar, yang perlu
// disiapkan, FAQ — lalu tombol ke Seller Center (login = daftar: nomor baru otomatis dibuatkan akun).

import (
	"context"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

type sellData struct {
	Rate      string // komisi global saat ini, mis. "8"
	Methods   string // cara daftar yang aktif sekarang, mis. "Telegram, Google or your phone number"
	SellerURL string // Seller Center dengan bahasa yang sama
}

// commissionRate: tarif komisi global yang berlaku (aturan admin → Komisi), bawaan 8%.
func (h *Handler) commissionRate(ctx context.Context) string {
	rate := decimal.NewFromInt(8)
	if h.Cat != nil && h.Cat.Pool != nil {
		var r decimal.Decimal
		if err := h.Cat.Pool.QueryRow(ctx, `SELECT rate_percent FROM commission_rules
			WHERE scope = 'global' AND valid_from <= now() AND (valid_to IS NULL OR valid_to > now())
			ORDER BY valid_from DESC LIMIT 1`).Scan(&r); err == nil {
			rate = r
		}
	}
	return strings.TrimSuffix(strings.TrimRight(rate.StringFixed(2), "0"), ".")
}

// loginMethods: cara login yang benar-benar aktif (mengikuti pengaturan admin), dalam bahasa halaman.
func (h *Handler) loginMethods(ctx context.Context, lang string) string {
	var m []string
	if h.TelegramLogin {
		m = append(m, "Telegram")
	}
	if h.Logins != nil && h.Logins.GoogleEnabled(ctx) {
		m = append(m, "Google")
	}
	if h.Logins != nil && h.Logins.FacebookEnabled(ctx) {
		m = append(m, "Facebook")
	}
	if h.Logins == nil || h.Logins.SMSEnabled(ctx) || len(m) == 0 {
		m = append(m, T(lang, "sl_phone"))
	}
	if len(m) == 1 {
		return m[0]
	}
	sep := ", "
	if lang == "zh" {
		sep = "、"
	}
	return strings.Join(m[:len(m)-1], sep) + " " + T(lang, "sl_or") + " " + m[len(m)-1]
}

func (h *Handler) sellPage(w http.ResponseWriter, r *http.Request, p *Page) (string, int, error) {
	p.Title, p.Description, p.Active = T(p.Lang, "sl_title")+" · Phsar", T(p.Lang, "sl_heroSub"), "sell"
	p.Data = sellData{Rate: h.commissionRate(r.Context()), Methods: h.loginMethods(r.Context(), p.Lang), SellerURL: "/seller/?lang=" + p.Lang}
	return "sell", http.StatusOK, nil
}

var sellStrs = map[string]map[string]string{
	"en": {
		"sl_phone": "your phone number", "sl_or": "or",
		"sl_title": "Sell on Phsar", "sl_badge": "For shops in Sihanoukville",
		"sl_hero": "Open your shop on Phsar — free", "sl_heroSub": "Reach buyers in Sihanoukville with photos and short videos. Buyers pay first, our courier picks up, you get paid to your bank.",
		"sl_cta": "Register as a seller", "sl_ctaSub": "Already a seller? The same button signs you in.",
		"sl_soon": "Phsar opens to buyers soon. Register now so your shop and products are ready on opening day.",
		"sl_whyT": "Why sell on Phsar",
		"sl_w1T":  "Free to open", "sl_w1B": "No sign-up fee, no monthly fee.",
		"sl_w2T": "Only {rate}% when you sell", "sl_w2B": "Commission is taken only from items sold. Nothing sold, nothing to pay. No payment fees for you.",
		"sl_w3T": "Buyers pay first", "sl_w3B": "Every order is paid by KHQR / ABA before you pack it — no cash-on-delivery risk.",
		"sl_w4T": "We pick up", "sl_w4B": "A Phsar courier collects orders from your pickup point and delivers them to the buyer.",
		"sl_w5T": "Paid to your bank", "sl_w5B": "Your earnings are transferred to your bank account on the payout schedule.",
		"sl_w6T": "Write once, 3 languages", "sl_w6B": "Write your product in one language — Phsar translates it to Khmer, English and Chinese.",
		"sl_stepsT": "How to start", "sl_stepsSub": "About 10–15 minutes. You can stop and continue later.",
		"sl_s1T": "Sign up in a minute", "sl_s1B": "Use what you already have — {methods}. Your account is created automatically.",
		"sl_s2T": "Create your shop", "sl_s2B": "Shop name, a short description and your logo.",
		"sl_s3T": "Verify your identity", "sl_s3B": "A photo of your ID card or passport (companies: business documents).",
		"sl_s4T": "Add your pickup point", "sl_s4B": "The place in Sihanoukville where our courier collects your orders.",
		"sl_s5T": "Send for review, add products", "sl_s5B": "Phsar checks your shop and tells you by notification. Meanwhile you can already prepare products.",
		"sl_prepT": "What to prepare",
		"sl_p1":    "One way to sign in: {methods}", "sl_p2": "Photo of your ID card or passport", "sl_p3": "A bank account in your own name or your company’s name",
		"sl_p4": "At least 3 photos per product — a short video is recommended", "sl_p5": "Your pickup address in Sihanoukville",
		"sl_p6":   "Companies: business registration and TIN (individuals without TIN can sell too)",
		"sl_faqT": "Questions",
		"sl_q1":   "Does it cost anything?", "sl_a1": "Opening a shop is free. Phsar takes {rate}% only from items you sell.",
		"sl_q2": "Where do you deliver?", "sl_a2": "For now in Sihanoukville. Our courier picks up from your pickup point and delivers to the buyer.",
		"sl_q3": "When do I get my money?", "sl_a3": "After the buyer receives the order, the money is added to your balance and transferred to your bank on the payout schedule.",
		"sl_q4": "Do I need a video?", "sl_a4": "No — 3 photos are enough. Products with a video usually sell better.",
		"sl_q5": "Can I sell used items?", "sl_a5": "Yes. New and used items are both welcome — just choose the condition.",
		"sl_q6": "How many shops can I open?", "sl_a6": "One account owns one shop. You can invite staff to help you run it.",
		"sl_endT": "Ready to start?", "sl_endB": "Sign up in a minute — your shop can be ready today.",
	},
	"km": {
		"sl_phone": "លេខទូរស័ព្ទ", "sl_or": "ឬ",
		"sl_title": "លក់លើ Phsar", "sl_badge": "សម្រាប់ហាងនៅក្រុងព្រះសីហនុ",
		"sl_hero": "បើកហាងរបស់អ្នកលើ Phsar — ឥតគិតថ្លៃ", "sl_heroSub": "ឈានដល់អ្នកទិញនៅក្រុងព្រះសីហនុ តាមរយៈរូបថត និងវីដេអូខ្លីៗ។ អ្នកទិញបង់ប្រាក់មុន អ្នកដឹករបស់យើងមកយក ហើយអ្នកទទួលប្រាក់ចូលគណនីធនាគារ។",
		"sl_cta": "ចុះឈ្មោះជាអ្នកលក់", "sl_ctaSub": "ជាអ្នកលក់រួចហើយ? ប៊ូតុងដដែលនេះ ក៏ប្រើសម្រាប់ចូលគណនីដែរ។",
		"sl_soon": "Phsar នឹងបើកសម្រាប់អ្នកទិញឆាប់ៗនេះ។ ចុះឈ្មោះឥឡូវ ដើម្បីឲ្យហាង និងផលិតផលរបស់អ្នករួចរាល់នៅថ្ងៃបើក។",
		"sl_whyT": "ហេតុអ្វីលក់លើ Phsar",
		"sl_w1T":  "បើកហាងឥតគិតថ្លៃ", "sl_w1B": "គ្មានថ្លៃចុះឈ្មោះ គ្មានថ្លៃប្រចាំខែ។",
		"sl_w2T": "តែ {rate}% ពេលលក់ដាច់", "sl_w2B": "កម្រៃជើងសារគិតតែលើទំនិញដែលលក់ដាច់ប៉ុណ្ណោះ។ មិនលក់ដាច់ មិនបង់អ្វីទេ។ គ្មានថ្លៃទូទាត់សម្រាប់អ្នក។",
		"sl_w3T": "អ្នកទិញបង់ប្រាក់មុន", "sl_w3B": "រាល់ការបញ្ជាទិញត្រូវបានបង់តាម KHQR / ABA មុនពេលអ្នកវេចខ្ចប់ — គ្មានហានិភ័យបង់ប្រាក់ពេលទទួលទំនិញ។",
		"sl_w4T": "យើងមកយកដល់កន្លែង", "sl_w4B": "អ្នកដឹករបស់ Phsar មកយកការបញ្ជាទិញពីទីតាំងទទួលកញ្ចប់របស់អ្នក ហើយដឹកជូនអ្នកទិញ។",
		"sl_w5T": "ប្រាក់ចូលគណនីធនាគារ", "sl_w5B": "ប្រាក់ចំណូលរបស់អ្នកត្រូវបានផ្ទេរចូលគណនីធនាគារតាមកាលវិភាគទូទាត់។",
		"sl_w6T": "សរសេរម្តង បាន 3 ភាសា", "sl_w6B": "សរសេរផលិតផលជាភាសាមួយ — Phsar បកប្រែជាភាសាខ្មែរ អង់គ្លេស និងចិនឲ្យ។",
		"sl_stepsT": "របៀបចាប់ផ្តើម", "sl_stepsSub": "ប្រហែល 10–15 នាទី។ អ្នកអាចឈប់ ហើយបន្តពេលក្រោយបាន។",
		"sl_s1T": "ចុះឈ្មោះក្នុងរយៈពេលមួយនាទី", "sl_s1B": "ប្រើអ្វីដែលអ្នកមានស្រាប់ — {methods}។ គណនីនឹងត្រូវបង្កើតដោយស្វ័យប្រវត្តិ។",
		"sl_s2T": "បង្កើតហាងរបស់អ្នក", "sl_s2B": "ឈ្មោះហាង ការពិពណ៌នាខ្លី និងឡូហ្គោ។",
		"sl_s3T": "ផ្ទៀងផ្ទាត់អត្តសញ្ញាណ", "sl_s3B": "រូបថតអត្តសញ្ញាណប័ណ្ណ ឬលិខិតឆ្លងដែន (ក្រុមហ៊ុន៖ ឯកសារអាជីវកម្ម)។",
		"sl_s4T": "បន្ថែមទីតាំងទទួលកញ្ចប់", "sl_s4B": "កន្លែងនៅក្រុងព្រះសីហនុ ដែលអ្នកដឹករបស់យើងមកយកការបញ្ជាទិញ។",
		"sl_s5T": "ផ្ញើឲ្យពិនិត្យ ហើយបន្ថែមផលិតផល", "sl_s5B": "Phsar ពិនិត្យហាងរបស់អ្នក ហើយជូនដំណឹងទៅអ្នក។ ក្នុងពេលនោះ អ្នកអាចរៀបចំផលិតផលបានហើយ។",
		"sl_prepT": "អ្វីដែលត្រូវរៀបចំ",
		"sl_p1":    "វិធីចូលមួយ៖ {methods}", "sl_p2": "រូបថតអត្តសញ្ញាណប័ណ្ណ ឬលិខិតឆ្លងដែន", "sl_p3": "គណនីធនាគារក្នុងឈ្មោះរបស់អ្នក ឬឈ្មោះក្រុមហ៊ុន",
		"sl_p4": "យ៉ាងតិច 3 រូបថតក្នុងមួយផលិតផល — ណែនាំឲ្យមានវីដេអូខ្លី", "sl_p5": "អាសយដ្ឋានទទួលកញ្ចប់នៅក្រុងព្រះសីហនុ",
		"sl_p6":   "ក្រុមហ៊ុន៖ ការចុះបញ្ជីអាជីវកម្ម និងលេខអត្តសញ្ញាណកម្មសារពើពន្ធ (បុគ្គលគ្មាន TIN ក៏អាចលក់បាន)",
		"sl_faqT": "សំណួរ",
		"sl_q1":   "តើត្រូវបង់ថ្លៃអ្វីទេ?", "sl_a1": "បើកហាងឥតគិតថ្លៃ។ Phsar យកតែ {rate}% ពីទំនិញដែលអ្នកលក់ដាច់។",
		"sl_q2": "ដឹកជញ្ជូនទៅណាខ្លះ?", "sl_a2": "បច្ចុប្បន្ននៅក្រុងព្រះសីហនុ។ អ្នកដឹករបស់យើងមកយកពីទីតាំងរបស់អ្នក ហើយដឹកជូនអ្នកទិញ។",
		"sl_q3": "ពេលណាខ្ញុំទទួលបានប្រាក់?", "sl_a3": "បន្ទាប់ពីអ្នកទិញទទួលបានទំនិញ ប្រាក់ត្រូវបានបញ្ចូលក្នុងសមតុល្យរបស់អ្នក ហើយផ្ទេរចូលធនាគារតាមកាលវិភាគទូទាត់។",
		"sl_q4": "តើត្រូវមានវីដេអូទេ?", "sl_a4": "មិនចាំបាច់ទេ — 3 រូបថតគឺគ្រប់គ្រាន់។ ផលិតផលដែលមានវីដេអូ ជាធម្មតាលក់ដាច់ជាង។",
		"sl_q5": "អាចលក់ទំនិញជជុះបានទេ?", "sl_a5": "បាន។ ទាំងទំនិញថ្មី និងជជុះ — គ្រាន់តែជ្រើសរើសស្ថានភាព។",
		"sl_q6": "អាចបើកហាងបានប៉ុន្មាន?", "sl_a6": "មួយគណនី មួយហាង។ អ្នកអាចអញ្ជើញបុគ្គលិកមកជួយគ្រប់គ្រងបាន។",
		"sl_endT": "រួចរាល់ហើយឬនៅ?", "sl_endB": "ចុះឈ្មោះក្នុងរយៈពេលមួយនាទី — ហាងរបស់អ្នកអាចរួចរាល់នៅថ្ងៃនេះ។",
	},
	"zh": {
		"sl_phone": "手机号", "sl_or": "或",
		"sl_title": "在 Phsar 开店", "sl_badge": "面向西哈努克港的商家",
		"sl_hero": "在 Phsar 免费开店", "sl_heroSub": "用照片和短视频触达西哈努克港的买家。买家先付款，我们的快递员上门取件，货款直接打到你的银行账户。",
		"sl_cta": "注册成为卖家", "sl_ctaSub": "已经是卖家？同一个按钮即可登录。",
		"sl_soon": "Phsar 即将向买家开放。现在注册，开业当天你的店铺和商品就能准备就绪。",
		"sl_whyT": "为什么在 Phsar 开店",
		"sl_w1T":  "免费开店", "sl_w1B": "没有注册费，没有月费。",
		"sl_w2T": "成交才收 {rate}%", "sl_w2B": "只对已售出的商品收取佣金。没卖出就不用付费，也不收支付手续费。",
		"sl_w3T": "买家先付款", "sl_w3B": "每笔订单在你打包前已通过 KHQR / ABA 付款——没有货到付款的风险。",
		"sl_w4T": "上门取件", "sl_w4B": "Phsar 快递员到你的取件点收货，并送到买家手中。",
		"sl_w5T": "货款打到银行", "sl_w5B": "你的收入按结算周期转入你的银行账户。",
		"sl_w6T": "写一次，三种语言", "sl_w6B": "用一种语言填写商品信息——Phsar 自动翻译成高棉语、英语和中文。",
		"sl_stepsT": "如何开始", "sl_stepsSub": "大约 10–15 分钟，可以随时暂停稍后继续。",
		"sl_s1T": "一分钟完成注册", "sl_s1B": "用你已有的方式——{methods}，账户会自动创建。",
		"sl_s2T": "创建店铺", "sl_s2B": "店铺名称、简短介绍和标志。",
		"sl_s3T": "身份验证", "sl_s3B": "身份证或护照照片（公司：营业文件）。",
		"sl_s4T": "添加取件点", "sl_s4B": "快递员在西哈努克港收取你订单的地点。",
		"sl_s5T": "提交审核，上架商品", "sl_s5B": "Phsar 审核你的店铺并通过通知告诉你。在此期间你已经可以准备商品。",
		"sl_prepT": "需要准备什么",
		"sl_p1":    "任一登录方式：{methods}", "sl_p2": "身份证或护照照片", "sl_p3": "本人或公司名下的银行账户",
		"sl_p4": "每件商品至少 3 张照片——建议加一段短视频", "sl_p5": "你在西哈努克港的取件地址",
		"sl_p6":   "公司：营业执照和税号 TIN（没有 TIN 的个人也可以开店）",
		"sl_faqT": "常见问题",
		"sl_q1":   "需要付费吗？", "sl_a1": "开店免费。Phsar 只对你卖出的商品收取 {rate}%。",
		"sl_q2": "配送到哪里？", "sl_a2": "目前仅限西哈努克港。快递员从你的取件点取货，送到买家手中。",
		"sl_q3": "什么时候能收到钱？", "sl_a3": "买家收到订单后，货款计入你的余额，并按结算周期转入你的银行账户。",
		"sl_q4": "必须有视频吗？", "sl_a4": "不必须——3 张照片就够了。有视频的商品通常卖得更好。",
		"sl_q5": "可以卖二手商品吗？", "sl_a5": "可以。全新和二手商品都欢迎——选择成色即可。",
		"sl_q6": "可以开几家店？", "sl_a6": "一个账户拥有一家店铺。你可以邀请员工帮忙管理。",
		"sl_endT": "准备好了吗？", "sl_endB": "一分钟完成注册——你的店铺今天就能准备好。",
	},
}

func init() {
	funcs["items"] = func(v ...any) []any { return v } // {{range (items "a" "b")}} di template
	for l, m := range sellStrs {
		for k, v := range m {
			strs[l][k] = v
		}
	}
}
