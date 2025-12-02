package session

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, s *Session) error
	FindByID(ctx context.Context, sessionID string) (*Session, error)
	Update(ctx context.Context, s *Session) error
}
