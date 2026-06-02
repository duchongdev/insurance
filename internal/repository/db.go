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
	); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}
