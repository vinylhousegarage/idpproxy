package callback

import (
	"context"
	"fmt"
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
		_ = c.Error(apierror.MissingGitHubCode(apierror.ErrMissingGitHubCode))

		return
	}

	if qState == "" {
		_ = c.Error(apierror.MissingState(apierror.ErrMissingState))

		return
	}

	cookie, err := c.Request.Cookie(stateCookieName)
	if err != nil || cookie == nil || cookie.Value == "" || cookie.Value != qState {
		query_state, query_err := apierror.FormatDetail("query_state", qState)
		if query_err != nil {
			h.OAuth.Logger.Error("failed to format query_state detail", zap.Error(query_err))
		}

		cookie_state, cookie_err := apierror.FormatDetail("cookie_state", safeCookieVal(cookie))
		if cookie_err != nil {
			h.OAuth.Logger.Error("failed to format cookie_state detail", zap.Error(cookie_err))
		}

		_ = c.Error(apierror.InvalidState(
			apierror.ErrInvalidState,
			apierror.APIInternal{
				Code: apierror.ErrorCodeInvalidQueryState,
				Err:  fmt.Errorf("query state is invalid: %s", query_state),
			},
			apierror.APIInternal{
				Code: apierror.ErrorCodeInvalidCookieState,
				Err:  fmt.Errorf("cookie state is invalid: %s", cookie_state),
			},
		))

		http.SetCookie(c.Writer, deleteStateCookie())

		return
	}

	http.SetCookie(c.Writer, deleteStateCookie())

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	req, err := githubtoken.BuildAccessTokenRequest(ctx, h.OAuth.Config, githubCode, qState)
	if err != nil {
		h.OAuth.Logger.Error("build github access token request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": apierror.ErrorCodeBuildAccessTokenRequest})

		return
	}

	resp, err := h.API.HTTPClient.Do(req)
	if err != nil {
		h.OAuth.Logger.Error("github access token request failed", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": apierror.ErrorCodeGitHubTokenRequest})

		return
	}

	defer resp.Body.Close()

	githubAccessToken, err := githubtoken.ExtractAccessTokenFromResponse(resp)
	if err != nil {
		h.OAuth.Logger.Warn("github access token response parse failed", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": apierror.ErrorCodeGitHubTokenExchange})

		return
	}

	githubUserReq, err := githubuser.NewGitHubUserRequest(ctx, githubAccessToken)
	if err != nil {
		h.OAuth.Logger.Error("build github /user request failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": apierror.ErrorCodeGitHubUserRequestBuild})

		return
	}

	githubUserResp, err := h.API.HTTPClient.Do(githubUserReq)
	if err != nil {
		h.OAuth.Logger.Error("failed to call github /user", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": apierror.ErrorCodeGitHubUserRequest})

		return
	}

	defer githubUserResp.Body.Close()

	githubUser, err := githubuser.DecodeGitHubUserResponse(githubUserResp)
	if err != nil {
		h.OAuth.Logger.Warn("failed to decode github /user response", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": apierror.ErrorCodeGitHubUserDecode})

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
