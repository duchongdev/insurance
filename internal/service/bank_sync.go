package service

import (
	"encoding/json"
	"strconv"

	"github.com/huaan/insurance-bridge/internal/huaan"
	"github.com/huaan/insurance-bridge/internal/model"
	"github.com/huaan/insurance-bridge/internal/repository"
)

// ReplaceBanksFromResponse 解析华安 getBankList 响应并全量覆盖 bank_info_t。
func ReplaceBanksFromResponse(banks *repository.BankRepo, respBytes []byte) error {
	if !huaan.ResponseOK(respBytes) {
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
		isAct := int(b.IsActBank)
		if isAct == 0 && b.Status == 1 {
			isAct = 1
		}
		items = append(items, map[string]interface{}{
			"bankName":     b.BankName,
			"bankCode":     b.BankCode,
			"debitCard":    cardFlagStr(b.DebitCard),
			"creditCard":   cardFlagStr(b.CreditCard),
			"payChannelId": b.PayChannelID,
			"isActBank":    isAct,
		})
	}
	return json.Marshal(map[string]interface{}{
		"code":    200,
		"message": "成功",
		"data":    items,
	})
}

func cardFlagStr(v int8) string {
	if v == 1 {
		return "1"
	}
	return "0"
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
		BankCode:     firstStr(m, "bankCode", "bank_code"),
		BankName:     firstStr(m, "bankName", "bank_name"),
		DebitCard:    firstInt8(m, "debitCard", "debit_card"),
		CreditCard:   firstInt8(m, "creditCard", "credit_card"),
		PayChannelID: firstStr(m, "payChannelId", "pay_channel_id"),
		IsActBank:    firstInt8(m, "isActBank", "is_act_bank"),
		Status:       firstInt8(m, "status"),
	}
	if _, ok := m["isActBank"]; ok {
		if b.IsActBank == 1 {
			b.Status = 1
		} else {
			b.Status = 2
		}
	} else if b.Status == 0 {
		b.Status = 1
		b.IsActBank = 1
	}
	if b.IsActBank == 0 && b.Status == 1 {
		b.IsActBank = 1
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

func str(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case json.Number:
		return t.String()
	default:
		return ""
	}
}
