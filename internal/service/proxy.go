// Package service 实现渠道请求代理、业务数据抽取及管理后台业务逻辑。
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
	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/pkg/cipher"
	"github.com/huaan/insurance-bridge/internal/pkg/pii"
	rediscache "github.com/huaan/insurance-bridge/internal/pkg/redis"
	"github.com/huaan/insurance-bridge/internal/pkg/sign"
	"github.com/huaan/insurance-bridge/internal/repository"
	"go.uber.org/zap"
)

// ProxyService 渠道 API 代理核心：验签 → 解密三要素 → 换华安 key 重签 → 转发 → 抽取业务 → 加密响应 → 落审计日志。
type ProxyService struct {
	cfg           *config.Config
	log           *zap.Logger
	channels      *repository.ChannelRepo // 渠道密钥查询
	banks         *repository.BankRepo      // 银行信息表
	extractor     *Extractor                // 从华安响应抽取业务实体
	pii           *pii.Transformer          // 三要素加解密
	bankListCache *rediscache.BankListCache // 银行列表 Redis 缓存
	httpClient    *http.Client              // 调用华安上游
}

// NewProxyService 构造代理服务，httpClient 超时取自配置 upstream_timeout。
func NewProxyService(
	cfg *config.Config,
	log *zap.Logger,
	channels *repository.ChannelRepo,
	banks *repository.BankRepo,
	extractor *Extractor,
	piiTransformer *pii.Transformer,
	bankListCache *rediscache.BankListCache,
) *ProxyService {
	return &ProxyService{
		cfg:           cfg,
		log:           log,
		channels:      channels,
		banks:         banks,
		extractor:     extractor,
		pii:           piiTransformer,
		bankListCache: bankListCache,
		httpClient: &http.Client{
			Timeout: cfg.Server.UpstreamTimeout,
		},
	}
}

// Forward 处理单次渠道 API 调用。
// 约定：HTTP 状态码恒为 200，业务成败由 JSON body 的 code 字段表达（与华安/渠道文档一致）。
// 流程：解析 JSON → 校验 channelCode/key/sign → 解密 PII → 替换华安 key 并重签 → POST 上游 → 抽取落库 → 加密响应字段 → 写服务日志。
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

	// getBankList：Redis 缓存 → 数据库 → 华安上游
	if apiPath == bankListAPIPath {
		if out, ok := s.respondBankListFromCache(ctx, traceID, channelCode, string(rawBody)); ok {
			return out, http.StatusOK, nil
		}
		if out, ok := s.respondBankListFromDB(ctx, traceID, channelCode, string(rawBody)); ok {
			return out, http.StatusOK, nil
		}
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
		s.logAPICall(traceID, channelCode, apiPath, string(rawBody), "", 0, duration, "upstream", err.Error())
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

	if apiPath == bankListAPIPath {
		s.cacheBankList(ctx, channelCode, respBytes)
	}

	// 基于明文响应抽取用户/保单等实体（异步失败不影响渠道响应）
	s.extractor.Process(channelCode, apiPath, body, respBody)

	// 返回渠道前对响应中的三要素重新加密
	_ = s.pii.EncryptResponse(respBody)
	out, _ := json.Marshal(respBody)

	s.logAPICall(traceID, channelCode, apiPath, maskLogBody(channelBody), string(respBytes), huaanCode, duration, "upstream", "")

	return out, http.StatusOK, nil
}

const bankListAPIPath = "/getBankList"

// GetCachedBankList 返回 Redis 中华安原始响应 JSON；未命中时 data 为 nil。
func (s *ProxyService) GetCachedBankList(ctx context.Context, channelCode string) ([]byte, error) {
	return s.bankListCache.Get(ctx, channelCode)
}

