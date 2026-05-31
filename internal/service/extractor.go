package service

import (
	"encoding/json"
	"strconv"

	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/pkg/cipher"
	"github.com/huaan/insurance-bridge/internal/repository"
)

// Extractor 从华安明文响应（及少量请求字段）中抽取业务实体并 Upsert 落库，供管理后台统计与对账。
// 抽取失败不向上游/渠道返回错误，由 repository 层静默处理。
type Extractor struct {
	biz    *repository.BusinessRepo // 用户/保单/签约等业务表
	cipher *cipher.AES256GCM        // 用户三要素入库前加密
}

// NewExtractor 构造业务抽取器。
func NewExtractor(biz *repository.BusinessRepo, c *cipher.AES256GCM) *Extractor {
	return &Extractor{biz: biz, cipher: c}
}

// Process 按 apiPath 路由到对应抽取逻辑；req 为已解密的渠道请求体，resp 为华安原始响应体。
func (e *Extractor) Process(channelCode, apiPath string, req, resp map[string]interface{}) {
	switch apiPath {
	case "/proInsurance":
		e.fromProInsurance(channelCode, req, resp)
	case "/verifyNoCode":
		e.fromVerify(channelCode, resp)
	case "/getSignUrl":
		e.fromSign(channelCode, req, resp)
	case "/getPolicyInfoByPhoneNo", "/getPolicyInfoByPolicyId":
		e.fromPolicyInfo(channelCode, resp)
	case "/getProductPricesByProductCode", "/getProductPricesByPolicyId":
		e.fromPrices(channelCode, apiPath, resp)
	case "/getProductInfoByChannel":
		e.fromProducts(channelCode, resp)
	case "/upGradeIns":
		e.fromUpgrade(channelCode, req, resp)
	}
}

// fromProInsurance 投保成功时抽取 policyId、userId、policyStatus 等写入保单表。
func (e *Extractor) fromProInsurance(channelCode string, req, resp map[string]interface{}) {
	data, _ := resp["data"].(map[string]interface{})
	if data == nil {
		return
	}
	p := &model.PolicyRecord{
		ChannelCode:  channelCode,
		PolicyID:     str(data["policyId"]),
		UserID:       str(data["userId"]),
		ProductCode:  str(req["productCode"]),
		PolicyStatus: str(data["policyStatus"]),
	}
	if p.PolicyID != "" {
		_ = e.biz.UpsertPolicy(p)
	}
}

// fromVerify 实名认证接口：将 phoneNo/name/idCard 加密后 Upsert 用户表。
func (e *Extractor) fromVerify(channelCode string, resp map[string]interface{}) {
	data, _ := resp["data"].(map[string]interface{})
	if data == nil {
		return
	}
	u := &model.UserRecord{
		ChannelCode: channelCode,
		UserID:      str(data["userid"]),
	}
	u.PhoneEnc, _ = e.cipher.Encrypt(str(data["phoneNo"]))
	u.NameEnc, _ = e.cipher.Encrypt(str(data["name"]))
	u.IDCardEnc, _ = e.cipher.Encrypt(str(data["idCard"]))
	if u.UserID != "" {
		_ = e.biz.UpsertUser(u)
	}
}

// fromSign 获取签约链接：记录 signId、policyId、签约 URL 等。
func (e *Extractor) fromSign(channelCode string, req, resp map[string]interface{}) {
	data, _ := resp["data"].(map[string]interface{})
	if data == nil {
		return
	}
	s := &model.SignRecord{
		ChannelCode:  channelCode,
		SignID:       str(data["signId"]),
		PolicyID:     str(req["policyId"]),
		BankCode:     str(req["bankCode"]),
		CardType:     str(req["cardType"]),
		PayChannelID: str(req["payChannelId"]),
		SignURL:      str(data["url"]),
	}
	if s.SignID != "" {
		_ = e.biz.UpsertSign(s)
	}
}

