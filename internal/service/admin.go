package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials 登录或 JWT 校验失败时返回。
var ErrInvalidCredentials = errors.New("invalid credentials")

// AdminService 管理后台业务：认证、渠道 CRUD 与银行信息查询。
type AdminService struct {
	admins    *repository.AdminRepo
	channels  *repository.ChannelRepo
	banks     *repository.BankRepo
	jwtSecret []byte
}

// NewAdminService 注入各仓储与 JWT 签名密钥。
func NewAdminService(
	admins *repository.AdminRepo,
	channels *repository.ChannelRepo,
	banks *repository.BankRepo,
	jwtSecret string,
) *AdminService {
	return &AdminService{
		admins:    admins,
		channels:  channels,
		banks:     banks,
		jwtSecret: []byte(jwtSecret),
	}
}

// Login 校验用户名密码，成功签发 HS256 JWT（sub=username，有效期 24 小时）。
func (s *AdminService) Login(username, password string) (string, error) {
	u, err := s.admins.GetByUsername(username)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	return token.SignedString(s.jwtSecret)
}

// ParseToken 解析并校验 JWT，返回 subject（用户名）。
func (s *AdminService) ParseToken(tokenStr string) (string, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !t.Valid {
		return "", ErrInvalidCredentials
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidCredentials
	}
	sub, _ := claims["sub"].(string)
	return sub, nil
}

// HashPassword 使用 bcrypt 生成密码哈希，供 seed 与改密使用。
func HashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(b), err
}

// ListChannels 分页查询渠道列表。
func (s *AdminService) ListChannels(page, size int) ([]model.Channel, int64, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	return s.channels.List((page-1)*size, size)
}

// CreateChannel 创建渠道；ChannelKey 为空时调用 GenerateChannelKey 自动生成。
func (s *AdminService) CreateChannel(ch *model.Channel) error {
	if err := ValidateChannelConfig(ch); err != nil {
		return err
	}
	if ch.HuaAnSettingID > 0 {
		exists, err := s.channels.ExistsByHuaAnSettingID(ch.HuaAnSettingID)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("%w: 该华安配置已关联渠道", ErrInvalidChannelConfig)
		}
	}
	if ch.Status == 0 {
		ch.Status = 1
	}
	if ch.Status == 1 {
		active, err := s.channels.HasActiveByCode(ch.ChannelCode, 0)
		if err != nil {
			return err
		}
		if active {
			return fmt.Errorf("%w: 渠道编码 %s 已有启用的渠道，请先禁用后再创建", ErrInvalidChannelConfig, ch.ChannelCode)
		}
	}
	if ch.ChannelKey == "" {
		ch.ChannelKey = GenerateChannelKey()
	}
	return s.channels.Create(ch)
}

// UpdateChannel 更新渠道可编辑字段；channelKey 为空时保留原值。
func (s *AdminService) UpdateChannel(id uint64, input *model.Channel) (*model.Channel, error) {
	existing, err := s.channels.GetByID(id)
	if err != nil {
		return nil, err
	}
	existing.ChannelName = strings.TrimSpace(input.ChannelName)
	existing.CallbackURL = strings.TrimSpace(input.CallbackURL)
	if input.Status == 0 || input.Status == 1 {
		existing.Status = input.Status
	}
	if strings.TrimSpace(input.ChannelKey) != "" {
		existing.ChannelKey = strings.TrimSpace(input.ChannelKey)
	}
	if err := ValidateChannelConfig(existing); err != nil {
		return nil, err
	}
	if existing.Status == 1 {
		active, err := s.channels.HasActiveByCode(existing.ChannelCode, existing.ID)
		if err != nil {
			return nil, err
		}
		if active {
			return nil, fmt.Errorf("%w: 渠道编码 %s 已有其他启用的渠道", ErrInvalidChannelConfig, existing.ChannelCode)
		}
	}
	if err := s.channels.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteChannel 按主键删除渠道。
func (s *AdminService) DeleteChannel(id uint64) error {
	return s.channels.Delete(id)
}

// GetChannel 按 ID 查询单条渠道。
func (s *AdminService) GetChannel(id uint64) (*model.Channel, error) {
	return s.channels.GetByID(id)
}

// ListBanks 分页查询银行信息；status 为 0 表示不限。
func (s *AdminService) ListBanks(status int8, page, size int) ([]model.BankInfo, int64, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	return s.banks.List(status, (page-1)*size, size)
}
