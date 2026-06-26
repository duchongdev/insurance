package service

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

	"github.com/google/uuid"
	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/pkg/sign"
	"go.uber.org/zap"
)

var (
	// ErrInvalidChannelConfig 渠道配置校验失败。
	ErrInvalidChannelConfig = errors.New("invalid channel config")
)

// insureNotifyRequired 华安回调必填字段（extra 可空）。
var insureNotifyRequired = []string{
	"productName", "productType", "insureTime", "policyNo", "policyId",
	"amount", "policyStartDate", "policyEndDate", "payStage", "channelCode", "timestamp",
}

// ValidateChannelConfig 校验渠道创建/更新时的必填项与回调地址格式。
func ValidateChannelConfig(ch *model.Channel) error {
	if ch == nil {
		return fmt.Errorf("%w: channel is nil", ErrInvalidChannelConfig)
	}
	if strings.TrimSpace(ch.ChannelCode) == "" {
		return fmt.Errorf("%w: channelCode required", ErrInvalidChannelConfig)
	}
	cb := strings.TrimSpace(ch.CallbackURL)
	if cb == "" {
		return fmt.Errorf("%w: callbackUrl required", ErrInvalidChannelConfig)
	}
	if !strings.HasPrefix(cb, "http://") && !strings.HasPrefix(cb, "https://") {
		return fmt.Errorf("%w: callbackUrl must start with http:// or https://", ErrInvalidChannelConfig)
	}
	return nil
}

// HandleInsureNotify 接收华安投保结果回调，转发至渠道 callbackUrl 并返回华安约定 JSON。
// reqSign 为华安 HTTP 头 sign（sign_enabled 时校验）；空则跳过验签。
func (s *ProxyService) HandleInsureNotify(ctx context.Context, rawBody []byte, reqSign string) ([]byte, int, error) {
	traceID := uuid.New().String()
	var body map[string]interface{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return s.callbackError(400, "invalid json"), http.StatusOK, nil
	}

	if err := validateInsureNotifyBody(body); err != nil {
		return s.callbackError(400, err.Error()), http.StatusOK, nil
	}

	channelCode, _ := body["channelCode"].(string)
	ch, err := s.channels.GetByCode(channelCode)
	if err != nil {
		s.log.Warn("insure notify channel not found",
			zap.String("traceId", traceID),
			zap.String("channelCode", channelCode),
		)
		return s.callbackError(403, "invalid channel"), http.StatusOK, nil
	}
	if strings.TrimSpace(ch.CallbackURL) == "" {
		return s.callbackError(500, "channel callbackUrl not configured"), http.StatusOK, nil
	}

	if s.cfg.HuaAn.SignEnabled && reqSign != "" {
		if !sign.VerifyHuaAnHeader(sign.MapFromJSON(body), s.cfg.HuaAn.Key, reqSign) {
			return s.callbackError(401, "sign verify failed"), http.StatusOK, nil
		}
	}

	outbound := sign.MapFromJSON(body)
	outbound["key"] = ch.ChannelKey
	outbound["sign"] = sign.Build(outbound, ch.ChannelKey)

	respBytes, err := s.postJSON(ctx, ch.CallbackURL, outbound)
	if err != nil {
		s.log.Error("insure notify forward failed",
			zap.String("traceId", traceID),
			zap.String("channelCode", channelCode),
			zap.String("callbackUrl", ch.CallbackURL),
			zap.Error(err),
		)
		return s.callbackError(502, "channel callback failed"), http.StatusOK, nil
	}

	var channelResp map[string]interface{}
	if err := json.Unmarshal(respBytes, &channelResp); err != nil {
		s.log.Warn("insure notify channel response not json",
			zap.String("traceId", traceID),
			zap.String("channelCode", channelCode),
		)
		return s.callbackError(502, "invalid channel response"), http.StatusOK, nil
	}

	code := parseIntCode(channelResp["code"])
	msg, _ := channelResp["message"].(string)
	s.log.Info("insure notify forwarded",
		zap.String("traceId", traceID),
		zap.String("channelCode", channelCode),
		zap.String("policyId", fmt.Sprint(body["policyId"])),
		zap.Int("channelRespCode", code),
	)
	if code != 200 {
		if msg == "" {
			msg = "channel callback rejected"
		}
		return s.callbackError(code, msg), http.StatusOK, nil
	}

	return s.callbackSuccess(), http.StatusOK, nil
}

func validateInsureNotifyBody(body map[string]interface{}) error {
	for _, k := range insureNotifyRequired {
		if _, ok := body[k]; !ok || body[k] == nil {
			return fmt.Errorf("%s required", k)
		}
		if str, ok := body[k].(string); ok && strings.TrimSpace(str) == "" {
			return fmt.Errorf("%s required", k)
		}
	}
	return nil
}

func (s *ProxyService) postJSON(ctx context.Context, url string, body map[string]interface{}) ([]byte, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	timeout := s.cfg.Server.UpstreamTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (s *ProxyService) callbackSuccess() []byte {
	b, _ := json.Marshal(map[string]interface{}{
		"code":    200,
		"message": "成功",
	})
	return b
}

func (s *ProxyService) callbackError(code int, message string) []byte {
	b, _ := json.Marshal(map[string]interface{}{
		"code":    code,
		"message": message,
	})
	return b
}

func parseIntCode(v interface{}) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case int64:
		return int(t)
	case json.Number:
		n, _ := t.Int64()
		return int(n)
	default:
		return 0
	}
}
