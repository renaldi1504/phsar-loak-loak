package translate

// Gemini (Google AI Studio, generateContent) — cadangan kalau Claude gagal (mis. HTTP 400 saldo habis).
// Kunci & model diatur super admin di halaman Integrasi; .env hanya cadangan.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Gemini struct {
	Key     string
	Model   string
	BaseURL string // kosong = https://generativelanguage.googleapis.com (diganti di tes)
}

// GeminiBaseURL: alamat API (hanya diganti tes).
var GeminiBaseURL = "https://generativelanguage.googleapis.com"

func (g *Gemini) base() string {
	if g.BaseURL != "" {
		return strings.TrimRight(g.BaseURL, "/")
	}
	return GeminiBaseURL
}

// Gen: hasil satu panggilan generateContent.
type Gen struct {
	Text, Model string
	In, Out     int // token (usageMetadata)
}

// Generate: satu panggilan generateContent; dipakai terjemahan dan asisten tanya-jawab.
func (g *Gemini) Generate(ctx context.Context, system, user string, wantJSON bool, maxTokens int) (string, string, error) {
	r, err := g.GenerateU(ctx, system, user, wantJSON, maxTokens)
	return r.Text, r.Model, err
}

// GenerateU: seperti Generate, plus jumlah token.
func (g *Gemini) GenerateU(ctx context.Context, system, user string, wantJSON bool, maxTokens int) (Gen, error) {
	cfg := map[string]any{"temperature": 0.2, "maxOutputTokens": maxTokens}
	if wantJSON {
		cfg["responseMimeType"] = "application/json"
	}
	body := map[string]any{
		"systemInstruction": map[string]any{"parts": []map[string]string{{"text": system}}},
		"contents":          []map[string]any{{"role": "user", "parts": []map[string]string{{"text": user}}}},
		"generationConfig":  cfg,
	}
	b, _ := json.Marshal(body)
	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.base()+"/v1beta/models/"+g.Model+":generateContent", bytes.NewReader(b))
	if err != nil {
		return Gen{}, err
	}
	hreq.Header.Set("content-type", "application/json")
	hreq.Header.Set("x-goog-api-key", g.Key)
	resp, err := httpClient.Do(hreq)
	if err != nil {
		return Gen{}, fmt.Errorf("gemini: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if resp.StatusCode != http.StatusOK {
		return Gen{}, geminiError(resp.StatusCode, raw)
	}
	var out struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text    string `json:"text"`
					Thought bool   `json:"thought"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
		ModelVersion  string `json:"modelVersion"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
			ThoughtsTokenCount   int `json:"thoughtsTokenCount"`
		} `json:"usageMetadata"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return Gen{}, fmt.Errorf("gemini: respons tidak terbaca: %w", err)
	}
	if len(out.Candidates) == 0 {
		return Gen{}, errors.New("gemini: tidak ada jawaban (mungkin diblokir filter keamanan)")
	}
	var sb strings.Builder
	for _, p := range out.Candidates[0].Content.Parts {
		if !p.Thought {
			sb.WriteString(p.Text)
		}
	}
	text := strings.TrimSpace(sb.String())
	if text == "" {
		return Gen{}, fmt.Errorf("gemini: jawaban kosong (finishReason=%s)", out.Candidates[0].FinishReason)
	}
	model := out.ModelVersion
	if model == "" {
		model = g.Model
	}
	u := out.UsageMetadata
	return Gen{Text: text, Model: model, In: u.PromptTokenCount, Out: u.CandidatesTokenCount + u.ThoughtsTokenCount}, nil
}

func (g *Gemini) Translate(ctx context.Context, req Request) (Result, error) {
	if len(req.Fields) == 0 || len(req.Targets) == 0 {
		return Result{Texts: map[string]map[string]string{}, Model: g.Model}, nil
	}
	src, _ := json.MarshalIndent(req.Fields, "", "  ")
	targetNames := []string{}
	for _, l := range req.Targets {
		targetNames = append(targetNames, langNames[l]+" ("+l+")")
	}
	keys := make([]string, 0, len(req.Fields))
	for k := range req.Fields {
		keys = append(keys, `"`+k+`"`)
	}
	user := fmt.Sprintf("Translate this %s listing from %s into %s.\n"+
		"Answer with ONLY a JSON object: one key per target language code (%s), each an object with exactly these keys: %s.\n<listing>\n%s\n</listing>",
		req.Kind, langNames[req.Source], strings.Join(targetNames, " and "), strings.Join(req.Targets, ", "), strings.Join(keys, ", "), src)
	sys := strings.Replace(systemPrompt, "Always answer by calling the save_translations tool.", "Always answer with the JSON object only.", 1)
	gen, err := g.GenerateU(ctx, sys, user, true, 4096)
	if err != nil {
		return Result{}, err
	}
	text, model := gen.Text, gen.Model
	text = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(text, "```json"), "```"), "```"))
	var parsed map[string]map[string]string
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		return Result{}, fmt.Errorf("gemini: JSON terjemahan tidak terbaca: %w", err)
	}
	res := Result{Texts: map[string]map[string]string{}, Model: model, In: gen.In, Out: gen.Out}
	for _, l := range req.Targets {
		m := parsed[l]
		for k := range req.Fields {
			if strings.TrimSpace(m[k]) == "" && strings.TrimSpace(req.Fields[k]) != "" {
				return Result{}, fmt.Errorf("gemini: terjemahan %s.%s kosong", l, k)
			}
		}
		res.Texts[l] = m
	}
	return res, nil
}

