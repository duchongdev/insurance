// 华安直连探测：按依赖顺序调用 18 个渠道路径，输出 JSON 结果供文档归档。
// 直连华安（不经渠道验签/PII 加密），三要素使用明文。
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/huaan/insurance-bridge/internal/config"
	"github.com/huaan/insurance-bridge/internal/huaan"
	"github.com/huaan/insurance-bridge/internal/pkg/sign"
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
	cfg     *config.Config
	client  *http.Client
	headers bool
}

func newHahealthCaller(cfg *config.Config, useHeaders bool) *hahealthCaller {
	return &hahealthCaller{
		cfg:     cfg,
		client:  &http.Client{Timeout: cfg.Server.UpstreamTimeout},
		headers: useHeaders,
	}
}

func (c *hahealthCaller) Call(ctx context.Context, apiPath string, body map[string]interface{}, headerTimestamp bool) (*huaan.CallResult, error) {
	body["key"] = c.cfg.HuaAn.Key
	if c.cfg.HuaAn.SignEnabled {
		delete(body, "sign")
		body["sign"] = sign.Build(body, c.cfg.HuaAn.Key)
	} else {
		body["sign"] = ""
	}

	upstreamPath := resolveHahealthPath(c.cfg.HuaAn.APIPath, apiPath)
	upstreamURL := strings.TrimRight(c.cfg.HuaAn.BaseURL, "/") + upstreamPath
	reqBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if c.headers {
		if cc, ok := body["channelCode"].(string); ok && cc != "" {
			req.Header.Set("channelCode", cc)
		}
		if headerTimestamp {
			if ts, ok := body["timestamp"].(string); ok && ts != "" {
				req.Header.Set("timestamp", ts)
			}
		}
	}

	resp, err := c.client.Do(req)
	duration := time.Since(start).Milliseconds()
	if err != nil {
		if ctx.Err() != nil {
			return nil, huaan.ErrTimeout
		}
		return nil, fmt.Errorf("%w: %v", huaan.ErrUnavailable, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return &huaan.CallResult{
		Body:       respBytes,
		DurationMs: duration,
		HuaAnCode:  huaan.ParseCode(respBytes),
	}, nil
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
	smsCode := os.Getenv("HUAAN_TEST_SMS_CODE")
	seedProduct := env("HUAAN_TEST_PRODUCT_CODE", "ZFHLW1040003")

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
	)

	call := func(idx int, name, path, depends string, fields map[string]interface{}, note string, hdrTs bool) stepResult {
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
			Index:        idx,
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

	// 阶段 1：无依赖
	r := call(1, "getBankList", huaan.BankListPath, "", nil, "无业务参数", false)
	appendResult(r)
	if r.BusinessOK {
		if pairs, err := huaan.ParseBankPayPairs(r.Response); err == nil && len(pairs) > 0 {
			b := pairs[0]
			bankCode, payChannel, cardType = b.BankCode, b.PayChannelID, b.DefaultCardType()
		}
	}

	r = call(2, "getProductInfoByChannel", "/getProductInfoByChannel", "", nil,
		"无业务参数；正常情况下 productCode 取自本接口", false)
	appendResult(r)
	if r.BusinessOK {
		if products, err := huaan.ParseChannelProducts(r.Response); err == nil && len(products) > 0 {
			productCode, productName = products[0].ProductCode, products[0].ProductName
		}
	}

	// 阶段 2：三要素 / 手机号（common 接口在 hahealth 上附加 timestamp 请求头）
	r = call(3, "product/info", huaan.ProductInfoPath, "", map[string]interface{}{
		"phoneNo": phone, "name": name, "idCard": idCard,
	}, "三要素明文", hahealth)
	appendResult(r)

	r = call(4, "sms/send", huaan.SmsSendPath, "", map[string]interface{}{
		"phoneNo": phone, "name": name, "idCard": idCard,
	}, "三要素明文", hahealth)
	appendResult(r)

	r = call(5, "sms/noValid", huaan.SmsNoValidPath, "", map[string]interface{}{
		"mobile": phone,
	}, "mobile 明文", hahealth)
	appendResult(r)

	r = call(6, "verifyNoCode", "/verifyNoCode", "", map[string]interface{}{
		"phoneNo": phone, "name": name, "idCard": idCard,
	}, "三要素明文", false)
	appendResult(r)

	r = call(7, "getPolicyInfoByPhoneNo", "/getPolicyInfoByPhoneNo", "", map[string]interface{}{
		"phoneNo": phone, "name": name, "idCard": idCard,
	}, "三要素明文", false)
	appendResult(r)

	r = call(8, "getUserInfoByPhoneNo", huaan.UserInfoByPhoneNoPath, "", map[string]interface{}{
		"phoneNo": phone,
	}, "phoneNo 明文", hahealth)
	appendResult(r)

	r = call(9, "policy/phone", huaan.PolicyByPhonePath, "", map[string]interface{}{
		"phoneNo": phone,
	}, "phoneNo 明文", hahealth)
	appendResult(r)

	// 阶段 3：productCode（优先 getProductInfoByChannel，否则用种子 productCode 走报价接口反查）
	if productCode == "" {
		productCode = seedProduct
	}
	prodDep := "getProductInfoByChannel → productCode=" + productCode
	if !out.Results[1].BusinessOK {
		prodDep = "getProductInfoByChannel 失败，暂用种子 productCode=" + productCode + "（来自 HUAAN_TEST_PRODUCT_CODE 或报价反查）"
	}

	r = call(10, "getProductPricesByProductCode", "/getProductPricesByProductCode", prodDep, map[string]interface{}{
		"productCode": productCode, "hasSocialSecurity": "1",
		"phoneNo": phone, "name": name, "idCard": idCard,
	}, "hasSocialSecurity 为字符串；含三要素", false)
	appendResult(r)
	if productName == "" && r.BusinessOK {
		var resp map[string]interface{}
		if json.Unmarshal(r.Response, &resp) == nil {
			if items, ok := resp["data"].([]interface{}); ok && len(items) > 0 {
				if m, ok := items[0].(map[string]interface{}); ok {
					if v, ok := m["productCode"].(string); ok && v != "" {
						productCode = v
					}
					if v, ok := m["productName"].(string); ok {
						productName = v
					}
				}
			}
		}
	}

	r = call(11, "priceByUser", huaan.PriceByUserPath, prodDep, map[string]interface{}{
		"productCode": productCode, "hasSocialSecurity": 1, "idCard": idCard,
	}, "hasSocialSecurity 为整数", hahealth)
	appendResult(r)

	r = call(12, "getLiabilitiesByProductId", huaan.LiabilitiesByProductIDPath, prodDep, map[string]interface{}{
		"productCode": productCode, "hasSocialSecurity": 1, "productType": 1, "idCard": idCard,
	}, "productType 1=体验版", hahealth)
	appendResult(r)

	if smsCode == "" {
		r = call(13, "sms/valid", huaan.SmsValidPath, "sms/send", map[string]interface{}{
			"mobile": phone, "smsCode": "",
		}, "需真实验证码（未设置 HUAAN_TEST_SMS_CODE）", hahealth)
	} else {
		r = call(13, "sms/valid", huaan.SmsValidPath, "sms/send", map[string]interface{}{
			"mobile": phone, "smsCode": smsCode,
		}, "smsCode 明文", hahealth)
	}
	appendResult(r)

	insFields := map[string]interface{}{
		"productCode": productCode, "hasSocialSecurity": 1, "isUpgrade": 0, "autoRenew": 1,
		"phoneNo": phone, "name": name, "idCard": idCard,
	}
	r = call(14, "proInsurance", huaan.ProInsurancePath, prodDep, insFields, "成功时返回 policyId/userId", false)
	appendResult(r)
	if r.BusinessOK {
		if parsed, err := huaan.ParseProInsuranceResponse(r.Response); err == nil {
			policyID, userID = parsed.PolicyID, parsed.UserID
		}
	}
	policyDep := "proInsurance → policyId=" + policyID
	if policyID == "" {
		policyDep = "proInsurance 未成功，policyId 为空"
	}

	r = call(15, "getPolicyInfoByPolicyId", huaan.PolicyInfoByPolicyIDPath, policyDep, map[string]interface{}{
		"policyId": policyID,
	}, "仅需 policyId", false)
	appendResult(r)

	r = call(16, "getProductPricesByPolicyId", huaan.ProductPricesByPolicyIDPath, policyDep, map[string]interface{}{
		"policyId": policyID,
	}, "仅需 policyId", false)
	appendResult(r)

	r = call(17, "upGradeIns", "/upGradeIns", policyDep, map[string]interface{}{
		"policyId": policyID,
	}, "仅需 policyId", false)
	appendResult(r)

	signDep := policyDep + "; getBankList → bankCode/payChannelId"
	r = call(18, "getSignUrl", "/getSignUrl", signDep, map[string]interface{}{
		"policyId": policyID, "userId": userID,
		"bankCode": bankCode, "payChannelId": payChannel, "cardType": cardType,
		"phoneNo": phone, "name": name, "idCard": idCard,
	}, fmt.Sprintf("product=%s bank=%s", productName, bankCode), false)
	appendResult(r)

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
