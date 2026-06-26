// Package bootstrap 应用首次启动时的数据初始化逻辑。
package bootstrap

import (
	"errors"

	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/repository"
	"github.com/huaan/insurance-bridge/internal/service"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed 确保内置管理员账号存在且密码与 DefaultAdminPassword 一致（bcrypt 哈希存储）。
func Seed(adminRepo *repository.AdminRepo, log *zap.Logger) error {
	hash, err := service.HashPassword(DefaultAdminPassword)
	if err != nil {
		return err
	}

	u, err := adminRepo.GetByUsername(DefaultAdminUsername)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := adminRepo.Create(&model.AdminUser{
				Username:     DefaultAdminUsername,
				PasswordHash: hash,
			}); err != nil {
				return err
			}
			log.Info("default admin user created", zap.String("username", DefaultAdminUsername))
			return nil
		}
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(DefaultAdminPassword)) != nil {
		if err := adminRepo.UpdatePasswordHash(u.ID, hash); err != nil {
			return err
		}
		log.Info("built-in admin password synced", zap.String("username", DefaultAdminUsername))
	}
	return nil
}
