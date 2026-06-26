// 模拟下游渠道依次调用本服务 upChannelApi 接口，原样输出请求/响应 JSON。
// 敏感信息通过环境变量传入，勿写入代码或提交到 Git。
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/huaan/insurance-bridge/internal/huaan"
	"github.com/huaan/insurance-bridge/internal/pkg/cipher"
	"github.com/huaan/insurance-bridge/internal/pkg/sign"
)

type caseDef struct {
	name            string
	path            string
	fields          map[string]interface{}
	needPII         bool
	needMobile      bool
	needIdCard      bool
	needPhoneOrUser bool
}

func main() {
	baseURL := env("BRIDGE_BASE_URL", "http://10.41.61.41")
	channelCode := mustEnv("CHANNEL_CODE")
	channelKey := mustEnv("CHANNEL_KEY")
	encKey := mustEnv("DATA_ENCRYPTION_KEY")

	phone := env("TEST_PHONE", "")
	name := env("TEST_NAME", "")
	idCard := env("TEST_ID_CARD", "")
	productCode := env("TEST_PRODUCT_CODE", "ZFHLW1040003")
	policyID := env("TEST_POLICY_ID", "91a51cc70e724d4885ac2e27530099ec")
	userID := env("TEST_USER_ID", "5c819c252ab343e2a2a6df98b3a014ab")
	smsCode := env("TEST_SMS_CODE", "1234")

	aes, err := cipher.New(encKey)
	if err != nil {
		fatal(err)
	}

	enc := func(plain string) string {
		if plain == "" {
			return ""
		}
		out, err := aes.Encrypt(plain)
		if err != nil {
			fatal(err)
		}
		return out
	}

	cases := []caseDef{
		{name: "getBankList", path: huaan.BankListPath},
		{name: "getProductInfoByChannel", path: "/getProductInfoByChannel"},
		{name: "productInfo", path: huaan.ProductInfoPath, needPII: true},
		{name: "smsSend", path: huaan.SmsSendPath, needPII: true},
		{name: "smsValid", path: huaan.SmsValidPath, needMobile: true, fields: map[string]interface{}{"smsCode": smsCode}},
		{name: "smsNoValid", path: huaan.SmsNoValidPath, needMobile: true},
		{name: "priceByUser", path: huaan.PriceByUserPath, needIdCard: true, fields: map[string]interface{}{
			"productCode": productCode, "hasSocialSecurity": 1,
		}},
		{name: "policyByPhone", path: huaan.PolicyByPhonePath, needPhoneOrUser: true},
		{name: "getUserInfoByPhoneNo", path: huaan.UserInfoByPhoneNoPath, needPhoneOrUser: true},
		{name: "liabilitiesByProductId", path: huaan.LiabilitiesByProductIDPath, needIdCard: true, fields: map[string]interface{}{
			"productCode": productCode, "hasSocialSecurity": 1, "productType": 1,
		}},
		{name: "getProductPricesByProductCode", path: "/getProductPricesByProductCode", needPII: true, fields: map[string]interface{}{
			"productCode": productCode, "hasSocialSecurity": "1",
		}},
		{name: "verifyNoCode", path: "/verifyNoCode", needPII: true},
		{name: "proInsurance", path: huaan.ProInsurancePath, needPII: true, fields: map[string]interface{}{
			"productCode": productCode, "hasSocialSecurity": 1, "isUpgrade": 0, "autoRenew": 1,
		}},
		{name: "upGradeIns", path: "/upGradeIns", fields: map[string]interface{}{"policyId": policyID}},
		{name: "getSignUrl", path: "/getSignUrl", needPII: true, fields: map[string]interface{}{
			"policyId": policyID, "bankCode": "BOC", "cardType": "1", "userId": userID,
		}},
		{name: "getPolicyInfoByPhoneNo", path: "/getPolicyInfoByPhoneNo", needPII: true},
		{name: "getProductPricesByPolicyId", path: huaan.ProductPricesByPolicyIDPath, fields: map[string]interface{}{
			"policyId": policyID,
		}},
		{name: "getPolicyInfoByPolicyId", path: huaan.PolicyInfoByPolicyIDPath, fields: map[string]interface{}{"policyId": policyID}},
	}

	client := &http.Client{Timeout: 35 * time.Second}
	sep := strings.Repeat("=", 72)

	for i, tc := range cases {
		fields := copyMap(tc.fields)
		if tc.needPII {
			if phone == "" || name == "" || idCard == "" {
				fmt.Printf("\n[%d/%d] %s — 跳过：未设置 TEST_PHONE/TEST_NAME/TEST_ID_CARD\n", i+1, len(cases), tc.name)
				continue
			}
			fields["phoneNo"] = enc(phone)
			fields["name"] = enc(name)
			fields["idCard"] = enc(idCard)
		}
		if tc.needMobile {
			if phone == "" {
				fmt.Printf("\n[%d/%d] %s — 跳过：未设置 TEST_PHONE\n", i+1, len(cases), tc.name)
				continue
			}
			fields["mobile"] = enc(phone)
		}
		if tc.needIdCard {
			if idCard == "" {
				fmt.Printf("\n[%d/%d] %s — 跳过：未设置 TEST_ID_CARD\n", i+1, len(cases), tc.name)
				continue
			}
			fields["idCard"] = enc(idCard)
		}
		if tc.needPhoneOrUser {
			if phone != "" {
				fields["phoneNo"] = enc(phone)
			} else if userID != "" {
				fields["userId"] = userID
			} else {
				fmt.Printf("\n[%d/%d] %s — 跳过：未设置 TEST_PHONE 或 TEST_USER_ID\n", i+1, len(cases), tc.name)
				continue
			}
		}

		body := huaan.BuildRequestBody(channelCode, fields)
		body["key"] = channelKey
		body["sign"] = sign.Build(sign.MapFromJSON(body), channelKey)

		reqJSON, _ := json.MarshalIndent(body, "", "  ")
		fmt.Printf("\n%s\n[%d/%d] %s\n%s\n", sep, i+1, len(cases), tc.name, sep)
		fmt.Printf(">>> POST %s/upChannelApi%s\n\n请求参数:\n%s\n\n", baseURL, tc.path, reqJSON)

		raw, _ := json.Marshal(body)
		resp, err := client.Post(baseURL+"/upChannelApi"+tc.path, "application/json; charset=utf-8", bytes.NewReader(raw))
		if err != nil {
			fmt.Printf("HTTP 错误: %v\n", err)
			continue
		}
		respBody, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		var parsed json.RawMessage
		if err := json.Unmarshal(respBody, &parsed); err != nil {
			fmt.Printf("HTTP %d\n响应(raw):\n%s\n", resp.StatusCode, string(respBody))
			continue
		}
		var pretty bytes.Buffer
		_ = json.Indent(&pretty, parsed, "", "  ")
		fmt.Printf("HTTP %d\n\n响应参数:\n%s\n", resp.StatusCode, pretty.String())

		var m map[string]interface{}
		if json.Unmarshal(respBody, &m) == nil {
			if tc.name == "proInsurance" && env("TEST_POLICY_ID", "") == "" {
				if data, ok := m["data"].(map[string]interface{}); ok {
					if pid, ok := data["policyId"].(string); ok && pid != "" {
						policyID = pid
						fmt.Printf("\n(已从 proInsurance 响应更新 policyId=%s)\n", policyID)
					}
				}
			}
			if tc.name == "verifyNoCode" && env("TEST_USER_ID", "") == "" {
				if data, ok := m["data"].(map[string]interface{}); ok {
					if uid, ok := data["userId"].(string); ok && uid != "" {
						userID = uid
						fmt.Printf("\n(已从 verifyNoCode 响应更新 userId=%s)\n", userID)
					}
				}
			}
		}
	}
}

func copyMap(in map[string]interface{}) map[string]interface{} {
	if in == nil {
		return map[string]interface{}{}
	}
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
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
		fmt.Fprintf(os.Stderr, "缺少环境变量 %s\n", k)
		os.Exit(1)
	}
	return v
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
