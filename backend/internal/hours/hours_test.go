package hours

import (
	"testing"
	"time"
)

func at(d, hm string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04", d+" "+hm, Loc)
	return t
}

func TestAdd(t *testing.T) {
	w := everyDay("08:00", "22:00")
	cases := []struct {
		name  string
		start time.Time
		mins  int
		want  time.Time
	}{
		{"di jam buka", at("2026-11-16", "10:00"), 120, at("2026-11-16", "12:00")},
		{"melewati jam tutup", at("2026-11-16", "21:30"), 120, at("2026-11-17", "09:30")},
		{"dibayar saat tutup (malam)", at("2026-11-16", "23:10"), 120, at("2026-11-17", "10:00")},
		{"dibayar saat tutup (subuh)", at("2026-11-17", "02:00"), 60, at("2026-11-17", "09:00")},
	}
	for _, c := range cases {
		if got := w.Add(c.start, c.mins); !got.Equal(c.want) {
			t.Errorf("%s: got %s want %s", c.name, got.In(Loc).Format("01-02 15:04"), c.want.Format("01-02 15:04"))
		}
	}
	// Lewat tengah malam (Tuktuk: 08:00–03:00).
	n := everyDay("08:00", "03:00")
	if got := n.Add(at("2026-11-17", "02:30"), 60); !got.Equal(at("2026-11-17", "08:30")) {
		t.Errorf("overnight: %s", got.In(Loc).Format("01-02 15:04"))
	}
	if got := n.Add(at("2026-11-16", "23:00"), 120); !got.Equal(at("2026-11-17", "01:00")) {
		t.Errorf("overnight2: %s", got.In(Loc).Format("01-02 15:04"))
	}
	// Tutup hari Minggu (16 Nov 2026 = Senin; 22 Nov = Minggu).
	s := everyDay("08:00", "22:00")
	s[0] = nil
	if got := s.Add(at("2026-11-21", "21:00"), 120); !got.Equal(at("2026-11-23", "09:00")) {
		t.Errorf("sunday closed: %s", got.In(Loc).Format("01-02 15:04"))
	}
	// 24 jam.
	f := everyDay("00:00", "24:00")
	if got := f.Add(at("2026-11-16", "23:30"), 120); !got.Equal(at("2026-11-17", "01:30")) {
		t.Errorf("24h: %s", got.In(Loc).Format("01-02 15:04"))
	}
	if open, next := w.OpenAt(at("2026-11-16", "23:00")); open || next == nil || !next.Equal(at("2026-11-17", "08:00")) {
		t.Errorf("OpenAt: %v %v", open, next)
	}
	var none Week
	if none.Validate("h") == nil {
		t.Error("minggu kosong harus ditolak")
	}
}
