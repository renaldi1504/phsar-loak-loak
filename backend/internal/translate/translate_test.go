package translate

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClaudeRequestAndParse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "k" || r.Header.Get("anthropic-version") == "" || r.URL.Path != "/v1/messages" {
			t.Errorf("header/path salah: %v %s", r.Header, r.URL.Path)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["model"] != "claude-haiku-5-5" {
			t.Errorf("model = %v", body["model"])
		}
		msg := body["messages"].([]any)[0].(map[string]any)["content"].(string)
		if !strings.Contains(msg, "<listing>") || !strings.Contains(msg, "Kaos merah") {
			t.Errorf("prompt tidak memuat listing: %s", msg)
		}
		_, _ = w.Write([]byte(`{"model":"claude-haiku-5-5","stop_reason":"tool_use","content":[{"type":"tool_use","name":"save_translations",
			"input":{"km":{"name":"អាវយឺតក្រហម"},"zh":{"name":"红色T恤"}}}]}`))
	}))
	defer srv.Close()
	c := &Claude{Key: "k", Model: "claude-haiku-5-5", BaseURL: srv.URL}
	res, err := c.Translate(context.Background(), Request{Kind: "product", Source: "en", Targets: []string{"km", "zh"},
		Fields: map[string]string{"name": "Kaos merah"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.Texts["zh"]["name"] != "红色T恤" || res.Texts["km"]["name"] == "" {
		t.Fatalf("hasil: %+v", res)
	}
}

func TestFallbackToGoogle(t *testing.T) {
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"overloaded"}`, 529)
	}))
	defer bad.Close()
	google := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Target string   `json:"target"`
			Q      []string `json:"q"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Target != "zh-CN" {
			t.Errorf("target = %s", body.Target)
		}
		_, _ = w.Write([]byte(`{"data":{"translations":[{"translatedText":"你好"}]}}`))
	}))
	defer google.Close()
	tr := WithFallback{&Claude{Key: "k", Model: "m", BaseURL: bad.URL}, &Google{Key: "g", BaseURL: google.URL}}
	res, err := tr.Translate(context.Background(), Request{Source: "en", Targets: []string{"zh"}, Fields: map[string]string{"name": "hello"}})
	if err != nil || res.Texts["zh"]["name"] != "你好" || res.Model != "google-translate-v2" {
		t.Fatalf("res=%+v err=%v", res, err)
	}
}

func TestNewDisabledWithoutKey(t *testing.T) {
	if New(Config{Provider: "claude"}) != nil {
		t.Fatal("tanpa key harus nil")
	}
	if _, ok := New(Config{Provider: "claude", AnthropicKey: "a", GoogleKey: "g", Fallback: "google"}).(WithFallback); !ok {
		t.Fatal("harus pakai cadangan")
	}
}

func TestClaudeErrorHints(t *testing.T) {
	cases := map[int]string{401: "API key salah", 404: "nama model", 529: "sibuk"}
	for status, want := range cases {
		err := claudeError(status, []byte(`{"type":"error","error":{"type":"x","message":"boom"}}`))
		if !strings.Contains(err.Error(), want) {
			t.Errorf("%d: %v (harusnya berisi %q)", status, err, want)
		}
	}
	err := claudeError(400, []byte(`{"error":{"type":"invalid_request_error","message":"Your credit balance is too low to access the Anthropic API."}}`))
	if !strings.Contains(err.Error(), "saldo API habis") {
		t.Errorf("credit: %v", err)
	}
}
