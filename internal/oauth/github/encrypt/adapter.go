package encrypt

import (
	"fmt"
	"strings"
	"time"

	"github.com/fernet/fernet-go"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/store"
)

const fernetPrefix = "fe1:"

type Key struct {
	KID string
	Key *fernet.Key
}

type FernetAdapter struct {
	keys map[string]*fernet.Key
	currentKID string
	ttl        uint
}

func NewFernetAdapter(key *fernet.Key, ttl uint) (*FernetAdapter, error) { // 互換用
    if key == nil { return nil, ErrNilKey }
    return NewFernetAdapterWithKeys([]Key{{KID:"default", Key:key}}, ttl)
}

func NewFernetAdapterWithKeys(keys []Key, ttl uint) (*FernetAdapter, error) {
    if len(keys) == 0 { return nil, ErrNilKeySet }
    m := make(map[string]*fernet.Key, len(keys))
    for _, k := range keys {
        if k.KID == "" || k.Key == nil { return nil, fmt.Errorf("bad key entry: %q", k.KID) }
        m[k.KID] = k.Key
    }
    return &FernetAdapter{ keys: m, currentKID: keys[len(keys)-1].KID, ttl: ttl }, nil
}

func (a *FernetAdapter) EncryptString(plain string) (store.Ciphertext, error) {
    key := a.keys[a.currentKID]
    if key == nil { return store.Ciphertext{}, ErrUnknownKID }
    token, err := fernet.EncryptAndSign([]byte(plain), key)
    if err != nil { return store.Ciphertext{}, fmt.Errorf("encrypt: %w", err) }
    return store.Ciphertext{ KID: a.currentKID, Blob: fernetPrefix + string(token) }, nil
}

func (a *FernetAdapter) DecryptString(ct store.Ciphertext) (string, error) {
    if ct.Blob == "" { return "", ErrEmptyBlob }
    if !strings.HasPrefix(ct.Blob, fernetPrefix) { return "", ErrBadFormat }
    tok := strings.TrimPrefix(ct.Blob, fernetPrefix)

    if key := a.keys[ct.KID]; key != nil {
        if b, err := a.verifyAndDecrypt(tok, key); err == nil { return string(b), nil }
    }
    for _, key := range a.keys {
        if b, err := a.verifyAndDecrypt(tok, key); err == nil { return string(b), nil }
    }
    return "", fmt.Errorf("%w: %s", ErrUnknownKID, ct.KID)
}

func (a *FernetAdapter) verifyAndDecrypt(token string, key *fernet.Key) ([]byte, error) {
    if a.ttl == 0 {
        b := fernet.VerifyAndDecrypt([]byte(token), 0, []*fernet.Key{key})
        if b == nil { return nil, ErrDecryptFailed }
        return b, nil
    }
    b := fernet.VerifyAndDecrypt([]byte(token), time.Duration(a.ttl)*time.Second, []*fernet.Key{key})
    if b == nil { return nil, ErrDecryptFailed }
    return b, nil
}

var _ store.TokenEncryptor = (*FernetAdapter)(nil)
