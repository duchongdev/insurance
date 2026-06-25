package service

import (
	"testing"

	"github.com/huaan/insurance-bridge/internal/model"
)

func TestValidateChannelConfig(t *testing.T) {
	if err := ValidateChannelConfig(&model.Channel{
		ChannelCode: "CH1",
		CallbackURL: "https://channel.example/callback",
	}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateChannelConfig(&model.Channel{ChannelCode: "CH1"}); err == nil {
		t.Fatal("expected callbackUrl error")
	}
	if err := ValidateChannelConfig(&model.Channel{
		ChannelCode: "CH1",
		CallbackURL: "ftp://bad",
	}); err == nil {
		t.Fatal("expected invalid callbackUrl")
	}
}

func TestValidateInsureNotifyBody(t *testing.T) {
	if err := validateInsureNotifyBody(map[string]interface{}{
		"productName": "综合意外险", "productType": 1, "insureTime": "2026-06-05 10:10:10",
		"policyNo": "P1", "policyId": "pid", "amount": 95.0,
		"policyStartDate": "2026-06-05", "policyEndDate": "2027-06-05",
		"payStage": 1, "channelCode": "H50001", "timestamp": 1780701234123,
	}); err != nil {
		t.Fatal(err)
	}
	if err := validateInsureNotifyBody(map[string]interface{}{"channelCode": "X"}); err == nil {
		t.Fatal("expected missing fields error")
	}
}
