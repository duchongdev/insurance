//go:build integration

package huaan_test

import (
	"context"
	"os"
	"testing"

	"github.com/huaan/insurance-bridge/internal/huaan"
)

// 华安直连集成测试：调用本服务真实的 huaan.Client，直连华安环境。
// 运行前设置环境变量，见 .env.huaan.example。
//
//	go test -tags=integration ./internal/huaan/ -v -count=1
const defaultHasSocialSecurity = "1" // 华安报价/投保类接口必填

func TestHuaAnDirect_AllPaths(t *testing.T) {
	client, env := huaan.NewTestClient(t)

	cases := []struct {
		name    string
		path    string
		needPII bool
		fields  map[string]interface{}
	}{
		{
			name:   "getBankList",
			path:   huaan.BankListPath,
			fields: nil,
		},
		{
			name:   "getProductInfoByChannel",
			path:   "/getProductInfoByChannel",
			fields: nil,
		},
		{
			name:    "getPolicyInfoByPhoneNo",
			path:    "/getPolicyInfoByPhoneNo",
			needPII: true,
			fields:  nil,
		},
		{
			name:    "verifyNoCode",
			path:    "/verifyNoCode",
			needPII: true,
			fields:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fields := tc.fields
			if tc.needPII {
				if env.Phone == "" || env.Name == "" || env.IDCard == "" {
					t.Skip("跳过：请设置 HUAAN_TEST_PHONE、HUAAN_TEST_NAME、HUAAN_TEST_ID_CARD")
				}
				if fields == nil {
					fields = map[string]interface{}{}
				}
				fields["phoneNo"] = env.Phone
				fields["name"] = env.Name
				fields["idCard"] = env.IDCard
			}

			body := huaan.BuildRequestBody(env.ChannelCode, fields)
			huaan.LogRequest(t, tc.path, body)
			result, err := client.Call(t.Context(), tc.path, body)
			if err != nil {
				t.Fatalf("Call %s: %v", tc.path, err)
			}
			huaan.LogResponse(t, tc.path, result)
		})
	}
}

// TestHuaAnDirect_ProductFromChannel 报价与投保：先 getProductInfoByChannel 取 productCode，再分别调两个接口。
//
//	go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_ProductFromChannel
func TestHuaAnDirect_ProductFromChannel(t *testing.T) {
	client, env := huaan.NewTestClient(t)
	if env.Phone == "" || env.Name == "" || env.IDCard == "" {
		t.Skip("跳过：请设置 HUAAN_TEST_PHONE、HUAAN_TEST_NAME、HUAAN_TEST_ID_CARD")
	}

	ctx := t.Context()
	products := fetchChannelProducts(t, ctx, client, env.ChannelCode)

	for _, prod := range products {
		t.Run(prod.ProductCode, func(t *testing.T) {
			t.Logf("选用产品（来自 getProductInfoByChannel）: productCode=%s productName=%s",
				prod.ProductCode, prod.ProductName)

			baseFields := insuranceFields(env, prod.ProductCode)

			t.Run("getProductPricesByProductCode", func(t *testing.T) {
				const path = "/getProductPricesByProductCode"
				body := huaan.BuildRequestBody(env.ChannelCode, baseFields)
				huaan.LogRequest(t, path, body)
				result, err := client.Call(ctx, path, body)
				if err != nil {
					t.Fatalf("Call %s: %v", path, err)
				}
				huaan.LogResponse(t, path, result)
			})

			t.Run("proInsurance", func(t *testing.T) {
				const path = "/proInsurance"
				body := huaan.BuildRequestBody(env.ChannelCode, baseFields)
				huaan.LogRequest(t, path, body)
				result, err := client.Call(ctx, path, body)
				if err != nil {
					t.Fatalf("Call %s: %v", path, err)
				}
				huaan.LogResponse(t, path, result)
			})
		})
	}
}

