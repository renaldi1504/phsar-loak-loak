package integ

// Kunci AI (Anthropic = utama, Gemini = cadangan) dan Google Maps, diatur super admin di halaman Integrasi.
// Prioritas: kunci tersimpan di admin; kalau belum ada, nilai dari .env (Defaults) tetap dipakai.

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
)

const (
	DefaultAnthropicModel = "claude-haiku-5-5"
	DefaultGeminiModel    = "gemini-3.5-flash"
)

var modelRe = regexp.MustCompile(`^[a-zA-Z0-9._:/-]{3,80}$`)

// AIKeys: kunci yang berlaku sekarang (rahasia — jangan pernah dikirim ke browser).
type AIKeys struct {
	AnthropicKey, AnthropicModel string
	GeminiKey, GeminiModel       string
}

// Ready: ada minimal satu penyedia AI.
func (k AIKeys) Ready() bool { return k.AnthropicKey != "" || k.GeminiKey != "" }

type aiState struct {
	keys    AIKeys
	srcA    string // admin · env · ""
	srcG    string
	mapsKey string
	mapID   string
	srcM    string
}

// readAI: bagian AI & peta dari database, dengan .env sebagai cadangan.
func (s *Store) readAI(ctx context.Context, q querier) (aiState, error) {
	st := aiState{keys: AIKeys{AnthropicModel: s.Def.AnthropicModel, GeminiModel: s.Def.GeminiModel}, mapID: s.Def.MapsMapID}
	if st.keys.AnthropicModel == "" {
		st.keys.AnthropicModel = DefaultAnthropicModel
	}
	if st.keys.GeminiModel == "" {
		st.keys.GeminiModel = DefaultGeminiModel
	}
	if s.Def.AnthropicKey != "" {
		st.keys.AnthropicKey, st.srcA = s.Def.AnthropicKey, "env"
	}
	if s.Def.GeminiKey != "" {
		st.keys.GeminiKey, st.srcG = s.Def.GeminiKey, "env"
	}
	if s.Def.MapsKey != "" {
		st.mapsKey, st.srcM = s.Def.MapsKey, "env"
	}
	for _, code := range []string{"anthropic", "gemini", "google_maps"} {
		r, err := readRow(ctx, q, code)
		if err != nil {
			return st, err
		}
		if !r.exists {
			continue
		}
		key, err := s.open(r, code)
		if err != nil {
			return st, err
		}
		switch code {
		case "anthropic":
			if m := r.config["model"]; m != "" {
				st.keys.AnthropicModel = m
			}
			if r.enabled && len(key) > 0 {
				st.keys.AnthropicKey, st.srcA = string(key), "admin"
			} else if !r.enabled {
				st.keys.AnthropicKey, st.srcA = "", "" // dimatikan di admin: .env juga tidak dipakai
			}
		case "gemini":
			if m := r.config["model"]; m != "" {
				st.keys.GeminiModel = m
			}
			if r.enabled && len(key) > 0 {
				st.keys.GeminiKey, st.srcG = string(key), "admin"
			} else if !r.enabled {
				st.keys.GeminiKey, st.srcG = "", ""
			}
		case "google_maps":
			if id := r.config["map_id"]; id != "" {
				st.mapID = id
			}
			if len(key) > 0 {
				st.mapsKey, st.srcM = string(key), "admin"
			}
		}
	}
	return st, nil
}

// AI: kunci AI yang berlaku (cache sama dengan pengaturan login/SMS).
func (s *Store) AI(ctx context.Context) AIKeys { return s.load(ctx).ai.keys }

// Maps: kunci Google Maps (browser) + Map ID.
func (s *Store) Maps(ctx context.Context) (key, mapID string) {
	st := s.load(ctx)
	return st.ai.mapsKey, st.ai.mapID
}

// ---- tampilan & perubahan oleh super admin ----

type AIView struct {
	Anthropic ProviderView `json:"anthropic"`
	Gemini    ProviderView `json:"gemini"`
	Order     []string     `json:"order"` // urutan yang dicoba
}

type ProviderView struct {
	Enabled    bool        `json:"enabled"`
	Active     bool        `json:"active"` // punya kunci & dipakai
	Source     string      `json:"source"` // admin · env · ""
	Model      string      `json:"model"`
	KeySet     bool        `json:"key_set"`
	KeyHint    string      `json:"key_hint"`
	EnvKey     bool        `json:"env_key"` // .env punya kunci cadangan
	VerifiedAt any         `json:"verified_at"`
	LastTest   *TestResult `json:"last_test"`
}

