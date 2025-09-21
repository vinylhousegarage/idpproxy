package jwt

import (
	"context"
)

type Signer interface {
	Sign(ctx context.Context, payload []byte) (token string, kid string, err error)
	Verify(ctx context.Context, token string) (payload []byte, kid string, err error)
	Alg() string
	KeyID() string
}

type Claims struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
	Gen int    `json:"gen"`
}
