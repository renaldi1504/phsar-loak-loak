package paygw

import "testing"

func TestAllowedBaseURL(t *testing.T) {
	for raw, want := range map[string]bool{
		"https://checkout.payway.com.kh":          true,
		"https://checkout-sandbox.payway.com.kh/": true,
		"https://payway.com.kh":                   true,
		"http://checkout.payway.com.kh":           false,
		"https://payway.com.kh.evil.com":          false,
		"https://evilpayway.com.kh":               false,
		"https://a@checkout.payway.com.kh":        false,
		"https://checkout.payway.com.kh/x":        false,
		"https://checkout.payway.com.kh?x=1":      false,
	} {
		if got := allowedBaseURL(raw); got != want {
			t.Errorf("%s: %v, want %v", raw, got, want)
		}
	}
}

func TestHintAndExplain(t *testing.T) {
	if hintOf("abcdefgh1234") != "1234" || hintOf("abc") != "••••" {
		t.Fatal("hint")
	}
	if m := explain("payway check-transaction: kode 1: wrong hash"); m[:9] != "Wrong API" {
		t.Fatalf("explain: %s", m)
	}
}
