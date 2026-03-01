package store

import (
	"context"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type Store interface {
	Save(ctx context.Context, proxyCode authcode.ProxyCode) error
	Consume(ctx context.Context, proxyCodeValue, clientID string) (string, error)
}
