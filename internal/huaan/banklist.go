package huaan

import (
	"encoding/json"
	"fmt"
)

// BankPayPair 银行编码与支付渠道 ID 的配套组合，须来自同一条 getBankList 记录。
type BankPayPair struct {
	BankCode     string
	PayChannelID string
	BankName     string
	DebitCard    string // "1" 支持储蓄卡
	CreditCard   string // "1" 支持信用卡
}

// ParseBankPayPairs 从 getBankList 华安响应 JSON 解析银行列表（仅保留同时含 bankCode、payChannelId 的项）。
func ParseBankPayPairs(respBytes []byte) ([]BankPayPair, error) {
	var resp map[string]interface{}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if ParseCode(respBytes) != 200 {
		return nil, fmt.Errorf("getBankList code=%d", ParseCode(respBytes))
	}
	items := bankItemsFromResponseData(resp["data"])
	out := make([]BankPayPair, 0, len(items))
	for _, m := range items {
		p := pairFromMap(m)
		if p.BankCode == "" || p.PayChannelID == "" {
			continue
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no bank with bankCode and payChannelId in getBankList data")
	}
	return out, nil
}

// SelectBankPayPair 按 bankCode 选取配套银行；bankCode 为空时取列表第一项。
func SelectBankPayPair(pairs []BankPayPair, bankCode string) (BankPayPair, error) {
	if len(pairs) == 0 {
		return BankPayPair{}, fmt.Errorf("empty bank list")
	}
	if bankCode == "" {
		return pairs[0], nil
	}
	for _, p := range pairs {
		if p.BankCode == bankCode {
			return p, nil
		}
	}
	return BankPayPair{}, fmt.Errorf("bankCode %q not found in getBankList (%d banks)", bankCode, len(pairs))
}

// DefaultCardType 根据银行能力返回 getSignUrl 用的 cardType：优先储蓄卡 "1"，否则信用卡 "2"。
func (p BankPayPair) DefaultCardType() string {
	if p.DebitCard == "1" {
		return "1"
	}
	if p.CreditCard == "1" {
		return "2"
	}
	return "1"
}

func bankItemsFromResponseData(data interface{}) []map[string]interface{} {
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

func pairFromMap(m map[string]interface{}) BankPayPair {
	return BankPayPair{
		BankCode:     strVal(m, "bankCode", "bank_code"),
		PayChannelID: strVal(m, "payChannelId", "pay_channel_id"),
		BankName:     strVal(m, "bankName", "bank_name"),
		DebitCard:    strVal(m, "debitCard", "debit_card"),
		CreditCard:   strVal(m, "creditCard", "credit_card"),
	}
}

func strVal(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			switch t := v.(type) {
			case string:
				if t != "" {
					return t
				}
			case float64:
				if t == float64(int64(t)) {
					return fmt.Sprintf("%d", int64(t))
				}
				return fmt.Sprintf("%v", t)
			case json.Number:
				return t.String()
			case bool:
				if t {
					return "1"
				}
				return "0"
			}
		}
	}
	return ""
}
