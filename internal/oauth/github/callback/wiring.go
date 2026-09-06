package callback

import (
	"context"
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	githubstore "github.com/vinylhousegarage/idpproxy/internal/oauth/github/store"
)

type GitHubTokenRepository interface {
	Upsert(
		ctx context.Context,
		rec *githubstore.GitHubTokenRecord,
	) error
}

type GitHubCallbackHandler struct {
	OAuth            *deps.GitHubOAuthDependencies
	API              *deps.GitHubAPIDependencies
	UserService      UserService
	ProxyCodeService ProxyCodeService
	TokenRepo        GitHubTokenRepository
	ClientID         string
}

func NewGitHubCallbackHandler(
	oauth *deps.GitHubOAuthDependencies,
	api *deps.GitHubAPIDependencies,
	userSvc UserService,
	proxyCodeSvc ProxyCodeService,
	tokenRepo GitHubTokenRepository,
	clientID string,
) *GitHubCallbackHandler {
	return &GitHubCallbackHandler{
		OAuth:            oauth,
		API:              api,
		UserService:      userSvc,
		ProxyCodeService: proxyCodeSvc,
		TokenRepo:        tokenRepo,
		ClientID:         clientID,
	}
}

func (h *GitHubCallbackHandler) ready() bool {
	return h != nil &&
		h.OAuth != nil &&
		h.OAuth.Config != nil &&
		h.OAuth.Logger != nil &&
		h.API != nil &&
		h.API.HTTPClient != nil &&
		h.UserService != nil &&
		h.ProxyCodeService != nil &&
		h.TokenRepo != nil &&
		h.ClientID != ""
}

func (h *GitHubCallbackHandler) notReady(w http.ResponseWriter) {
	if h != nil && h.OAuth != nil && h.OAuth.Logger != nil {
		h.OAuth.Logger.Error("handler dependencies not satisfied")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	_, _ = w.Write([]byte(`{"error":"server not ready"}`))
}
