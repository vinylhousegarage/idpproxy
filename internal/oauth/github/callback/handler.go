package callback

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"
	githubtoken "github.com/vinylhousegarage/idpproxy/internal/oauth/github/token"
	githubuser "github.com/vinylhousegarage/idpproxy/internal/oauth/github/user"
	"go.uber.org/zap"
)

func callbackSuccessLocation(proxyCode, qState string) string {
	return "/oauth/github/callback/success?code=" + url.QueryEscape(proxyCode) +
		"&state=" + url.QueryEscape(qState)
}

func (h *GitHubCallbackHandler) Serve(c *gin.Context) {
	if !h.ready() {
		h.notReady(c.Writer)

		return
	}

	githubCode := c.Query("code")
	qState := c.Query("state")

	if githubCode == "" {
		h.OAuth.Logger.Warn("missing githubCode")
		c.JSON(http.StatusBadRequest, gin.H{"error": apierror.ErrorCodeMissingGitHubCode})

		return
	}

	if qState == "" {
		h.OAuth.Logger.Warn("missing state")
		c.JSON(http.StatusBadRequest, gin.H{"error": apierror.ErrorCodeMissingState})

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
		c.JSON(http.StatusBadRequest, gin.H{"error": apierror.ErrorCodeInvalidState})

		return
	}

	http.SetCookie(c.Writer, deleteStateCookie())

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	req, err := githubtoken.BuildAccessTokenRequest(ctx, h.OAuth.Config, githubCode, qState)
	if err != nil {
		h.OAuth.Logger.Error("build github access token request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": githubtoken.ErrorCodeBuildAccessTokenRequest})

		return
	}

	resp, err := h.API.HTTPClient.Do(req)
	if err != nil {
		h.OAuth.Logger.Error("github access token request failed", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": githubtoken.ErrorCodeGitHubTokenRequest})

		return
	}

	githubAccessToken, err := githubtoken.ExtractAccessTokenFromResponse(resp)
	if err != nil {
		h.OAuth.Logger.Warn("github access token response parse failed", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": githubtoken.ErrorCodeGitHubTokenExchange})

		return
	}

	githubUserReq, err := githubuser.NewGitHubUserRequest(ctx, githubAccessToken)
	if err != nil {
		h.OAuth.Logger.Error("build github /user request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": githubuser.ErrorCodeGitHubUserRequestBuild})

		return
	}

	githubUserResp, err := h.API.HTTPClient.Do(githubUserReq)
	if err != nil {
		h.OAuth.Logger.Error("failed to call github /user", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": githubuser.ErrorCodeGitHubUserRequest})

		return
	}

	githubUser, err := githubuser.DecodeGitHubUserResponse(githubUserResp)
	if err != nil {
		h.OAuth.Logger.Warn("failed to decode github /user response", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": githubuser.ErrorCodeGitHubUserDecode})

		return
	}

	internalUserID, err := h.UserService.UpsertFromGitHub(
		ctx,
		githubUser.ID,
		githubUser.Login,
		githubUser.Email,
	)
	if err != nil {
		h.OAuth.Logger.Error("failed to upsert github user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": apierror.ErrorCodeUserUpsert})

		return
	}

	proxyCode, err := h.ProxyCodeService.Issue(
		ctx,
		internalUserID,
		h.ClientID,
	)
	if err != nil {
		h.OAuth.Logger.Error("failed to issue proxy code", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": apierror.ErrorCodeProxyCodeIssue})

		return
	}

	c.Redirect(
		http.StatusFound,
		callbackSuccessLocation(proxyCode, qState),
	)
}
