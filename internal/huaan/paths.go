// Package huaan 封装对华安上游 API 的 HTTP 调用，与渠道层（验签、PII）解耦。
package huaan

// BankListPath 渠道侧银行列表路径（挂载在 huaan.api_path 路由组下）。
const BankListPath = "/getBankList"

// BankListUpstreamPath 华安侧银行列表路径（相对 base_url，非 /upChannelApi 前缀）。
const BankListUpstreamPath = "/common/channel/api/getBankList"

// ProInsurancePath 渠道侧预投保路径（挂载在 huaan.api_path 路由组下）。
const ProInsurancePath = "/proInsurance"

// ProInsuranceUpstreamPath 华安侧预投保路径（相对 base_url，经 proxy 前缀）。
const ProInsuranceUpstreamPath = "/proxy/upChannelApi/proInsurance"

// ProductInfoByChannelPath 已废弃：华安侧不再提供/使用该接口，请改用 ProductInfoPath（product/info）。
// 保留常量仅作历史引用，勿注册路由、勿在探测/集成测试中调用。
const ProductInfoByChannelPath = "/getProductInfoByChannel"

// ProductInfoPath 渠道侧获取渠道产品信息路径（挂载在 huaan.api_path 路由组下）。
const ProductInfoPath = "/product/info"

// ProductInfoUpstreamPath 华安侧获取渠道产品信息路径（相对 base_url，非 /upChannelApi 前缀）。
const ProductInfoUpstreamPath = "/common/channel/api/product/info"

// SmsSendPath 渠道侧发送短信验证码路径（挂载在 huaan.api_path 路由组下）。
const SmsSendPath = "/sms/send"

// SmsSendUpstreamPath 华安侧发送短信验证码路径（相对 base_url，非 /upChannelApi 前缀）。
const SmsSendUpstreamPath = "/common/channel/api/sms/send"

// SmsValidPath 渠道侧短信验证码校验路径（挂载在 huaan.api_path 路由组下）。
// 联调状态：未完成联调（华安上游验签/参数绑定规则待确认，见 docs/TESTING.md §7）。
const SmsValidPath = "/sms/valid"

// SmsValidUpstreamPath 华安侧短信验证码校验路径（相对 base_url，非 /upChannelApi 前缀）。
const SmsValidUpstreamPath = "/common/channel/api/sms/valid"

// SmsNoValidPath 渠道侧免短信验证码注册登录路径（挂载在 huaan.api_path 路由组下）。
const SmsNoValidPath = "/sms/noValid"

// SmsNoValidUpstreamPath 华安侧免短信验证码注册登录路径（相对 base_url，非 /upChannelApi 前缀）。
const SmsNoValidUpstreamPath = "/common/channel/api/sms/noValid"

// PriceByUserPath 渠道侧查询产品价格路径（挂载在 huaan.api_path 路由组下）。
const PriceByUserPath = "/priceByUser"

// PriceByUserUpstreamPath 华安侧查询产品价格路径（相对 base_url，非 /upChannelApi 前缀）。
const PriceByUserUpstreamPath = "/common/channel/api/priceByUser"

// PolicyByPhonePath 渠道侧通过手机号查询用户投保情况路径（挂载在 huaan.api_path 路由组下）。
const PolicyByPhonePath = "/policy/phone"

// PolicyByPhoneUpstreamPath 华安侧通过手机号查询用户投保情况路径（相对 base_url，非 /upChannelApi 前缀）。
const PolicyByPhoneUpstreamPath = "/common/channel/api/policy/phone"

// UserInfoByPhoneNoPath 渠道侧根据手机号/用户 ID 查询用户信息路径（挂载在 huaan.api_path 路由组下）。
const UserInfoByPhoneNoPath = "/getUserInfoByPhoneNo"

// UserInfoByPhoneNoUpstreamPath 华安侧查询用户信息路径（相对 base_url，非 /upChannelApi 前缀）。
const UserInfoByPhoneNoUpstreamPath = "/common/channel/api/getUserInfoByPhoneNo"

// LiabilitiesByProductIDPath 渠道侧根据产品查询可选责任列表路径（挂载在 huaan.api_path 路由组下）。
const LiabilitiesByProductIDPath = "/getLiabilitiesByProductId"

// LiabilitiesByProductIDUpstreamPath 华安侧根据产品查询可选责任列表路径（相对 base_url，非 /upChannelApi 前缀）。
const LiabilitiesByProductIDUpstreamPath = "/common/channel/api/getLiabilitiesByProductId"

