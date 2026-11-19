// Package shipping: ongkir per zona dan kurs riel yang diatur admin (menu Shipping & FX rates).
//
//	wilayah layanan (service_areas, poligon di peta) ──► zona ongkir (shipping_zones)
//	ongkir dibayar per toko dalam satu checkout; subtotal toko ≥ free_over → gratis.
//	Default saat launch: Sihanoukville, ongkir 0 (gratis).
//
// Kurs USD→KHR: admin memasukkan kurs baru; riwayat disimpan, yang terbaru dipakai di semua tampilan riel.
package shipping

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"phsar/internal/geo"
	"phsar/internal/platform/db"
)

// FieldError: input tidak valid (422).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

var (
	ErrNotFound = errors.New("tidak ditemukan")
	ErrLastArea = errors.New("minimal satu wilayah layanan harus aktif")
	ErrExists   = errors.New("kode sudah dipakai")
)

// Querier: pool atau transaksi.
type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Zone: aturan ongkir satu zona.
type Zone struct {
	Code        string            `json:"code"`
	Name        map[string]string `json:"name"`
	FeePerStore string            `json:"fee_per_store"`
	FreeOver    *string           `json:"free_over"`
	CourierFee  string            `json:"courier_fee"`
	Areas       int               `json:"areas"`
	UpdatedAt   time.Time         `json:"updated_at"`

	fee, freeOver decimal.Decimal
	hasFree       bool
}

// FeeFor: ongkir untuk satu toko dengan subtotal ini.
func (z Zone) FeeFor(subtotal decimal.Decimal) decimal.Decimal {
	if !subtotal.IsPositive() || !z.fee.IsPositive() {
		return decimal.Zero
	}
	if z.hasFree && subtotal.GreaterThanOrEqual(z.freeOver) {
		return decimal.Zero
	}
	return z.fee
}

// Free: tanpa ongkir sama sekali (tampilan "Gratis").
func (z Zone) Free() bool { return !z.fee.IsPositive() }

func (z *Zone) parse() {
	z.fee, _ = decimal.NewFromString(z.FeePerStore)
	if z.FreeOver != nil {
		z.freeOver, _ = decimal.NewFromString(*z.FreeOver)
		z.hasFree = true
	}
}

// ZoneFor: zona untuk kode wilayah. area "" (keranjang belum ada alamat) → zona termurah di antara wilayah aktif (perkiraan).
func ZoneFor(ctx context.Context, q Querier, area string) (Zone, error) {
	var z Zone
	var name []byte
	sql := `SELECT z.code, z.name, z.fee_per_store::text, z.free_over::text, z.courier_fee::text FROM shipping_zones z `
	var err error
	if area != "" {
		err = q.QueryRow(ctx, sql+`JOIN service_areas a ON a.shipping_zone_code = z.code WHERE a.code = $1`, area).
			Scan(&z.Code, &name, &z.FeePerStore, &z.FreeOver, &z.CourierFee)
	} else {
		err = q.QueryRow(ctx, sql+`WHERE EXISTS (SELECT 1 FROM service_areas a WHERE a.shipping_zone_code = z.code AND a.is_active)
			ORDER BY z.fee_per_store, z.code LIMIT 1`).Scan(&z.Code, &name, &z.FeePerStore, &z.FreeOver, &z.CourierFee)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return Zone{FeePerStore: "0.00", CourierFee: "0.00"}, nil // tanpa zona: gratis (perilaku lama)
	}
	if err != nil {
		return z, err
	}
	_ = json.Unmarshal(name, &z.Name)
	z.parse()
	return z, nil
}

// Info: ringkasan untuk tampilan website (halaman produk, banner) — zona default.
type Info struct {
	Free     bool
	Fee      string
	FreeOver string
}

func GetInfo(ctx context.Context, q Querier) Info {
	z, err := ZoneFor(ctx, q, "")
	if err != nil {
		return Info{Free: true, Fee: "0.00"}
	}
	in := Info{Free: z.Free(), Fee: z.fee.StringFixed(2)}
	if z.hasFree && !z.Free() {
		in.FreeOver = z.freeOver.StringFixed(2)
	}
	return in
}

// ---------------------------------------------------------------------------
// Admin
// ---------------------------------------------------------------------------

type Store struct{ Pool *pgxpool.Pool }

