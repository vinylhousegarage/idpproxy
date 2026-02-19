package callback

import (
	"context"
)

type UserService interface {
	UpsertFromGitHub(ctx context.Context, ghID int64, login, email string) (string, error)
}

type AuthCodeService interface {
	Issue(ctx context.Context, userID string, clientID string) (string, error)
}
