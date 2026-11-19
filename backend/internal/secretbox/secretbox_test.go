package secretbox

import "testing"

func TestSealOpen(t *testing.T) {
	b, err := New("ZGV2LW9ubHktZGF0YS1rZXktMzItYnl0ZXMtLS0tLS0=")
	if err != nil {
		t.Fatal(err)
	}
	s, err := b.Seal([]byte("000123456"), []byte("store-1"))
	if err != nil {
		t.Fatal(err)
	}
	if got, err := b.Open(s, []byte("store-1")); err != nil || string(got) != "000123456" {
		t.Fatalf("open: %q %v", got, err)
	}
	if _, err := b.Open(s, []byte("store-2")); err == nil {
		t.Fatal("data yang dipindah ke toko lain harus gagal dibuka")
	}
	if _, err := New("short"); err == nil {
		t.Fatal("kunci pendek harus ditolak")
	}
}
