package storefront

// pricehistory: transparansi flash sale — "Harga terendah 30 hari: $X" di halaman produk saat produk sedang flash.
// X = harga normal terendah 30 hari sebelum sesi mulai (acuan aturan harga flash; fungsi SQL variant_low_price).

import "context"

var phStrs = map[string]map[string]string{
	"en": {"phLow30": "Lowest price in the 30 days before this sale: ${a}"},
	"km": {"phLow30": "តម្លៃទាបបំផុតក្នុង 30 ថ្ងៃមុនការលក់នេះ៖ ${a}"},
	"zh": {"phLow30": "本次抢购前 30 天最低价：${a}"},
}

func init() {
	for l, m := range phStrs {
		for k, v := range m {
			strs[l][k] = v
		}
	}
}

// phLow30: isi p.Low30 untuk varian yang sedang memakai harga flash (yang terendah di antaranya).
func (c *Catalog) phLow30(ctx context.Context, p *Product) {
	ids := []string{}
	for _, v := range p.Variants {
		if v.Flash {
			ids = append(ids, v.ID)
		}
	}
	if len(ids) == 0 {
		return
	}
	var low *string
	if err := c.Pool.QueryRow(ctx, `SELECT min(variant_low_price(f.variant_id, f.starts_at))::text FROM flash_live f
		WHERE f.variant_id = ANY($1::uuid[]) AND f.sold < f.qty_cap`, ids).Scan(&low); err == nil && low != nil {
		p.Low30 = *low
	}
}
