package callback

import (
	"context"

	"github.com/vinylhousegarage/idpproxy/internal/idtoken"
)

type UserService interface {
	UpsertFromGitHub(ctx context.Context, ghID int64, login, email string) (string, error)
}

type IDTokenIssuer interface {
	Issue(ctx context.Context, in *idtoken.IDTokenInput) (jwt string, kid string, err error)
}
