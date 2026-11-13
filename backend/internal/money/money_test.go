package money

import (
	"testing"

	"github.com/shopspring/decimal"
)

func d(s string) decimal.Decimal { return decimal.RequireFromString(s) }

func TestRound2HalfUp(t *testing.T) {
	cases := map[string]string{
		"12.5960": "12.60",
		"12.594":  "12.59",
		"12.595":  "12.60",
		"0.005":   "0.01",
		"7.00":    "7.00",
	}
	for in, want := range cases {
		if got := Round2(d(in)).StringFixed(2); got != want {
			t.Errorf("Round2(%s) = %s, mau %s", in, got, want)
		}
	}
}

func TestSplitCommissionAlwaysSumsToTotal(t *testing.T) {
	for _, total := range []string{"12.99", "32.00", "0.01", "189.00", "53.37"} {
		c, s := SplitCommission(d(total), d("8"))
		if !c.Add(s).Equal(d(total)) {
			t.Fatalf("komisi %s + seller %s ≠ total %s", c, s, total)
		}
	}
	c, s := SplitCommission(d("12.99"), d("8"))
	if c.StringFixed(2) != "1.04" || s.StringFixed(2) != "11.95" {
		t.Fatalf("12.99 @8%% = %s / %s, mau 1.04 / 11.95", c, s)
	}
}

func TestRielDisplay(t *testing.T) {
	if got := ToRiel(d("32.00"), d("4000")); got != 128000 {
		t.Fatalf("$32 → %d", got)
	}
	if got := ToRiel(d("12.99"), d("4100")); got != 53300 { // 53,259 → 53,300
		t.Fatalf("$12.99 @4100 → %d", got)
	}
	if got := FormatRiel(212000); got != "212,000 ៛" {
		t.Fatalf("FormatRiel = %q", got)
	}
	if got := FormatUSD(d("1234.5")); got != "$1,234.50" {
		t.Fatalf("FormatUSD = %q", got)
	}
}
