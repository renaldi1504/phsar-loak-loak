package storefront

// promos: harga flash sale di kartu produk & halaman produk (harga normal jadi harga coret + hitung mundur).
// Harga yang benar-benar dibayar tetap dihitung server saat keranjang/checkout (paket shop).

import (
	"context"
	"strconv"
	"time"
)

type flashRow struct {
	variant string
	price   string
	ends    time.Time
}

// activeFlash: flash sale yang berlaku & masih ada kuota, per varian (harga terendah).
func (c *Catalog) activeFlash(ctx context.Context, where string, arg any) (map[string]flashRow, map[string]string, error) {
	rows, err := c.Pool.Query(ctx, `SELECT DISTINCT ON (f.variant_id) f.variant_id, v.product_id::text, f.sale_price::text, f.ends_at
		FROM flash_live f JOIN product_variants v ON v.id = f.variant_id
		WHERE `+where+` AND v.deleted_at IS NULL AND v.status = 'active' AND f.sold < f.qty_cap
		ORDER BY f.variant_id, f.sale_price`, arg)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	out, product := map[string]flashRow{}, map[string]string{}
	for rows.Next() {
		var f flashRow
		var pid string
		if err := rows.Scan(&f.variant, &pid, &f.price, &f.ends); err != nil {
			return nil, nil, err
		}
		out[f.variant] = f
		product[f.variant] = pid
	}
	return out, product, rows.Err()
}

func fmt2(v float64) string { return strconv.FormatFloat(v, 'f', 2, 64) }

// flashCards: harga termurah kartu ikut flash sale (kalau lebih murah).
func (c *Catalog) flashCards(ctx context.Context, cards []Card) {
	if len(cards) == 0 {
		return
	}
	ids := make([]string, len(cards))
	for i, k := range cards {
		ids[i] = k.ID
	}
	fl, prod, err := c.activeFlash(ctx, `v.product_id = ANY($1::uuid[])`, ids)
	if err != nil || len(fl) == 0 {
		return
	}
	idx := map[string]int{}
	for i, k := range cards {
		idx[k.ID] = i
	}
	for vid, f := range fl {
		i, ok := idx[prod[vid]]
		if !ok {
			continue
		}
		k := &cards[i]
		if num(f.price) < num(k.MinPrice) {
			if k.Compare == "" || num(k.Compare) < num(k.MinPrice) {
				k.Compare = k.MinPrice
			}
			k.MinPrice = fmt2(num(f.price))
			if num(k.MaxPrice) < num(k.MinPrice) {
				k.MaxPrice = k.MinPrice
			}
		}
		if k.FlashEnds == nil || f.ends.Before(*k.FlashEnds) {
			e := f.ends
			k.FlashEnds = &e
		}
	}
}

// flashProduct: harga varian ikut flash sale; harga normal jadi harga coret.
func (c *Catalog) flashProduct(ctx context.Context, p *Product) {
	if p == nil || len(p.Variants) == 0 {
		return
	}
	fl, _, err := c.activeFlash(ctx, `v.product_id = $1::uuid`, p.ID)
	if err != nil || len(fl) == 0 {
		return
	}
	min := -1.0
	for i := range p.Variants {
		v := &p.Variants[i]
		if f, ok := fl[v.ID]; ok && num(f.price) < num(v.Price) {
			if v.Compare == "" || num(v.Compare) < num(v.Price) {
				v.Compare = v.Price
			}
			v.Price = fmt2(num(f.price))
			v.Flash = true
			if p.FlashEnds == nil || f.ends.Before(*p.FlashEnds) {
				e := f.ends
				p.FlashEnds = &e
			}
			if num(v.Compare) > num(p.Compare) {
				p.Compare = v.Compare
			}
		}
		if min < 0 || num(v.Price) < min {
			min = num(v.Price)
		}
	}
	if min >= 0 && min < num(p.MinPrice) {
		p.MinPrice = fmt2(min)
	}
	c.phLow30(ctx, p) // pricehistory
}
