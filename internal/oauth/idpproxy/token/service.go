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

	ac, err := s.Store.Consume(ctx, req.Code, req.ClientID)
	if err != nil {
		return nil, ErrInvalidGrant
	}

	if s.Clock.Now().After(ac.ExpiresAt) {
		return nil, ErrInvalidGrant
	}

	idToken := "idtoken-for-" + ac.UserID

	return &TokenResponse{
		IDToken: idToken,
	}, nil
}
