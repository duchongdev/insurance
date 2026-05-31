// Package service 实现渠道请求代理、业务数据抽取及管理后台业务逻辑。
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/pkg/cipher"
	"github.com/huaan/insurance-bridge/internal/pkg/pii"
	"github.com/huaan/insurance-bridge/internal/pkg/sign"
	"github.com/huaan/insurance-bridge/internal/repository"
	"go.uber.org/zap"
)

// ProxyService 渠道 API 代理核心：验签 → 解密三要素 → 换华安 key 重签 → 转发 → 抽取业务 → 加密响应 → 落审计日志。
type ProxyService struct {
	cfg        *config.Config
	log        *zap.Logger
	channels   *repository.ChannelRepo // 渠道密钥查询
	logs       *repository.LogRepo       // 接口审计日志
	extractor  *Extractor                // 从华安响应抽取业务实体
	pii        *pii.Transformer          // 三要素加解密
	httpClient *http.Client              // 调用华安上游
}

// NewProxyService 构造代理服务，httpClient 超时取自配置 upstream_timeout。
func NewProxyService(
	cfg *config.Config,
	log *zap.Logger,
	channels *repository.ChannelRepo,
	logs *repository.LogRepo,
	extractor *Extractor,
	piiTransformer *pii.Transformer,
) *ProxyService {
	return &ProxyService{
		cfg:       cfg,
		log:       log,
		channels:  channels,
		logs:      logs,
		extractor: extractor,
		pii:       piiTransformer,
		httpClient: &http.Client{
			Timeout: cfg.Server.UpstreamTimeout,
		},
	}
}

// Forward 处理单次渠道 API 调用。
// 约定：HTTP 状态码恒为 200，业务成败由 JSON body 的 code 字段表达（与华安/渠道文档一致）。
// 流程：解析 JSON → 校验 channelCode/key/sign → 解密 PII → 替换华安 key 并重签 → POST 上游 → 抽取落库 → 加密响应字段 → 写审计日志。
func (s *ProxyService) Forward(ctx context.Context, apiPath string, rawBody []byte) ([]byte, int, error) {
	traceID := uuid.New().String()
	var body map[string]interface{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return s.errorResponse(400, "invalid json"), http.StatusOK, nil
	}

	channelCode, _ := body["channelCode"].(string)
	if channelCode == "" {
		return s.errorResponse(400, "channelCode required"), http.StatusOK, nil
	}

	ch, err := s.channels.GetByCode(channelCode)
	if err != nil {
		return s.errorResponse(403, "invalid channel"), http.StatusOK, nil
	}

	// 校验渠道分配的 key（与 sign 中参与签名的 key 一致）
	reqKey, _ := body["key"].(string)
	if reqKey != ch.ChannelKey {
		return s.errorResponse(401, "invalid channel key"), http.StatusOK, nil
	}
	params := sign.MapFromJSON(body)
	if !sign.Verify(params, ch.ChannelKey) {
		return s.errorResponse(401, "sign verify failed"), http.StatusOK, nil
	}

	// 保留加密态副本，供日志落库脱敏（避免在 audit 表存明文三要素）
	channelBody := cloneMap(body)
	if err := s.pii.DecryptRequest(body); err != nil {
		s.log.Warn("pii decrypt failed", zap.String("traceId", traceID), zap.Error(err))
		return s.errorResponse(400, "pii decrypt failed"), http.StatusOK, nil
	}

	// 替换为华安密钥并重新签名后转发
	body["key"] = ch.HuaAnKey
	delete(body, "sign")
	body["sign"] = sign.Build(body, ch.HuaAnKey)

	upstreamURL := s.cfg.HuaAn.BaseURL + s.cfg.HuaAn.APIPath + apiPath
	reqBytes, _ := json.Marshal(body)

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(reqBytes))
	if err != nil {
		return s.errorResponse(500, "build upstream request failed"), http.StatusOK, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := s.httpClient.Do(req)
	duration := time.Since(start).Milliseconds()
	if err != nil {
		s.saveLog(traceID, channelCode, apiPath, string(rawBody), "", 0, duration, err.Error())
		if ctx.Err() != nil {
			return s.errorResponse(504, "upstream timeout"), http.StatusOK, nil
		}
		return s.errorResponse(502, "upstream unavailable"), http.StatusOK, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.errorResponse(500, "read upstream response failed"), http.StatusOK, err
	}

	var respBody map[string]interface{}
	_ = json.Unmarshal(respBytes, &respBody)
	huaanCode := 0
	if c, ok := respBody["code"].(float64); ok {
		huaanCode = int(c)
	}

	// 基于明文响应抽取用户/保单等实体（异步失败不影响渠道响应）
	s.extractor.Process(channelCode, apiPath, body, respBody)

	// 返回渠道前对响应中的三要素重新加密
	_ = s.pii.EncryptResponse(respBody)
	out, _ := json.Marshal(respBody)

	// 请求侧存脱敏副本，响应侧存华安原始 JSON
	s.saveLog(traceID, channelCode, apiPath, maskLogBody(channelBody), string(respBytes), huaanCode, duration, "")

	s.log.Info("proxy done",
		zap.String("traceId", traceID),
		zap.String("channelCode", channelCode),
		zap.String("apiPath", apiPath),
		zap.Int("huaanCode", huaanCode),
		zap.Int64("durationMs", duration),
	)

	return out, http.StatusOK, nil
}

// saveLog 写入 api_request_logs，错误被忽略以免阻塞主链路。
func (s *ProxyService) saveLog(traceID, channelCode, apiPath, req, resp string, code int, duration int64, errMsg string) {
	_ = s.logs.Create(&model.APIRequestLog{
		TraceID:      traceID,
		ChannelCode:  channelCode,
		APIPath:      apiPath,
		RequestBody:  req,
		ResponseBody: resp,
		HuaAnCode:    code,
		DurationMS:   duration,
		ErrorMessage: errMsg,
	})
}

// errorResponse 构造桥接层错误 JSON（code/message/data），供渠道解析。
func (s *ProxyService) errorResponse(code int, message string) []byte {
	b, _ := json.Marshal(map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    nil,
	})
	return b
}

// cloneMap 浅拷贝 map，用于在解密前保留渠道原始请求副本。
func cloneMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// maskLogBody 对请求体中的三要素字段做 Mask 脱敏后序列化，供审计日志存储。
func maskLogBody(m map[string]interface{}) string {
	cp := cloneMap(m)
	for _, f := range pii.RequestFields {
		if v, ok := cp[f].(string); ok && v != "" {
			cp[f] = cipher.Mask(v)
		}
	}
	b, _ := json.Marshal(cp)
	return string(b)
}

// GenerateChannelKey 为新渠道生成 32 位十六进制密钥（UUID 去横线）。
func GenerateChannelKey() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
