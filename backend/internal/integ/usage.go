package integ

// Pemakaian layanan pihak ketiga (Anthropic, Gemini, PlasGate) + pengingat saldo.
// Penyedia tidak memberi API saldo yang bisa dipakai di sini, jadi: super admin mengisi saldo yang terlihat di console
// penyedia, server mencatat setiap panggilan + perkiraan biaya, lalu menghitung perkiraan sisa dan memberi peringatan.
// Error "saldo habis" (mis. Anthropic HTTP 400 credit balance) langsung jadi peringatan merah.

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
)

// Usage: satu panggilan ke layanan luar.
type Usage struct {
	Provider string // anthropic · gemini · plasgate
	Purpose  string // translate · askbot · test · sms
	Model    string
	In, Out  int // token
	Units    int // SMS
	Err      error
}

// ClassifyErr: jenis kegagalan untuk peringatan.
func ClassifyErr(err error) string {
	if err == nil {
		return ""
	}
	m := strings.ToLower(err.Error())
	switch {
	case strings.Contains(m, "credit balance") || strings.Contains(m, "billing") || strings.Contains(m, "saldo") || strings.Contains(m, "insufficient"):
		return "billing"
	case strings.Contains(m, "http 429") || strings.Contains(m, "resource_exhausted") || strings.Contains(m, "quota") || strings.Contains(m, "rate limit"):
		return "quota"
	case strings.Contains(m, "http 401") || strings.Contains(m, "http 403") || strings.Contains(m, "api key") || strings.Contains(m, "permission"):
		return "auth"
	}
	return "other"
}

// Budget: harga & saldo per penyedia (diisi super admin). Angka sebagai teks desimal.
type Budget struct {
	BalanceUSD   string     `json:"balance_usd"`    // saldo saat diisi ("" = belum diisi)
	BalanceAt    *time.Time `json:"balance_at"`     // kapan saldo diisi; biaya sesudahnya dikurangkan
	WarnBelowUSD string     `json:"warn_below_usd"` // peringatan kalau perkiraan sisa di bawah ini
	PriceIn      string     `json:"price_in"`       // USD per 1 juta token input
	PriceOut     string     `json:"price_out"`      // USD per 1 juta token output
	PricePerSMS  string     `json:"price_per_sms"`  // PlasGate
	DailyLimit   int        `json:"daily_limit"`    // Gemini gratis: batas permintaan per hari yang terlihat di AI Studio (0 = tidak diisi)
}

var usageProviders = []string{"anthropic", "gemini", "plasgate"}

func defaultBudgets() map[string]Budget {
	return map[string]Budget{
		"anthropic": {WarnBelowUSD: "2", PriceIn: "1", PriceOut: "5"},
		"gemini":    {PriceIn: "0", PriceOut: "0"},
		"plasgate":  {WarnBelowUSD: "2", PricePerSMS: "0.03"},
	}
}

