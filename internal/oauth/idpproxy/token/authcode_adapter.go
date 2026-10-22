package token

import (
	"context"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
)

type proxyCodeConsumer interface {
	Consume(
		ctx context.Context,
		proxyCode string,
		clientID string,
	) (authcode.ProxyCode, error)
}

type AuthCodeServiceAdapter struct {
	consumer proxyCodeConsumer
}

func NewAuthCodeServiceAdapter(
	consumer proxyCodeConsumer,
) *AuthCodeServiceAdapter {
	return &AuthCodeServiceAdapter{
		consumer: consumer,
	}
}

func (a *AuthCodeServiceAdapter) Consume(
	ctx context.Context,
	code string,
	clientID string,
) (*AuthCode, error) {
	proxyCode, err := a.consumer.Consume(ctx, code, clientID)
	if err != nil {
		return nil, err
	}

	return &AuthCode{
		UserID:    proxyCode.UserID,
		ClientID:  proxyCode.ClientID,
		ExpiresAt: proxyCode.ExpiresAt,
	}, nil
}
