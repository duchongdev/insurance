package sign

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// BuildHuaAn 华安上游签名：非空参数按 key 字典序 k=v& 拼接，末尾 &key=密钥，MD5 后 32 位大写。
// key/sign 不参与排序段；密钥不写入请求体，仅用于签名。
func BuildHuaAn(params map[string]interface{}, secretKey string) string {
	plain := buildHuaAnPlaintext(params, secretKey)
	sum := md5.Sum([]byte(plain))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// VerifyHuaAnHeader 校验 HTTP 头 sign 与请求体参数是否匹配华安规则。
func VerifyHuaAnHeader(params map[string]interface{}, secretKey, gotSign string) bool {
	if gotSign == "" {
		return false
	}
	want := BuildHuaAn(params, secretKey)
	return strings.EqualFold(gotSign, want)
}

func buildHuaAnPlaintext(params map[string]interface{}, secretKey string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || k == "key" {
			continue
		}
		if isEmptyValue(params[k]) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(formatHuaAnValue(params[k]))
	}
	if secretKey != "" {
		if b.Len() > 0 {
			b.WriteByte('&')
		}
		b.WriteString("key=")
		b.WriteString(secretKey)
	}
	return b.String()
}

func formatHuaAnValue(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return ""
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%v", val)
	case json.Number:
		return val.String()
	case map[string]interface{}:
		return formatHuaAnObject(val)
	case []interface{}:
		return formatHuaAnArray(val)
	default:
		// 兼容 json.Unmarshal 到 map[string]interface{} 之外的 map 类型
		if m, ok := toStringKeyMap(v); ok {
			return formatHuaAnObject(m)
		}
		return fmt.Sprintf("%v", val)
	}
}

func formatHuaAnObject(m map[string]interface{}) string {
	keys := make([]string, 0, len(m))
	for k, v := range m {
		if isEmptyValue(v) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+formatHuaAnValue(m[k]))
	}
	return strings.Join(parts, "&")
}

func formatHuaAnArray(items []interface{}) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		if isEmptyValue(item) {
			continue
		}
		switch v := item.(type) {
		case map[string]interface{}:
			parts = append(parts, formatHuaAnObject(v))
		default:
			if m, ok := toStringKeyMap(v); ok {
				parts = append(parts, formatHuaAnObject(m))
			} else {
				parts = append(parts, formatHuaAnValue(v))
			}
		}
	}
	return strings.Join(parts, ",")
}

func isEmptyValue(v interface{}) bool {
	if v == nil {
		return true
	}
	switch t := v.(type) {
	case string:
		return t == ""
	case []interface{}:
		return len(t) == 0
	case map[string]interface{}:
		return len(t) == 0
	case json.Number:
		return t.String() == ""
	}
	return false
}

func toStringKeyMap(v interface{}) (map[string]interface{}, bool) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, false
	}
	var out map[string]interface{}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, false
	}
	return out, true
}