func (s *Store) budgets(ctx context.Context, q querier) (map[string]Budget, error) {
	out := defaultBudgets()
	var raw []byte
	err := q.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'usage_budget'`).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, nil
	}
	if err != nil {
		return out, err
	}
	var saved map[string]Budget
	if json.Unmarshal(raw, &saved) == nil {
		for k, v := range saved {
			out[k] = v
		}
	}
	return out, nil
}

func num(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

// RecordUsage mencatat satu panggilan. Tidak pernah menggagalkan pemanggil (error hanya dicatat di log).
func (s *Store) RecordUsage(ctx context.Context, u Usage) {
	b, err := s.budgets(ctx, s.Pool)
	if err != nil && s.Log != nil {
		s.Log.Warn("pemakaian: harga tidak terbaca", "err", err)
	}
	p := b[u.Provider]
	cost := 0.0
	if u.Err == nil {
		cost = float64(u.In)/1e6*num(p.PriceIn) + float64(u.Out)/1e6*num(p.PriceOut) + float64(u.Units)*num(p.PricePerSMS)
	}
	kind := ClassifyErr(u.Err)
	var errText *string
	if u.Err != nil {
		t := trunc(u.Err.Error())
		if len(t) > 300 {
			t = t[:300]
		}
		errText = &t
	}
	if _, err := s.Pool.Exec(context.WithoutCancel(ctx), `INSERT INTO integration_usage (provider, purpose, ok, error_kind, error_text, model, input_tokens, output_tokens, units, cost_usd)
		VALUES ($1, $2, $3, nullif($4, ''), $5, nullif($6, ''), $7, $8, $9, $10)`,
		u.Provider, u.Purpose, u.Err == nil, kind, errText, u.Model, u.In, u.Out, u.Units, cost); err != nil && s.Log != nil {
		s.Log.Warn("pemakaian tidak tercatat", "provider", u.Provider, "err", err)
	}
}

// ---- laporan ----

type Period struct {
	Requests int     `json:"requests"`
	Failed   int     `json:"failed"`
	In       int64   `json:"input_tokens"`
	Out      int64   `json:"output_tokens"`
	Units    int64   `json:"units"`
	CostUSD  float64 `json:"cost_usd"`
}

type ProviderUsage struct {
	Provider     string     `json:"provider"`
	Today        Period     `json:"today"`
	Week         Period     `json:"week"`
	Month        Period     `json:"month"`
	Budget       Budget     `json:"budget"`
	SpentSince   float64    `json:"spent_since"`   // biaya sejak saldo diisi
	RemainingUSD *float64   `json:"remaining_usd"` // perkiraan sisa (null kalau saldo belum diisi)
	LastError    *LastError `json:"last_error"`
	LastOK       *time.Time `json:"last_ok"`
	Fallbacks    int        `json:"fallbacks_7d"` // Gemini: dipakai karena Anthropic gagal (perkiraan: panggilan Gemini sukses dalam 7 hari)
}

type LastError struct {
	Kind string    `json:"kind"`
	Text string    `json:"text"`
	At   time.Time `json:"at"`
}

type DayUsage struct {
	Day      string  `json:"day"` // YYYY-MM-DD (waktu Kamboja)
	Provider string  `json:"provider"`
	Requests int     `json:"requests"`
	Failed   int     `json:"failed"`
	CostUSD  float64 `json:"cost_usd"`
}

type Alert struct {
	Provider string `json:"provider"`
	Level    string `json:"level"` // bad · warn
	Code     string `json:"code"`  // billing · auth · low_balance · quota · daily_limit · no_ai
	Value    string `json:"value,omitempty"`
}

type UsageReport struct {
	Providers []ProviderUsage `json:"providers"`
	Days      []DayUsage      `json:"days"`
	Alerts    []Alert         `json:"alerts"`
}

const dayStart = `(date_trunc('day', now() AT TIME ZONE 'Asia/Phnom_Penh') AT TIME ZONE 'Asia/Phnom_Penh')`

func (s *Store) period(ctx context.Context, provider, since string) (Period, error) {
	var p Period
	err := s.Pool.QueryRow(ctx, `SELECT count(*), count(*) FILTER (WHERE NOT ok), coalesce(sum(input_tokens), 0), coalesce(sum(output_tokens), 0),
		coalesce(sum(units), 0), coalesce(sum(cost_usd), 0)::float8 FROM integration_usage WHERE provider = $1 AND created_at >= `+since, provider).
		Scan(&p.Requests, &p.Failed, &p.In, &p.Out, &p.Units, &p.CostUSD)
	return p, err
}

func (s *Store) Usage(ctx context.Context) (UsageReport, error) {
	rep := UsageReport{Providers: []ProviderUsage{}, Days: []DayUsage{}, Alerts: []Alert{}}
	b, err := s.budgets(ctx, s.Pool)
	if err != nil {
		return rep, err
	}
	for _, prov := range usageProviders {
		pu := ProviderUsage{Provider: prov, Budget: b[prov]}
		if pu.Today, err = s.period(ctx, prov, dayStart); err != nil {
			return rep, err
		}
		if pu.Week, err = s.period(ctx, prov, `now() - interval '7 days'`); err != nil {
			return rep, err
		}
		if pu.Month, err = s.period(ctx, prov, `now() - interval '30 days'`); err != nil {
			return rep, err
		}
		if pu.Budget.BalanceAt != nil && pu.Budget.BalanceUSD != "" {
			_ = s.Pool.QueryRow(ctx, `SELECT coalesce(sum(cost_usd), 0)::float8 FROM integration_usage WHERE provider = $1 AND created_at >= $2`,
				prov, pu.Budget.BalanceAt).Scan(&pu.SpentSince)
			rem := num(pu.Budget.BalanceUSD) - pu.SpentSince
			pu.RemainingUSD = &rem
		}
		var le LastError
		var errAt *time.Time
		_ = s.Pool.QueryRow(ctx, `SELECT error_kind, coalesce(error_text, ''), created_at FROM integration_usage WHERE provider = $1 AND NOT ok
			ORDER BY created_at DESC LIMIT 1`, prov).Scan(&le.Kind, &le.Text, &errAt)
		if errAt != nil {
			le.At = *errAt
			pu.LastError = &le
		}
		_ = s.Pool.QueryRow(ctx, `SELECT max(created_at) FROM integration_usage WHERE provider = $1 AND ok`, prov).Scan(&pu.LastOK)
		if prov == "gemini" {
			pu.Fallbacks = pu.Week.Requests - pu.Week.Failed
		}
		rep.Providers = append(rep.Providers, pu)
	}
	rows, err := s.Pool.Query(ctx, `SELECT to_char(created_at AT TIME ZONE 'Asia/Phnom_Penh', 'YYYY-MM-DD') AS d, provider, count(*),
		count(*) FILTER (WHERE NOT ok), coalesce(sum(cost_usd), 0)::float8
		FROM integration_usage WHERE created_at >= `+dayStart+` - interval '13 days' GROUP BY 1, 2 ORDER BY 1 DESC, 2`)
	if err != nil {
		return rep, err
	}
	defer rows.Close()
	for rows.Next() {
		var d DayUsage
		if err := rows.Scan(&d.Day, &d.Provider, &d.Requests, &d.Failed, &d.CostUSD); err != nil {
			return rep, err
		}
		rep.Days = append(rep.Days, d)
	}
	if err := rows.Err(); err != nil {
		return rep, err
	}
	rep.Alerts = s.alertsFrom(ctx, rep.Providers)
	return rep, nil
}

// alertsFrom: peringatan dari laporan (dipakai halaman Pemakaian & saldo dan dasbor).
func (s *Store) alertsFrom(ctx context.Context, list []ProviderUsage) []Alert {
	out := []Alert{}
	ai := s.AI(ctx)
	if !ai.Ready() {
		out = append(out, Alert{Provider: "ai", Level: "warn", Code: "no_ai"})
	}
	smsOn := s.load(ctx).smsMode == SMSPlasGate
	for _, pu := range list {
		active := map[string]bool{"anthropic": ai.AnthropicKey != "", "gemini": ai.GeminiKey != "", "plasgate": smsOn}[pu.Provider]
		if !active {
			continue
		}
		// Error terakhir belum "terobati" oleh panggilan sukses sesudahnya, dan masih baru (24 jam).
		if le := pu.LastError; le != nil && time.Since(le.At) < 24*time.Hour && (pu.LastOK == nil || pu.LastOK.Before(le.At)) {
			switch le.Kind {
			case "billing":
				out = append(out, Alert{Provider: pu.Provider, Level: "bad", Code: "billing"})
			case "auth":
				out = append(out, Alert{Provider: pu.Provider, Level: "bad", Code: "auth"})
			case "quota":
				out = append(out, Alert{Provider: pu.Provider, Level: "warn", Code: "quota"})
			}
		}
		if pu.RemainingUSD != nil && pu.Budget.WarnBelowUSD != "" && *pu.RemainingUSD < num(pu.Budget.WarnBelowUSD) {
			out = append(out, Alert{Provider: pu.Provider, Level: map[bool]string{true: "bad", false: "warn"}[*pu.RemainingUSD <= 0],
				Code: "low_balance", Value: strconv.FormatFloat(*pu.RemainingUSD, 'f', 2, 64)})
		}
		if pu.Budget.DailyLimit > 0 && pu.Today.Requests >= pu.Budget.DailyLimit*8/10 {
			out = append(out, Alert{Provider: pu.Provider, Level: "warn", Code: "daily_limit",
				Value: strconv.Itoa(pu.Today.Requests) + "/" + strconv.Itoa(pu.Budget.DailyLimit)})
		}
	}
	return out
}

// Alerts: hanya peringatan (untuk dasbor & badge menu).
func (s *Store) Alerts(ctx context.Context) []Alert {
	rep, err := s.Usage(ctx)
	if err != nil {
		if s.Log != nil {
			s.Log.Warn("peringatan integrasi tidak bisa dihitung", "err", err)
		}
		return []Alert{}
	}
	return rep.Alerts
}

// BudgetInput: saldo & harga dari super admin. balance_usd diisi = saldo baru mulai dihitung sekarang.
type BudgetInput struct {
	Provider     string  `json:"provider"`
	BalanceUSD   *string `json:"balance_usd"` // null = tidak diubah · "" = hapus saldo
	WarnBelowUSD string  `json:"warn_below_usd"`
	PriceIn      string  `json:"price_in"`
	PriceOut     string  `json:"price_out"`
	PricePerSMS  string  `json:"price_per_sms"`
	DailyLimit   int     `json:"daily_limit"`
}

func validMoney(v string, max float64) bool {
	if strings.TrimSpace(v) == "" {
		return true
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	return err == nil && f >= 0 && f <= max
}

func (s *Store) SaveBudget(ctx context.Context, actorID string, in BudgetInput) error {
	known := false
	for _, p := range usageProviders {
		known = known || p == in.Provider
	}
	if !known {
		return FieldError{"provider", "unknown provider"}
	}
	switch {
	case in.BalanceUSD != nil && !validMoney(*in.BalanceUSD, 1e6):
		return FieldError{"balance_usd", "enter the balance in US dollars, e.g. 10.00"}
	case !validMoney(in.WarnBelowUSD, 1e6):
		return FieldError{"warn_below_usd", "enter an amount in US dollars"}
	case !validMoney(in.PriceIn, 1000):
		return FieldError{"price_in", "price per million tokens, e.g. 1.00"}
	case !validMoney(in.PriceOut, 1000):
		return FieldError{"price_out", "price per million tokens, e.g. 5.00"}
	case !validMoney(in.PricePerSMS, 10):
		return FieldError{"price_per_sms", "price per SMS, e.g. 0.03"}
	case in.DailyLimit < 0 || in.DailyLimit > 1000000:
		return FieldError{"daily_limit", "requests per day, 0 = not set"}
	}
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		all, err := s.budgets(ctx, tx)
		if err != nil {
			return err
		}
		b := all[in.Provider]
		if in.BalanceUSD != nil {
			b.BalanceUSD = strings.TrimSpace(*in.BalanceUSD)
			now := time.Now().UTC()
			b.BalanceAt = &now
			if b.BalanceUSD == "" {
				b.BalanceAt = nil
			}
		}
		b.WarnBelowUSD, b.PriceIn, b.PriceOut, b.PricePerSMS, b.DailyLimit = strings.TrimSpace(in.WarnBelowUSD),
			strings.TrimSpace(in.PriceIn), strings.TrimSpace(in.PriceOut), strings.TrimSpace(in.PricePerSMS), in.DailyLimit
		all[in.Provider] = b
		raw, _ := json.Marshal(all)
		_, err = tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by, updated_at) VALUES ('usage_budget', $1, $2, now())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, raw, actorID)
		return err
	})
}
