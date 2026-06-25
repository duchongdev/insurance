package huaan

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/huaan/insurance-bridge/internal/config"
)

func TestParseProInsuranceResponse(t *testing.T) {
	raw := []byte(`{
		"code": 200,
		"message": "成功",
		"data": {
			"policyId": "9697d878b6474c619fa43ffa9ca3b4b0",
			"policyNo": "PN2025001",
			"policyStatus": "0",
			"userId": "5c819c252ab343e2a2a6df98b3a014ab"
		}
	}`)
	got, err := ParseProInsuranceResponse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.PolicyID != "9697d878b6474c619fa43ffa9ca3b4b0" {
		t.Fatalf("policyId=%q", got.PolicyID)
	}
	if got.PolicyNo != "PN2025001" {
		t.Fatalf("policyNo=%q", got.PolicyNo)
	}
}

func TestClient_Call_ProInsuranceUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"成功","data":{"policyId":"p1","policyNo":"n1"}}`))
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

	body := BuildRequestBody("CH001", map[string]interface{}{
		"productCode":       "PROD2025001",
		"phoneNo":           "13800138000",
		"name":              "张三",
		"idCard":            "110101199003071234",
		"hasSocialSecurity": 1,
		"isUpgrade":         0,
		"autoRenew":         1,
	})
	if _, err := client.Call(t.Context(), ProInsurancePath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != ProInsuranceUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, ProInsuranceUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", ProInsurancePath); got != ProInsuranceUpstreamPath {
		t.Fatalf("resolve path=%q", got)
	}
}
