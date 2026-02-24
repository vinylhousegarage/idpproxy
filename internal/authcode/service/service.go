package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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
) (string, error) {

	authCode, err := generateAuthCode()
	if err != nil {
		return "", err
	}

	ac := authcode.AuthCode{
		Code:      authCode,
		UserID:    userID,
		ClientID:  clientID,
		ExpiresAt: time.Now().Add(60 * time.Second),
	}

	if err := s.store.Save(ctx, ac); err != nil {
		return "", err
	}

	return authCode, nil
}

func (s *Service) Consume(
	ctx context.Context,
	authCode string,
	clientID string,
) (string, error) {
	return s.store.Consume(ctx, authCode, clientID)
}

func generateAuthCode() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
