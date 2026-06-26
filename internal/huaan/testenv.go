package huaan

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/huaan/insurance-bridge/internal/config"
)

// EnvTestConfig 华安直连集成测试所需的环境变量。
type EnvTestConfig struct {
	BaseURL     string
	APIPath     string
	ChannelCode string
	HuaAnKey    string
	Phone       string
	Name        string
	IDCard      string
}

// LoadEnvTestConfig 从环境变量加载华安测试配置；必填项缺失时 skip 测试。
func LoadEnvTestConfig(t *testing.T) EnvTestConfig {
	t.Helper()
	cfg := EnvTestConfig{
		BaseURL:     os.Getenv("HUAAN_BASE_URL"),
		APIPath:     os.Getenv("HUAAN_API_PATH"),
		ChannelCode: os.Getenv("HUAAN_CHANNEL_CODE"),
		HuaAnKey:    os.Getenv("HUAAN_KEY"),
		Phone:       os.Getenv("HUAAN_TEST_PHONE"),
		Name:        os.Getenv("HUAAN_TEST_NAME"),
		IDCard:      os.Getenv("HUAAN_TEST_ID_CARD"),
	}
	if cfg.BaseURL == "" || cfg.ChannelCode == "" {
		t.Skip("跳过华安直连测试：请设置 HUAAN_BASE_URL、HUAAN_CHANNEL_CODE")
	}
	if cfg.APIPath == "" {
		cfg.APIPath = "/upChannelApi"
	}
	return cfg
}

// ToAppConfig 转为应用 Config，供 huaan.Client 使用。
func (e EnvTestConfig) ToAppConfig() *config.Config {
	signEnabled := false
	if v := os.Getenv("HUAAN_SIGN_ENABLED"); v != "" {
		signEnabled = v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(v, "yes")
	}
	return &config.Config{
		HuaAn: config.HuaAnConfig{
			BaseURL:     e.BaseURL,
			APIPath:     e.APIPath,
			Key:         e.HuaAnKey,
			SignEnabled: signEnabled,
		},
		Server: config.ServerConfig{UpstreamTimeout: 30 * time.Second},
	}
}

// NewTestClient 基于环境变量构造华安客户端。
func NewTestClient(t *testing.T) (*Client, EnvTestConfig) {
	t.Helper()
	env := LoadEnvTestConfig(t)
	return NewClient(env.ToAppConfig(), nil, nil), env
}

// LogRequest 输出即将发往华安的请求体（Call 注入 sign 请求头或空 key/sign 之前的状态）。
func LogRequest(t *testing.T, apiPath string, body map[string]interface{}) {
	t.Helper()
	t.Logf("=== 请求 %s ===", apiPath)
	logJSON(t, "请求体", body)
}

// LogResponse 原样输出华安 JSON 响应（明文 PII，不做加密）。
func LogResponse(t *testing.T, apiPath string, result *CallResult) {
	t.Helper()
	t.Logf("=== 响应 %s ===", apiPath)
	t.Logf("durationMs: %d, huaanCode: %d", result.DurationMs, result.HuaAnCode)
	logJSONBytes(t, "华安响应", result.Body)
}

func logJSON(t *testing.T, label string, v interface{}) {
	t.Helper()
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Logf("%s: %+v", label, v)
		return
	}
	t.Logf("%s:\n%s", label, b)
}

func logJSONBytes(t *testing.T, label string, raw []byte) {
	t.Helper()
	var pretty json.RawMessage
	if err := json.Unmarshal(raw, &pretty); err != nil {
		t.Logf("%s(raw): %s", label, raw)
		return
	}
	indented, err := json.MarshalIndent(pretty, "", "  ")
	if err != nil {
		t.Logf("%s: %s", label, raw)
		return
	}
	t.Logf("%s:\n%s", label, indented)
}
