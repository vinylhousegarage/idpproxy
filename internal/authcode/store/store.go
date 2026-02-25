package store

import (
	"context"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type Store interface {
	Save(ctx context.Context, authCode authcode.AuthCode) error
	Consume(ctx context.Context, authCodeValue, clientID string) (string, error)
}
