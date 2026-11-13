package geo

import "testing"

// Kotak kasar di sekitar kota Sihanoukville (sama dengan seed migrasi).
const box = `{"type":"Polygon","coordinates":[[[103.44,10.55],[103.62,10.55],[103.62,10.72],[103.44,10.72],[103.44,10.55]]]}`

func TestContains(t *testing.T) {
	p, err := ParsePolygon([]byte(box))
	if err != nil {
		t.Fatal(err)
	}
	if !p.Contains(10.627, 103.522) { // pusat kota
		t.Fatal("pusat Sihanoukville seharusnya di dalam")
	}
	if p.Contains(11.556, 104.928) { // Phnom Penh
		t.Fatal("Phnom Penh seharusnya di luar")
	}
}
