// 华安直连探测：按依赖顺序调用渠道路径，输出 JSON 结果供文档归档。
// 直连华安（不经渠道验签/PII 加密），三要素使用明文。
// productCode 取自 product/info；bankCode/payChannelId 取自 getBankList；
// sms 系列默认不测；getProductInfoByChannel、verifyNoCode、getPolicyInfoByPhoneNo、upGradeIns 已废弃不测。
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/huaan"
)

type stepResult struct {
	Index        int                    `json:"index"`
	Name         string                 `json:"name"`
	ChannelPath  string                 `json:"channelPath"`
	UpstreamPath string                 `json:"upstreamPath"`
	DependsOn    string                 `json:"dependsOn,omitempty"`
	Request      map[string]interface{} `json:"request"`
	Response     json.RawMessage        `json:"response,omitempty"`
	HuaAnCode    int                    `json:"huaanCode"`
	DurationMs   int64                  `json:"durationMs"`
	HTTPError    string                 `json:"httpError,omitempty"`
	BusinessOK   bool                   `json:"businessOK"`
	Note         string                 `json:"note,omitempty"`
}

type probeOutput struct {
	BaseURL     string       `json:"baseURL"`
	ChannelCode string       `json:"channelCode"`
	TestedAt    string       `json:"testedAt"`
	Hahealth    bool         `json:"hahealthMode"`
	Results     []stepResult `json:"results"`
}

// hahealthCaller 在 huaan.Client 基础上适配 ins.api.hahealth.ink：
// - 请求头 channelCode（及 common 接口可选 timestamp）
// - /upChannelApi/* 404 时改用 /proxy/upChannelApi/*
type hahealthCaller struct {
	client *huaan.Client
}

func newHahealthCaller(cfg *config.Config, _ bool) *hahealthCaller {
	return &hahealthCaller{
		client: huaan.NewClient(cfg, nil, nil),
	}
}

func (c *hahealthCaller) Call(ctx context.Context, apiPath string, body map[string]interface{}, headerTimestamp bool) (*huaan.CallResult, error) {
	_ = headerTimestamp // Client 对 /common/channel/api 自动附加 timestamp 头
	return c.client.Call(ctx, apiPath, body)
}

// resolveHahealthPath 在标准 ResolveUpstreamPath 基础上，将仍走 apiPath 前缀的路径改为 /proxy 前缀。
func resolveHahealthPath(apiPathPrefix, channelPath string) string {
	p := huaan.ResolveUpstreamPath(apiPathPrefix, channelPath)
	if strings.HasPrefix(p, "/common/channel/api/") {
		return p
	}
	if strings.HasPrefix(p, "/proxy/") {
		return p
	}
	if strings.HasPrefix(p, apiPathPrefix+"/") {
		suffix := strings.TrimPrefix(p, apiPathPrefix)
		return "/proxy" + apiPathPrefix + suffix
	}
	return p
}

