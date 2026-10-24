package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vinylhousegarage/idpproxy/internal/auth/idtoken"
)

type AuthCode struct {
	UserID    string
	ClientID  string
	ExpiresAt time.Time
}

type AuthCodeStore interface {
	Consume(
		ctx context.Context,
		code string,
		clientID string,
	) (*AuthCode, error)
}

type Clock interface {
	Now() time.Time
}

type IDTokenIssuer interface {
	Issue(
		ctx context.Context,
		in *idtoken.IDTokenInput,
	) (token string, kid string, err error)
}

type Service struct {
	Store         AuthCodeStore
	Clock         Clock
	IDTokenIssuer IDTokenIssuer
	IDTokenTTL    time.Duration
}

func (s *Service) Exchange(
	ctx context.Context,
	req TokenRequest,
) (*TokenResponse, error) {
	if s == nil ||
		s.Store == nil ||
		s.Clock == nil ||
		s.IDTokenIssuer == nil ||
		s.IDTokenTTL <= 0 {
		return nil, errors.New("token: invalid service configuration")
	}

	if req.GrantType != "authorization_code" {
		return nil, ErrUnsupportedGrantType
	}

	ac, err := s.Store.Consume(ctx, req.Code, req.ClientID)
	if err != nil {
		return nil, ErrInvalidGrant
	}

	now := s.Clock.Now()

	if ac == nil ||
		ac.UserID == "" ||
		ac.ClientID == "" ||
		ac.ClientID != req.ClientID ||
		!ac.ExpiresAt.After(now) {
		return nil, ErrInvalidGrant
	}

	idToken, _, err := s.IDTokenIssuer.Issue(ctx, &idtoken.IDTokenInput{
		UserID:   ac.UserID,
		ClientID: ac.ClientID,
		Now:      now,
		TTL:      s.IDTokenTTL,
	})
	if err != nil {
		return nil, fmt.Errorf("issue ID token: %w", err)
	}

	return &TokenResponse{
		IDToken: idToken,
	}, nil
}
