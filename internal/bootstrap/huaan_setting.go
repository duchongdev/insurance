package bootstrap

import (
	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/repository"
	"go.uber.org/zap"
)

// SyncHuaAnSetting 启动时迁移/初始化华安配置，并将首条记录写入运行时 Config 作为兜底。
func SyncHuaAnSetting(cfg *config.Config, repo *repository.HuaAnSettingRepo, log *zap.Logger) error {
	list, err := repo.List()
	if err != nil {
		return err
	}

	if len(list) == 0 {
		s := &model.HuaAnSetting{
			Name:          "默认测试环境",
			EnvType:       model.HuaAnEnvTest,
			BaseURL:       cfg.HuaAn.BaseURL,
			ChannelCode:   cfg.HuaAn.ChannelCode,
			ChannelSecret: cfg.HuaAn.Key,
			IsActive:      false,
		}
		if err := repo.Create(s); err != nil {
			return err
		}
		log.Info("huaan settings initialized from config file")
		applyHuaAnSetting(cfg, s)
		return nil
	}

	for i := range list {
		s := &list[i]
		if s.EnvType == "" {
			s.EnvType = model.HuaAnEnvTest
			if err := repo.Update(s); err != nil {
				return err
			}
		}
	}

	applyHuaAnSetting(cfg, &list[0])
	return nil
}

// applyHuaAnSetting 将库内华安配置写入运行时 Config。
func applyHuaAnSetting(cfg *config.Config, s *model.HuaAnSetting) {
	cfg.HuaAn.BaseURL = s.BaseURL
	cfg.HuaAn.ChannelCode = s.ChannelCode
	cfg.HuaAn.Key = s.ChannelSecret
}
