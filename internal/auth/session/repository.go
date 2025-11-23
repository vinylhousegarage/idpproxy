package session

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, s *Session) error
}
