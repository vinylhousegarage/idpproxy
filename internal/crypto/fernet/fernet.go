package fernet

import (
	"fmt"

	"github.com/fernet/fernet-go"
)

func Encrypt(key *fernet.Key, plaintext []byte) (string, error) {
	token := fernet.EncryptAndSign(plaintext, key)
	return string(token), nil
}

func DecryptWithTTL(key *fernet.Key, token string, ttl uint) ([]byte, error) {
	if key == nil {
		return nil, fmt.Errorf("fernet key is nil")
	}

	plaintext := fernet.VerifyAndDecrypt([]byte(token), ttl, []*fernet.Key{key})
	if plaintext == nil {
		return nil, fmt.Errorf("invalid or expired token")
	}

	return plaintext, nil
}
