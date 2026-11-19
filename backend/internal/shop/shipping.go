package shop

// shipping: ongkir per toko dari zona wilayah (paket shipping, diatur admin → Shipping & FX rates).
// Dipanggil setelah quote(): menambah ongkir ke tiap toko dan ke total. area "" = keranjang tanpa alamat → perkiraan.

import (
	"context"

	"github.com/shopspring/decimal"

	"phsar/internal/shipping"
)

func (s *Service) applyShipping(ctx context.Context, q shipping.Querier, qt *Quote, area string) error {
	z, err := shipping.ZoneFor(ctx, q, area)
	if err != nil {
		return err
	}
	total := decimal.Zero
	for i := range qt.Stores {
		g := &qt.Stores[i]
		g.shipping = z.FeeFor(g.subtotal)
		g.ShippingFee = g.shipping.StringFixed(2)
		total = total.Add(g.shipping)
	}
	qt.shipping = total
	qt.ShippingFee = total.StringFixed(2)
	qt.total = qt.total.Add(total)
	qt.Total = qt.total.StringFixed(2)
	qt.ShippingEstimate = area == "" && total.IsPositive()
	return nil
}
