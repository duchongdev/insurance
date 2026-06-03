package huaan

import "testing"

func TestParseProInsuranceResponse(t *testing.T) {
	raw := []byte(`{
		"code": 200,
		"message": "操作成功",
		"data": {
			"policyId": "9697d878b6474c619fa43ffa9ca3b4b0",
			"policyStatus": "0",
			"userId": "5c819c252ab343e2a2a6df98b3a014ab"
		}
	}`)
	got, err := ParseProInsuranceResponse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.PolicyID != "9697d878b6474c619fa43ffa9ca3b4b0" || got.UserID == "" {
		t.Fatalf("unexpected: %+v", got)
	}
}
