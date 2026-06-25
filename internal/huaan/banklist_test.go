package huaan

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/huaan/insurance-bridge/internal/config"
)

func TestParseBankPayPairs(t *testing.T) {
	raw := []byte(`{
		"code": 200,
		"message": "ok",
		"data": [
			{"bankCode": "BOC", "bankName": "中国银行", "debitCard": "1", "payChannelId": "aaa"},
			{"bankCode": "CCB", "bankName": "建设银行", "debitCard": "1", "creditCard": "1", "payChannelId": "bbb"}
		]
	}`)
	pairs, err := ParseBankPayPairs(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(pairs) != 2 || pairs[0].BankCode != "BOC" || pairs[0].PayChannelID != "aaa" {
		t.Fatalf("unexpected pairs: %+v", pairs)
	}

	got, err := SelectBankPayPair(pairs, "CCB")
	if err != nil || got.PayChannelID != "bbb" {
		t.Fatalf("select CCB: %+v err=%v", got, err)
	}
	if pairs[0].DefaultCardType() != "1" {
		t.Fatalf("cardType=%s", pairs[0].DefaultCardType())
	}
}

func TestSelectBankPayPair_notFound(t *testing.T) {
	pairs := []BankPayPair{{BankCode: "BOC", PayChannelID: "x"}}
	if _, err := SelectBankPayPair(pairs, "XYZ"); err == nil {
		t.Fatal("expected error")
	}
}

func TestClient_Call_BankListUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"成功","data":[]}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		HuaAn: config.HuaAnConfig{
			BaseURL:     srv.URL,
			APIPath:     "/upChannelApi",
			Key:         "test-key",
			SignEnabled: false,
		},
		Server: config.ServerConfig{UpstreamTimeout: 5 * time.Second},
	}
	client := NewClient(cfg, nil, srv.Client())

	body := BuildRequestBody("CH001", nil)
	if _, err := client.Call(t.Context(), BankListPath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != BankListUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, BankListUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", BankListPath); got != BankListUpstreamPath {
		t.Fatalf("resolve path=%q", got)
	}
}
