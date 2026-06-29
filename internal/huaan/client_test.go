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
	testClientCall(t, true, "test-huaan-key")
}

func TestClient_Call_NoSign(t *testing.T) {
	testClientCall(t, false, "")
}

func testClientCall(t *testing.T, signEnabled bool, cfgKey string) {
	var gotPath string
	var gotBody map[string]interface{}
	var gotSign string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotSign = r.Header.Get("sign")
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"success","data":{"phoneNo":"13800138000"}}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		HuaAn: config.HuaAnConfig{
			BaseURL:     srv.URL,
			APIPath:     "/upChannelApi",
			Key:         cfgKey,
			SignEnabled: signEnabled,
		},
		Server: config.ServerConfig{UpstreamTimeout: 5 * time.Second},
	}
	client := NewClient(cfg, nil, srv.Client())

	body := BuildRequestBody("CH001", map[string]interface{}{
		"phoneNo": "13800138000",
		"name":    "张三",
		"key":     "channel-key-should-be-stripped",
		"sign":    "old-sign-should-be-stripped",
	})
	result, err := client.Call(t.Context(), "/getProductPricesByProductCode", body)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/upChannelApi/getProductPricesByProductCode" {
		t.Fatalf("path=%s", gotPath)
	}
	if signEnabled {
		if _, ok := gotBody["key"]; ok {
			t.Fatalf("key should not be in body when sign enabled: %v", gotBody["key"])
		}
		if !sign.VerifyHuaAnHeader(sign.MapFromJSON(gotBody), cfgKey, gotSign) {
			t.Fatalf("upstream header sign verify failed, sign=%s", gotSign)
		}
	} else {
		if gotBody["key"] != "" || gotBody["sign"] != "" {
			t.Fatalf("key=%v sign=%v want empty", gotBody["key"], gotBody["sign"])
		}
		if gotSign != "" {
			t.Fatalf("header sign=%q want empty", gotSign)
		}
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

func TestClient_Call_ProductInfoUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"msg":"请求成功","data":{"productCode":"PRO2026001"}}`))
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

	body := BuildRequestBody("CH001", map[string]interface{}{"phoneNo": "13333333333"})
	if _, err := client.Call(t.Context(), ProductInfoPath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != ProductInfoUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, ProductInfoUpstreamPath)
	}
}

func TestClient_Call_SmsSendUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"data":"验证码发送成功","message":"操作成功"}`))
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

	body := BuildRequestBody("CH001", map[string]interface{}{"phoneNo": "13333333333"})
	if _, err := client.Call(t.Context(), SmsSendPath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != SmsSendUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, SmsSendUpstreamPath)
	}
}

func TestClient_Call_SmsValidUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"请求成功","data":{"userId":"u1","phoneNo":"13333333333"}}`))
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
		"phoneNo": "13333333333",
		"code":    "1234",
	})
	if _, err := client.Call(t.Context(), SmsValidPath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != SmsValidUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, SmsValidUpstreamPath)
	}
}

func TestClient_Call_SmsNoValidUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"请求成功","data":{"userId":"u1","phoneNo":"13333333333"}}`))
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

	body := BuildRequestBody("CH001", map[string]interface{}{"mobile": "13333333333"})
	if _, err := client.Call(t.Context(), SmsNoValidPath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != SmsNoValidUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, SmsNoValidUpstreamPath)
	}
}

func TestClient_Call_PriceByUserUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"成功","data":{"productCode":"PROD2025001","price":95.00}}`))
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
		"idCard":            "110101199003071234",
		"hasSocialSecurity": 1,
	})
	if _, err := client.Call(t.Context(), PriceByUserPath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != PriceByUserUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, PriceByUserUpstreamPath)
	}
}

func TestClient_Call_PolicyByPhoneUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"操作成功","data":[{"productCode":"Z041001","productName":"百万医疗险-体验版"}]}`))
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

	body := BuildRequestBody("CH001", map[string]interface{}{"phoneNo": "13800138000"})
	if _, err := client.Call(t.Context(), PolicyByPhonePath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != PolicyByPhoneUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, PolicyByPhoneUpstreamPath)
	}
}

