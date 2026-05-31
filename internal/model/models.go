// Package model 定义 GORM 持久化实体，对应渠道配置、接口审计日志及从华安响应抽取的业务快照。
package model

import "time"

// Channel 渠道方配置。每个渠道拥有独立的验签密钥与华安上游密钥，代理转发时完成 key 替换与重签。
type Channel struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`                                        // 主键
	ChannelCode string    `gorm:"uniqueIndex;size:64;not null" json:"channelCode"`             // 渠道编码，请求体 channelCode 与之匹配
	ChannelName string    `gorm:"size:128" json:"channelName"`                                 // 渠道显示名称
	ChannelKey  string    `gorm:"size:128;not null" json:"channelKey"`                         // 本服务分配给渠道的密钥，用于校验渠道请求 sign
	HuaAnKey    string    `gorm:"size:128;not null" json:"huaAnKey"`                           // 华安分配的密钥，转发上游时替换 body.key 并按此重签
	Status      int8      `gorm:"default:1" json:"status"`                                     // 1=启用（可代理），0=禁用
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

// APIRequestLog 单次渠道 API 调用的审计记录，与代理链路 trace_id 对应，供后台查询与排障。
type APIRequestLog struct {
	ID           uint64    `gorm:"primaryKey"`
	TraceID      string    `gorm:"index;size:64;not null"`  // 链路 ID（代理内生成 UUID，可与 X-Trace-Id 关联）
	ChannelCode  string    `gorm:"index;size:64"`           // 渠道编码
	APIPath      string    `gorm:"index;size:128;not null"` // 如 /proInsurance
	RequestBody  string    `gorm:"type:longtext"`           // 脱敏后的请求 JSON（三要素已 Mask）
	ResponseBody string    `gorm:"type:longtext"`           // 华安原始响应 JSON（明文，便于对账）
	HuaAnCode    int       `gorm:"index"`                   // 华安响应 body.code
	DurationMS   int64                                     // 上游往返耗时（毫秒）
	ErrorMessage string    `gorm:"size:512"`                // 上游失败或超时时的错误摘要
	CreatedAt    time.Time `gorm:"index"`
}

// UserRecord 用户三要素快照（库内 AES 加密），由 /verifyNoCode 等接口响应抽取。
type UserRecord struct {
	ID          uint64    `gorm:"primaryKey"`
	ChannelCode string    `gorm:"uniqueIndex:idx_channel_user;size:64;not null"` // 渠道 + 用户唯一
	UserID      string    `gorm:"uniqueIndex:idx_channel_user;size:64"`          // 华安 userid
	PhoneEnc    string    `gorm:"size:512"`                                      // 手机号密文
	NameEnc     string    `gorm:"size:512"`                                      // 姓名密文
	IDCardEnc   string    `gorm:"size:512"`                                      // 身份证号密文
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PolicyRecord 保单信息快照，由投保、查询、升级等接口响应抽取。
type PolicyRecord struct {
	ID              uint64    `gorm:"primaryKey"`
	ChannelCode     string    `gorm:"index;size:64;not null"`
	PolicyID        string    `gorm:"uniqueIndex;size:64;not null"` // 华安 policyId，全局唯一
	PolicyNo        string    `gorm:"size:64"`                      // 保单号
	UserID          string    `gorm:"index;size:64"`
	ProductCode     string    `gorm:"size:64"`
	ProductID       string    `gorm:"size:64"`
	PolicyStatus    string    `gorm:"size:8"`   // 保单状态码
	PolicyStartDate string    `gorm:"size:32"`  // 起保日
	PolicyEndDate   string    `gorm:"size:32"`  // 终保日
	PayPremium      string    `gorm:"size:32"`  // 保费（字符串保持与上游一致）
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// SignRecord 签约/代扣链接记录，由 /getSignUrl 响应抽取。
type SignRecord struct {
	ID           uint64    `gorm:"primaryKey"`
	ChannelCode  string    `gorm:"index;size:64;not null"`
	SignID       string    `gorm:"index;size:64"`    // 签约流水号
	PolicyID     string    `gorm:"index;size:64"`    // 关联保单
	BankCode     string    `gorm:"size:32"`          // 银行编码
	CardType     string    `gorm:"size:32"`          // 卡类型
	PayChannelID string    `gorm:"size:64"`          // 支付渠道
	SignURL      string    `gorm:"size:1024"`        // 签约 H5 地址
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ProductRecord 渠道产品目录快照，由产品列表/报价接口抽取。
type ProductRecord struct {
	ID          uint64    `gorm:"primaryKey"`
	ChannelCode string    `gorm:"uniqueIndex:idx_channel_product;size:64;not null"`
	ProductCode string    `gorm:"uniqueIndex:idx_channel_product;size:64;not null"`
	ProductName string    `gorm:"size:128"`
	ProductID   string    `gorm:"size:64"`
	ProductType string    `gorm:"size:32"`
	Price       string    `gorm:"size:32"` // 最近报价
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PaymentRecord 报价/支付流水（只增），用于统计各接口产生的价格记录。
type PaymentRecord struct {
	ID          uint64    `gorm:"primaryKey"`
	ChannelCode string    `gorm:"index;size:64;not null"`
	PolicyID    string    `gorm:"index;size:64"`
	SignID      string    `gorm:"index;size:64"`
	ProductCode string    `gorm:"size:64"`
	Price       string    `gorm:"size:32"`
	SourceAPI   string    `gorm:"size:128"` // 产生记录的接口路径
	CreatedAt   time.Time
}
