// Package money menangani nominal USD dan tampilan riel.
//
// Aturan Phsar:
//   - Semua nominal disimpan NUMERIC(12,2) USD dan diproses dengan decimal, tidak pernah float64.
//   - Pembulatan biasa (half-up): 12.5960 → 12.60, 12.594 → 12.59.
//   - Riel hanya untuk tampilan, dibulatkan ke ៛100 terdekat.
package money

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

var hundred = decimal.NewFromInt(100)

// Round2 membulatkan ke 2 desimal. decimal.Round memakai "half away from zero",
// yang untuk nominal positif sama dengan pembulatan biasa (≥ 5 naik).
func Round2(d decimal.Decimal) decimal.Decimal { return d.Round(2) }

// Percent menghitung pct% dari amount, dibulatkan ke sen. Contoh komisi 8% dari $12.99 = $1.04.
func Percent(amount, pct decimal.Decimal) decimal.Decimal {
	return Round2(amount.Mul(pct).Div(hundred))
}

// SplitCommission mengembalikan (komisi, bagian seller). Bagian seller = total − komisi,
// supaya jumlahnya selalu tepat sama dengan total.
func SplitCommission(total, pct decimal.Decimal) (commission, seller decimal.Decimal) {
	commission = Percent(total, pct)
	return commission, total.Sub(commission)
}

// ToRiel mengubah USD ke riel untuk tampilan, dibulatkan ke ៛100 terdekat.
func ToRiel(usd, rielPerUSD decimal.Decimal) int64 {
	return usd.Mul(rielPerUSD).Div(hundred).Round(0).Mul(hundred).IntPart()
}

// FormatUSD → "$1,234.50".
func FormatUSD(d decimal.Decimal) string {
	s := Round2(d).StringFixed(2)
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	intPart, frac, _ := strings.Cut(s, ".")
	out := groupThousands(intPart) + "." + frac
	if neg {
		return "−$" + out
	}
	return "$" + out
}

// FormatRiel → "128,000 ៛" (format yang dipakai website Phsar).
func FormatRiel(r int64) string {
	return fmt.Sprintf("%s ៛", groupThousands(fmt.Sprint(r)))
}

func groupThousands(s string) string {
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	pre := len(s) % 3
	if pre > 0 {
		b.WriteString(s[:pre])
	}
	for i := pre; i < len(s); i += 3 {
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		b.WriteString(s[i : i+3])
	}
	return b.String()
}
