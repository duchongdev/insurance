package huaan

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/pkg/sign"
	"go.uber.org/zap"
)

// ErrTimeout 调用华安超时（context 取消或 HTTP 客户端超时）。
var ErrTimeout = errors.New("upstream timeout")

// ErrUnavailable 华安不可达或连接失败。
var ErrUnavailable = errors.New("upstream unavailable")

// Client 华安上游 HTTP 客户端：写入 key/sign、POST、返回原始 JSON。
type Client struct {
	cfg        *config.Config
	log        *zap.Logger
	httpClient *http.Client
}

// NewClient 构造华安客户端；httpClient 为 nil 时使用 cfg.Server.UpstreamTimeout 创建默认客户端。
func NewClient(cfg *config.Config, log *zap.Logger, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: cfg.Server.UpstreamTimeout}
	}
	if log == nil {
		log = zap.NewNop()
	}
	return &Client{cfg: cfg, log: log, httpClient: httpClient}
}

// CallResult 单次华安调用的结果。
type CallResult struct {
	Body       []byte // 华安原始响应 JSON（PII 为明文）
	DurationMs int64
	HuaAnCode  int // 响应 JSON 中 code 字段，解析失败时为 0
}

// Call 向华安 POST 请求。body 会被原地修改：始终写入 config 中的 key；开启签名时计算 sign，否则 sign 为空字符串。
func (c *Client) Call(ctx context.Context, apiPath string, body map[string]interface{}) (*CallResult, error) {
	body["key"] = c.cfg.HuaAn.Key
	if c.cfg.HuaAn.SignEnabled {
		delete(body, "sign")
		body["sign"] = sign.Build(body, c.cfg.HuaAn.Key)
	} else {
		body["sign"] = ""
	}

	upstreamPath := ResolveUpstreamPath(c.cfg.HuaAn.APIPath, apiPath)
	upstreamURL := strings.TrimRight(c.cfg.HuaAn.BaseURL, "/") + upstreamPath
	reqBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	duration := time.Since(start).Milliseconds()
	if err != nil {
		if ctx.Err() != nil {
			return nil, ErrTimeout
		}
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	result := &CallResult{
		Body:       respBytes,
		DurationMs: duration,
		HuaAnCode:  ParseCode(respBytes),
	}
	return result, nil
}

// BuildRequestBody 构造带 timestamp、channelCode 的华安请求体（key/sign 由 Call 注入）。
func BuildRequestBody(channelCode string, fields map[string]interface{}) map[string]interface{} {
	body := map[string]interface{}{
		"timestamp":   fmt.Sprintf("%d", time.Now().UnixMilli()),
		"channelCode": channelCode,
	}
	for k, v := range fields {
		body[k] = v
	}
	return body
}

// ParseCode 从华安响应 JSON 解析 code 字段。
func ParseCode(respBytes []byte) int {
	var resp map[string]interface{}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return 0
	}
	if code, ok := resp["code"].(float64); ok {
		return int(code)
	}
	return 0
}

// ResponseOK 判断华安响应是否为业务成功（code==200）。
func ResponseOK(respBytes []byte) bool {
	return ParseCode(respBytes) == 200
}