// ProductPricesByPolicyIDPath 渠道侧按保单 ID 查价格路径（挂载在 huaan.api_path 路由组下）。
const ProductPricesByPolicyIDPath = "/getProductPricesByPolicyId"

// ProductPricesByPolicyIDUpstreamPath 华安侧按保单 ID 查价格路径（相对 base_url，非 /upChannelApi 前缀）。
const ProductPricesByPolicyIDUpstreamPath = "/common/channel/api/getProductPricesByPolicyId"

// PolicyInfoByPolicyIDPath 渠道侧按保单 ID 查详情路径（挂载在 huaan.api_path 路由组下）。
const PolicyInfoByPolicyIDPath = "/getPolicyInfoByPolicyId"

// PolicyInfoByPolicyIDUpstreamPath 华安侧按保单 ID 查详情路径（相对 base_url，非 /upChannelApi 前缀）。
const PolicyInfoByPolicyIDUpstreamPath = "/common/channel/api/getPolicyInfoByPolicyId"

// GetPhoneByTokenPath 渠道侧一键登录解密手机号路径（挂载在 huaan.api_path 路由组下）。
const GetPhoneByTokenPath = "/getPhoneByToken"

// GetPhoneByTokenUpstreamPath 华安侧一键登录解密手机号路径（相对 base_url，非 /upChannelApi 前缀）。
const GetPhoneByTokenUpstreamPath = "/common/channel/api/getPhoneByToken"

// UpGradeInsPath 已废弃：华安侧不再使用，勿注册路由、勿在探测/集成测试中调用。
const UpGradeInsPath = "/upGradeIns"

// VerifyNoCodePath 已废弃：华安侧不再使用，勿注册路由、勿在探测/集成测试中调用。
const VerifyNoCodePath = "/verifyNoCode"

// PolicyInfoByPhoneNoPath 已废弃：华安侧不再使用，勿注册路由、勿在探测/集成测试中调用。
const PolicyInfoByPhoneNoPath = "/getPolicyInfoByPhoneNo"

// upstreamPathOverrides 渠道路径 → 华安完整路径后缀（相对 base_url）；未命中时拼 api_path + 渠道路径。
var upstreamPathOverrides = map[string]string{
	BankListPath:                    BankListUpstreamPath,
	ProInsurancePath:                ProInsuranceUpstreamPath,
	ProductInfoPath:                 ProductInfoUpstreamPath,
	SmsSendPath:                     SmsSendUpstreamPath,
	SmsValidPath:                    SmsValidUpstreamPath,
	SmsNoValidPath:                  SmsNoValidUpstreamPath,
	PriceByUserPath:                 PriceByUserUpstreamPath,
	PolicyByPhonePath:               PolicyByPhoneUpstreamPath,
	UserInfoByPhoneNoPath:           UserInfoByPhoneNoUpstreamPath,
	LiabilitiesByProductIDPath:      LiabilitiesByProductIDUpstreamPath,
	ProductPricesByPolicyIDPath:     ProductPricesByPolicyIDUpstreamPath,
	PolicyInfoByPolicyIDPath:        PolicyInfoByPolicyIDUpstreamPath,
	GetPhoneByTokenPath:             GetPhoneByTokenUpstreamPath,
}

// ResolveUpstreamPath 返回调用华安时的 URL 路径后缀（相对 base_url）。
func ResolveUpstreamPath(apiPathPrefix, channelPath string) string {
	if p, ok := upstreamPathOverrides[channelPath]; ok {
		return p
	}
	return apiPathPrefix + channelPath
}

// APIPaths 渠道 POST 接口路径（挂载在 config huaan.api_path 路由组下），与 Register 一一对应。
var APIPaths = []string{
	ProInsurancePath,                 // 预投保
	"/getSignUrl",                    // 获取签约链接
	"/getProductPricesByProductCode", // 按产品编码报价
	BankListPath,                     // 银行列表
	ProductInfoPath,                  // 获取渠道产品信息（productCode 取自本接口）
	SmsSendPath,                      // 发送短信验证码
	SmsValidPath,                     // 短信验证码校验
	SmsNoValidPath,                   // 免短信验证码注册登录
	PriceByUserPath,                  // 查询产品价格
	PolicyByPhonePath,                // 通过手机号查询用户投保情况
	UserInfoByPhoneNoPath,            // 根据手机号/用户 ID 查询用户信息
	LiabilitiesByProductIDPath,       // 根据产品查询可选责任列表
	ProductPricesByPolicyIDPath,      // 按保单 ID 查价格（收银台展示支付金额）
	PolicyInfoByPolicyIDPath,         // 按保单 ID 查详情
	GetPhoneByTokenPath,              // 一键登录解密手机号
}
