// Package huaan 封装对华安上游 API 的 HTTP 调用，与渠道层（验签、PII）解耦。
package huaan

// BankListPath 华安银行列表接口路径。
const BankListPath = "/getBankList"

// APIPaths 华安 /upChannelApi 下全部 POST 接口，与渠道侧路由一一对应。
var APIPaths = []string{
	"/proInsurance",                  // 投保
	"/upGradeIns",                    // 升级险种
	"/getSignUrl",                    // 获取签约链接
	"/verifyNoCode",                  // 无验证码实名
	"/getPolicyInfoByPhoneNo",        // 按手机号查保单
	"/getProductPricesByProductCode", // 按产品编码报价
	BankListPath,                     // 银行列表
	"/getProductInfoByChannel",       // 渠道产品列表
	"/getProductPricesByPolicyId",    // 按保单报价
	"/getPolicyInfoByPolicyId",       // 按保单 ID 查询
}
