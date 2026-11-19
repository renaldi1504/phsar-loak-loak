package askbot

// Jawaban berbasis aturan (tanpa Claude): kenali maksud pertanyaan dari kata kunci, lalu jawab
// HANYA dari data (stok, harga, varian, garansi, retur, kondisi, jam buka) atau FAQ yang cocok.
// Tidak ada yang cocok → "seller akan mengecek".

import (
	"fmt"
	"strings"
	"unicode"
)

var intentWords = map[string][]string{
	"price":     {"price", "how much", "cost", "cheap", "discount", "$", "តម្លៃ", "ប៉ុន្មាន", "ថ្លៃ", "价格", "价钱", "多少钱", "多少", "便宜", "优惠"},
	"stock":     {"stock", "available", "availability", "in stock", "sold out", "left", "still have", "do you have", "ស្តុក", "ស្ដុក", "នៅមាន", "មានទេ", "អស់", "库存", "有货", "还有", "有没有", "现货", "卖完"},
	"variant":   {"size", "color", "colour", "variant", "option", "ទំហំ", "ពណ៌", "សាយ", "ម៉ូដ", "尺码", "尺寸", "颜色", "款式", "型号", "规格"},
	"warranty":  {"warranty", "guarantee", "guaranty", "ធានា", "保修", "质保", "保固"},
	"return":    {"return", "refund", "exchange", "send back", "ប្តូរ", "ប្ដូរ", "សងប្រាក់", "យកមកវិញ", "退货", "退款", "换货", "退换"},
	"hours":     {"open", "opening", "hours", "close", "closing", "what time", "បើក", "ម៉ោង", "បិទ", "营业", "几点", "开门", "关门", "营业时间"},
	"condition": {"condition", "new", "used", "second hand", "secondhand", "original", "scratch", "ថ្មី", "មួយទឹក", "ស្ថានភាព", "ជជុះ", "新的", "全新", "二手", "成色", "划痕"},
}

var intentOrder = []string{"stock", "price", "variant", "condition", "warranty", "return", "hours"}

type phrases struct {
	hello, thanks, check, more                                       string
	price1, priceRange, stockYes, stockNo, stockFew, notForSale      string
	options, optionLine                                              string
	warrantyNone, warrantyStore, warrantyOfficial                    string
	returnYes, returnNo, condition, usedFor, hoursIntro, closed, sep string
	preorder                                                         string
	days                                                             [7]string
}

var texts = map[string]phrases{
	"en": {hello: "Hello!", thanks: "Thank you for your question!", check: "We will check and reply to you soon.", more: "Feel free to ask if you have other questions.",
		price1: "The price is $%s.", priceRange: "The price is $%s to $%s, depending on the option.",
		stockYes: "It is in stock (%d available).", stockNo: "Sorry, it is sold out at the moment.", stockFew: "Only %d left.",
		notForSale: "Sorry, this item is not available for sale right now.",
		options:    "Available options:", optionLine: "%s: $%s (%d in stock)",
		warrantyNone: "This item has no warranty.", warrantyStore: "It comes with a %d-month store warranty.", warrantyOfficial: "It comes with a %d-month official warranty.",
		returnYes: "You can request a return within %d days after delivery.", returnNo: "This item cannot be returned.",
		condition: "Condition: %s.", usedFor: "Used for about %d months.", hoursIntro: "Our opening hours:", closed: "closed", sep: "; ",
		preorder: "This is a pre-order item: it ships within %d days.",
		days:     [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}},
	"km": {hello: "សួស្តី!", thanks: "អរគុណសម្រាប់សំណួររបស់អ្នក!", check: "យើងនឹងពិនិត្យ ហើយឆ្លើយតបអ្នកឆាប់ៗ។", more: "បើមានសំណួរផ្សេងទៀត សូមសួរមកបាន។",
		price1: "តម្លៃ $%s។", priceRange: "តម្លៃពី $%s ដល់ $%s អាស្រ័យលើជម្រើស។",
		stockYes: "នៅមានក្នុងស្តុក (%d)។", stockNo: "សូមអភ័យទោស ពេលនេះអស់ពីស្តុកហើយ។", stockFew: "នៅសល់តែ %d ប៉ុណ្ណោះ។",
		notForSale: "សូមអភ័យទោស ទំនិញនេះមិនទាន់លក់នៅពេលនេះទេ។",
		options:    "ជម្រើសដែលមាន៖", optionLine: "%s៖ $%s (មាន %d)",
		warrantyNone: "ទំនិញនេះគ្មានការធានាទេ។", warrantyStore: "មានការធានាដោយហាង %d ខែ។", warrantyOfficial: "មានការធានាផ្លូវការ %d ខែ។",
		returnYes: "អ្នកអាចស្នើប្តូរវិញក្នុងរយៈពេល %d ថ្ងៃ បន្ទាប់ពីទទួលបាន។", returnNo: "ទំនិញនេះមិនទទួលប្តូរវិញទេ។",
		condition: "ស្ថានភាព៖ %s។", usedFor: "បានប្រើប្រហែល %d ខែ។", hoursIntro: "ម៉ោងបើករបស់យើង៖", closed: "បិទ", sep: "; ",
		preorder: "ទំនិញនេះជាការកម្មង់ទុកមុន៖ ផ្ញើក្នុងរយៈពេល %d ថ្ងៃ។",
		days:     [7]string{"អាទិត្យ", "ច័ន្ទ", "អង្គារ", "ពុធ", "ព្រហស្បតិ៍", "សុក្រ", "សៅរ៍"}},
	"zh": {hello: "您好！", thanks: "感谢您的提问！", check: "我们会尽快确认并回复您。", more: "如有其他问题，欢迎随时询问。",
		price1: "价格是 $%s。", priceRange: "价格为 $%s 至 $%s，视款式而定。",
		stockYes: "有现货（%d 件）。", stockNo: "抱歉，目前已售罄。", stockFew: "仅剩 %d 件。",
		notForSale: "抱歉，该商品目前暂不销售。",
		options:    "可选款式：", optionLine: "%s：$%s（库存 %d）",
		warrantyNone: "该商品没有保修。", warrantyStore: "提供 %d 个月店铺保修。", warrantyOfficial: "提供 %d 个月官方保修。",
		returnYes: "收货后 %d 天内可申请退换。", returnNo: "该商品不支持退换。",
		condition: "成色：%s。", usedFor: "已使用约 %d 个月。", hoursIntro: "我们的营业时间：", closed: "休息", sep: "；",
		preorder: "该商品为预售，%d 天内发货。",
		days:     [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}},
}

