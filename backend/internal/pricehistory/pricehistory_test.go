package pricehistory

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"

	"phsar/internal/i18n"
)

func TestMsg(t *testing.T) {
	low := decimal.RequireFromString("18")
	for l, want := range map[i18n.Lang]string{i18n.EN: "lowest price of the last 30 days ($18.00) — at most $17.99", i18n.KM: "30 ថ្ងៃ", i18n.ZH: "近 30 天最低价（$18.00）"} {
		if got := Msg(l, "Kettle", low); !strings.Contains(got, want) || !strings.HasPrefix(got, "Kettle") {
			t.Errorf("%s: %q", l, got)
		}
	}
}