type MapsView struct {
	Active  bool   `json:"active"`
	Source  string `json:"source"`
	MapID   string `json:"map_id"`
	KeySet  bool   `json:"key_set"`
	KeyHint string `json:"key_hint"`
	EnvKey  bool   `json:"env_key"`
}

func (s *Store) aiView(ctx context.Context, st *state) (AIView, MapsView, error) {
	var v AIView
	var m MapsView
	for _, p := range []struct {
		code string
		pv   *ProviderView
		env  string
		src  string
		key  string
		mod  string
	}{
		{"anthropic", &v.Anthropic, s.Def.AnthropicKey, st.ai.srcA, st.ai.keys.AnthropicKey, st.ai.keys.AnthropicModel},
		{"gemini", &v.Gemini, s.Def.GeminiKey, st.ai.srcG, st.ai.keys.GeminiKey, st.ai.keys.GeminiModel},
	} {
		r, err := readRow(ctx, s.Pool, p.code)
		if err != nil {
			return v, m, err
		}
		p.pv.Enabled = !r.exists || r.enabled
		p.pv.Active, p.pv.Source, p.pv.Model = p.key != "", p.src, p.mod
		p.pv.KeySet, p.pv.KeyHint, p.pv.EnvKey, p.pv.LastTest = r.hint != "", r.hint, p.env != "", r.last
		if r.verified != nil {
			p.pv.VerifiedAt = r.verified
		}
	}
	if v.Anthropic.Active {
		v.Order = append(v.Order, "anthropic")
	}
	if v.Gemini.Active {
		v.Order = append(v.Order, "gemini")
	}
	r, err := readRow(ctx, s.Pool, "google_maps")
	if err != nil {
		return v, m, err
	}
	m.Active, m.Source, m.MapID = st.ai.mapsKey != "", st.ai.srcM, st.ai.mapID
	m.KeySet, m.KeyHint, m.EnvKey = r.hint != "", r.hint, s.Def.MapsKey != ""
	return v, m, nil
}

// AIInput: kunci kosong/null = tidak diubah; Clear = hapus kunci tersimpan (kembali ke .env).
type AIInput struct {
	Provider string  `json:"provider"` // anthropic · gemini
	Enabled  bool    `json:"enabled"`
	Model    string  `json:"model"`
	Key      *string `json:"key"`
	Clear    bool    `json:"clear"`
}

var aiKeyRe = map[string]*regexp.Regexp{
	"anthropic": regexp.MustCompile(`^sk-ant-[A-Za-z0-9_-]{20,200}$`),
	// Gemini: kunci lama "AIza…" (standard key) atau kunci baru "AQ.…" (auth key, bawaan Google AI Studio sejak 2026).
	"gemini": regexp.MustCompile(`^(AIza[A-Za-z0-9_-]{30,80}|AQ\.[A-Za-z0-9._-]{20,300})$`),
}

func (s *Store) SaveAI(ctx context.Context, actorID string, in AIInput) error {
	defer s.Invalidate()
	re, ok := aiKeyRe[in.Provider]
	if !ok {
		return FieldError{"provider", "unknown provider"}
	}
	in.Model = strings.TrimSpace(in.Model)
	if in.Model == "" {
		in.Model = map[string]string{"anthropic": DefaultAnthropicModel, "gemini": DefaultGeminiModel}[in.Provider]
	}
	if !modelRe.MatchString(in.Model) {
		return FieldError{"model", "model name looks wrong"}
	}
	key := strings.TrimSpace(deref(in.Key))
	if key != "" && !re.MatchString(key) {
		if in.Provider == "anthropic" {
			return FieldError{"key", "Anthropic keys start with sk-ant- (console.anthropic.com → API keys)"}
		}
		return FieldError{"key", "Gemini key looks wrong — it starts with AQ. or AIza (Google AI Studio → Get API key)"}
	}
	return s.saveSecret(ctx, actorID, in.Provider, in.Enabled, map[string]string{"model": in.Model}, key, in.Clear)
}