func main() {
	baseURL := env("HUAAN_BASE_URL", "https://ins.api.hahealth.ink/")
	apiPath := env("HUAAN_API_PATH", "/upChannelApi")
	channelCode := mustEnv("HUAAN_CHANNEL_CODE")
	huaAnKey := os.Getenv("HUAAN_KEY")
	hahealth := envBool("HUAAN_HAHEALTH_MODE", strings.Contains(baseURL, "hahealth.ink"))

	phone := env("HUAAN_TEST_PHONE", "13811045503")
	name := env("HUAAN_TEST_NAME", "杜冲")
	idCard := env("HUAAN_TEST_ID_CARD", "13068319940517031X")
	skipSMS := envBool("HUAAN_PROBE_SKIP_SMS", true)

	cfg := &config.Config{
		HuaAn: config.HuaAnConfig{
			BaseURL:     baseURL,
			APIPath:     apiPath,
			Key:         huaAnKey,
			SignEnabled: envBool("HUAAN_SIGN_ENABLED", false),
		},
		Server: config.ServerConfig{UpstreamTimeout: 35 * time.Second},
	}

	var caller interface {
		Call(ctx context.Context, apiPath string, body map[string]interface{}, headerTimestamp bool) (*huaan.CallResult, error)
	}
	if hahealth {
		caller = newHahealthCaller(cfg, true)
	} else {
		caller = &legacyCaller{client: huaan.NewClient(cfg, nil, nil)}
	}

	ctx := context.Background()
	out := probeOutput{
		BaseURL:     strings.TrimRight(baseURL, "/"),
		ChannelCode: channelCode,
		TestedAt:    time.Now().Format(time.RFC3339),
		Hahealth:    hahealth,
	}

	var (
		productCode string
		productName string
		policyID    string
		userID      string
		bankCode    string
		payChannel  string
		cardType    string
		stepIdx     int
	)

	nextIdx := func() int {
		stepIdx++
		return stepIdx
	}

	piiFields := func() map[string]interface{} {
		return map[string]interface{}{
			"phoneNo": phone, "name": name, "idCard": idCard,
		}
	}

	call := func(name, path, depends string, fields map[string]interface{}, note string, hdrTs bool) stepResult {
		if fields == nil {
			fields = map[string]interface{}{}
		}
		body := huaan.BuildRequestBody(channelCode, fields)
		reqCopy := copyMap(body)
		upstream := resolveHahealthPath(apiPath, path)
		if !hahealth {
			upstream = huaan.ResolveUpstreamPath(apiPath, path)
		}

		res := stepResult{
			Index:        nextIdx(),
			Name:         name,
			ChannelPath:  apiPath + path,
			UpstreamPath: strings.TrimRight(baseURL, "/") + upstream,
			DependsOn:    depends,
			Request:      reqCopy,
			Note:         note,
		}

		result, err := caller.Call(ctx, path, body, hdrTs)
		if err != nil {
			res.HTTPError = err.Error()
			return res
		}
		res.DurationMs = result.DurationMs
		res.HuaAnCode = result.HuaAnCode
		res.BusinessOK = huaan.ResponseOK(result.Body)
		res.Response = json.RawMessage(result.Body)
		return res
	}

	appendResult := func(r stepResult) {
		out.Results = append(out.Results, r)
	}

	// 阶段 1：基础数据（银行、产品）
	r := call("getBankList", huaan.BankListPath, "", nil, "无业务参数；bankCode/payChannelId 取自本接口", false)
	appendResult(r)
	if r.BusinessOK {
		if pairs, err := huaan.ParseBankPayPairs(r.Response); err == nil && len(pairs) > 0 {
			b := pairs[0]
			bankCode, payChannel, cardType = b.BankCode, b.PayChannelID, b.DefaultCardType()
		}
	}

	r = call("product/info", huaan.ProductInfoPath, "", piiFields(), "三要素明文；productCode 取自本接口", hahealth)
	appendResult(r)
	if r.BusinessOK {
		if products, err := huaan.ParseChannelProducts(r.Response); err == nil && len(products) > 0 {
			productCode, productName = products[0].ProductCode, products[0].ProductName
		}
	}
	prodDep := "product/info → productCode=" + productCode
	if productCode == "" {
		prodDep = "product/info 未成功，productCode 为空"
	}

	// 阶段 2：查询类（不含已废弃 verifyNoCode、getPolicyInfoByPhoneNo）
	r = call("getUserInfoByPhoneNo", huaan.UserInfoByPhoneNoPath, "", map[string]interface{}{
		"phoneNo": phone,
	}, "phoneNo 明文", hahealth)
	appendResult(r)

	r = call("policy/phone", huaan.PolicyByPhonePath, "", map[string]interface{}{
		"phoneNo": phone,
	}, "phoneNo 明文", hahealth)
	appendResult(r)

	if productCode != "" {
		r = call("getProductPricesByProductCode", "/getProductPricesByProductCode", prodDep, mergeFields(piiFields(), map[string]interface{}{
			"productCode": productCode, "hasSocialSecurity": "1",
		}), "hasSocialSecurity 为字符串；含三要素", false)
		appendResult(r)

		r = call("priceByUser", huaan.PriceByUserPath, prodDep, map[string]interface{}{
			"productCode": productCode, "hasSocialSecurity": 1, "idCard": idCard,
		}, "hasSocialSecurity 为整数", hahealth)
		appendResult(r)

		r = call("getLiabilitiesByProductId", huaan.LiabilitiesByProductIDPath, prodDep, map[string]interface{}{
			"productCode": productCode, "hasSocialSecurity": 1, "productType": 1, "idCard": idCard,
		}, "productType 1=体验版", hahealth)
		appendResult(r)

		insFields := mergeFields(piiFields(), map[string]interface{}{
			"productCode": productCode, "hasSocialSecurity": 1, "isUpgrade": 0, "autoRenew": 1,
		})
		r = call("proInsurance", huaan.ProInsurancePath, prodDep, insFields, "成功时返回 policyId/userId", false)
		appendResult(r)
		if r.BusinessOK {
			if parsed, err := huaan.ParseProInsuranceResponse(r.Response); err == nil {
				policyID, userID = parsed.PolicyID, parsed.UserID
			}
		}
	}

	policyDep := "proInsurance → policyId=" + policyID
	if policyID == "" {
		policyDep = "proInsurance 未成功，跳过 policyId 相关请求"
	}

	if policyID != "" {
		r = call("getPolicyInfoByPolicyId", huaan.PolicyInfoByPolicyIDPath, policyDep, map[string]interface{}{
			"policyId": policyID,
		}, "仅需 policyId", false)
		appendResult(r)

		r = call("getProductPricesByPolicyId", huaan.ProductPricesByPolicyIDPath, policyDep, map[string]interface{}{
			"policyId": policyID,
		}, "仅需 policyId", false)
		appendResult(r)
	}

	if policyID != "" && userID != "" && bankCode != "" {
		signDep := policyDep + "; getBankList → bankCode=" + bankCode
		signFields := mergeFields(piiFields(), map[string]interface{}{
			"policyId": policyID, "userId": userID,
			"bankCode": bankCode, "payChannelId": payChannel, "cardType": cardType,
		})
		r = call("getSignUrl", "/getSignUrl", signDep, signFields,
			fmt.Sprintf("product=%s bank=%s", productName, bankCode), false)
		appendResult(r)
	}

	if !skipSMS {
		// sms 系列需验证码或单独联调，默认跳过；设 HUAAN_PROBE_SKIP_SMS=false 可启用
		smsCode := os.Getenv("HUAAN_TEST_SMS_CODE")
		r = call("sms/send", huaan.SmsSendPath, "", piiFields(), "三要素明文", hahealth)
		appendResult(r)
		r = call("sms/noValid", huaan.SmsNoValidPath, "", map[string]interface{}{"mobile": phone}, "mobile 明文", hahealth)
		appendResult(r)
		if smsCode != "" {
			r = call("sms/valid", huaan.SmsValidPath, "sms/send", map[string]interface{}{
				"mobile": phone, "smsCode": smsCode,
			}, "smsCode 明文", hahealth)
			appendResult(r)
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "encode: %v\n", err)
		os.Exit(1)
	}
}

type legacyCaller struct {
	client *huaan.Client
}

func (l *legacyCaller) Call(ctx context.Context, apiPath string, body map[string]interface{}, _ bool) (*huaan.CallResult, error) {
	return l.client.Call(ctx, apiPath, body)
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		fmt.Fprintf(os.Stderr, "missing env %s\n", k)
		os.Exit(1)
	}
	return v
}

func envBool(k string, def bool) bool {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(v, "yes")
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func mergeFields(base map[string]interface{}, extra map[string]interface{}) map[string]interface{} {
	out := copyMap(base)
	for k, v := range extra {
		out[k] = v
	}
	return out
}
