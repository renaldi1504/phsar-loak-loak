package seller

import "testing"

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Urban Finds":          "urban-finds",
		"  Grace Beauty Shop!": "grace-beauty-shop",
		"TukTuk & Co. 2026":    "tuktuk-co-2026",
		"ហាងស្អាត":             "",
	}
	for in, want := range cases {
		if got := Slugify(in); got != want {
			t.Errorf("Slugify(%q) = %q, mau %q", in, got, want)
		}
	}
}
