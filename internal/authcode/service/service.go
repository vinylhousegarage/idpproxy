package service

import (
	"context"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/authcode"
	"github.com/vinylhousegarage/idpproxy/internal/authcode/store"
)

type Service struct {
	store store.Store
	now   func() time.Time
}

func (s *Service) Issue(
	ctx context.Context,
	userID string,
	clientID string,
) error {

	now := s.now
	if now == nil {
		now = time.Now
	}

	code := authcode.AuthCode{
		Code:      "dummy",
		UserID:    userID,
		ClientID:  clientID,
		ExpiresAt: now().Add(5 * time.Minute),
	}

	return s.store.Save(ctx, code)
}

func (s *Service) Consume(
	ctx context.Context,
	code string,
	clientID string,
) (string, error) {

	now := s.now
	if now == nil {
		now = time.Now
	}

	ac, err := s.store.Get(ctx, code, clientID)
	if err != nil {
		return "", err
	}

	if ac.ClientID != clientID {
		return "", ErrClientMismatch
	}

	if ac.ExpiresAt.Before(now()) {
		return "", ErrExpired
	}

	return s.store.Consume(ctx, code, clientID)
}
