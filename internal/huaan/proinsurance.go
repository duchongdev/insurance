package huaan

import (
	"context"
	"encoding/json"
	"fmt"
)

const proInsurancePath = "/proInsurance"

// ProInsuranceResult proInsurance 成功响应中的业务字段。
type ProInsuranceResult struct {
	PolicyID     string
	UserID       string
	PolicyStatus string
}

// ParseProInsuranceResponse 从 proInsurance 华安响应 JSON 解析 policyId、userId 等。
func ParseProInsuranceResponse(respBytes []byte) (ProInsuranceResult, error) {
	var resp map[string]interface{}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return ProInsuranceResult{}, fmt.Errorf("unmarshal response: %w", err)
	}
	if ParseCode(respBytes) != 200 {
		return ProInsuranceResult{}, fmt.Errorf("proInsurance code=%d", ParseCode(respBytes))
	}
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		return ProInsuranceResult{}, fmt.Errorf("proInsurance data missing or invalid")
	}
	out := ProInsuranceResult{
		PolicyID:     strVal(data, "policyId", "policy_id"),
		UserID:       strVal(data, "userId", "user_id"),
		PolicyStatus: strVal(data, "policyStatus", "policy_status"),
	}
	if out.PolicyID == "" {
		return ProInsuranceResult{}, fmt.Errorf("proInsurance data.policyId empty")
	}
	return out, nil
}

// CallProInsurance 发起预投保；fields 须含 productCode、hasSocialSecurity 及三要素等。
func CallProInsurance(ctx context.Context, client *Client, channelCode string, fields map[string]interface{}) (*CallResult, ProInsuranceResult, error) {
	body := BuildRequestBody(channelCode, fields)
	result, err := client.Call(ctx, proInsurancePath, body)
	if err != nil {
		return nil, ProInsuranceResult{}, err
	}
	parsed, err := ParseProInsuranceResponse(result.Body)
	if err != nil {
		return result, ProInsuranceResult{}, err
	}
	return result, parsed, nil
}
