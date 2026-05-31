// Package pii 定义三要素字段名及响应加密路径，并提供渠道侧 AES 加解密转换（见 transform.go）。
package pii

// RequestFields 渠道请求体顶层需解密的字段（转发华安前转为明文）。
var RequestFields = []string{"phoneNo", "name", "idCard"}

// ResponseEncryptPaths 文档约定的响应内 PII 路径（点分表示，[] 表示数组下标），供参考；
// 实际加密由 transform.encryptObjectPII 按字段名递归处理 data / insuredList。
var ResponseEncryptPaths = []string{
	"data.phoneNo",
	"data.name",
	"data.idCard",
	"data.userid",
	"data.insuredList[].insuredCardNo",
	"data.insuredList[].insuredName",
	"data.insuredList[].phoneNo",
	// getPolicyInfoByPhoneNo 可能返回列表
	"data[].phoneNo",
	"data[].insuredList[].insuredCardNo",
	"data[].insuredList[].insuredName",
	"data[].insuredList[].phoneNo",
}
