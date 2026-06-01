package service

import (
	"encoding/json"
	"testing"
)

func TestBankItemsFromData(t *testing.T) {
	arr := bankItemsFromData([]interface{}{
		map[string]interface{}{"bankCode": "ICBC", "bankName": "工商银行"},
	})
	if len(arr) != 1 || arr[0]["bankCode"] != "ICBC" {
		t.Fatalf("unexpected array parse: %+v", arr)
	}

	nested := bankItemsFromData(map[string]interface{}{
		"bankList": []interface{}{
			map[string]interface{}{"bankCode": "ABC", "bankName": "农业银行"},
		},
	})
	if len(nested) != 1 || nested[0]["bankCode"] != "ABC" {
		t.Fatalf("unexpected nested parse: %+v", nested)
	}
}

func TestBankFromMap(t *testing.T) {
	b := bankFromMap(map[string]interface{}{
		"bankCode":   "CCB",
		"bankName":   "建设银行",
		"debitCard":  float64(1),
		"creditCard": float64(0),
	})
	if b.BankCode != "CCB" || b.DebitCard != 1 || b.Status != 1 {
		t.Fatalf("unexpected bank: %+v", b)
	}
}

func TestBuildBankListResponseBytes(t *testing.T) {
	raw, err := BuildBankListResponseBytes(banksFromData([]interface{}{
		map[string]interface{}{"bankCode": "ICBC", "bankName": "工商银行", "debitCard": float64(1)},
	}))
	if err != nil {
		t.Fatal(err)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatal(err)
	}
	if int(resp["code"].(float64)) != 200 {
		t.Fatalf("unexpected code: %v", resp["code"])
	}
	data, ok := resp["data"].([]interface{})
	if !ok || len(data) != 1 {
		t.Fatalf("unexpected data: %v", resp["data"])
	}
}
