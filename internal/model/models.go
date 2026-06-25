// Package model 定义 GORM 持久化实体：渠道配置、管理员账号与银行信息。
package model

import "time"

// Channel 渠道方配置。每个渠道拥有独立的验签密钥与华安上游密钥，代理转发时完成 key 替换与重签。
type Channel struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`                                        // 主键
	ChannelCode string    `gorm:"uniqueIndex;size:64;not null" json:"channelCode"`             // 渠道编码，请求体 channelCode 与之匹配
	ChannelName string    `gorm:"size:128" json:"channelName"`                                 // 渠道显示名称
	ChannelKey  string    `gorm:"size:128;not null" json:"channelKey"`                         // 本服务分配给渠道的密钥，用于校验渠道请求 sign
	HuaAnKey     string    `gorm:"size:128;not null" json:"huaAnKey"`                           // 华安分配的密钥，转发上游时替换 body.key 并按此重签
	CallbackURL  string    `gorm:"column:callback_url;size:512;not null" json:"callbackUrl"`    // 投保结果回调地址，本服务转发华安 notify 时使用
	Status       int8      `gorm:"default:1" json:"status"`                                     // 1=启用（可代理），0=禁用
	Remark      string    `gorm:"size:512" json:"remark"`                                      // 备注
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// AdminUser 管理后台登录账号，密码以 bcrypt 哈希存储。
type AdminUser struct {
	ID           uint64    `gorm:"primaryKey"`
	Username     string    `gorm:"uniqueIndex;size:64;not null"` // 登录名
	PasswordHash string    `gorm:"size:255;not null"`            // bcrypt 哈希，不可逆
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// BankInfo 银行信息，对应 bank_info_t；由华安 getBankList 响应同步。
type BankInfo struct {
	ID           uint64    `gorm:"primaryKey;column:id" json:"id"`
	BankCode     string    `gorm:"column:bank_code;uniqueIndex;size:64;not null" json:"bankCode"`
	BankName     string    `gorm:"column:bank_name;size:128;not null" json:"bankName"`
	DebitCard    int8      `gorm:"column:debit_card;not null;default:0" json:"debitCard"`
	CreditCard   int8      `gorm:"column:credit_card;not null;default:0" json:"creditCard"`
	PayChannelID string    `gorm:"column:pay_channel_id;size:64" json:"payChannelId"`
	IsActBank    int8      `gorm:"column:is_act_bank;not null;default:1" json:"isActBank"` // 1=活跃银行 0=非活跃
	Status       int8      `gorm:"column:status;not null;default:1;index" json:"status"`     // 1=启用 2=停用（与 isActBank 同步）
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

// TableName 指定 GORM 表名。
func (BankInfo) TableName() string { return "bank_info_t" }
