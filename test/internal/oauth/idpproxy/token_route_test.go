package idpproxy_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	authidtoken "github.com/vinylhousegarage/idpproxy/internal/auth/idtoken"
	authcodeservice "github.com/vinylhousegarage/idpproxy/internal/authcode/service"
	authcodestore "github.com/vinylhousegarage/idpproxy/internal/authcode/store"
	githubcallback "github.com/vinylhousegarage/idpproxy/internal/oauth/github/callback"
	githubme "github.com/vinylhousegarage/idpproxy/internal/oauth/github/me"
	idpproxytoken "github.com/vinylhousegarage/idpproxy/internal/oauth/idpproxy/token"
	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/public"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

type fakeIDTokenIssuer struct {
	called bool
	input  *authidtoken.IDTokenInput
	token  string
	err    error
}

func (f *fakeIDTokenIssuer) Issue(
	_ context.Context,
	in *authidtoken.IDTokenInput,
) (string, string, error) {
	f.called = true
	f.input = in

	if f.err != nil {
		return "", "", f.err
	}

	return f.token, "test-key-id", nil
}

func TestTokenRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("exchanges a proxy code for an ID token", func(t *testing.T) {
		r, proxyCodeService, issuer := newTokenRouteRouter(t)

		proxyCode, err := proxyCodeService.Issue(
			t.Context(),
			"user-123",
			"client-123",
		)
		require.NoError(t, err)

		body, err := json.Marshal(idpproxytoken.TokenRequest{
			GrantType: "authorization_code",
			Code:      proxyCode,
			ClientID:  "client-123",
		})
		require.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodPost,
			"/oauth/token",
			bytes.NewReader(body),
		)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response idpproxytoken.TokenResponse
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&response))
		require.Equal(t, "signed-id-token", response.IDToken)

		require.True(t, issuer.called)
		require.NotNil(t, issuer.input)
		require.Equal(t, "user-123", issuer.input.UserID)
		require.Equal(t, "client-123", issuer.input.ClientID)
		require.Equal(t, time.Hour, issuer.input.TTL)
	})

	t.Run("rejects a previously consumed proxy code", func(t *testing.T) {
		r, proxyCodeService, issuer := newTokenRouteRouter(t)

		proxyCode, err := proxyCodeService.Issue(
			t.Context(),
			"user-123",
			"client-123",
		)
		require.NoError(t, err)

		body, err := json.Marshal(idpproxytoken.TokenRequest{
			GrantType: "authorization_code",
			Code:      proxyCode,
			ClientID:  "client-123",
		})
		require.NoError(t, err)

		firstRequest := httptest.NewRequest(
			http.MethodPost,
			"/oauth/token",
			bytes.NewReader(body),
		)
		firstRequest.Header.Set("Content-Type", "application/json")
		firstRecorder := httptest.NewRecorder()

		r.ServeHTTP(firstRecorder, firstRequest)

		require.Equal(t, http.StatusOK, firstRecorder.Code)
		require.True(t, issuer.called)

		secondRequest := httptest.NewRequest(
			http.MethodPost,
			"/oauth/token",
			bytes.NewReader(body),
		)
		secondRequest.Header.Set("Content-Type", "application/json")
		secondRecorder := httptest.NewRecorder()

		r.ServeHTTP(secondRecorder, secondRequest)

		require.Equal(t, http.StatusBadRequest, secondRecorder.Code)

		var response map[string]string
		require.NoError(
			t,
			json.NewDecoder(secondRecorder.Body).Decode(&response),
		)
		require.Equal(t, "invalid_grant", response["error"])
	})
}

func newTokenRouteRouter(
	t *testing.T,
) (
	*gin.Engine,
	*authcodeservice.Service,
	*fakeIDTokenIssuer,
) {
	t.Helper()

	logger := zap.NewNop()

	proxyCodeService := authcodeservice.NewService(
		authcodestore.NewMemoryStore(),
	)

	issuer := &fakeIDTokenIssuer{
		token: "signed-id-token",
	}

	tokenHandler := idpproxytoken.NewHandler(
		&idpproxytoken.Service{
			Store: idpproxytoken.NewAuthCodeServiceAdapter(
				proxyCodeService,
			),
			Clock:         idpproxytoken.SystemClock{},
			IDTokenIssuer: issuer,
			IDTokenTTL:    time.Hour,
		},
		logger,
	)

	d := router.NewRouterDeps(
		public.PublicFS,
		testhelpers.NewMockGitHubAPIDeps(logger),
		testhelpers.NewMockGitHubOAuthDeps(logger),
		&githubcallback.GitHubCallbackHandler{},
		&githubme.GitHubMeHandler{},
		testhelpers.NewMockSessionValidator(),
		tokenHandler,
		testhelpers.NewMockGoogleDeps(logger),
		logger,
		testhelpers.NewMockSystemDeps(logger),
	)

	return router.NewRouter(d), proxyCodeService, issuer
}
