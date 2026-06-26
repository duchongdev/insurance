package repository

import (
	"errors"

	"github.com/huaan/insurance-bridge/internal/model"
	"gorm.io/gorm"
)

// AdminRepo 管理后台用户表访问。
type AdminRepo struct {
	db *gorm.DB
}

// NewAdminRepo 构造管理员仓储。
func NewAdminRepo(db *gorm.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

// GetByUsername 按登录名查询管理员，不存在返回 gorm.ErrRecordNotFound。
func (r *AdminRepo) GetByUsername(username string) (*model.AdminUser, error) {
	var u model.AdminUser
	err := r.db.Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &u, err
}

// Create 创建管理员账号。
func (r *AdminRepo) Create(u *model.AdminUser) error {
	return r.db.Create(u).Error
}

// UpdatePasswordHash 更新管理员密码哈希。
func (r *AdminRepo) UpdatePasswordHash(id uint64, hash string) error {
	return r.db.Model(&model.AdminUser{}).Where("id = ?", id).Update("password_hash", hash).Error
}