// RefreshBankList 管理后台手动请求华安 /getBankList，使用调用方提供的渠道码与华安密钥。
func (s *ProxyService) RefreshBankList(ctx context.Context, channelCode, huaAnKey string) ([]byte, error) {
	if channelCode == "" || huaAnKey == "" {
		return nil, errors.New("channelCode and huaAnKey required")
	}

	traceID := uuid.New().String()
	body := map[string]interface{}{
		"timestamp":   fmt.Sprintf("%d", time.Now().UnixMilli()),
		"channelCode": channelCode,
		"key":         huaAnKey,
	}
	body["sign"] = sign.Build(body, huaAnKey)

	upstreamURL := s.cfg.HuaAn.BaseURL + s.cfg.HuaAn.APIPath + bankListAPIPath
	reqBytes, _ := json.Marshal(body)

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("build upstream request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := s.httpClient.Do(req)
	duration := time.Since(start).Milliseconds()
	if err != nil {
		s.logAPICall(traceID, channelCode, bankListAPIPath, maskAdminLogBody(body), "", 0, duration, "admin-upstream", err.Error())
		if ctx.Err() != nil {
			return nil, errors.New("upstream timeout")
		}
		return nil, fmt.Errorf("upstream unavailable: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read upstream response: %w", err)
	}

	var respBody map[string]interface{}
	_ = json.Unmarshal(respBytes, &respBody)
	huaanCode := 0
	if c, ok := respBody["code"].(float64); ok {
		huaanCode = int(c)
	}

	s.logAPICall(traceID, channelCode, bankListAPIPath, maskAdminLogBody(body), string(respBytes), huaanCode, duration, "admin-upstream", "")

	s.cacheBankList(ctx, channelCode, respBytes)

	return respBytes, nil
}

func (s *ProxyService) cacheBankList(ctx context.Context, channelCode string, respBytes []byte) {
	if !huaAnResponseOK(respBytes) {
		return
	}
	if err := s.bankListCache.Set(ctx, channelCode, respBytes); err != nil {
		s.log.Warn("bank list cache set failed",
			zap.String("channelCode", channelCode),
			zap.Error(err),
		)
	}
	if err := ReplaceBanksFromResponse(s.banks, respBytes); err != nil {
		s.log.Warn("bank list db replace failed",
			zap.String("channelCode", channelCode),
			zap.Error(err),
		)
	}
}

// respondBankListFromDB 命中数据库时构造渠道响应、回填 Redis 并写审计日志。
func (s *ProxyService) respondBankListFromDB(ctx context.Context, traceID, channelCode, reqLog string) ([]byte, bool) {
	banks, err := s.banks.ListEnabled()
	if err != nil {
		s.log.Warn("bank list db query failed",
			zap.String("traceId", traceID),
			zap.String("channelCode", channelCode),
			zap.Error(err),
		)
		return nil, false
	}
	if len(banks) == 0 {
		return nil, false
	}

	raw, err := BuildBankListResponseBytes(banks)
	if err != nil {
		return nil, false
	}
	if err := s.bankListCache.Set(ctx, channelCode, raw); err != nil {
		s.log.Warn("bank list cache set from db failed",
			zap.String("channelCode", channelCode),
			zap.Error(err),
		)
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(raw, &respBody); err != nil {
		return nil, false
	}
	_ = s.pii.EncryptResponse(respBody)
	out, _ := json.Marshal(respBody)

	s.logAPICall(traceID, channelCode, bankListAPIPath, reqLog, string(raw), 200, 0, "db", "")
	return out, true
}

// respondBankListFromCache 命中缓存时构造渠道响应并写服务日志。
func (s *ProxyService) respondBankListFromCache(ctx context.Context, traceID, channelCode, reqLog string) ([]byte, bool) {
	cached, err := s.bankListCache.Get(ctx, channelCode)
	if err != nil {
		s.log.Warn("bank list cache get failed",
			zap.String("traceId", traceID),
			zap.String("channelCode", channelCode),
			zap.Error(err),
		)
		return nil, false
	}
	if len(cached) == 0 {
		return nil, false
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(cached, &respBody); err != nil {
		return nil, false
	}
	huaanCode := 0
	if c, ok := respBody["code"].(float64); ok {
		huaanCode = int(c)
	}
	_ = s.pii.EncryptResponse(respBody)
	out, _ := json.Marshal(respBody)

	s.logAPICall(traceID, channelCode, bankListAPIPath, reqLog, string(cached), huaanCode, 0, "cache", "")
	return out, true
}

func huaAnResponseOK(respBytes []byte) bool {
	var resp map[string]interface{}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return false
	}
	code, ok := resp["code"].(float64)
	return ok && int(code) == 200
}

// maskAdminLogBody 管理后台直连请求日志中掩码密钥字段。
func maskAdminLogBody(m map[string]interface{}) string {
	cp := cloneMap(m)
	if _, ok := cp["key"]; ok {
		cp["key"] = "***"
	}
	b, _ := json.Marshal(cp)
	return string(b)
}

// logAPICall 将渠道 API 调用详情写入服务日志；摘要 info，请求/响应 debug。
func (s *ProxyService) logAPICall(traceID, channelCode, apiPath, req, resp string, code int, duration int64, source, errMsg string) {
	fields := []zap.Field{
		zap.String("traceId", traceID),
		zap.String("channelCode", channelCode),
		zap.String("apiPath", apiPath),
		zap.Int("huaanCode", code),
		zap.Int64("durationMs", duration),
		zap.String("source", source),
	}
	if errMsg != "" {
		fields = append(fields, zap.String("error", errMsg))
	}
	s.log.Info("api call", fields...)
	s.log.Debug("api call detail",
		zap.String("traceId", traceID),
		zap.String("channelCode", channelCode),
		zap.String("apiPath", apiPath),
		zap.String("source", source),
		zap.String("request", req),
		zap.String("response", resp),
	)
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
