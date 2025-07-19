package fernet

import (
	"fmt"
	"time"

	"github.com/fernet/fernet-go"
)

func Encrypt(key *fernet.Key, plaintext []byte) (string, error) {
	token, err := fernet.EncryptAndSign(plaintext, key)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func DecryptWithTTL(key *fernet.Key, token string, ttl uint) ([]byte, error) {
	if key == nil {
		return nil, fmt.Errorf("fernet key is nil")
	}

	ttlDuration := time.Duration(ttl) * time.Second

	plaintext := fernet.VerifyAndDecrypt(
		[]byte(token),
		ttlDuration,
		[]*fernet.Key{key},
	)
	if plaintext == nil {
		return nil, fmt.Errorf("invalid or expired token")
	}

	return plaintext, nil
}