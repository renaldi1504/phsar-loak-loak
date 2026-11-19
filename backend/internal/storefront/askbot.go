package storefront

// askbot: blok "Tanya jawab" di halaman produk — FAQ yang ditulis seller untuk produk ini.

import "context"

// QA: satu tanya-jawab. Lang = bahasa yang ditulis seller (atribut lang di HTML).
type QA struct{ Q, A, Lang string }

// FAQSource: FAQ per produk (paket askbot). nil = blok tidak tampil.
type FAQSource interface {
	ProductFAQ(ctx context.Context, lang, productID string) ([]QA, error)
}

// productFAQ: kalau gagal, halaman produk tetap tampil tanpa blok tanya jawab.
func (h *Handler) productFAQ(ctx context.Context, lang, productID string) []QA {
	if h.FAQ == nil {
		return nil
	}
	list, err := h.FAQ.ProductFAQ(ctx, lang, productID)
	if err != nil {
		if h.Log != nil {
			h.Log.Warn("askbot: FAQ produk gagal dimuat", "err", err)
		}
		return nil
	}
	return list
}

func init() {
	page := map[string]map[string]string{
		"en": {"faqTitle": "Questions & answers", "faqBy": "Answered by the seller"},
		"km": {"faqTitle": "សំណួរ និងចម្លើយ", "faqBy": "ឆ្លើយដោយអ្នកលក់"},
		"zh": {"faqTitle": "问答", "faqBy": "卖家解答"},
	}
	for l, m := range page {
		for k, v := range m {
			strs[l][k] = v
		}
	}
}