// MapsInput: kunci browser Google Maps (Maps JavaScript + Places + Geocoding).
type MapsInput struct {
	MapID string  `json:"map_id"`
	Key   *string `json:"key"`
	Clear bool    `json:"clear"`
}

var mapsKeyRe = regexp.MustCompile(`^AIza[A-Za-z0-9_-]{30,60}$`)
var mapIDRe = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

func (s *Store) SaveMaps(ctx context.Context, actorID string, in MapsInput) error {
	defer s.Invalidate()
	in.MapID = strings.TrimSpace(in.MapID)
	if in.MapID != "" && !mapIDRe.MatchString(in.MapID) {
		return FieldError{"map_id", "Map ID looks wrong"}
	}
	key := strings.TrimSpace(deref(in.Key))
	if key != "" && !mapsKeyRe.MatchString(key) {
		return FieldError{"key", "Google Maps keys start with AIza (Google Cloud → APIs & Services → Credentials)"}
	}
	return s.saveSecret(ctx, actorID, "google_maps", true, map[string]string{"map_id": in.MapID}, key, in.Clear)
}

func (s *Store) saveSecret(ctx context.Context, actorID, code string, enabled bool, cfg map[string]string, key string, clear bool) error {
	var sealed []byte
	var hint string
	if key != "" {
		var err error
		if sealed, err = s.Box.Seal([]byte(key), aad(code)); err != nil {
			return err
		}
		hint = hintOf(key)
	}
	raw, _ := json.Marshal(cfg)
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `INSERT INTO integrations (code, enabled, config, secret_enc, secret_hint, updated_by)
			VALUES ($1, $2, $3, $4, nullif($5, ''), $6)
			ON CONFLICT (code) DO UPDATE SET enabled = EXCLUDED.enabled, config = EXCLUDED.config,
				secret_enc = CASE WHEN $7 THEN NULL ELSE coalesce(EXCLUDED.secret_enc, integrations.secret_enc) END,
				secret_hint = CASE WHEN $7 THEN NULL ELSE coalesce(EXCLUDED.secret_hint, integrations.secret_hint) END,
				verified_at = CASE WHEN $7 OR EXCLUDED.secret_enc IS NOT NULL OR integrations.config <> EXCLUDED.config THEN NULL ELSE integrations.verified_at END,
				last_test = CASE WHEN $7 OR EXCLUDED.secret_enc IS NOT NULL THEN NULL ELSE integrations.last_test END,
				updated_by = EXCLUDED.updated_by`, code, enabled, raw, sealed, hint, actorID, clear); err != nil {
			return err
		}
		return nil
	})
}

// Tester: menguji satu penyedia AI dengan kunci tertentu (diisi paket httpx → translate, supaya integ tidak bergantung ke translate).
type Tester func(ctx context.Context, provider, key, model string) error

// TestAI: uji kunci yang berlaku untuk penyedia itu (terjemahan pendek), simpan hasilnya.
func (s *Store) TestAI(ctx context.Context, actorID, provider string, test Tester) (TestResult, error) {
	if _, ok := aiKeyRe[provider]; !ok {
		return TestResult{}, FieldError{"provider", "unknown provider"}
	}
	s.Invalidate()
	k := s.AI(ctx)
	key, model := k.AnthropicKey, k.AnthropicModel
	if provider == "gemini" {
		key, model = k.GeminiKey, k.GeminiModel
	}
	if key == "" {
		return TestResult{}, FieldError{"key", "enter the API key first"}
	}
	res := TestResult{At: time.Now().UTC()}
	terr := test(ctx, provider, key, model)
	if terr != nil {
		terr = redact(terr, key)
	}
	s.RecordUsage(ctx, Usage{Provider: provider, Purpose: "test", Model: model, Err: terr})
	if terr != nil {
		res.Message = terr.Error()
	} else {
		res.OK, res.Message = true, "Test translation worked ("+model+")."
	}
	// Baris mungkin belum ada (kunci dari .env): buat dulu supaya hasil uji tersimpan.
	if err := db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		cfg, _ := json.Marshal(map[string]string{"model": model})
		_, err := tx.Exec(ctx, `INSERT INTO integrations (code, enabled, config, updated_by) VALUES ($1, true, $2, $3) ON CONFLICT (code) DO NOTHING`,
			provider, cfg, actorID)
		return err
	}); err != nil {
		return res, err
	}
	return res, s.saveTest(ctx, actorID, provider, res)
}
