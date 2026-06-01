package repository

import (
	"github.com/huaan/insurance-bridge/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BusinessRepo 用户/保单/签约/产品/报价等业务快照表访问，抽取逻辑通过 Upsert 保持幂等。
type BusinessRepo struct {
	db *gorm.DB
}

// NewBusinessRepo 构造业务数据仓储。
func NewBusinessRepo(db *gorm.DB) *BusinessRepo {
	return &BusinessRepo{db: db}
}

// UpsertUser 按 (channel_code, user_id) 冲突时更新三要素密文。
func (r *BusinessRepo) UpsertUser(u *model.UserRecord) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "channel_code"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"phone_enc", "name_enc", "id_card_enc", "updated_at"}),
	}).Create(u).Error
}

// UpsertPolicy 按 policy_id 冲突时更新保单字段。
func (r *BusinessRepo) UpsertPolicy(p *model.PolicyRecord) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "policy_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"policy_no", "policy_status", "product_code", "user_id", "pay_premium", "policy_start_date", "policy_end_date", "updated_at"}),
	}).Create(p).Error
}

// UpsertSign 按 sign_id 冲突时更新签约相关字段。
func (r *BusinessRepo) UpsertSign(s *model.SignRecord) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "sign_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"policy_id", "sign_url", "bank_code", "updated_at"}),
	}).Create(s).Error
}

// UpsertProduct 按 (channel_code, product_code) 冲突时更新产品信息。
func (r *BusinessRepo) UpsertProduct(p *model.ProductRecord) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "channel_code"}, {Name: "product_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"product_name", "product_id", "price", "updated_at"}),
	}).Create(p).Error
}

// CreatePayment 追加一条报价/支付流水（不做 Upsert，保留历史）。
func (r *BusinessRepo) CreatePayment(p *model.PaymentRecord) error {
	return r.db.Create(p).Error
}

// ListPolicies 分页查询保单，可选按 channelCode 过滤。
func (r *BusinessRepo) ListPolicies(channelCode string, offset, limit int) ([]model.PolicyRecord, int64, error) {
	var list []model.PolicyRecord
	var total int64
	q := r.db.Model(&model.PolicyRecord{})
	if channelCode != "" {
		q = q.Where("channel_code = ?", channelCode)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset(offset).Limit(limit).Order("id desc").Find(&list).Error
	return list, total, err
}

// ListUsers 分页查询用户记录。
func (r *BusinessRepo) ListUsers(channelCode string, offset, limit int) ([]model.UserRecord, int64, error) {
	var list []model.UserRecord
	var total int64
	q := r.db.Model(&model.UserRecord{})
	if channelCode != "" {
		q = q.Where("channel_code = ?", channelCode)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset(offset).Limit(limit).Order("id desc").Find(&list).Error
	return list, total, err
}

// ListSigns 分页查询签约记录。
func (r *BusinessRepo) ListSigns(channelCode string, offset, limit int) ([]model.SignRecord, int64, error) {
	var list []model.SignRecord
	var total int64
	q := r.db.Model(&model.SignRecord{})
	if channelCode != "" {
		q = q.Where("channel_code = ?", channelCode)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset(offset).Limit(limit).Order("id desc").Find(&list).Error
	return list, total, err
}

// Stats 统计保单、用户、签约条数，供管理后台仪表盘使用。
func (r *BusinessRepo) Stats(channelCode string) (map[string]interface{}, error) {
	stats := map[string]interface{}{}
	pq := r.db.Model(&model.PolicyRecord{})
	uq := r.db.Model(&model.UserRecord{})
	sq := r.db.Model(&model.SignRecord{})
	if channelCode != "" {
		pq = pq.Where("channel_code = ?", channelCode)
		uq = uq.Where("channel_code = ?", channelCode)
		sq = sq.Where("channel_code = ?", channelCode)
	}
	var policies, users, signs int64
	_ = pq.Count(&policies).Error
	_ = uq.Count(&users).Error
	_ = sq.Count(&signs).Error
	stats["policyCount"] = policies
	stats["userCount"] = users
	stats["signCount"] = signs
	return stats, nil
}
