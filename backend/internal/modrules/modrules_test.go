package modrules

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func d(s string) decimal.Decimal { return decimal.RequireFromString(s) }
func dp(s string) *decimal.Decimal {
	v := d(s)
	return &v
}

func reasons(hits []PriceHit) []string {
	out := []string{}
	for _, h := range hits {
		out = append(out, h.Reason)
	}
	return out
}

func TestCheckPrice(t *testing.T) {
	floors := []Floor{{"iphone 15", d("250")}, {"ipad", d("80")}, {"ipad pro", d("200")}, {"rolex", d("1000")}}
	r := Defaults()
	cases := []struct {
		name string
		in   PriceCheck
		want []string
	}{
		{"iphone murah", PriceCheck{Name: "iPhone 15 Pro Max 256GB", Prices: []decimal.Decimal{d("120")}}, []string{ReasonPriceFloor}},
		{"iphone wajar", PriceCheck{Name: "iPhone 15 128GB", Prices: []decimal.Decimal{d("480")}}, []string{}},
		{"case iphone", PriceCheck{Name: "Case iPhone 15 silicone", Prices: []decimal.Decimal{d("3")}}, []string{}},
		{"iphone 150 bukan model", PriceCheck{Name: "iphone 150 sticker pack", Prices: []decimal.Decimal{d("2")}}, []string{}},
		{"ipad pro pakai floor terpanjang", PriceCheck{Name: "iPad Pro 11", Prices: []decimal.Decimal{d("150")}}, []string{ReasonPriceFloor}},
		{"ipad biasa", PriceCheck{Name: "iPad 9th gen", Prices: []decimal.Decimal{d("150")}}, []string{}},
		{"diskon 80%", PriceCheck{Name: "Kaos polos", Prices: []decimal.Decimal{d("2")}, Compares: []*decimal.Decimal{dp("10")}}, []string{ReasonBigDiscount}},
		{"diskon 50%", PriceCheck{Name: "Kaos polos", Prices: []decimal.Decimal{d("5")}, Compares: []*decimal.Decimal{dp("10")}}, []string{}},
		{"harga anjlok", PriceCheck{Name: "Blender", Prices: []decimal.Decimal{d("9")}, ApprovedMin: dp("20")}, []string{ReasonPriceDrop}},
		{"turun sedikit", PriceCheck{Name: "Blender", Prices: []decimal.Decimal{d("15")}, ApprovedMin: dp("20")}, []string{}},
		{"varian termurah dipakai", PriceCheck{Name: "Rolex Submariner", Prices: []decimal.Decimal{d("5000"), d("300")}}, []string{ReasonPriceFloor}},
	}
	for _, c := range cases {
		got := reasons(CheckPrice(c.in, floors, r))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
}

func TestDecide(t *testing.T) {
	r := Defaults()
	cases := []struct {
		name string
		s    Situation
		want []string
	}{
		{"toko lulus, bersih", Situation{}, nil},
		{"masa percobaan", Situation{ProbationLeft: true}, []string{ReasonProbation}},
		{"kategori terbatas + kata", Situation{Restricted: true, KeywordReview: true}, []string{ReasonRestricted, ReasonKeyword}},
		{"strike", Situation{StrikePoints: 1}, []string{ReasonStrikes}},
		{"ditahan laporan, seller edit", Situation{PrevModeration: "pending_review", Reported: true, ContentChanged: true}, []string{ReasonReported}},
		{"kirim ulang setelah diminta perbaiki", Situation{PrevModeration: "changes_requested"}, []string{ReasonChanges}},
		{"tayang, edit kecil, toko terbatas", Situation{PrevModeration: "approved", ContentChanged: true, Restricted: true, ProbationLeft: true}, nil},
		{"tayang, edit kecil, kata baru", Situation{PrevModeration: "approved", ContentChanged: true, KeywordReview: true}, []string{ReasonKeyword}},
		{"tayang, harga saja turun", Situation{PrevModeration: "approved", PriceHits: []PriceHit{{Reason: ReasonPriceDrop}}}, []string{ReasonPriceDrop}},
		{"tayang, edit besar, kategori terbatas", Situation{PrevModeration: "approved", MajorEdit: true, ContentChanged: true, Restricted: true}, []string{ReasonRestricted}},
	}
	for _, c := range cases {
		if got := Decide(c.s, r); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
	off := r
	off.AutoPublish = false
	if got := Decide(Situation{}, off); !reflect.DeepEqual(got, []string{ReasonManualOnly}) {
		t.Errorf("auto publish off: %v", got)
	}
}

func TestContainsWord(t *testing.T) {
	if ContainsWord("burgundy dress", "gun") || !ContainsWord("toy gun", "gun") || !ContainsWord("ស្រោមទូរសព្ទ", "ស្រោម") {
		t.Fatal("ContainsWord")
	}
}
