package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/httperror"
)

type GoogleLoginHandler struct {
	Deps *deps.Dependencies
}

func NewGoogleLoginHandler(di *deps.Dependencies) *GoogleLoginHandler {
	return &GoogleLoginHandler{
		Deps: di,
	}
}

func (h *GoogleLoginHandler) Serve(c *gin.Context) {
	state := GenerateState()
	http.SetCookie(c.Writer, BuildStateCookie(state))

	endpoint, err := GetGoogleLoginURL(h.Deps.MetadataURL, h.Deps.HTTPClient, h.Deps.Logger)
	if err != nil {
		httperror.WriteErrorResponse(c.Writer, err, h.Logger)
		return
	}

	loginURL, err := BuildGoogleLoginURL(h.Deps.Config, endpoint, state)
	if err != nil {
		httperror.WriteErrorResponse(c.Writer, err, h.Logger)
		return
	}

	h.Deps.Logger.Info("redirecting to Google login",
		zap.String("url", loginURL),
		zap.String("state", state),
	)

	c.Redirect(http.StatusFound, loginURL)
}
