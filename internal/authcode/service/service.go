package service

import (
	"context"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
	"github.com/vinylhousegarage/idpproxy/internal/authcode/store"
)

type Service struct {
	store store.Store
}

func (s *Service) Issue(
	ctx context.Context,
	userID string,
	clientID string,
) error {

	code := authcode.AuthCode{
		Code:      "dummy",
		UserID:    userID,
		ClientID:  clientID,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return s.store.Save(ctx, code)
}