// TestHuaAnDirect_PolicyFromProInsurance 依赖 policyId 的接口：先投保取 policyId，再查保单/报价/升级。
// policyId 来自 proInsurance 响应（与生产一致：下游请求携带，本服务仅转发）。
//
//	go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_PolicyFromProInsurance
func TestHuaAnDirect_PolicyFromProInsurance(t *testing.T) {
	client, env := huaan.NewTestClient(t)
	if env.Phone == "" || env.Name == "" || env.IDCard == "" {
		t.Skip("跳过：请设置 HUAAN_TEST_PHONE、HUAAN_TEST_NAME、HUAAN_TEST_ID_CARD")
	}

	ctx := t.Context()
	products := fetchChannelProducts(t, ctx, client, env.ChannelCode)
	prod := products[0]
	t.Logf("选用产品: productCode=%s productName=%s", prod.ProductCode, prod.ProductName)

	insResult, insParsed, err := huaan.CallProInsurance(ctx, client, env.ChannelCode, insuranceFields(env, prod.ProductCode))
	if err != nil {
		t.Fatalf("proInsurance: %v", err)
	}
	huaan.LogResponse(t, "/proInsurance", insResult)
	t.Logf("选用保单（来自 proInsurance）: policyId=%s userId=%s policyStatus=%s",
		insParsed.PolicyID, insParsed.UserID, insParsed.PolicyStatus)

	policyID := insParsed.PolicyID

	t.Run("getPolicyInfoByPolicyId", func(t *testing.T) {
		const path = "/getPolicyInfoByPolicyId"
		body := huaan.BuildRequestBody(env.ChannelCode, map[string]interface{}{
			"policyId": policyID,
		})
		huaan.LogRequest(t, path, body)
		result, err := client.Call(ctx, path, body)
		if err != nil {
			t.Fatalf("Call %s: %v", path, err)
		}
		huaan.LogResponse(t, path, result)
	})

	t.Run("getProductPricesByPolicyId", func(t *testing.T) {
		const path = "/getProductPricesByPolicyId"
		body := huaan.BuildRequestBody(env.ChannelCode, map[string]interface{}{
			"policyId":          policyID,
			"hasSocialSecurity": defaultHasSocialSecurity,
			"phoneNo":           env.Phone,
			"name":              env.Name,
			"idCard":            env.IDCard,
		})
		huaan.LogRequest(t, path, body)
		result, err := client.Call(ctx, path, body)
		if err != nil {
			t.Fatalf("Call %s: %v", path, err)
		}
		huaan.LogResponse(t, path, result)
	})

	t.Run("upGradeIns", func(t *testing.T) {
		const path = "/upGradeIns"
		body := huaan.BuildRequestBody(env.ChannelCode, map[string]interface{}{
			"policyId": policyID,
		})
		huaan.LogRequest(t, path, body)
		result, err := client.Call(ctx, path, body)
		if err != nil {
			t.Fatalf("Call %s: %v", path, err)
		}
		huaan.LogResponse(t, path, result)
	})
}

