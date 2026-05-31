package cipher

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	c, err := New("0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}
	plain := "张三"
	enc, err := c.Encrypt(plain)
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.Decrypt(enc)
	if err != nil {
		t.Fatal(err)
	}
	if got != plain {
		t.Fatalf("got %s want %s", got, plain)
	}
}

func TestInvalidKey(t *testing.T) {
	if _, err := New("short"); err != ErrInvalidKey {
		t.Fatal("expected invalid key error")
	}
}
