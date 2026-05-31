// Package bootstrap 应用首次启动时的数据初始化逻辑。
package bootstrap

import (
	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/repository"
	"github.com/huaan/insurance-bridge/internal/service"
	"go.uber.org/zap"
)

// Seed 当 admin_users 表为空时，使用配置中的默认用户名/密码创建首个管理员（bcrypt 哈希存储）。
// 生产环境部署后应立即登录并修改密码。
func Seed(cfg *config.Config, adminRepo *repository.AdminRepo, log *zap.Logger) error {
	n, err := adminRepo.Count()
	if err != nil {
		return err
	}
	if n == 0 {
		hash, err := service.HashPassword(cfg.Admin.DefaultPassword)
		if err != nil {
			return err
		}
		if err := adminRepo.Create(&model.AdminUser{
			Username:     cfg.Admin.DefaultUsername,
			PasswordHash: hash,
		}); err != nil {
			return err
		}
		log.Info("default admin user created", zap.String("username", cfg.Admin.DefaultUsername))
	}
	return nil
}
