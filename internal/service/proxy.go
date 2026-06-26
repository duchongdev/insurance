// Package service 实现渠道请求代理及管理后台业务逻辑。
package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/huaan"
	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/pkg/pii"
	rediscache "github.com/huaan/insurance-bridge/internal/pkg/redis"
	"github.com/huaan/insurance-bridge/internal/pkg/sign"
	"github.com/huaan/insurance-bridge/internal/repository"
	"go.uber.org/zap"
)

// ProxyService 渠道 API 代理：渠道层（验签、PII）编排 + 华安层（HuaAnClient）转发。
type ProxyService struct {
	cfg           *config.Config
	log           *zap.Logger
	channels      *repository.ChannelRepo
	banks         *repository.BankRepo
	huaanSettings *repository.HuaAnSettingRepo
	pii           *pii.Transformer
	bankListCache *rediscache.BankListCache
	huaan         *huaan.Client
}

// NewProxyService 构造代理服务；huaanClient 为 nil 时自动创建默认华安客户端。
func NewProxyService(
	cfg *config.Config,
	log *zap.Logger,
	channels *repository.ChannelRepo,
	banks *repository.BankRepo,
	huaanSettings *repository.HuaAnSettingRepo,
	piiTransformer *pii.Transformer,
	bankListCache *rediscache.BankListCache,
	huaanClient *huaan.Client,
) *ProxyService {
	if huaanClient == nil {
		huaanClient = huaan.NewClient(cfg, log, nil)
	}
	return &ProxyService{
		cfg:           cfg,
		log:           log,
		channels:      channels,
		banks:         banks,
		huaanSettings: huaanSettings,
		pii:           piiTransformer,
		bankListCache: bankListCache,
		huaan:         huaanClient,
	}
}

// HuaAnClient 返回华安上游客户端，供集成测试或管理任务直接调用。
func (s *ProxyService) HuaAnClient() *huaan.Client {
	return s.huaan
}

// Forward 处理单次渠道 API 调用。
// 约定：HTTP 状态码恒为 200，业务成败由 JSON body 的 code 字段表达（与华安/渠道文档一致）。
// 流程：解析 JSON → 校验 channelCode/key/sign → 解密 PII → 华安 Client.Call（key/sign 由 huaan 配置注入）→ 加密响应字段。
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

	reqKey, _ := body["key"].(string)
	if reqKey != ch.ChannelKey {
		return s.errorResponse(401, "invalid channel key"), http.StatusOK, nil
	}
	params := sign.MapFromJSON(body)
	if !sign.Verify(params, ch.ChannelKey) {
		return s.errorResponse(401, "sign verify failed"), http.StatusOK, nil
	}

	if apiPath == huaan.BankListPath {
		if out, ok := s.respondBankListFromCache(ctx, traceID, channelCode); ok {
			return out, http.StatusOK, nil
		}
		if out, ok := s.respondBankListFromDB(ctx, traceID, channelCode); ok {
			return out, http.StatusOK, nil
		}
	}

	if err := s.pii.DecryptRequest(body); err != nil {
		s.log.Warn("pii decrypt failed", zap.String("traceId", traceID), zap.Error(err))
		return s.errorResponse(400, "pii decrypt failed"), http.StatusOK, nil
	}

	result, err := s.callHuaAn(ctx, ch, apiPath, body)
	if err != nil {
		s.log.Warn("upstream failed",
			zap.String("traceId", traceID),
			zap.String("channelCode", channelCode),
			zap.String("apiPath", apiPath),
			zap.Error(err),
		)
		if errors.Is(err, huaan.ErrTimeout) {
			return s.errorResponse(504, "upstream timeout"), http.StatusOK, nil
		}
		return s.errorResponse(502, "upstream unavailable"), http.StatusOK, err
	}

	var respBody map[string]interface{}
	_ = json.Unmarshal(result.Body, &respBody)

	if apiPath == huaan.BankListPath {
		s.cacheBankList(ctx, channelCode, result.Body)
	}

	_ = s.pii.EncryptResponse(respBody)
	out, _ := json.Marshal(respBody)

	s.log.Info("proxy done",
		zap.String("traceId", traceID),
		zap.String("channelCode", channelCode),
		zap.String("apiPath", apiPath),
		zap.Int("huaanCode", result.HuaAnCode),
		zap.Int64("durationMs", result.DurationMs),
	)

	return out, http.StatusOK, nil
}

