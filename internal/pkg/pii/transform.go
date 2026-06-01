package pii

import (
	"github.com/huaan/insurance-bridge/internal/pkg/cipher"
)

// Transformer 在渠道与华安之间转换三要素：入站解密（渠道密文 → 明文给华安），出站加密（华安明文 → 密文给渠道）。
type Transformer struct {
	c *cipher.AES256GCM
}

// NewTransformer 使用与库内字段相同的 AES-256-GCM 实例构造 PII 转换器。
func NewTransformer(c *cipher.AES256GCM) *Transformer {
	return &Transformer{c: c}
}

// DecryptRequest 解密渠道请求体顶层三要素及 insuredList 内嵌字段，原地修改 body。
// 若某字段为明文（非 Base64 密文），Decrypt 失败则跳过该字段，兼容混合传参。
func (t *Transformer) DecryptRequest(body map[string]interface{}) error {
	for _, field := range RequestFields {
		v, ok := body[field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		plain, err := t.c.Decrypt(s)
		if err != nil {
			return err
		}
		body[field] = plain
	}
	return t.decryptNested(body)
}

// decryptNested 处理请求体 insuredList 数组中被渠道加密的被保人字段。
func (t *Transformer) decryptNested(body map[string]interface{}) error {
	if list, ok := body["insuredList"].([]interface{}); ok {
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				for _, f := range []string{"insuredCardNo", "insuredName", "phoneNo"} {
					if v, ok := m[f].(string); ok && v != "" {
						// 渠道侧若加密则解密；明文则 Decrypt 会失败，跳过
						if plain, err := tryDecrypt(t, v); err == nil {
							m[f] = plain
						}
					}
				}
			}
		}
	}
	return nil
}

// tryDecrypt 尝试解密单字段，失败时由调用方决定是否忽略。
func tryDecrypt(t *Transformer, v string) (string, error) {
	return t.c.Decrypt(v)
}

// EncryptResponse 对返回渠道的响应体 data 对象（及 data 数组元素）中的 PII 字段重新加密。
func (t *Transformer) EncryptResponse(body map[string]interface{}) error {
	encryptMapFields(body, "data", t.c)
	if dataArr, ok := body["data"].([]interface{}); ok {
		for _, item := range dataArr {
			if m, ok := item.(map[string]interface{}); ok {
				encryptObjectPII(m, t.c)
			}
		}
	}
	return nil
}

// encryptMapFields 当 root[key] 为 map 时，对其执行 PII 加密。
func encryptMapFields(root map[string]interface{}, key string, c *cipher.AES256GCM) {
	data, ok := root[key]
	if !ok {
		return
	}
	if m, ok := data.(map[string]interface{}); ok {
		encryptObjectPII(m, c)
	}
}

// encryptObjectPII 加密单个 JSON 对象内的 phoneNo/name/idCard/userid 及 insuredList 嵌套项。
func encryptObjectPII(m map[string]interface{}, c *cipher.AES256GCM) {
	fields := []string{"phoneNo", "name", "idCard", "userid", "insuredCardNo", "insuredName"}
	for _, f := range fields {
		if v, ok := m[f].(string); ok && v != "" {
			if enc, err := c.Encrypt(v); err == nil {
				m[f] = enc
			}
		}
	}
	if list, ok := m["insuredList"].([]interface{}); ok {
		for _, item := range list {
			if im, ok := item.(map[string]interface{}); ok {
				encryptObjectPII(im, c)
			}
		}
	}
}