type Area struct {
	Code      string            `json:"code"`
	Name      map[string]string `json:"name"`
	Boundary  json.RawMessage   `json:"boundary"`
	Zone      string            `json:"zone"`
	Active    bool              `json:"active"`
	Stores    int               `json:"stores"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Rate struct {
	ID         int       `json:"id"`
	RateDate   string    `json:"rate_date"`
	RielPerUSD string    `json:"riel_per_usd"`
	Source     string    `json:"source"`
	Note       *string   `json:"note"`
	By         *string   `json:"by"`
	At         time.Time `json:"at"`
}

type View struct {
	Zones   []Zone `json:"zones"`
	Areas   []Area `json:"areas"`
	Rates   []Rate `json:"rates"`
	Current string `json:"riel_per_usd"`
	Courier struct {
		Deliveries int    `json:"deliveries"`
		Cost       string `json:"cost"`
		FeesIn     string `json:"fees_collected"`
	} `json:"courier_30d"`
}

func (s *Store) View(ctx context.Context) (*View, error) {
	v := &View{Zones: []Zone{}, Areas: []Area{}, Rates: []Rate{}}
	rows, err := s.Pool.Query(ctx, `SELECT z.code, z.name, z.fee_per_store::text, z.free_over::text, z.courier_fee::text, z.updated_at,
		(SELECT count(*) FROM service_areas a WHERE a.shipping_zone_code = z.code) FROM shipping_zones z ORDER BY z.code`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var z Zone
		var name []byte
		if err := rows.Scan(&z.Code, &name, &z.FeePerStore, &z.FreeOver, &z.CourierFee, &z.UpdatedAt, &z.Areas); err != nil {
			rows.Close()
			return nil, err
		}
		_ = json.Unmarshal(name, &z.Name)
		v.Zones = append(v.Zones, z)
	}
	rows.Close()
	rows, err = s.Pool.Query(ctx, `SELECT a.code, a.name, a.boundary, a.shipping_zone_code, a.is_active, a.updated_at,
		(SELECT count(*) FROM stores st WHERE st.service_area_code = a.code AND st.deleted_at IS NULL) FROM service_areas a ORDER BY a.code`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a Area
		var name, b []byte
		if err := rows.Scan(&a.Code, &name, &b, &a.Zone, &a.Active, &a.UpdatedAt, &a.Stores); err != nil {
			rows.Close()
			return nil, err
		}
		_ = json.Unmarshal(name, &a.Name)
		a.Boundary = b
		v.Areas = append(v.Areas, a)
	}
	rows.Close()
	rows, err = s.Pool.Query(ctx, `SELECT r.id, r.rate_date::text, r.riel_per_usd::text, r.source, r.note, coalesce(nullif(u.name, ''), u.phone), r.fetched_at
		FROM exchange_rates r LEFT JOIN users u ON u.id = r.created_by ORDER BY r.rate_date DESC, r.id DESC LIMIT 30`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r Rate
		if err := rows.Scan(&r.ID, &r.RateDate, &r.RielPerUSD, &r.Source, &r.Note, &r.By, &r.At); err != nil {
			rows.Close()
			return nil, err
		}
		v.Rates = append(v.Rates, r)
	}
	rows.Close()
	v.Current = "4000"
	if len(v.Rates) > 0 {
		v.Current = v.Rates[0].RielPerUSD
	}
	// Biaya kurir 30 hari (perkiraan: antaran selesai × biaya kurir zona) vs ongkir yang dibayar pembeli.
	_ = s.Pool.QueryRow(ctx, `SELECT count(*), coalesce(sum(z.courier_fee), 0)::numeric(12,2)::text, coalesce(sum(o.shipping_fee), 0)::numeric(12,2)::text
		FROM deliveries d JOIN orders o ON o.id = d.order_id JOIN checkouts c ON c.id = o.checkout_id
		LEFT JOIN service_areas a ON a.code = c.service_area_code LEFT JOIN shipping_zones z ON z.code = a.shipping_zone_code
		WHERE d.status = 'delivered' AND d.delivered_at > now() - interval '30 days'`).Scan(&v.Courier.Deliveries, &v.Courier.Cost, &v.Courier.FeesIn)
	return v, nil
}

var codeRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

type ZoneInput struct {
	Code        string            `json:"code"`
	Name        map[string]string `json:"name"`
	FeePerStore string            `json:"fee_per_store"`
	FreeOver    string            `json:"free_over"` // kosong = tidak ada gratis ongkir
	CourierFee  string            `json:"courier_fee"`
}

func names(in map[string]string) (map[string]string, error) {
	out := map[string]string{}
	for _, l := range []string{"en", "km", "zh"} {
		v := strings.TrimSpace(in[l])
		if n := len([]rune(v)); n > 60 {
			return nil, FieldError{"name_" + l, "too long"}
		}
		if v != "" {
			out[l] = v
		}
	}
	if out["en"] == "" {
		return nil, FieldError{"name_en", "required"}
	}
	return out, nil
}

func money(field, v string, max int64, allowEmpty bool) (*decimal.Decimal, error) {
	v = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(v), "$"))
	if v == "" {
		if allowEmpty {
			return nil, nil
		}
		v = "0"
	}
	d, err := decimal.NewFromString(v)
	if err != nil || d.IsNegative() || d.GreaterThan(decimal.NewFromInt(max)) || d.Exponent() < -2 {
		return nil, FieldError{field, "enter an amount between 0 and " + decimal.NewFromInt(max).String()}
	}
	return &d, nil
}

// SaveZone membuat (create=true) atau mengubah zona.
func (s *Store) SaveZone(ctx context.Context, actorID string, in ZoneInput, create bool) error {
	in.Code = strings.TrimSpace(in.Code)
	if !codeRe.MatchString(in.Code) || len(in.Code) > 32 {
		return FieldError{"code", "lowercase letters, numbers and dashes"}
	}
	nm, err := names(in.Name)
	if err != nil {
		return err
	}
	fee, err := money("fee_per_store", in.FeePerStore, 50, false)
	if err != nil {
		return err
	}
	free, err := money("free_over", in.FreeOver, 100000, true)
	if err != nil {
		return err
	}
	if free != nil && !free.IsPositive() {
		free = nil
	}
	cf, err := money("courier_fee", in.CourierFee, 50, false)
	if err != nil {
		return err
	}
	name, _ := json.Marshal(nm)
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		if create {
			_, err := tx.Exec(ctx, `INSERT INTO shipping_zones (code, name, fee_per_store, free_over, courier_fee, is_estimate, updated_by)
				VALUES ($1, $2, $3, $4, $5, false, $6)`, in.Code, name, *fee, free, *cf, actorID)
			if err != nil && strings.Contains(err.Error(), "23505") {
				return ErrExists
			}
			return err
		}
		tag, err := tx.Exec(ctx, `UPDATE shipping_zones SET name = $2, fee_per_store = $3, free_over = $4, courier_fee = $5, updated_by = $6, updated_at = now()
			WHERE code = $1`, in.Code, name, *fee, free, *cf, actorID)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
}

type AreaInput struct {
	Code     string            `json:"code"`
	Name     map[string]string `json:"name"`
	Boundary json.RawMessage   `json:"boundary"`
	Zone     string            `json:"zone"`
	Active   bool              `json:"active"`
}

// SaveArea membuat atau mengubah wilayah layanan (poligon GeoJSON [lng, lat]).
func (s *Store) SaveArea(ctx context.Context, actorID string, in AreaInput, create bool) error {
	in.Code = strings.TrimSpace(in.Code)
	if !codeRe.MatchString(in.Code) || len(in.Code) > 32 {
		return FieldError{"code", "lowercase letters, numbers and dashes"}
	}
	nm, err := names(in.Name)
	if err != nil {
		return err
	}
	if len(in.Boundary) == 0 || len(in.Boundary) > 200<<10 {
		return FieldError{"boundary", "draw or paste the area polygon (GeoJSON)"}
	}
	if _, err := geo.ParsePolygon(in.Boundary); err != nil {
		return FieldError{"boundary", `GeoJSON {"type":"Polygon","coordinates":[[[lng,lat],…]]} with at least 3 points`}
	}
	name, _ := json.Marshal(nm)
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		var ok bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM shipping_zones WHERE code = $1)`, in.Zone).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return FieldError{"zone", "choose a delivery zone"}
		}
		if create {
			_, err := tx.Exec(ctx, `INSERT INTO service_areas (code, name, boundary, shipping_zone_code, is_active) VALUES ($1, $2, $3, $4, $5)`,
				in.Code, name, []byte(in.Boundary), in.Zone, in.Active)
			if err != nil && strings.Contains(err.Error(), "23505") {
				return ErrExists
			}
			return err
		}
		tag, err := tx.Exec(ctx, `UPDATE service_areas SET name = $2, boundary = $3, shipping_zone_code = $4, is_active = $5, updated_at = now() WHERE code = $1`,
			in.Code, name, []byte(in.Boundary), in.Zone, in.Active)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		var active int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM service_areas WHERE is_active`).Scan(&active); err != nil {
			return err
		}
		if active == 0 {
			return ErrLastArea // pembeli tidak akan bisa checkout sama sekali
		}
		return nil
	})
}

// AddRate menyimpan kurs baru (berlaku sekarang).
func (s *Store) AddRate(ctx context.Context, actorID, rate, note string) error {
	d, err := decimal.NewFromString(strings.ReplaceAll(strings.TrimSpace(rate), ",", ""))
	if err != nil || d.LessThan(decimal.NewFromInt(3000)) || d.GreaterThan(decimal.NewFromInt(6000)) || d.Exponent() < -4 {
		return FieldError{"riel_per_usd", "between 3,000 and 6,000 riel per dollar"}
	}
	note = strings.TrimSpace(note)
	if len([]rune(note)) > 200 {
		return FieldError{"note", "too long"}
	}
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO exchange_rates (rate_date, riel_per_usd, source, created_by, note)
			VALUES ((now() AT TIME ZONE 'Asia/Phnom_Penh')::date, $1, 'manual', $2, nullif($3, ''))`, d, actorID, note)
		return err
	})
}
