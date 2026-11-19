package storefront

// shipping: teks ongkir berbayar (dipakai kalau admin mengisi ongkir > 0 di Shipping & FX rates).
func init() {
	page := map[string]map[string]string{
		"en": {"trust1Paid": "Fast delivery in Sihanoukville", "deliveryNotePaid": "Delivered by a Phsar courier within Sihanoukville. Fee per store.",
			"shipFreeOver": "free over {a}", "shipFreeOverStore": "Free when you spend {a} or more in one store."},
		"km": {"trust1Paid": "ដឹកជញ្ជូនរហ័សក្នុងក្រុងព្រះសីហនុ", "deliveryNotePaid": "ដឹកដោយអ្នកដឹក Phsar ក្នុងក្រុងព្រះសីហនុ។ គិតថ្លៃតាមហាង។",
			"shipFreeOver": "ឥតគិតថ្លៃពី {a}", "shipFreeOverStore": "ឥតគិតថ្លៃ ពេលទិញចាប់ពី {a} ក្នុងហាងតែមួយ។"},
		"zh": {"trust1Paid": "西哈努克市内快速配送", "deliveryNotePaid": "由 Phsar 快递员在西哈努克市内配送，按店铺计费。",
			"shipFreeOver": "满 {a} 免运费", "shipFreeOverStore": "单店消费满 {a} 免运费。"},
	}
	shop := map[string]map[string]string{
		"en": {"shipEstimate": "(estimate)"}, "km": {"shipEstimate": "(ប៉ាន់ស្មាន)"}, "zh": {"shipEstimate": "（预估）"},
	}
	for l, m := range page {
		for k, v := range m {
			strs[l][k] = v
		}
	}
	for l, m := range shop {
		for k, v := range m {
			shopStrs[l][k] = v
		}
	}
}
