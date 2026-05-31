package repository

import (
	"time"

	"github.com/huaan/insurance-bridge/internal/model"
	"gorm.io/gorm"
)

// LogRepo 接口审计日志 api_request_logs 表访问。
type LogRepo struct {
	db *gorm.DB
}

// NewLogRepo 构造日志仓储。
func NewLogRepo(db *gorm.DB) *LogRepo {
	return &LogRepo{db: db}
}

// Create 插入一条代理调用审计记录。
func (r *LogRepo) Create(log *model.APIRequestLog) error {
	return r.db.Create(log).Error
}

// List 分页查询日志，可按 channelCode、apiPath 过滤。
func (r *LogRepo) List(channelCode, apiPath string, offset, limit int) ([]model.APIRequestLog, int64, error) {
	var list []model.APIRequestLog
	var total int64
	q := r.db.Model(&model.APIRequestLog{})
	if channelCode != "" {
		q = q.Where("channel_code = ?", channelCode)
	}
	if apiPath != "" {
		q = q.Where("api_path = ?", apiPath)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset(offset).Limit(limit).Order("id desc").Find(&list).Error
	return list, total, err
}

// DeleteBefore 物理删除 created_at 早于 t 的记录，返回删除行数（供定时清理任务使用）。
func (r *LogRepo) DeleteBefore(t time.Time) (int64, error) {
	res := r.db.Where("created_at < ?", t).Delete(&model.APIRequestLog{})
	return res.RowsAffected, res.Error
}
