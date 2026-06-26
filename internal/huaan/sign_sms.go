package huaan

import "github.com/huaan/insurance-bridge/internal/pkg/sign"

// SmsSignExcludeKeys 华安 sms 系列接口签名时 body 中不参与签名的字段（仍随请求发送）。
func SmsSignExcludeKeys(apiPath string) []string {
	switch apiPath {
	case SmsSendPath:
		return []string{"idCard", "name"}
	case SmsNoValidPath:
		return []string{"mobile"}
	default:
		return nil
	}
}

func bodyForHuaAnSign(apiPath string, body map[string]interface{}) map[string]interface{} {
	exclude := SmsSignExcludeKeys(apiPath)
	if len(exclude) == 0 {
		return body
	}
	skip := make(map[string]struct{}, len(exclude))
	for _, k := range exclude {
		skip[k] = struct{}{}
	}
	out := make(map[string]interface{}, len(body))
	for k, v := range body {
		if _, ok := skip[k]; ok {
			continue
		}
		out[k] = v
	}
	return out
}

// huaanSignSecretKey 返回参与华安签名的密钥；空字符串表示明文末尾不拼接 &key=。
// SmsValidPath：临时规则（mobile+smsCode 参与签名、不含 key 后缀），未完成联调，后续可能调整。
func huaanSignSecretKey(apiPath string, secretKey string) string {
	if apiPath == SmsValidPath {
		return ""
	}
	return secretKey
}

func buildHuaAnSign(apiPath string, body map[string]interface{}, secretKey string) string {
	return sign.BuildHuaAn(bodyForHuaAnSign(apiPath, body), huaanSignSecretKey(apiPath, secretKey))
}
