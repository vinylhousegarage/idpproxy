package idtoken

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type hmacSigner interface {
	Sign(
		ctx context.Context,
		payload []byte,
	) (token string, kid string, err error)

	Alg() string
}

type HMACSignerAdapter struct {
	signer hmacSigner
}

func NewHMACSignerAdapter(
	signer hmacSigner,
) *HMACSignerAdapter {
	return &HMACSignerAdapter{
		signer: signer,
	}
}

func (a *HMACSignerAdapter) SignJWT(
	ctx context.Context,
	payload map[string]any,
) (string, string, error) {
	if a == nil || a.signer == nil {
		return "", "", errors.New("idtoken: nil HMAC signer")
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("marshal JWT payload: %w", err)
	}

	return a.signer.Sign(ctx, payloadJSON)
}

func (a *HMACSignerAdapter) Algorithm() string {
	if a == nil || a.signer == nil {
		return ""
	}

	return a.signer.Alg()
}
