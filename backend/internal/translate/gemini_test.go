package translate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Anthropic menjawab 400 (saldo habis) → Gemini dipakai sebagai cadangan.
func TestDynamicFallbackToGemini(t *testing.T) {
	claude := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"type":"error","error":{"type":"invalid_request_error","message":"Your credit balance is too low to access the Anthropic API."}}`))
	}))
	defer claude.Close()
	var gotKey, gotPath string
	gemini := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey, gotPath = r.Header.Get("x-goog-api-key"), r.URL.Path
		_, _ = w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"{\"km\":{\"name\":\"អាវយឺតក្រហម\"}}"}]},"finishReason":"STOP"}],"modelVersion":"gemini-3.5-flash","usageMetadata":{"promptTokenCount":120,"candidatesTokenCount":30}}`))
	}))
	defer gemini.Close()
	old := GeminiBaseURL
	GeminiBaseURL = gemini.URL
	defer func() { GeminiBaseURL = old }()

	var calls []string
	d := Dynamic{Keys: func(context.Context) Keys {
		return Keys{AnthropicKey: "sk-ant-x", AnthropicModel: "m", AnthropicBaseURL: claude.URL, GeminiKey: "g-key", GeminiModel: "gemini-3.5-flash"}
	}, Record: func(_ context.Context, provider, model string, in, out int, err error) {
		calls = append(calls, fmt.Sprintf("%s|%s|%d|%d|%v", provider, model, in, out, err != nil))
	}}
	res, err := d.Translate(context.Background(), Request{Kind: "product", Source: "en", Targets: []string{"km"}, Fields: map[string]string{"name": "Red T-shirt"}})
	if err != nil {
		t.Fatalf("cadangan Gemini harus dipakai: %v", err)
	}
	if res.Texts["km"]["name"] != "អាវយឺតក្រហម" || res.Model != "gemini-3.5-flash" {
		t.Fatalf("hasil salah: %+v", res)
	}
	// Pemakaian tercatat per percobaan: Anthropic gagal, Gemini berhasil dengan jumlah token.
	if len(calls) != 2 || calls[0] != "anthropic|m|0|0|true" || calls[1] != "gemini|gemini-3.5-flash|120|30|false" {
		t.Fatalf("catatan pemakaian salah: %v", calls)
	}
	if gotKey != "g-key" || !strings.HasSuffix(gotPath, "/models/gemini-3.5-flash:generateContent") {
		t.Fatalf("permintaan Gemini salah: key=%q path=%q", gotKey, gotPath)
	}

	// Gemini juga gagal → kedua pesan error ikut.
	GeminiBaseURL = claude.URL
	_, err = d.Translate(context.Background(), Request{Kind: "product", Source: "en", Targets: []string{"km"}, Fields: map[string]string{"name": "x"}})
	if err == nil || !strings.Contains(err.Error(), "credit balance") || !strings.Contains(err.Error(), "gemini") {
		t.Fatalf("error gabungan: %v", err)
	}

	// Tanpa kunci: ErrNoKeys dan Ready=false (worker menunggu).
	none := Dynamic{Keys: func(context.Context) Keys { return Keys{} }}
	if none.Ready(context.Background()) {
		t.Fatal("tanpa kunci tidak boleh Ready")
	}
	if _, err := none.Translate(context.Background(), Request{Kind: "product", Source: "en", Targets: []string{"km"}, Fields: map[string]string{"name": "x"}}); !errors.Is(err, ErrNoKeys) {
		t.Fatalf("harus ErrNoKeys: %v", err)
	}
}
