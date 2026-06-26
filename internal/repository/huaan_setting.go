package repository

import (
	"errors"

	"github.com/huaan/insurance-bridge/internal/model"
	"gorm.io/gorm"
)

// HuaAnSettingRepo 华安上游配置表访问。
type HuaAnSettingRepo struct {
	db *gorm.DB
}

// NewHuaAnSettingRepo 构造华安配置仓储。
func NewHuaAnSettingRepo(db *gorm.DB) *HuaAnSettingRepo {
	return &HuaAnSettingRepo{db: db}
}

// List 列出全部华安配置，生产环境在前、同类型按 id 升序。
func (r *HuaAnSettingRepo) List() ([]model.HuaAnSetting, error) {
	var list []model.HuaAnSetting
	err := r.db.Order("env_type asc, id asc").Find(&list).Error
	return list, err
}

// GetByID 按主键查询。
func (r *HuaAnSettingRepo) GetByID(id uint64) (*model.HuaAnSetting, error) {
	var s model.HuaAnSetting
	err := r.db.First(&s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &s, err
}

// GetActive 返回当前启用的华安配置。
func (r *HuaAnSettingRepo) GetActive() (*model.HuaAnSetting, error) {
	var s model.HuaAnSetting
	err := r.db.Where("is_active = ?", true).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &s, err
}

// Create 插入新配置。
func (r *HuaAnSettingRepo) Create(s *model.HuaAnSetting) error {
	return r.db.Create(s).Error
}

// Update 按主键 Save 更新。
func (r *HuaAnSettingRepo) Update(s *model.HuaAnSetting) error {
	return r.db.Save(s).Error
}

// Delete 按主键删除。
func (r *HuaAnSettingRepo) Delete(id uint64) error {
	return r.db.Delete(&model.HuaAnSetting{}, id).Error
}

// SetActive 将指定配置设为唯一启用项。
func (r *HuaAnSettingRepo) SetActive(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.HuaAnSetting{}).Where("is_active = ?", true).Update("is_active", false).Error; err != nil {
			return err
		}
		res := tx.Model(&model.HuaAnSetting{}).Where("id = ?", id).Update("is_active", true)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

// Count 返回配置总数。
func (r *HuaAnSettingRepo) Count() (int64, error) {
	var n int64
	err := r.db.Model(&model.HuaAnSetting{}).Count(&n).Error
	return n, err
}
