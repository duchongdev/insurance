package cipher

import "crypto/rand"

// cipherReadRand 从 crypto/rand 读取随机字节，供 GCM nonce 使用；单独函数便于测试注入。
func cipherReadRand(b []byte) (int, error) {
	return rand.Read(b)
}
