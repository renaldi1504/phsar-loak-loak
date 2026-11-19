// Package askbot membantu seller menjawab pertanyaan pembeli tentang produknya.
//
// Dari pertanyaan pembeli + data produk (nama, deskripsi, varian/harga/stok, kondisi, garansi, retur)
// + jam buka toko + FAQ toko, bot membuat DRAF jawaban singkat dan sopan dalam bahasa pertanyaan
// (en/km/zh). Seller tetap yang mengirim. Utama: Claude (ANTHROPIC_API_KEY). Tanpa key / Claude gagal:
// jawaban berbasis aturan dari data. Bot tidak pernah mengarang fakta — kalau datanya tidak ada,
// draf bilang seller akan mengecek.
package askbot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/hours"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrDisabled    = errors.New("assist is turned off")
	ErrRateLimited = errors.New("too many requests")
)

// FieldError → 422 {"error":{"code":"validation","field":…}}.
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

const (
	MaxQuestion    = 1000 // karakter pertanyaan pembeli
	MaxFAQQuestion = 300
	MaxFAQAnswer   = 2000
	MaxFAQs        = 100 // per toko
	DefaultLimit   = 60  // draf per toko per jam
)

// Generator membuat draf dari fakta. Claude adalah implementasi utama; nil = hanya aturan.
type Generator interface {
	Draft(ctx context.Context, lang, question string, f Facts) (Draft, error)
}

// Draft: hasil generator.
type Draft struct {
	Reply        string
	Model        string
	InputTokens  int
	OutputTokens int
}

type Service struct {
	Pool *pgxpool.Pool
	Log  *slog.Logger
	Gen  Generator // nil kalau ANTHROPIC_API_KEY kosong
}

// ---- Pengaturan (super admin) ----

type Settings struct {
	Enabled     bool `json:"enabled"`
	HourlyLimit int  `json:"hourly_limit"`
}

func (s Settings) Check() error {
	if s.HourlyLimit < 1 || s.HourlyLimit > 1000 {
		return FieldError{"hourly_limit", "must be between 1 and 1000"}
	}
	return nil
}

