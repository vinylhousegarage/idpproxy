package token

import (
	"context"
	"time"
)

type AuthCode struct {
	UserID    string
	ClientID  string
	ExpiresAt time.Time
}

type AuthCodeStore interface {
	Consume(ctx context.Context, code string, clientID string) (*AuthCode, error)
}

type Clock interface {
	Now() time.Time
}

type Service struct {
	Store AuthCodeStore
	Clock Clock
}

func validateClientSecret(clientID, secret string) bool {
	return secret == "secret"
}

func generateAccessToken(userID string) string {
	return "access-" + userID
}

func generateRefreshToken() string {
	return "refresh-token"
}

const accessTokenTTL = 3600

func (s *Service) Exchange(
	ctx context.Context,
	req TokenRequest,
) (*TokenResponse, error) {

	if s.Store == nil || s.Clock == nil {
		return nil, ErrInvalidGrant
	}

	if req.GrantType != "authorization_code" {
		return nil, ErrUnsupportedGrantType
	}

	if !validateClientSecret(req.ClientID, req.ClientSecret) {
		return nil, ErrInvalidClient
	}

	ac, err := s.Store.Consume(ctx, req.Code, req.ClientID)
	if err != nil {
		return nil, ErrInvalidGrant
	}

	if s.Clock.Now().After(ac.ExpiresAt) {
		return nil, ErrInvalidGrant
	}

	accessToken := generateAccessToken(ac.UserID)
	refreshToken := generateRefreshToken()

	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    accessTokenTTL,
		RefreshToken: refreshToken,
	}, nil
}
