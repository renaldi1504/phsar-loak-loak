// Package hours: jam buka toko + aturan batas waktu pesanan (diatur super admin).
//
// Batas waktu dihitung dalam "menit buka": kalau toko tutup, hitungan berhenti dan lanjut saat toko buka lagi.
// Contoh: batas 120 menit, toko buka 08:00–22:00, pesanan dibayar 21:30 → 30 menit terpakai hari itu,
// sisa 90 menit mulai 08:00 besok → batas 09:30.
package hours

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// Loc: waktu Kamboja (UTC+7, tanpa musim panas) — tidak bergantung tzdata di server.
var Loc = time.FixedZone("ICT", 7*3600)

// Day: jam buka satu hari, "HH:MM". Close ≤ Open berarti lewat tengah malam (mis. 08:00–03:00).
// 00:00–24:00 = buka 24 jam.
type Day struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

// Week: indeks 0 = Minggu … 6 = Sabtu; nil = tutup.
type Week [7]*Day

func parseHM(s string, allow24 bool) (int, bool) {
	h, m, ok := strings.Cut(strings.TrimSpace(s), ":")
	if !ok || len(h) == 0 || len(h) > 2 || len(m) != 2 {
		return 0, false
	}
	hh, err1 := strconv.Atoi(h)
	mm, err2 := strconv.Atoi(m)
	if err1 != nil || err2 != nil || mm < 0 || mm > 59 || hh < 0 {
		return 0, false
	}
	v := hh*60 + mm
	if v > 24*60 || (v == 24*60 && !allow24) || (hh == 24 && mm != 0) {
		return 0, false
	}
	return v, true
}

// FieldError: isian yang salah (field mengikuti JSON, mis. "hours.1.close").
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

// Validate: format jam benar dan minimal satu hari buka (kalau tidak, batas waktu tidak pernah habis).
func (w Week) Validate(field string) error {
	open := false
	for i, d := range w {
		if d == nil {
			continue
		}
		o, ok := parseHM(d.Open, false)
		if !ok {
			return FieldError{fmt.Sprintf("%s.%d.open", field, i), "use HH:MM, e.g. 08:00"}
		}
		c, ok := parseHM(d.Close, true)
		if !ok {
			return FieldError{fmt.Sprintf("%s.%d.close", field, i), "use HH:MM, e.g. 22:00"}
		}
		if o == c {
			return FieldError{fmt.Sprintf("%s.%d.close", field, i), "closing time must differ from opening time (use 00:00–24:00 for 24 hours)"}
		}
		d.Open, d.Close = fmt.Sprintf("%02d:%02d", o/60, o%60), fmt.Sprintf("%02d:%02d", c/60, c%60)
		open = true
	}
	if !open {
		return FieldError{field, "open at least one day a week"}
	}
	return nil
}

type span struct{ from, to time.Time }

// spans: rentang buka dari (hari sebelum `from`) sampai `days` hari ke depan, terurut & digabung.
func (w Week) spans(from time.Time, days int) []span {
	f := from.In(Loc)
	day0 := time.Date(f.Year(), f.Month(), f.Day(), 0, 0, 0, 0, Loc)
	var out []span
	for i := -1; i <= days; i++ {
		date := day0.AddDate(0, 0, i)
		d := w[int(date.Weekday())]
		if d == nil {
			continue
		}
		o, ok1 := parseHM(d.Open, false)
		c, ok2 := parseHM(d.Close, true)
		if !ok1 || !ok2 || o == c {
			continue
		}
		if c < o {
			c += 24 * 60 // lewat tengah malam
		}
		out = append(out, span{date.Add(time.Duration(o) * time.Minute), date.Add(time.Duration(c) * time.Minute)})
	}
	sort.Slice(out, func(a, b int) bool { return out[a].from.Before(out[b].from) })
	merged := out[:0]
	for _, s := range out {
		if n := len(merged); n > 0 && !s.from.After(merged[n-1].to) {
			if s.to.After(merged[n-1].to) {
				merged[n-1].to = s.to
			}
			continue
		}
		merged = append(merged, s)
	}
	return merged
}