func (s *Service) Settings(ctx context.Context) (Settings, error) {
	st := Settings{Enabled: true, HourlyLimit: DefaultLimit}
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'askbot'`).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return st, nil
	}
	if err != nil {
		return st, err
	}
	_ = json.Unmarshal(raw, &st)
	if st.HourlyLimit < 1 {
		st.HourlyLimit = DefaultLimit
	}
	return st, nil
}

// ---- Fakta untuk bot ----

// Facts: semua yang boleh dipakai bot. Yang tidak ada di sini tidak boleh dijawab sebagai fakta.
type Facts struct {
	Store   StoreFacts    `json:"store"`
	Product *ProductFacts `json:"product,omitempty"`
	FAQs    []FAQFact     `json:"faqs,omitempty"`
}

type StoreFacts struct {
	Name  string   `json:"name"`
	Hours []string `json:"opening_hours,omitempty"` // "Sunday: 08:00-22:00" / "Monday: closed"; kosong = tidak diketahui
	week  *hours.Week
}

type ProductFacts struct {
	Name           string        `json:"name"`
	Description    string        `json:"description,omitempty"`
	ForSale        bool          `json:"currently_for_sale"`
	Condition      string        `json:"condition"`
	ConditionNote  string        `json:"condition_note,omitempty"`
	UsageMonths    *int          `json:"used_for_months,omitempty"`
	WarrantyType   string        `json:"warranty_type"` // none · store · official
	WarrantyMonths int           `json:"warranty_months"`
	ReturnDays     int           `json:"return_days"` // 0 = tidak menerima retur
	PreorderDays   *int          `json:"preorder_days,omitempty"`
	Options        []string      `json:"option_names,omitempty"`
	Variants       []VariantFact `json:"variants"`
	CurrencyNote   string        `json:"currency"`
}

type VariantFact struct {
	Label string `json:"label,omitempty"` // mis. "Black / M"; kosong = produk tanpa varian
	Price string `json:"price_usd"`
	Stock int    `json:"in_stock"`
}

type FAQFact struct {
	Question string `json:"q"`
	Answer   string `json:"a"`
	Lang     string `json:"lang"`
	Product  bool   `json:"about_this_product,omitempty"`
}

var conditionNames = map[string]string{"new": "new", "like_new": "like new", "used_good": "used - good", "used_fair": "used - fair", "refurbished": "refurbished"}
var dayNames = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// LoadFacts: data toko + produk (kalau ada) + FAQ. Produk harus milik toko itu.
func (s *Service) LoadFacts(ctx context.Context, storeID, productID, lang string) (Facts, error) {
	var f Facts
	if err := s.Pool.QueryRow(ctx, `SELECT name FROM stores WHERE id = $1 AND deleted_at IS NULL`, storeID).Scan(&f.Store.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return f, ErrNotFound
		}
		return f, err
	}
	if w, err := hours.StoreWeek(ctx, s.Pool, storeID); err == nil && w != nil {
		f.Store.week = w
		for i, d := range w {
			if d == nil {
				f.Store.Hours = append(f.Store.Hours, dayNames[i]+": closed")
			} else {
				f.Store.Hours = append(f.Store.Hours, dayNames[i]+": "+d.Open+"-"+d.Close)
			}
		}
	}
	if productID != "" {
		p, err := s.loadProduct(ctx, storeID, productID, lang)
		if err != nil {
			return f, err
		}
		f.Product = p
	}
	rows, err := s.Pool.Query(ctx, `SELECT question, answer, lang::text, product_id IS NOT NULL FROM store_faqs
		WHERE store_id = $1 AND (product_id IS NULL OR product_id = $2::uuid)
		ORDER BY product_id IS NULL, (lang = $3::lang_code) DESC, position, created_at LIMIT 60`, storeID, nullable(productID), lang)
	if err != nil {
		return f, err
	}
	defer rows.Close()
	for rows.Next() {
		var q FAQFact
		if err := rows.Scan(&q.Question, &q.Answer, &q.Lang, &q.Product); err != nil {
			return f, err
		}
		f.FAQs = append(f.FAQs, q)
	}
	return f, rows.Err()
}

func (s *Service) loadProduct(ctx context.Context, storeID, productID, lang string) (*ProductFacts, error) {
	p := &ProductFacts{CurrencyNote: "USD"}
	var src, status, moderation string
	var usage, preorder *int16
	var wm, rd int16
	err := s.Pool.QueryRow(ctx, `SELECT p.source_lang::text, p.status::text, p.moderation_status::text, p.condition::text,
		coalesce(p.condition_note, ''), p.usage_months, p.warranty_type, p.warranty_months, p.return_days, p.preorder_days,
		coalesce(t.name, o.name, ''), coalesce(nullif(t.description, ''), o.description, '')
		FROM products p
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = $3::lang_code
		LEFT JOIN product_translations o ON o.product_id = p.id AND o.lang = p.source_lang
		WHERE p.id = $1 AND p.store_id = $2 AND p.deleted_at IS NULL`, productID, storeID, lang).
		Scan(&src, &status, &moderation, &p.Condition, &p.ConditionNote, &usage, &p.WarrantyType, &wm, &rd, &preorder, &p.Name, &p.Description)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p.ForSale = status == "active" && moderation == "approved"
	p.WarrantyMonths, p.ReturnDays = int(wm), int(rd)
	if c, ok := conditionNames[p.Condition]; ok {
		p.Condition = c
	}
	if usage != nil {
		u := int(*usage)
		p.UsageMonths = &u
	}
	if preorder != nil {
		d := int(*preorder)
		p.PreorderDays = &d
	}
	p.Description = clip(p.Description, 3000)

	rows, err := s.Pool.Query(ctx, `SELECT coalesce(name_i18n->>$2, name_i18n->>$3, '') FROM product_options WHERE product_id = $1 ORDER BY position`,
		productID, lang, src)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			rows.Close()
			return nil, err
		}
		p.Options = append(p.Options, n)
	}
	rows.Close()

	rows, err = s.Pool.Query(ctx, `
		SELECT v.price::text, greatest(v.stock_qty - v.reserved_qty, 0),
		       coalesce(string_agg(coalesce(ov.value_i18n->>$2, ov.value_i18n->>$3, ''), ' / ' ORDER BY o.position) FILTER (WHERE ov.id IS NOT NULL), '')
		FROM product_variants v
		LEFT JOIN variant_option_values vov ON vov.variant_id = v.id
		LEFT JOIN product_option_values ov ON ov.id = vov.option_value_id
		LEFT JOIN product_options o ON o.id = vov.option_id
		WHERE v.product_id = $1 AND v.deleted_at IS NULL AND v.status = 'active'
		GROUP BY v.id ORDER BY v.is_default DESC, v.price, v.sku LIMIT 40`, productID, lang, src)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var v VariantFact
		if err := rows.Scan(&v.Price, &v.Stock, &v.Label); err != nil {
			return nil, err
		}
		p.Variants = append(p.Variants, v)
	}
	return p, rows.Err()
}

// ---- Draf jawaban ----

type SuggestInput struct {
	Question  string `json:"question"`
	ProductID string `json:"product_id,omitempty"`
}

type SuggestResult struct {
	Reply string `json:"reply"`
	Mode  string `json:"-"` // claude · rules
	Lang  string `json:"-"`
}

// Suggest: validasi → fitur aktif? → batas per jam → fakta → Claude (atau aturan) → catat pemakaian.
func (s *Service) Suggest(ctx context.Context, storeID, userID string, in SuggestInput) (SuggestResult, error) {
	q := strings.TrimSpace(in.Question)
	pid := strings.TrimSpace(in.ProductID)
	switch {
	case q == "":
		return SuggestResult{}, FieldError{"question", "required"}
	case utf8.RuneCountInString(q) > MaxQuestion:
		return SuggestResult{}, FieldError{"question", fmt.Sprintf("at most %d characters", MaxQuestion)}
	case pid != "" && !isUUID(pid):
		return SuggestResult{}, FieldError{"product_id", "invalid id"}
	}
	st, err := s.Settings(ctx)
	if err != nil {
		return SuggestResult{}, err
	}
	if !st.Enabled {
		return SuggestResult{}, ErrDisabled
	}
	var used int
	if err := s.Pool.QueryRow(ctx, `SELECT count(*) FROM askbot_usage WHERE store_id = $1 AND created_at > now() - interval '1 hour'`, storeID).Scan(&used); err != nil {
		return SuggestResult{}, err
	}
	if used >= st.HourlyLimit {
		return SuggestResult{}, ErrRateLimited
	}
	lang := DetectLang(q)
	f, err := s.LoadFacts(ctx, storeID, pid, lang)
	if err != nil {
		return SuggestResult{}, err
	}

	res := SuggestResult{Lang: lang, Mode: "rules"}
	var d Draft
	var genErr error
	if s.Gen != nil {
		d, genErr = s.Gen.Draft(ctx, lang, q, f)
		noAI := errors.Is(genErr, ErrNoAI) // belum ada kunci AI: langsung jawaban berbasis aturan, bukan error
		if noAI {
			genErr = nil
		}
		if noAI {
		} else if genErr == nil && strings.TrimSpace(d.Reply) != "" {
			res.Mode, res.Reply = "claude", clip(strings.TrimSpace(d.Reply), 2000)
		} else if genErr == nil {
			genErr = errors.New("claude: jawaban kosong")
		}
		if genErr != nil && s.Log != nil {
			s.Log.Warn("askbot: Claude gagal, pakai jawaban berbasis aturan", "err", genErr)
		}
	}
	if res.Mode == "rules" {
		res.Reply = RuleReply(lang, q, f)
	}
	var errText *string
	if genErr != nil {
		e := clip(genErr.Error(), 300)
		errText = &e
	}
	var model *string
	if res.Mode == "claude" && d.Model != "" {
		model = &d.Model
	}
	if _, err := s.Pool.Exec(ctx, `INSERT INTO askbot_usage (store_id, user_id, product_id, lang, mode, model, input_tokens, output_tokens, error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, storeID, nullable(userID), nullable(pid), lang, res.Mode, model, d.InputTokens, d.OutputTokens, errText); err != nil {
		return SuggestResult{}, err
	}
	return res, nil
}

// DetectLang: bahasa pertanyaan dari aksaranya — Khmer → km, Han → zh, selain itu en.
func DetectLang(s string) string {
	var km, zh int
	for _, r := range s {
		switch {
		case r >= 0x1780 && r <= 0x17FF, r >= 0x19E0 && r <= 0x19FF:
			km++
		case unicode.Is(unicode.Han, r):
			zh++
		}
	}
	switch {
	case km > 0 && km >= zh:
		return "km"
	case zh > 0:
		return "zh"
	}
	return "en"
}

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func clip(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	r := []rune(s)
	return string(r[:n]) + "…"
}

func isUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
				return false
			}
		}
	}
	return true
}
