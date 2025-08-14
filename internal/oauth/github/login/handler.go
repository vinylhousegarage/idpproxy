package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

type GitHubLoginHandler struct {
	Deps *deps.GitHubOAuthDependencies
}

func NewGitHubLoginHandler(githubDeps *deps.GitHubOAuthDependencies) *GitHubLoginHandler {
	return &GitHubLoginHandler{
		Deps: githubDeps,
	}
}

func (h *GitHubLoginHandler) Serve(c *gin.Context) {
	state := GenerateState()
	http.SetCookie(c.Writer, BuildStateCookie(state))

	loginURL := BuildGitHubLoginURL(h.Deps.Config, state)

	h.Deps.Logger.Info("redirecting to GitHub login",
		zap.String("url", loginURL),
		zap.String("state", state),
	)

	c.Redirect(http.StatusFound, loginURL)
}
