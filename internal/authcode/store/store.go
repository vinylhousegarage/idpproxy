package store

import (
	"context"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type Store interface {
	Save(ctx context.Context, code authcode.AuthCode) error
	Consume(ctx context.Context, code string, clientID string) (string, error)
}
