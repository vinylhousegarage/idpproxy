package callback

import (
	"context"

	"github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

type UserService interface {
	UpsertFromGitHub(
		ctx context.Context,
		githubID int64,
		login string,
		email string,
	) (string, error)
}

type ProxyCodeService interface {
	Issue(
		ctx context.Context,
		userID string,
		clientID string,
	) (string, error)
}

type SessionService interface {
	Start(
		ctx context.Context,
		userID string,
	) (*session.Session, error)
}