var conditionLocal = map[string]map[string]string{
	"km": {"new": "ថ្មី", "like new": "ដូចថ្មី", "used - good": "ប្រើរួច - ល្អ", "used - fair": "ប្រើរួច - មធ្យម", "refurbished": "ជួសជុលឡើងវិញ"},
	"zh": {"new": "全新", "like new": "几乎全新", "used - good": "二手 - 良好", "used - fair": "二手 - 一般", "refurbished": "翻新"},
}

// Intents: maksud yang dikenali dari pertanyaan (urutan tetap).
func Intents(q string) []string {
	low := strings.ToLower(q)
	out := []string{}
	for _, k := range intentOrder {
		for _, w := range intentWords[k] {
			if containsWord(low, w) {
				out = append(out, k)
				break
			}
		}
	}
	return out
}

// containsWord: kata Latin harus utuh ("new" tidak cocok dengan "news"); aksara lain cukup substring.
func containsWord(s, w string) bool {
	if !isLatin(w) {
		return strings.Contains(s, w)
	}
	for i := 0; ; {
		j := strings.Index(s[i:], w)
		if j < 0 {
			return false
		}
		start, end := i+j, i+j+len(w)
		okL := start == 0 || !isWordByte(s[start-1])
		okR := end == len(s) || !isWordByte(s[end])
		if okL && okR {
			return true
		}
		i = start + 1
	}
}

