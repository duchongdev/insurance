package huaan

import (
	"encoding/json"
	"os"
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
	return &config.Config{
		HuaAn: config.HuaAnConfig{
			BaseURL:     e.BaseURL,
			APIPath:     e.APIPath,
			SignEnabled: false,
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

// LogResponse 原样输出华安 JSON 响应（明文 PII，不做加密）。
func LogResponse(t *testing.T, apiPath string, result *CallResult) {
	t.Helper()
	t.Logf("=== %s ===", apiPath)
	t.Logf("durationMs: %d, huaanCode: %d", result.DurationMs, result.HuaAnCode)
	var pretty json.RawMessage
	if err := json.Unmarshal(result.Body, &pretty); err != nil {
		t.Logf("华安响应(raw): %s", result.Body)
		return
	}
	indented, err := json.MarshalIndent(pretty, "", "  ")
	if err != nil {
		t.Logf("华安响应: %s", result.Body)
		return
	}
	t.Logf("华安响应:\n%s", indented)
}
