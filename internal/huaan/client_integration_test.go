//go:build integration

package huaan_test

import (
	"os"
	"testing"

	"github.com/huaan/insurance-bridge/internal/huaan"
)

// 华安直连集成测试：调用本服务真实的 huaan.Client，直连华安环境（默认不携带 key/sign）。
// 运行前设置环境变量，见 .env.huaan.example。
//
//	go test -tags=integration ./internal/huaan/ -v -count=1
func TestHuaAnDirect_AllPaths(t *testing.T) {
	client, env := huaan.NewTestClient(t)

	cases := []struct {
		name       string
		path       string
		needPII    bool
		needFields bool // 为 true 且未配置 HUAAN_TEST_* 业务参数时 skip
		fields     map[string]interface{}
	}{
		{
			name:       "getBankList",
			path:       huaan.BankListPath,
			fields:     nil,
		},
		{
			name:   "getProductInfoByChannel",
			path:   "/getProductInfoByChannel",
			fields: nil,
		},
		{
			name:       "getProductPricesByProductCode",
			path:       "/getProductPricesByProductCode",
			needFields: true,
			fields:     map[string]interface{}{"productCode": os.Getenv("HUAAN_TEST_PRODUCT_CODE")},
		},
		{
			name:       "getPolicyInfoByPhoneNo",
			path:       "/getPolicyInfoByPhoneNo",
			needPII:    true,
			fields:     nil,
		},
		{
			name:       "getPolicyInfoByPolicyId",
			path:       "/getPolicyInfoByPolicyId",
			needFields: true,
			fields:     map[string]interface{}{"policyId": os.Getenv("HUAAN_TEST_POLICY_ID")},
		},
		{
			name:       "getProductPricesByPolicyId",
			path:       "/getProductPricesByPolicyId",
			needFields: true,
			fields:     map[string]interface{}{"policyId": os.Getenv("HUAAN_TEST_POLICY_ID")},
		},
		{
			name:       "proInsurance",
			path:       "/proInsurance",
			needPII:    true,
			needFields: true,
			fields:     map[string]interface{}{"productCode": os.Getenv("HUAAN_TEST_PRODUCT_CODE")},
		},
		{
			name:       "upGradeIns",
			path:       "/upGradeIns",
			needFields: true,
			fields:     map[string]interface{}{"policyId": os.Getenv("HUAAN_TEST_POLICY_ID")},
		},
		{
			name:    "getSignUrl",
			path:    "/getSignUrl",
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
			if tc.needFields {
				for k, v := range fields {
					if s, ok := v.(string); ok && s == "" {
						t.Skipf("跳过：请设置业务参数（如 HUAAN_TEST_PRODUCT_CODE、HUAAN_TEST_POLICY_ID），缺少 %s", k)
					}
				}
			}

			body := huaan.BuildRequestBody(env.ChannelCode, fields)
			result, err := client.Call(t.Context(), tc.path, body, env.HuaAnKey)
			if err != nil {
				t.Fatalf("Call %s: %v", tc.path, err)
			}
			huaan.LogResponse(t, tc.path, result)
		})
	}
}

// TestHuaAnDirect_APIPathsComplete 确保集成测试覆盖的路径与 APIPaths 一致。
func TestHuaAnDirect_APIPathsComplete(t *testing.T) {
	covered := map[string]bool{
		huaan.BankListPath:                  true,
		"/getProductInfoByChannel":          true,
		"/getProductPricesByProductCode":    true,
		"/getPolicyInfoByPhoneNo":           true,
		"/getPolicyInfoByPolicyId":          true,
		"/getProductPricesByPolicyId":       true,
		"/proInsurance":                     true,
		"/upGradeIns":                       true,
		"/getSignUrl":                       true,
		"/verifyNoCode":                     true,
	}
	for _, p := range huaan.APIPaths {
		if !covered[p] {
			t.Fatalf("integration test missing path: %s", p)
		}
	}
}
