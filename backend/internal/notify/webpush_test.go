package notify

import (
	"crypto/ecdh"
	"testing"
)

// Vektor uji RFC 8291 Lampiran A.
func TestEncryptRFC8291(t *testing.T) {
	asD, _ := decodeB64("yfWPiYE-n46HLnH0KqZOF1fJJU3MYrct3AELtAQ-oRw")
	as, err := ecdh.P256().NewPrivateKey(asD)
	if err != nil {
		t.Fatal(err)
	}
	salt, _ := decodeB64("DGv6ra1nlYgDCS1FRnbzlw")
	sub := Subscription{P256dh: "BCVxsr7N_eNgVRqvHtD0zTZsEc6-VV-JvLexhqUzORcxaOzi6-AYWXvTBHm4bjyPjs7Vd8pZGH6SRpkNtoIAiw4", Auth: "BTBZMqHH6r4Tts7J_aSIgg"}
	got, err := Encrypt([]byte("When I grow up, I want to be a watermelon"), sub, salt, as)
	if err != nil {
		t.Fatal(err)
	}
	want := "DGv6ra1nlYgDCS1FRnbzlwAAEABBBP4z9KsN6nGRTbVYI_c7VJSPQTBtkgcy27mlmlMoZIIgDll6e3vCYLocInmYWAmS6TlzAC8wEqKK6PBru3jl7A_yl95bQpu6cVPTpK4Mqgkf1CXztLVBSt2Ks3oZwbuwXPXLWyouBWLVWGNWQexSgSxsj_Qulcy4a-fN"
	if b64.EncodeToString(got) != want {
		t.Fatalf("hasil enkripsi tidak sama dengan RFC 8291:\n got %s\nwant %s", b64.EncodeToString(got), want)
	}
}