// GetCachedBankList 返回 Redis 中华安原始响应 JSON；未命中时 data 为 nil。
func (s *ProxyService) GetCachedBankList(ctx context.Context, channelCode string) ([]byte, error) {
	return s.bankListCache.Get(ctx, channelCode)
}

// RefreshBankList 使用指定华安配置请求上游 getBankList。
func (s *ProxyService) RefreshBankList(ctx context.Context, huaAnSettingID uint64) ([]byte, error) {
	if huaAnSettingID == 0 {
		return nil, errors.New("huaanSettingId required")
	}
	setting, err := s.huaanSettings.GetByID(huaAnSettingID)
	if err != nil {
		return nil, errors.New("huaan setting not found")
	}

	body := huaan.BuildRequestBody(setting.ChannelCode, nil)
	result, err := s.callHuaAnWithSetting(ctx, setting, huaan.BankListPath, body)
	if err != nil {
		if errors.Is(err, huaan.ErrTimeout) {
			return nil, errors.New("upstream timeout")
		}
		return nil, err
	}

	s.cacheBankList(ctx, setting.ChannelCode, result.Body)
	return result.Body, nil
}

func (s *ProxyService) callHuaAn(ctx context.Context, ch *model.Channel, apiPath string, body map[string]interface{}) (*huaan.CallResult, error) {
	if ch != nil && ch.HuaAnSettingID > 0 {
		setting, err := s.huaanSettings.GetByID(ch.HuaAnSettingID)
		if err == nil {
			return s.callHuaAnWithSetting(ctx, setting, apiPath, body)
		}
	}
	return s.huaan.Call(ctx, apiPath, body)
}

func (s *ProxyService) callHuaAnWithSetting(ctx context.Context, setting *model.HuaAnSetting, apiPath string, body map[string]interface{}) (*huaan.CallResult, error) {
	cfg := *s.cfg
	cfg.HuaAn = config.HuaAnConfig{
		BaseURL:     setting.BaseURL,
		APIPath:     s.cfg.HuaAn.APIPath,
		ChannelCode: setting.ChannelCode,
		Key:         setting.ChannelSecret,
		SignEnabled: s.cfg.HuaAn.SignEnabled,
	}
	client := huaan.NewClient(&cfg, s.log, nil)
	return client.Call(ctx, apiPath, body)
}

func (s *ProxyService) cacheBankList(ctx context.Context, channelCode string, respBytes []byte) {
	if !huaan.ResponseOK(respBytes) {
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

// respondBankListFromDB 命中数据库时构造渠道响应并回填 Redis。
func (s *ProxyService) respondBankListFromDB(ctx context.Context, traceID, channelCode string) ([]byte, bool) {
	banks, err := s.banks.ListAll()
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

	s.log.Info("bank list db hit",
		zap.String("traceId", traceID),
		zap.String("channelCode", channelCode),
		zap.Int("count", len(banks)),
	)
	return out, true
}

// respondBankListFromCache 命中缓存时构造渠道响应。
func (s *ProxyService) respondBankListFromCache(ctx context.Context, traceID, channelCode string) ([]byte, bool) {
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
	huaanCode := huaan.ParseCode(cached)
	_ = s.pii.EncryptResponse(respBody)
	out, _ := json.Marshal(respBody)

	s.log.Info("bank list cache hit",
		zap.String("traceId", traceID),
		zap.String("channelCode", channelCode),
		zap.Int("huaanCode", huaanCode),
	)
	return out, true
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

// GenerateChannelKey 为新渠道生成 32 位十六进制密钥（UUID 去横线）。
func GenerateChannelKey() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
