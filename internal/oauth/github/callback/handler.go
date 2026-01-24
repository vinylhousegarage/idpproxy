package callback

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	githubtoken "github.com/vinylhousegarage/idpproxy/internal/oauth/github/token"
	githubuser "github.com/vinylhousegarage/idpproxy/internal/oauth/github/user"
	"go.uber.org/zap"
)

func callbackSuccessLocation(code, qState string) string {
	return "/oauth/github/callback/success?code=" + url.QueryEscape(code) +
		"&state=" + url.QueryEscape(qState)
}

func (h *GitHubCallbackHandler) Serve(c *gin.Context) {
	if !h.ready() {
		h.notReady(c.Writer)

		return
	}

	code := c.Query("code")
	qState := c.Query("state")

	if code == "" {
		h.OAuth.Logger.Warn("missing code")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})

		return
	}
	if qState == "" {
		h.OAuth.Logger.Warn("missing state")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing state"})

		return
	}

	cookie, err := c.Request.Cookie(stateCookieName)
	if err != nil || cookie == nil || cookie.Value == "" || cookie.Value != qState {
		h.OAuth.Logger.Warn("state mismatch or missing cookie",
			zap.String("query_state", qState),
			zap.String("cookie_state", safeCookieVal(cookie)),
			zap.Error(err),
		)
		http.SetCookie(c.Writer, deleteStateCookie())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})

		return
	}

	http.SetCookie(c.Writer, deleteStateCookie())

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	req, err := githubtoken.BuildAccessTokenRequest(ctx, h.OAuth.Config, code, qState)
	if err != nil {
		h.OAuth.Logger.Error("build token request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "build request failed"})

		return
	}

	resp, err := h.API.HTTPClient.Do(req)
	if err != nil {
		h.OAuth.Logger.Error("token request failed", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "token request failed"})

		return
	}

	accessToken, err := githubtoken.ExtractAccessTokenFromResponse(resp)
	if err != nil {

		h.OAuth.Logger.Warn("token response parse failed", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "token exchange failed"})

		return
	}

	userReq, err := githubuser.NewGitHubUserRequest(ctx, accessToken)
	if err != nil {
		h.OAuth.Logger.Error("build /user request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build user request"})

		return
	}

	userResp, err := h.API.HTTPClient.Do(userReq)
	if err != nil {
		h.OAuth.Logger.Error("call /user failed", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to call GitHub /user"})

		return
	}

	ghUser, err := githubuser.DecodeGitHubUserResponse(userResp)
	if err != nil {
		h.OAuth.Logger.Warn("decode /user failed", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to decode GitHub /user"})

		return
	}

	_, err = h.UserService.UpsertFromGitHub(ctx, ghUser.ID, ghUser.Login, ghUser.Email)
	if err != nil {
		h.OAuth.Logger.Error("upsert user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upsert user"})

		return
	}

	c.Redirect(http.StatusFound, callbackSuccessLocation(code, qState))
}
