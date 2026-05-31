// Package cipher 提供 AES-256-GCM 加解密及日志脱敏工具，用于渠道三要素传输与库内敏感字段存储。
package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	// ErrInvalidKey 密钥长度不是 32 字节。
	ErrInvalidKey = errors.New("encryption key must be 32 bytes")
	// ErrInvalidInput 密文格式非法或认证标签校验失败。
	ErrInvalidInput = errors.New("invalid ciphertext")
)

// AES256GCM 封装 AES-256-GCM，密文格式：Base64(nonce || ciphertext+tag)。
type AES256GCM struct {
	gcm cipher.AEAD
}

// New 使用 32 字节字符串密钥构造加解密器（配置 security.data_encryption_key）。
func New(key string) (*AES256GCM, error) {
	k := []byte(key)
	if len(k) != 32 {
		return nil, ErrInvalidKey
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}
	return &AES256GCM{gcm: gcm}, nil
}

// Encrypt 将明文加密为 Base64 字符串；空串直接返回空。
func (c *AES256GCM) Encrypt(plain string) (string, error) {
	if plain == "" {
		return "", nil
	}
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := cipherReadRand(nonce); err != nil {
		return "", err
	}
	ciphertext := c.gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 将 Base64 密文解密为明文；空串直接返回空。
func (c *AES256GCM) Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrInvalidInput
	}
	nonceSize := c.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrInvalidInput
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plain, err := c.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrInvalidInput
	}
	return string(plain), nil
}

// Mask 对字符串做日志脱敏（保留首尾各 2 字符，中间替换为 ****）。
func Mask(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}