func isLatin(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func isWordByte(b byte) bool { return b >= 'a' && b <= 'z' || b >= '0' && b <= '9' }

// RuleReply: draf tanpa AI.
func RuleReply(lang, q string, f Facts) string {
	t, ok := texts[lang]
	if !ok {
		t, lang = texts["en"], "en"
	}
	parts := []string{}
	answered := false
	if faq := matchFAQ(lang, q, f.FAQs); faq != "" {
		parts = append(parts, faq)
		answered = true
	}
	p := f.Product
	intents := Intents(q)
	if p != nil && !p.ForSale && len(intents) > 0 && hasAny(intents, "stock", "price", "variant") {
		parts = append(parts, t.notForSale)
		answered = true
		intents = without(intents, "stock", "price", "variant")
	}
	listed := false // daftar varian sudah ditulis (mencakup harga + stok)
	for _, in := range intents {
		switch in {
		case "stock", "price", "variant":
			if p == nil || len(p.Variants) == 0 || listed {
				continue
			}
			if len(p.Variants) > 1 && p.Variants[0].Label != "" {
				if (in == "variant" || len(p.Variants) <= 6) && len(p.Variants) <= 12 {
					lines := []string{}
					for _, v := range p.Variants {
						lines = append(lines, fmt.Sprintf(t.optionLine, v.Label, v.Price, v.Stock))
					}
					parts = append(parts, t.options+" "+strings.Join(lines, t.sep))
					listed = true
					answered = true
					continue
				}
			}
			if in == "price" || in == "variant" {
				lo, hi := priceRange(p.Variants)
				if lo == hi {
					parts = append(parts, fmt.Sprintf(t.price1, lo))
				} else {
					parts = append(parts, fmt.Sprintf(t.priceRange, lo, hi))
				}
				answered = true
			}
			if in == "stock" {
				n := 0
				for _, v := range p.Variants {
					n += v.Stock
				}
				switch {
				case n <= 0:
					parts = append(parts, t.stockNo)
				case n <= 3:
					parts = append(parts, fmt.Sprintf(t.stockFew, n))
				default:
					parts = append(parts, fmt.Sprintf(t.stockYes, n))
				}
				answered = true
			}
		case "condition":
			if p == nil {
				continue
			}
			c := p.Condition
			if m := conditionLocal[lang]; m != nil && m[c] != "" {
				c = m[c]
			}
			s := fmt.Sprintf(t.condition, c)
			if p.UsageMonths != nil && *p.UsageMonths > 0 {
				s += " " + fmt.Sprintf(t.usedFor, *p.UsageMonths)
			}
			if p.ConditionNote != "" {
				s += " " + clip(p.ConditionNote, 300)
			}
			parts = append(parts, s)
			answered = true
		case "warranty":
			if p == nil {
				continue
			}
			switch {
			case p.WarrantyType == "none" || p.WarrantyMonths <= 0:
				parts = append(parts, t.warrantyNone)
			case p.WarrantyType == "official":
				parts = append(parts, fmt.Sprintf(t.warrantyOfficial, p.WarrantyMonths))
			default:
				parts = append(parts, fmt.Sprintf(t.warrantyStore, p.WarrantyMonths))
			}
			answered = true
		case "return":
			if p == nil {
				continue
			}
			if p.ReturnDays > 0 {
				parts = append(parts, fmt.Sprintf(t.returnYes, p.ReturnDays))
			} else {
				parts = append(parts, t.returnNo)
			}
			answered = true
		case "hours":
			if f.Store.week == nil {
				continue
			}
			days := []string{}
			for i, d := range f.Store.week {
				if d == nil {
					days = append(days, t.days[i]+" "+t.closed)
				} else {
					days = append(days, t.days[i]+" "+d.Open+"–"+d.Close)
				}
			}
			parts = append(parts, t.hoursIntro+" "+strings.Join(days, t.sep))
			answered = true
		}
	}
	if p != nil && p.PreorderDays != nil && hasAny(intents, "stock") {
		parts = append(parts, fmt.Sprintf(t.preorder, *p.PreorderDays))
	}
	if !answered {
		return t.hello + " " + t.thanks + " " + t.check
	}
	return t.hello + " " + strings.Join(parts, " ") + " " + t.more
}

func priceRange(vs []VariantFact) (string, string) {
	lo, hi := vs[0].Price, vs[0].Price
	for _, v := range vs[1:] {
		if num(v.Price) < num(lo) {
			lo = v.Price
		}
		if num(v.Price) > num(hi) {
			hi = v.Price
		}
	}
	return lo, hi
}

func num(s string) float64 {
	var f float64
	_, _ = fmt.Sscanf(s, "%g", &f)
	return f
}

func hasAny(list []string, keys ...string) bool {
	for _, a := range list {
		for _, k := range keys {
			if a == k {
				return true
			}
		}
	}
	return false
}

func without(list []string, keys ...string) []string {
	out := []string{}
	for _, a := range list {
		if !hasAny([]string{a}, keys...) {
			out = append(out, a)
		}
	}
	return out
}

// matchFAQ: FAQ berbahasa sama yang paling mirip (bigram karakter), kalau cukup mirip.
func matchFAQ(lang, q string, faqs []FAQFact) string {
	qb := bigrams(q)
	if len(qb) == 0 {
		return ""
	}
	best, bestScore := "", 0.0
	for _, f := range faqs {
		if f.Lang != lang {
			continue
		}
		fb := bigrams(f.Question)
		if len(fb) == 0 {
			continue
		}
		common := 0
		for g := range fb {
			if qb[g] {
				common++
			}
		}
		// Seberapa banyak isi pertanyaan FAQ muncul di pertanyaan pembeli (dan sebaliknya).
		score := float64(common) / float64(min(len(fb), len(qb)))
		if f.Product {
			score += 0.05
		}
		if score > bestScore {
			best, bestScore = f.Answer, score
		}
	}
	if bestScore >= 0.6 {
		return best
	}
	return ""
}

func bigrams(s string) map[string]bool {
	rs := []rune{}
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r) {
			rs = append(rs, r)
		} else if len(rs) > 0 && rs[len(rs)-1] != ' ' {
			rs = append(rs, ' ')
		}
	}
	out := map[string]bool{}
	for i := 0; i+1 < len(rs); i++ {
		if rs[i] == ' ' || rs[i+1] == ' ' {
			continue
		}
		out[string(rs[i:i+2])] = true
	}
	return out
}
