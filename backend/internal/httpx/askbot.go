package httpx

// askbot: bot pembantu seller menjawab pertanyaan pembeli + FAQ toko.

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/askbot"
	"phsar/internal/config"
	"phsar/internal/storefront"
)

// NewAskbot: draf AI dengan kunci dari Admin → Integrasi (Anthropic, cadangan Gemini); tanpa kunci → jawaban berbasis aturan.
func NewAskbot(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool) *askbot.Service {
	ist := NewIntegrations(cfg, log, pool)
	return &askbot.Service{Pool: pool, Log: log, Gen: askbot.DynamicGen{Keys: AIKeys(cfg, ist), Record: UsageRecorder(ist, "askbot")}}
}

// askbotFAQ: FAQ per produk untuk blok "Tanya jawab" di halaman produk.
type askbotFAQ struct{ s *askbot.Service }

func (a askbotFAQ) ProductFAQ(ctx context.Context, lang, productID string) ([]storefront.QA, error) {
	list, err := a.s.ProductFAQ(ctx, lang, productID)
	if err != nil {
		return nil, err
	}
	out := make([]storefront.QA, 0, len(list))
	for _, q := range list {
		out = append(out, storefront.QA{Q: q.Question, A: q.Answer, Lang: q.Lang})
	}
	return out, nil
}