func TestClient_Call_UserInfoByPhoneNoUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"请求成功","data":{"userId":"u1","phoneNo":"13800138000","name":"张三","idCard":"110101199001011234","channelCode":"CH001"}}`))
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

	body := BuildRequestBody("CH001", map[string]interface{}{"userId": "u1"})
	if _, err := client.Call(t.Context(), UserInfoByPhoneNoPath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != UserInfoByPhoneNoUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, UserInfoByPhoneNoUpstreamPath)
	}
}

func TestClient_Call_LiabilitiesByProductIDUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"操作成功","data":[{"liabilityName":"高保险金","price":3.16,"kindCode":"101","uwCount":20}]}`))
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
		"idCard":            "110101199003071234",
		"hasSocialSecurity": 1,
		"productType":       1,
	})
	if _, err := client.Call(t.Context(), LiabilitiesByProductIDPath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != LiabilitiesByProductIDUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, LiabilitiesByProductIDUpstreamPath)
	}
}

func TestClient_Call_GetPhoneByTokenUpstreamPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":"操作成功","data":"18212345678"}`))
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
		"userInformation": "one-click-info",
		"token":           "one-click-token",
	})
	if _, err := client.Call(t.Context(), GetPhoneByTokenPath, body); err != nil {
		t.Fatal(err)
	}
	if gotPath != GetPhoneByTokenUpstreamPath {
		t.Fatalf("path=%s want=%s", gotPath, GetPhoneByTokenUpstreamPath)
	}
}

func TestResolveUpstreamPath(t *testing.T) {
	if got := ResolveUpstreamPath("/upChannelApi", ProductInfoPath); got != ProductInfoUpstreamPath {
		t.Fatalf("product info path=%q want=%q", got, ProductInfoUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", SmsSendPath); got != SmsSendUpstreamPath {
		t.Fatalf("sms send path=%q want=%q", got, SmsSendUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", SmsValidPath); got != SmsValidUpstreamPath {
		t.Fatalf("sms valid path=%q want=%q", got, SmsValidUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", SmsNoValidPath); got != SmsNoValidUpstreamPath {
		t.Fatalf("sms noValid path=%q want=%q", got, SmsNoValidUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", PriceByUserPath); got != PriceByUserUpstreamPath {
		t.Fatalf("priceByUser path=%q want=%q", got, PriceByUserUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", PolicyByPhonePath); got != PolicyByPhoneUpstreamPath {
		t.Fatalf("policyByPhone path=%q want=%q", got, PolicyByPhoneUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", UserInfoByPhoneNoPath); got != UserInfoByPhoneNoUpstreamPath {
		t.Fatalf("getUserInfoByPhoneNo path=%q want=%q", got, UserInfoByPhoneNoUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", LiabilitiesByProductIDPath); got != LiabilitiesByProductIDUpstreamPath {
		t.Fatalf("liabilitiesByProductId path=%q want=%q", got, LiabilitiesByProductIDUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", BankListPath); got != BankListUpstreamPath {
		t.Fatalf("bank list path=%q want=%q", got, BankListUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", ProductPricesByPolicyIDPath); got != ProductPricesByPolicyIDUpstreamPath {
		t.Fatalf("productPricesByPolicyId path=%q want=%q", got, ProductPricesByPolicyIDUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", PolicyInfoByPolicyIDPath); got != PolicyInfoByPolicyIDUpstreamPath {
		t.Fatalf("policyInfoByPolicyId path=%q want=%q", got, PolicyInfoByPolicyIDUpstreamPath)
	}
	if got := ResolveUpstreamPath("/upChannelApi", GetPhoneByTokenPath); got != GetPhoneByTokenUpstreamPath {
		t.Fatalf("getPhoneByToken path=%q want=%q", got, GetPhoneByTokenUpstreamPath)
	}
}