func geminiError(status int, raw []byte) error {
	var e struct {
		Error struct {
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"error"`
	}
	_ = json.Unmarshal(raw, &e)
	msg := e.Error.Message
	if msg == "" {
		msg = truncate(string(raw), 200)
	}
	hint := ""
	switch {
	case status == 400 && strings.Contains(strings.ToLower(msg), "api key"):
		hint = "API key Gemini salah — buat di aistudio.google.com → Get API key, isi di Admin → Integrasi"
	case status == 403:
		hint = "key tidak diizinkan (API belum aktif / dibatasi wilayah)"
	case status == 404:
		hint = "nama model tidak dikenal — cek model Gemini di Admin → Integrasi"
	case status == 429:
		hint = "kuota gratis Gemini habis untuk saat ini — otomatis dicoba lagi"
	case status >= 500:
		hint = "server Gemini sedang sibuk — otomatis dicoba lagi"
	}
	if hint != "" {
		return fmt.Errorf("gemini: HTTP %d %s: %s (%s)", status, e.Error.Status, truncate(msg, 200), hint)
	}
	return fmt.Errorf("gemini: HTTP %d %s: %s", status, e.Error.Status, truncate(msg, 200))
}

// ---- Kunci dinamis: Anthropic utama → Gemini cadangan → (opsional) Google Translate dari .env ----

// Keys: kunci yang berlaku saat ini (dibaca ulang tiap job, jadi perubahan di admin langsung terpakai).
type Keys struct {
	AnthropicKey, AnthropicModel, AnthropicBaseURL string
	GeminiKey, GeminiModel                         string
}

type Dynamic struct {
	Keys   func(ctx context.Context) Keys
	Extra  Translator // cadangan terakhir (Google Translate .env), boleh nil
	Record Recorder   // laporan pemakaian (Admin → Integrasi → Pemakaian & saldo), boleh nil
}

// Recorder: dipanggil sekali per percobaan ke penyedia AI (berhasil atau gagal).
type Recorder func(ctx context.Context, provider, model string, in, out int, err error)

// ErrNoKeys: belum ada kunci AI sama sekali — job menunggu, tidak dihitung gagal.
var ErrNoKeys = errors.New("belum ada API key AI — isi kunci Anthropic atau Gemini di Admin → Integrasi")

// Ready: ada penyedia yang bisa dipakai.
func (d Dynamic) Ready(ctx context.Context) bool {
	k := d.Keys(ctx)
	return k.AnthropicKey != "" || k.GeminiKey != "" || d.Extra != nil
}

func (d Dynamic) Translate(ctx context.Context, req Request) (Result, error) {
	k := d.Keys(ctx)
	type step struct {
		name, model string
		t           Translator
	}
	var chain []step
	if k.AnthropicKey != "" {
		chain = append(chain, step{"anthropic", k.AnthropicModel, &Claude{Key: k.AnthropicKey, Model: k.AnthropicModel, BaseURL: k.AnthropicBaseURL}})
	}
	if k.GeminiKey != "" {
		chain = append(chain, step{"gemini", k.GeminiModel, &Gemini{Key: k.GeminiKey, Model: k.GeminiModel}})
	}
	if d.Extra != nil {
		chain = append(chain, step{"", "", d.Extra})
	}
	if len(chain) == 0 {
		return Result{}, ErrNoKeys
	}
	var errs []string
	for _, st := range chain {
		res, err := st.t.Translate(ctx, req)
		if d.Record != nil && st.name != "" {
			model := res.Model
			if model == "" {
				model = st.model
			}
			d.Record(ctx, st.name, model, res.In, res.Out, err)
		}
		if err == nil {
			return res, nil
		}
		// Claude gagal (400 saldo habis, 401, 429, 5xx, jaringan…) → coba penyedia berikutnya.
		errs = append(errs, err.Error())
	}
	return Result{}, errors.New(strings.Join(errs, "; cadangan: "))
}
