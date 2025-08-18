package google_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/public"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

func newMockGoogleDepsWithFunc(logger *zap.Logger, fn func(ctx context.Context, idToken string) (*firebaseauth.Token, error)) *deps.GoogleDependencies {
	return &deps.GoogleDependencies{
		Logger: logger,
		Verifier: &testhelpers.MockVerifier{
			VerifyFunc: fn,
		},
	}
}

func TestMeRoute_Returns200AndResponse(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	googleDeps := newMockGoogleDepsWithFunc(logger, func(ctx context.Context, idToken string) (*firebaseauth.Token, error) {
		return &firebaseauth.Token{
			UID: "test-sub",
			Claims: map[string]interface{}{
				"iss": "https://issuer.example.com",
				"aud": "test-audience",
				"exp": float64(1234567890),
			},
		}, nil
	})

	githubDeps := testhelpers.NewMockGitHubDeps(logger)
	systemDeps := testhelpers.NewMockSystemDeps(logger)
	r := router.NewRouter(githubDeps, githubAPIDeps, googleDeps, systemDeps, http.FS(public.PublicFS))

	req, err := http.NewRequest(http.MethodGet, "/me", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer dummy.token.value")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)

	require.Equal(t, "test-sub", body["sub"])
	require.Equal(t, "https://issuer.example.com", body["iss"])
	require.Equal(t, "test-audience", body["aud"])
	require.EqualValues(t, 1234567890, body["exp"])
}
