package repository

import (
	"github.com/huaan/insurance-bridge/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BankRepo 银行信息表 bank_info_t 访问。
type BankRepo struct {
	db *gorm.DB
}

// NewBankRepo 构造银行仓储。
func NewBankRepo(db *gorm.DB) *BankRepo {
	return &BankRepo{db: db}
}

// List 分页查询银行，status 非 0 时按状态过滤。
func (r *BankRepo) List(status int8, offset, limit int) ([]model.BankInfo, int64, error) {
	var list []model.BankInfo
	var total int64
	q := r.db.Model(&model.BankInfo{})
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset(offset).Limit(limit).Order("id asc").Find(&list).Error
	return list, total, err
}

// ListAll 返回全部银行，按 id 升序。
func (r *BankRepo) ListAll() ([]model.BankInfo, error) {
	var list []model.BankInfo
	err := r.db.Order("id asc").Find(&list).Error
	return list, err
}

// ListEnabled 返回启用中的银行（status=1）。
func (r *BankRepo) ListEnabled() ([]model.BankInfo, error) {
	var list []model.BankInfo
	err := r.db.Where("status = ?", 1).Order("id asc").Find(&list).Error
	return list, err
}

// ReplaceAll 全量覆盖 bank_info_t：先清空再批量写入。
func (r *BankRepo) ReplaceAll(banks []model.BankInfo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("1 = 1").Delete(&model.BankInfo{}).Error; err != nil {
			return err
		}
		if len(banks) == 0 {
			return nil
		}
		return tx.CreateInBatches(banks, 100).Error
	})
}

// Upsert 按 bank_code 冲突时更新名称、卡类型与状态。
func (r *BankRepo) Upsert(b *model.BankInfo) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "bank_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"bank_name", "debit_card", "credit_card", "status", "update_time"}),
	}).Create(b).Error
}