// TestHuaAnDirect_GetSignUrl 获取签约链接：getProductInfoByChannel → proInsurance（policyId/userId）
// → getBankList（bankCode+payChannelId）→ getSignUrl。
//
//	go test -tags=integration ./internal/huaan/ -v -count=1 -run TestHuaAnDirect_GetSignUrl
func TestHuaAnDirect_GetSignUrl(t *testing.T) {
	client, env := huaan.NewTestClient(t)
	if env.Phone == "" || env.Name == "" || env.IDCard == "" {
		t.Skip("跳过：请设置 HUAAN_TEST_PHONE、HUAAN_TEST_NAME、HUAAN_TEST_ID_CARD")
	}

	ctx := t.Context()
	products := fetchChannelProducts(t, ctx, client, env.ChannelCode)
	prod := products[0]

	insResult, insParsed, err := huaan.CallProInsurance(ctx, client, env.ChannelCode, insuranceFields(env, prod.ProductCode))
	if err != nil {
		t.Fatalf("proInsurance: %v", err)
	}
	huaan.LogResponse(t, "/proInsurance", insResult)
	policyID := insParsed.PolicyID
	userID := insParsed.UserID
	if userID == "" {
		userID = os.Getenv("HUAAN_TEST_USER_ID")
	}
	if userID == "" {
		t.Skip("proInsurance 未返回 userId，且未设置 HUAAN_TEST_USER_ID")
	}
	t.Logf("policyId=%s userId=%s（policyId 来自 proInsurance）", policyID, userID)

	bankBody := huaan.BuildRequestBody(env.ChannelCode, nil)
	huaan.LogRequest(t, huaan.BankListPath, bankBody)
	bankResult, err := client.Call(ctx, huaan.BankListPath, bankBody)
	if err != nil {
		t.Fatalf("getBankList: %v", err)
	}
	huaan.LogResponse(t, huaan.BankListPath, bankResult)

	pairs, err := huaan.ParseBankPayPairs(bankResult.Body)
	if err != nil {
		t.Fatalf("parse bank list: %v", err)
	}
	bank, err := huaan.SelectBankPayPair(pairs, os.Getenv("HUAAN_TEST_BANK_CODE"))
	if err != nil {
		t.Fatalf("select bank: %v", err)
	}
	t.Logf("选用银行: bankCode=%s payChannelId=%s bankName=%s",
		bank.BankCode, bank.PayChannelID, bank.BankName)

	cardType := os.Getenv("HUAAN_TEST_CARD_TYPE")
	if cardType == "" {
		cardType = bank.DefaultCardType()
	}

	signFields := map[string]interface{}{
		"policyId":     policyID,
		"userId":       userID,
		"bankCode":     bank.BankCode,
		"payChannelId": bank.PayChannelID,
		"cardType":     cardType,
		"phoneNo":      env.Phone,
		"name":         env.Name,
		"idCard":       env.IDCard,
	}
	signBody := huaan.BuildRequestBody(env.ChannelCode, signFields)
	const signPath = "/getSignUrl"
	huaan.LogRequest(t, signPath, signBody)
	signResult, err := client.Call(ctx, signPath, signBody)
	if err != nil {
		t.Fatalf("getSignUrl: %v", err)
	}
	huaan.LogResponse(t, signPath, signResult)
}

func fetchChannelProducts(t *testing.T, ctx context.Context, client *huaan.Client, channelCode string) []huaan.ChannelProduct {
	t.Helper()
	const channelPath = "/getProductInfoByChannel"
	channelBody := huaan.BuildRequestBody(channelCode, nil)
	huaan.LogRequest(t, channelPath, channelBody)
	channelResult, err := client.Call(ctx, channelPath, channelBody)
	if err != nil {
		t.Fatalf("Call %s: %v", channelPath, err)
	}
	huaan.LogResponse(t, channelPath, channelResult)
	products, err := huaan.ParseChannelProducts(channelResult.Body)
	if err != nil {
		t.Fatalf("parse products: %v", err)
	}
	return products
}

func insuranceFields(env huaan.EnvTestConfig, productCode string) map[string]interface{} {
	return map[string]interface{}{
		"productCode":       productCode,
		"hasSocialSecurity": defaultHasSocialSecurity,
		"phoneNo":           env.Phone,
		"name":              env.Name,
		"idCard":            env.IDCard,
	}
}

// TestHuaAnDirect_APIPathsComplete 确保集成测试覆盖的路径与 APIPaths 一致。
func TestHuaAnDirect_APIPathsComplete(t *testing.T) {
	covered := map[string]bool{
		huaan.BankListPath:               true,
		"/getProductInfoByChannel":       true,
		"/getProductPricesByProductCode": true, // TestHuaAnDirect_ProductFromChannel
		"/getPolicyInfoByPhoneNo":        true,
		"/getPolicyInfoByPolicyId":       true, // TestHuaAnDirect_PolicyFromProInsurance
		"/getProductPricesByPolicyId":    true, // TestHuaAnDirect_PolicyFromProInsurance
		"/proInsurance":                  true, // TestHuaAnDirect_ProductFromChannel
		"/upGradeIns":                    true, // TestHuaAnDirect_PolicyFromProInsurance
		"/getSignUrl":                    true, // TestHuaAnDirect_GetSignUrl
		"/verifyNoCode":                  true,
	}
	for _, p := range huaan.APIPaths {
		if !covered[p] {
			t.Fatalf("integration test missing path: %s", p)
		}
	}
}
