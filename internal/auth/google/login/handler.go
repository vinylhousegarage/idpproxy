package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
	"github.com/vinylhousegarage/idpproxy/internal/response"
)

type GoogleLoginHandler struct {
	Config config.GoogleConfig
	Client httpclient.HTTPClient
	Logger *zap.Logger
}

func NewGoogleLoginHandler(cfg config.GoogleConfig, cli httpclient.HTTPClient, logger *zap.Logger) *GoogleLoginHandler {
	return &GoogleLoginHandler{
		Config: cfg,
		Client: cli,
		Logger: logger,
	}
}

func (h *GoogleLoginHandler) Serve(c *gin.Context) {
	state := GenerateState()
	http.SetCookie(c.Writer, BuildStateCookie(state))

	endpoint, err := GetGoogleLoginURL(h.Config.MetadataURL, h.Client, h.Logger)
	if err != nil {
		response.WriteErrorResponse(c.Writer, err, h.Logger)
		return
	}

	loginURL, err := BuildGoogleLoginURL(h.Config, endpoint, state)
	if err != nil {
		response.WriteErrorResponse(c.Writer, err, h.Logger)
		return
	}

	h.Logger.Info("redirecting to Google login",
		zap.String("url", loginURL),
		zap.String("state", state),
	)

	c.Redirect(http.StatusFound, loginURL)
}