// fromPolicyInfo 保单查询：data 可能为单对象或数组，仅 isPolicy=1 的条目落库。
func (e *Extractor) fromPolicyInfo(channelCode string, resp map[string]interface{}) {
	data := resp["data"]
	switch v := data.(type) {
	case map[string]interface{}:
		e.savePolicyFromMap(channelCode, v)
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				e.savePolicyFromMap(channelCode, m)
			}
		}
	}
}

// savePolicyFromMap 将单条保单 map 转为 PolicyRecord 并 Upsert。
func (e *Extractor) savePolicyFromMap(channelCode string, m map[string]interface{}) {
	if str(m["isPolicy"]) != "1" {
		return
	}
	p := &model.PolicyRecord{
		ChannelCode:     channelCode,
		PolicyID:        str(m["policyId"]),
		PolicyNo:        str(m["policyNo"]),
		ProductCode:     str(m["productCode"]),
		ProductID:       str(m["productId"]),
		PolicyStatus:    str(m["policyStatus"]),
		PolicyStartDate: str(m["policyStartDate"]),
		PolicyEndDate:   str(m["policyEndDate"]),
		PayPremium:      str(m["payPremium"]),
	}
	if p.PolicyID != "" {
		_ = e.biz.UpsertPolicy(p)
	}
}

// fromPrices 产品/保单报价：data 为数组或单对象时分别写入 payment 与 product 快照。
func (e *Extractor) fromPrices(channelCode, apiPath string, resp map[string]interface{}) {
	data, _ := resp["data"].([]interface{})
	if data == nil {
		if m, ok := resp["data"].(map[string]interface{}); ok {
			e.savePayment(channelCode, apiPath, m)
		}
		return
	}
	for _, item := range data {
		if m, ok := item.(map[string]interface{}); ok {
			e.savePayment(channelCode, apiPath, m)
		}
	}
}

// savePayment 记录报价流水，并在有 productCode 时同步更新产品目录价格。
func (e *Extractor) savePayment(channelCode, apiPath string, m map[string]interface{}) {
	p := &model.PaymentRecord{
		ChannelCode: channelCode,
		PolicyID:    str(m["policyId"]),
		SignID:      str(m["signId"]),
		ProductCode: str(m["productCode"]),
		Price:       str(m["price"]),
		SourceAPI:   apiPath,
	}
	if p.Price != "" || p.SignID != "" {
		_ = e.biz.CreatePayment(p)
	}
	if p.ProductCode != "" {
		_ = e.biz.UpsertProduct(&model.ProductRecord{
			ChannelCode: channelCode,
			ProductCode: p.ProductCode,
			Price:       p.Price,
		})
	}
}

// fromProducts 渠道产品列表：批量 Upsert 产品名称等信息。
func (e *Extractor) fromProducts(channelCode string, resp map[string]interface{}) {
	data, _ := resp["data"].([]interface{})
	for _, item := range data {
		if m, ok := item.(map[string]interface{}); ok {
			_ = e.biz.UpsertProduct(&model.ProductRecord{
				ChannelCode: channelCode,
				ProductCode: str(m["productCode"]),
				ProductName: str(m["productName"]),
			})
		}
	}
}

// fromUpgrade 升级险种成功（code=200）时更新保单状态为有效。
func (e *Extractor) fromUpgrade(channelCode string, req, resp map[string]interface{}) {
	code, _ := resp["code"].(float64)
	if int(code) != 200 {
		return
	}
	p := &model.PolicyRecord{
		ChannelCode:  channelCode,
		PolicyID:     str(req["policyId"]),
		PolicyStatus: "1",
	}
	if p.PolicyID != "" {
		_ = e.biz.UpsertPolicy(p)
	}
}

// str 将 JSON 解析后的 interface{} 安全转为字符串（兼容 string/float64/json.Number）。
func str(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case json.Number:
		return t.String()
	default:
		return ""
	}
}
