package callback

import (
	"net/http"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

type GitHubCallbackHandler struct {
	OAuth         *deps.GitHubOAuthDependencies
	API           *deps.GitHubAPIDependencies
	UserService   UserService
	IDTokenIssuer IDTokenIssuer
	ClientID      string
}

func NewGitHubCallbackHandler(
	oauth *deps.GitHubOAuthDependencies,
	api *deps.GitHubAPIDependencies,
	userSvc UserService,
	issuer IDTokenIssuer,
	clientID string,
) *GitHubCallbackHandler {
	return &GitHubCallbackHandler{
		OAuth:         oauth,
		API:           api,
		UserService:   userSvc,
		IDTokenIssuer: issuer,
		ClientID:      clientID,
	}
}

func (h *GitHubCallbackHandler) ready() bool {
	return h != nil &&
		h.OAuth != nil && h.OAuth.Config != nil && h.OAuth.Logger != nil &&
		h.API != nil && h.API.HTTPClient != nil &&
		h.UserService != nil && h.IDTokenIssuer != nil &&
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
