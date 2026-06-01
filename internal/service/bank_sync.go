package service

import (
	"encoding/json"

	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/repository"
)

// ReplaceBanksFromResponse 解析华安 getBankList 响应并全量覆盖 bank_info_t。
func ReplaceBanksFromResponse(banks *repository.BankRepo, respBytes []byte) error {
	if !huaAnResponseOK(respBytes) {
		return nil
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil
	}
	list := banksFromData(resp["data"])
	return banks.ReplaceAll(list)
}

// BuildBankListResponseBytes 将库内银行列表构造成华安标准响应 JSON。
func BuildBankListResponseBytes(banks []model.BankInfo) ([]byte, error) {
	items := make([]map[string]interface{}, 0, len(banks))
	for _, b := range banks {
		items = append(items, map[string]interface{}{
			"bankCode":   b.BankCode,
			"bankName":   b.BankName,
			"debitCard":  b.DebitCard,
			"creditCard": b.CreditCard,
			"status":     b.Status,
		})
	}
	return json.Marshal(map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    items,
	})
}

func banksFromData(data interface{}) []model.BankInfo {
	items := bankItemsFromData(data)
	out := make([]model.BankInfo, 0, len(items))
	for _, item := range items {
		b := bankFromMap(item)
		if b.BankCode == "" {
			continue
		}
		out = append(out, *b)
	}
	return out
}

func bankItemsFromData(data interface{}) []map[string]interface{} {
	switch v := data.(type) {
	case []interface{}:
		return mapsFromSlice(v)
	case map[string]interface{}:
		for _, key := range []string{"bankList", "list", "banks"} {
			if arr, ok := v[key].([]interface{}); ok {
				return mapsFromSlice(arr)
			}
		}
		if _, ok := v["bankCode"]; ok {
			return []map[string]interface{}{v}
		}
	}
	return nil
}

func mapsFromSlice(items []interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

func bankFromMap(m map[string]interface{}) *model.BankInfo {
	b := &model.BankInfo{
		BankCode:   firstStr(m, "bankCode", "bank_code"),
		BankName:   firstStr(m, "bankName", "bank_name"),
		DebitCard:  firstInt8(m, "debitCard", "debit_card"),
		CreditCard: firstInt8(m, "creditCard", "credit_card"),
		Status:     firstInt8(m, "status"),
	}
	if b.Status == 0 {
		b.Status = 1
	}
	return b
}

func firstStr(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v := str(m[k]); v != "" {
			return v
		}
	}
	return ""
}

func firstInt8(m map[string]interface{}, keys ...string) int8 {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			return toInt8(v)
		}
	}
	return 0
}

func toInt8(v interface{}) int8 {
	switch t := v.(type) {
	case float64:
		return int8(t)
	case int:
		return int8(t)
	case int64:
		return int8(t)
	case json.Number:
		n, _ := t.Int64()
		return int8(n)
	case string:
		if t == "1" || t == "true" {
			return 1
		}
	}
	return 0
}
