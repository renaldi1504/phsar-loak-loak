package askbot

// DynamicGen: draf AI dengan kunci dari Admin → Integrasi. Anthropic dulu; kalau gagal (mis. HTTP 400 saldo habis)
// pakai Gemini. Tanpa kunci sama sekali → ErrNoAI (Service memakai jawaban berbasis aturan tanpa mencatat error).

import (
	"context"
	"errors"
	"fmt"

	"phsar/internal/translate"
)

var ErrNoAI = errors.New("askbot: belum ada API key AI")

type DynamicGen struct {
	Keys   func(ctx context.Context) translate.Keys
	Record translate.Recorder // laporan pemakaian, boleh nil
}

func (d DynamicGen) rec(ctx context.Context, provider, model string, in, out int, err error) {
	if d.Record != nil {
		d.Record(ctx, provider, model, in, out, err)
	}
}

// Ready: ada kunci AI (untuk tanda "ai" di halaman seller/admin).
func (d DynamicGen) Ready(ctx context.Context) bool {
	k := d.Keys(ctx)
	return k.AnthropicKey != "" || k.GeminiKey != ""
}

// aiOn: generator AI benar-benar bisa dipakai sekarang.
func (s *Service) aiOn(ctx context.Context) bool {
	if g, ok := s.Gen.(interface{ Ready(context.Context) bool }); ok {
		return g.Ready(ctx)
	}
	return s.Gen != nil
}

func (d DynamicGen) Draft(ctx context.Context, lang, question string, f Facts) (Draft, error) {
	k := d.Keys(ctx)
	if k.AnthropicKey == "" && k.GeminiKey == "" {
		return Draft{}, ErrNoAI
	}
	var first error
	if k.AnthropicKey != "" {
		dr, err := (&Claude{Key: k.AnthropicKey, Model: k.AnthropicModel, BaseURL: k.AnthropicBaseURL}).Draft(ctx, lang, question, f)
		model := dr.Model
		if model == "" {
			model = k.AnthropicModel
		}
		d.rec(ctx, "anthropic", model, dr.InputTokens, dr.OutputTokens, err)
		if err == nil {
			return dr, nil
		}
		first = err
	}
	if k.GeminiKey != "" {
		g := &translate.Gemini{Key: k.GeminiKey, Model: k.GeminiModel}
		gen, err := g.GenerateU(ctx, systemPrompt, draftPrompt(lang, question, f), false, 600)
		model := gen.Model
		if model == "" {
			model = k.GeminiModel
		}
		d.rec(ctx, "gemini", model, gen.In, gen.Out, err)
		if err == nil {
			return Draft{Reply: gen.Text, Model: gen.Model, InputTokens: gen.In, OutputTokens: gen.Out}, nil
		}
		if first != nil {
			return Draft{}, fmt.Errorf("%v; cadangan: %w", first, err)
		}
		return Draft{}, err
	}
	return Draft{}, first
}
