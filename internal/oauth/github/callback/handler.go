package callback

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/apierror"
	githubstore "github.com/vinylhousegarage/idpproxy/internal/oauth/github/store"
	githubtoken "github.com/vinylhousegarage/idpproxy/internal/oauth/github/token"
	githubuser "github.com/vinylhousegarage/idpproxy/internal/oauth/github/user"
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
		apiErr := apierror.InvalidState(apierror.ErrInvalidState).
			AddInternal(apierror.ErrorCodeInvalidQueryState, "query_state", qState).
			AddInternal(
				apierror.ErrorCodeInvalidCookieState,
				"cookie_state",
				safeCookieVal(cookie),
			)

		_ = c.Error(apiErr)

		http.SetCookie(c.Writer, deleteStateCookie())

		return
	}

	http.SetCookie(c.Writer, deleteStateCookie())

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	req, err := githubtoken.BuildAccessTokenRequest(
		ctx,
		h.OAuth.Config,
		githubCode,
		qState,
	)
	if err != nil {
		_ = c.Error(
			apierror.GitHubAccessTokenRequestError(
				apierror.ErrGitHubAccessTokenRequest,
			),
		)

		return
	}

	resp, err := h.API.HTTPClient.Do(req)
	if err != nil {
		_ = c.Error(
			apierror.GitHubTokenRequestError(
				apierror.ErrGitHubTokenRequest,
			),
		)

		return
	}
	defer resp.Body.Close()

	githubAccessToken, err := githubtoken.ExtractAccessTokenFromResponse(resp)
	if err != nil {
		_ = c.Error(
			apierror.GitHubTokenExchangeError(
				apierror.ErrGitHubTokenExchange,
			),
		)

		return
	}

	githubUserReq, err := githubuser.NewGitHubUserRequest(
		ctx,
		githubAccessToken,
	)
	if err != nil {
		_ = c.Error(
			apierror.GitHubUserRequestBuildError(
				apierror.ErrGitHubUserRequestBuild,
			),
		)

		return
	}

	githubUserResp, err := h.API.HTTPClient.Do(githubUserReq)
	if err != nil {
		_ = c.Error(
			apierror.GitHubUserRequestError(
				apierror.ErrGitHubUserRequest,
			),
		)

		return
	}
	defer githubUserResp.Body.Close()

	githubUser, err := githubuser.DecodeGitHubUserResponse(githubUserResp)
	if err != nil {
		_ = c.Error(
			apierror.GitHubUserDecodeError(
				apierror.ErrGitHubUserDecode,
			),
		)

		return
	}

	internalUserID, err := h.UserService.UpsertFromGitHub(
		ctx,
		githubUser.ID,
		githubUser.Login,
		githubUser.Email,
	)
	if err != nil {
		_ = c.Error(
			apierror.UserUpsertError(apierror.ErrUserUpsert),
		)

		return
	}

	err = h.TokenRepo.Upsert(ctx, &githubstore.GitHubTokenRecord{
		GitHubID:    strconv.FormatInt(githubUser.ID, 10),
		Provider:    "github",
		FirebaseUID: internalUserID,
		Login:       githubUser.Login,
		AccessToken: githubAccessToken,
		LastUsedAt:  time.Now(),
	})
	if err != nil {
		_ = c.Error(
			apierror.GitHubTokenUpsertError(
				apierror.ErrGitHubTokenUpsert,
			),
		)

		return
	}

	session, err := h.SessionService.Start(ctx, internalUserID)
	if err != nil || session == nil || session.SessionID == "" {
		_ = c.Error(
			apierror.SessionStartError(apierror.ErrSessionStart),
		)

		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     authsession.SessionCookieName,
		Value:    session.SessionID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	proxyCode, err := h.ProxyCodeService.Issue(
		ctx,
		internalUserID,
		h.ClientID,
	)
	if err != nil {
		_ = c.Error(
			apierror.ProxyCodeIssueError(
				apierror.ErrProxyCodeIssue,
			),
		)

		return
	}

	c.Redirect(
		http.StatusFound,
		callbackSuccessLocation(proxyCode, qState),
	)
}
