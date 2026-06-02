package huaan

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/pkg/sign"
)

func TestClient_Call(t *testing.T) {
	const huaAnKey = "test-huaan-key"
	var gotPath string
	var gotBody map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"success","data":{"phoneNo":"13800138000"}}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		HuaAn: config.HuaAnConfig{
			BaseURL: srv.URL,
			APIPath: "/upChannelApi",
		},
		Server: config.ServerConfig{UpstreamTimeout: 5 * time.Second},
	}
	client := NewClient(cfg, nil, srv.Client())

	body := BuildRequestBody("CH001", map[string]interface{}{
		"phoneNo": "13800138000",
		"name":    "张三",
	})
	result, err := client.Call(t.Context(), "/getProductInfoByChannel", body, huaAnKey)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/upChannelApi/getProductInfoByChannel" {
		t.Fatalf("path=%s", gotPath)
	}
	if gotBody["key"] != huaAnKey {
		t.Fatalf("key=%v", gotBody["key"])
	}
	if !sign.Verify(sign.MapFromJSON(gotBody), huaAnKey) {
		t.Fatal("upstream sign verify failed")
	}
	if result.HuaAnCode != 200 {
		t.Fatalf("huaanCode=%d", result.HuaAnCode)
	}
	if !ResponseOK(result.Body) {
		t.Fatalf("response: %s", result.Body)
	}
}

func TestBuildRequestBody(t *testing.T) {
	body := BuildRequestBody("CH001", map[string]interface{}{"productCode": "P1"})
	if body["channelCode"] != "CH001" || body["productCode"] != "P1" {
		t.Fatalf("unexpected body: %+v", body)
	}
	if _, ok := body["timestamp"]; !ok {
		t.Fatal("timestamp required")
	}
}
