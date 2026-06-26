// Package repository 封装 GORM 数据访问：连接初始化、表迁移及各领域仓储。
package repository

import (
	"fmt"

	"github.com/huaan/insurance-bridge/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB 打开 MySQL 连接并 AutoMigrate 全部实体表；日志级别为 Warn，减少 SQL 噪音。
func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.AutoMigrate(
		&model.Channel{},
		&model.AdminUser{},
		&model.BankInfo{},
		&model.HuaAnSetting{},
	); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	if err := migrateChannelSchema(db); err != nil {
		return nil, fmt.Errorf("migrate channel schema: %w", err)
	}
	return db, nil
}

// migrateChannelSchema 移除 channel_code 唯一约束，允许多环境共用同一渠道编码。
func migrateChannelSchema(db *gorm.DB) error {
	m := db.Migrator()
	for _, name := range []string{"channel_code", "ChannelCode", "idx_channels_channel_code"} {
		if m.HasIndex(&model.Channel{}, name) {
			if err := m.DropIndex(&model.Channel{}, name); err != nil {
				return err
			}
		}
	}
	return nil
}
