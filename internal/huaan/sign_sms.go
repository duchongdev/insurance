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
	if apiPath == SmsValidPath {
		return smsValidSignParams(body)
	}
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

// smsValidSignParams 构造 sms/valid 签名参数字段：phoneNo、code 及公共字段。
func smsValidSignParams(body map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for _, k := range []string{"channelCode", "timestamp", "phoneNo", "code"} {
		if v, ok := body[k]; ok {
			out[k] = v
		}
	}
	return out
}

// normalizeSmsValidUpstreamBody 华安 sms/valid 上游字段：phoneNo、code（兼容旧 mobile、smsCode）。
func normalizeSmsValidUpstreamBody(body map[string]interface{}) {
	if v, ok := body["mobile"]; ok {
		if phone, has := body["phoneNo"].(string); !has || phone == "" {
			body["phoneNo"] = v
		}
		delete(body, "mobile")
	}
	if v, ok := body["smsCode"]; ok {
		if code, has := body["code"].(string); !has || code == "" {
			body["code"] = v
		}
		delete(body, "smsCode")
	}
}

// huaanSignSecretKey 返回参与华安签名的密钥。
func huaanSignSecretKey(apiPath string, secretKey string) string {
	return secretKey
}

func buildHuaAnSign(apiPath string, body map[string]interface{}, secretKey string) string {
	return sign.BuildHuaAn(bodyForHuaAnSign(apiPath, body), huaanSignSecretKey(apiPath, secretKey))
}
