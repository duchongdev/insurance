package repository

import (
	"errors"

	"github.com/huaan/insurance-bridge/internal/model"
	"gorm.io/gorm"
)

// ErrChannelNotFound 渠道不存在或已禁用。
var ErrChannelNotFound = errors.New("channel not found")

// ChannelRepo 渠道配置表访问。
type ChannelRepo struct {
	db *gorm.DB
}

// NewChannelRepo 构造渠道仓储。
func NewChannelRepo(db *gorm.DB) *ChannelRepo {
	return &ChannelRepo{db: db}
}

// GetByCode 按渠道编码查询启用中的渠道（status=1），供代理验签使用。
func (r *ChannelRepo) GetByCode(code string) (*model.Channel, error) {
	var ch model.Channel
	err := r.db.Where("channel_code = ? AND status = 1", code).First(&ch).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrChannelNotFound
	}
	return &ch, err
}

// List 分页列出全部渠道（含禁用），按 id 降序。
func (r *ChannelRepo) List(offset, limit int) ([]model.Channel, int64, error) {
	var list []model.Channel
	var total int64
	q := r.db.Model(&model.Channel{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset(offset).Limit(limit).Order("id desc").Find(&list).Error
	return list, total, err
}

// Create 插入新渠道。
func (r *ChannelRepo) Create(ch *model.Channel) error {
	return r.db.Create(ch).Error
}

// Update 按主键 Save 更新渠道（全字段覆盖）。
func (r *ChannelRepo) Update(ch *model.Channel) error {
	return r.db.Save(ch).Error
}

// Delete 按主键硬删除渠道。
func (r *ChannelRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Channel{}, id).Error
}

// GetByID 按主键查询渠道。
func (r *ChannelRepo) GetByID(id uint64) (*model.Channel, error) {
	var ch model.Channel
	err := r.db.First(&ch, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrChannelNotFound
	}
	return &ch, err
}
