// Package model 定义 GORM 持久化实体：渠道配置、管理员账号与银行信息。
package model

import "time"

// Channel 渠道方配置。每个渠道拥有独立的验签密钥与华安上游密钥，代理转发时完成 key 替换与重签。
type Channel struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`                                        // 主键
	ChannelCode    string    `gorm:"size:64;not null;index" json:"channelCode"`                   // 渠道编码，与华安配置 channel_code 一致（允许多环境同码，同时仅一条启用）
	ChannelName    string    `gorm:"size:128" json:"channelName"`                                 // 渠道显示名称
	ChannelKey     string    `gorm:"size:128;not null" json:"channelKey"`                         // 本服务分配给渠道的密钥，用于校验渠道请求 sign
	HuaAnKey       string    `gorm:"size:128;not null" json:"huaAnKey"`                           // 同步自关联华安配置的 channel_secret
	HuaAnSettingID uint64    `gorm:"column:huaan_setting_id;not null;default:0;index" json:"huaanSettingId"` // 关联 huaan_settings.id，每条华安配置最多对应一个渠道
	CallbackURL    string    `gorm:"column:callback_url;size:512;not null" json:"callbackUrl"`    // 投保结果回调地址
	EnvType        string    `gorm:"-" json:"envType,omitempty"`                                // 关联华安配置的环境类型，列表接口填充
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

// HuaAnSetting 华安上游连接配置，支持多条；env_type 区分生产/测试，is_active 标记当前运行时使用的配置。
type HuaAnSetting struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:128" json:"name"`
	EnvType       string    `gorm:"column:env_type;size:16;not null;index" json:"envType"`
	BaseURL       string    `gorm:"column:base_url;size:512;not null" json:"baseUrl"`
	ChannelCode   string    `gorm:"column:channel_code;size:64;not null" json:"channelCode"`
	ChannelSecret string    `gorm:"column:channel_secret;size:128;not null" json:"channelSecret"`
	IsActive      bool      `gorm:"column:is_active;not null;default:0" json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// 华安配置环境类型。
const (
	HuaAnEnvProd = "prod"
	HuaAnEnvTest = "test"
)

// TableName 指定 GORM 表名。
func (HuaAnSetting) TableName() string { return "huaan_settings" }