// Add: waktu setelah `mins` menit buka terhitung sejak `start`.
func (w Week) Add(start time.Time, mins int) time.Time {
	left := time.Duration(mins) * time.Minute
	for _, s := range w.spans(start, 30) {
		if !s.to.After(start) {
			continue
		}
		from := s.from
		if start.After(from) {
			from = start
		}
		if avail := s.to.Sub(from); avail >= left {
			return from.Add(left)
		} else {
			left -= avail
		}
	}
	return start.Add(time.Duration(mins) * time.Minute) // jadwal tidak masuk akal → hitung biasa
}

// OpenAt: apakah buka pada t; kalau tutup, kapan buka lagi.
func (w Week) OpenAt(t time.Time) (bool, *time.Time) {
	for _, s := range w.spans(t, 8) {
		if !t.Before(s.from) && t.Before(s.to) {
			return true, nil
		}
		if s.from.After(t) {
			n := s.from
			return false, &n
		}
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// Aturan batas waktu pesanan (platform_settings 'order_timing', diatur super admin)
// ---------------------------------------------------------------------------

type Timing struct {
	AcceptMinutes int  `json:"accept_minutes"` // seller harus menerima pesanan dibayar dalam … menit
	PackEnabled   bool `json:"pack_enabled"`   // batas kemas aktif?
	PackMinutes   int  `json:"pack_minutes"`   // setelah diterima, harus "siap dijemput" dalam … menit
	BusinessHours bool `json:"business_hours"` // hitung hanya saat toko buka
	DefaultHours  Week `json:"default_hours"`  // dipakai kalau toko belum mengisi jam buka
}

func everyDay(open, close string) Week {
	var w Week
	for i := range w {
		w[i] = &Day{open, close}
	}
	return w
}

func DefaultTiming() Timing {
	return Timing{AcceptMinutes: 120, PackEnabled: true, PackMinutes: 720, BusinessHours: true, DefaultHours: everyDay("08:00", "22:00")}
}

type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func LoadTiming(ctx context.Context, q Querier) (Timing, error) {
	t := DefaultTiming()
	var raw []byte
	err := q.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'order_timing'`).Scan(&raw)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return t, err
	}
	if len(raw) > 0 {
		var x Timing
		if json.Unmarshal(raw, &x) == nil {
			if x.DefaultHours.Validate("default_hours") != nil {
				x.DefaultHours = t.DefaultHours
			}
			if x.AcceptMinutes <= 0 {
				x.AcceptMinutes = t.AcceptMinutes
			}
			if x.PackMinutes <= 0 {
				x.PackMinutes = t.PackMinutes
			}
			t = x
		}
	}
	return t, nil
}

// Check: batas wajar (supaya salah ketik tidak membatalkan semua pesanan / membuat pesanan menggantung).
func (t *Timing) Check() error {
	switch {
	case t.AcceptMinutes < 15 || t.AcceptMinutes > 2880:
		return FieldError{"accept_minutes", "between 15 minutes and 48 hours"}
	case t.PackMinutes < 30 || t.PackMinutes > 10080:
		return FieldError{"pack_minutes", "between 30 minutes and 7 days"}
	}
	return t.DefaultHours.Validate("default_hours")
}

// StoreWeek: jam buka toko (nil kalau belum diisi).
func StoreWeek(ctx context.Context, q Querier, storeID string) (*Week, error) {
	var raw []byte
	if err := q.QueryRow(ctx, `SELECT opening_hours FROM stores WHERE id = $1`, storeID).Scan(&raw); err != nil {
		return nil, err
	}
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var w Week
	if err := json.Unmarshal(raw, &w); err != nil || w.Validate("opening_hours") != nil {
		return nil, nil
	}
	return &w, nil
}

// Effective: jam yang dipakai menghitung batas toko ini.
func (t Timing) Effective(w *Week) Week {
	if w != nil {
		return *w
	}
	return t.DefaultHours
}

// Deadline: from + mins (menit buka kalau BusinessHours aktif).
func (t Timing) Deadline(w *Week, from time.Time, mins int) time.Time {
	if !t.BusinessHours {
		return from.Add(time.Duration(mins) * time.Minute)
	}
	return t.Effective(w).Add(from, mins)
}

// StoreDeadline: baca jam buka toko lalu hitung batas.
func StoreDeadline(ctx context.Context, q Querier, t Timing, storeID string, from time.Time, mins int) (time.Time, error) {
	w, err := StoreWeek(ctx, q, storeID)
	if err != nil {
		return time.Time{}, err
	}
	return t.Deadline(w, from, mins), nil
}
