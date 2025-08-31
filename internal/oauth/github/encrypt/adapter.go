package encrypt

import (
	"fmt"

	"github.com/fernet/fernet-go"

	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/store"
)

type FernetAdapter struct {
	key *fernet.Key
	ttl uint
}

func NewFernetAdapter(key *fernet.Key, ttl uint) (*FernetAdapter, error) {
	if key == nil {
		return nil, ErrNilKey
	}

	return &FernetAdapter{key: key, ttl: ttl}, nil
}

func (a *FernetAdapter) EncryptString(plain string) (string, error) {
	return Encrypt(a.key, []byte(plain))
}

func (a *FernetAdapter) DecryptString(token string) (string, error) {
	if token == "" {
		return "", ErrEmptyToken
	}

	b, err := DecryptWithTTL(a.key, token, a.ttl)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(b), nil
}

var _ store.TokenEncryptor = (*FernetAdapter)(nil)
