package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/repository"
	"gorm.io/gorm"
)

// ErrInvalidHuaAnConfig 华安配置校验失败。
var ErrInvalidHuaAnConfig = errors.New("invalid huaan config")

// ErrHuaAnConfigNotFound 华安配置不存在。
var ErrHuaAnConfigNotFound = errors.New("huaan config not found")

// HuaAnConfigDTO 管理后台华安配置读写结构。
type HuaAnConfigDTO struct {
	ID            uint64    `json:"id,omitempty"`
	Name          string    `json:"name"`
	EnvType       string    `json:"envType"`
	BaseURL       string    `json:"baseUrl"`
	ChannelCode   string    `json:"channelCode"`
	ChannelSecret string    `json:"channelSecret"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
}

// HuaAnConnectivityResult 华安连通性探测结果。
type HuaAnConnectivityResult struct {
	OK          bool   `json:"ok"`
	DurationMs  int64  `json:"durationMs"`
	UpstreamURL string `json:"upstreamUrl"`
	HTTPStatus  int    `json:"httpStatus,omitempty"`
	Message     string `json:"message"`
}

// HuaAnConfigService 华安上游配置管理：读库、写库并同步内存 Config。
type HuaAnConfigService struct {
	cfg  *config.Config
	repo *repository.HuaAnSettingRepo
}

// NewHuaAnConfigService 构造华安配置服务。
func NewHuaAnConfigService(cfg *config.Config, repo *repository.HuaAnSettingRepo) *HuaAnConfigService {
	return &HuaAnConfigService{cfg: cfg, repo: repo}
}

// GetByID 按主键查询华安配置。
func (s *HuaAnConfigService) GetByID(id uint64) (*HuaAnConfigDTO, error) {
	setting, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHuaAnConfigNotFound
		}
		return nil, err
	}
	return toHuaAnConfigDTO(setting), nil
}

// ResolveChannelCredentials 根据华安配置 ID 返回渠道编码与密钥。
func (s *HuaAnConfigService) ResolveChannelCredentials(huaAnSettingID uint64) (channelCode, channelSecret string, err error) {
	dto, err := s.GetByID(huaAnSettingID)
	if err != nil {
		return "", "", err
	}
	return dto.ChannelCode, dto.ChannelSecret, nil
}

// List 返回全部华安配置。
func (s *HuaAnConfigService) List() ([]HuaAnConfigDTO, error) {
	list, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	out := make([]HuaAnConfigDTO, len(list))
	for i := range list {
		out[i] = *toHuaAnConfigDTO(&list[i])
	}
	return out, nil
}

// GetActive 返回 is_active=1 的华安配置（历史字段，新逻辑不再依赖）。
func (s *HuaAnConfigService) GetActive() (*HuaAnConfigDTO, error) {
	setting, err := s.repo.GetActive()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHuaAnConfigNotFound
		}
		return nil, err
	}
	dto := toHuaAnConfigDTO(setting)
	return dto, nil
}

// Create 新建华安配置。
func (s *HuaAnConfigService) Create(dto *HuaAnConfigDTO) (*HuaAnConfigDTO, error) {
	if err := validateHuaAnConfigInput(dto, true); err != nil {
		return nil, err
	}
	setting := &model.HuaAnSetting{
		Name:          strings.TrimSpace(dto.Name),
		EnvType:       dto.EnvType,
		BaseURL:       strings.TrimSpace(dto.BaseURL),
		ChannelCode:   strings.TrimSpace(dto.ChannelCode),
		ChannelSecret: strings.TrimSpace(dto.ChannelSecret),
		IsActive:      false,
	}
	if err := s.repo.Create(setting); err != nil {
		return nil, err
	}
	return toHuaAnConfigDTO(setting), nil
}

// Update 更新华安配置；channelSecret 为空时保留原值。
func (s *HuaAnConfigService) Update(id uint64, dto *HuaAnConfigDTO) (*HuaAnConfigDTO, error) {
	if err := validateHuaAnConfigInput(dto, false); err != nil {
		return nil, err
	}
	current, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHuaAnConfigNotFound
		}
		return nil, err
	}
	secret := strings.TrimSpace(dto.ChannelSecret)
	if secret == "" {
		secret = current.ChannelSecret
	}
	if secret == "" {
		return nil, ErrInvalidHuaAnConfig
	}
	current.Name = strings.TrimSpace(dto.Name)
	current.EnvType = dto.EnvType
	current.BaseURL = strings.TrimSpace(dto.BaseURL)
	current.ChannelCode = strings.TrimSpace(dto.ChannelCode)
	current.ChannelSecret = secret
	if err := s.repo.Update(current); err != nil {
		return nil, err
	}
	return toHuaAnConfigDTO(current), nil
}

// Delete 删除华安配置。
func (s *HuaAnConfigService) Delete(id uint64) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrHuaAnConfigNotFound
		}
		return err
	}
	return nil
}

// TestConnectivity 对华安接口根地址做 HTTP 探测（类似 curl），仅验证网络是否可达。
func (s *HuaAnConfigService) TestConnectivity(ctx context.Context, baseURL string) (*HuaAnConnectivityResult, error) {
	probeURL, err := normalizeProbeURL(baseURL)
	if err != nil {
		return nil, err
	}

	timeout := s.cfg.Server.UpstreamTimeout
	if timeout <= 0 {
		timeout = 25 * time.Second
	}
	client := &http.Client{Timeout: timeout}

	start := time.Now()
	status, err := probeHTTP(ctx, client, probeURL)
	duration := time.Since(start).Milliseconds()
	if err != nil {
		return &HuaAnConnectivityResult{
			OK:          false,
			UpstreamURL: probeURL,
			Message:     probeErrorMessage(err),
		}, nil
	}

	return &HuaAnConnectivityResult{
		OK:          true,
		DurationMs:  duration,
		UpstreamURL: probeURL,
		HTTPStatus:  status,
		Message:     fmt.Sprintf("网络通畅，HTTP %d", status),
	}, nil
}

// TestConnectivityByID 按配置 ID 探测其接口地址。
func (s *HuaAnConfigService) TestConnectivityByID(ctx context.Context, id uint64) (*HuaAnConnectivityResult, error) {
	setting, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHuaAnConfigNotFound
		}
		return nil, err
	}
	return s.TestConnectivity(ctx, setting.BaseURL)
}

func validateHuaAnConfigInput(dto *HuaAnConfigDTO, requireSecret bool) error {
	dto.EnvType = strings.TrimSpace(dto.EnvType)
	dto.BaseURL = strings.TrimSpace(dto.BaseURL)
	dto.ChannelCode = strings.TrimSpace(dto.ChannelCode)
	dto.ChannelSecret = strings.TrimSpace(dto.ChannelSecret)

	if dto.EnvType != model.HuaAnEnvProd && dto.EnvType != model.HuaAnEnvTest {
		return ErrInvalidHuaAnConfig
	}
	if dto.BaseURL == "" || dto.ChannelCode == "" {
		return ErrInvalidHuaAnConfig
	}
	if requireSecret && dto.ChannelSecret == "" {
		return ErrInvalidHuaAnConfig
	}
	return nil
}

func normalizeProbeURL(baseURL string) (string, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return "", ErrInvalidHuaAnConfig
	}
	u, err := url.Parse(baseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", ErrInvalidHuaAnConfig
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", ErrInvalidHuaAnConfig
	}
	return strings.TrimRight(baseURL, "/") + "/", nil
}

func probeHTTP(ctx context.Context, client *http.Client, probeURL string) (int, error) {
	status, err := doProbe(ctx, client, probeURL, http.MethodHead)
	if err != nil {
		return 0, err
	}
	if status == http.StatusMethodNotAllowed || status == http.StatusNotImplemented {
		return doProbe(ctx, client, probeURL, http.MethodGet)
	}
	return status, nil
}

func doProbe(ctx context.Context, client *http.Client, probeURL, method string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, method, probeURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

func probeErrorMessage(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "连接超时，请检查接口地址或网络连通性"
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "连接超时，请检查接口地址或网络连通性"
	}
	return "无法连接上游，请检查接口地址是否正确"
}

func toHuaAnConfigDTO(s *model.HuaAnSetting) *HuaAnConfigDTO {
	return &HuaAnConfigDTO{
		ID:            s.ID,
		Name:          s.Name,
		EnvType:       s.EnvType,
		BaseURL:       s.BaseURL,
		ChannelCode:   s.ChannelCode,
		ChannelSecret: s.ChannelSecret,
		IsActive:      s.IsActive,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
