// Package sign 实现华安接口文档约定的 MD5 签名：参数按 key 的 ASCII 升序拼接为 k=v&...，再 MD5 为 32 位小写 hex。
package sign

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Build 按华安文档规则生成 MD5 签名。
// 规则：排除 sign 字段；其余参数按 key 排序拼接；若 params 无 key 则追加 key=secretKey；对拼接串做 MD5。
func Build(params map[string]interface{}, secretKey string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
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
		b.WriteString(formatValue(params[k]))
	}
	// 文档示例中 key 参与签名；若请求体未带 key 则补上密钥
	if _, ok := params["key"]; !ok && secretKey != "" {
		if b.Len() > 0 {
			b.WriteByte('&')
		}
		b.WriteString("key=")
		b.WriteString(secretKey)
	}

	sum := md5.Sum([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

// Verify 校验请求中的 sign 是否与使用 secretKey 重新计算的结果一致。
func Verify(params map[string]interface{}, secretKey string) bool {
	got, ok := params["sign"].(string)
	if !ok || got == "" {
		return false
	}
	return Build(params, secretKey) == got
}

// formatValue 将 JSON 反序列化后的值格式化为签名字符串（与华安服务端规则对齐）。
func formatValue(v interface{}) string {
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
	default:
		return fmt.Sprintf("%v", val)
	}
}

// MapFromJSON 浅拷贝 JSON 对象为 map，避免在验签过程中修改原始 body。
func MapFromJSON(raw map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(raw))
	for k, v := range raw {
		out[k] = v
	}
	return out
}
