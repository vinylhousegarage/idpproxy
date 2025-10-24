package idtoken

import "context"

type Signer interface {
	SignJWT(ctx context.Context, payload map[string]any) (jwt string, kid string, err error)
}
