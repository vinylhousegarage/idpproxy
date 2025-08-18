package google_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/router"
	"github.com/vinylhousegarage/idpproxy/public"
	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

func TestLoginfirebaseRoute_Returns200AndIDToken(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	githubDeps := testhelpers.NewMockGitHubDeps(logger)
	githubAPIDeps := testhelpers.NewMockGitHubAPIDeps(logger)
	googleDeps := testhelpers.NewMockGoogleDeps(logger)
	systemDeps := testhelpers.NewMockSystemDeps(logger)

	r := router.NewRouter(githubDeps, githubAPIDeps, googleDeps, systemDeps, http.FS(public.PublicFS))

	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"id_token":"dummy.token.value"}`)
	req, err := http.NewRequest(http.MethodPost, "/google/login/firebase", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	cookies := w.Result().Cookies()
	require.NotEmpty(t, cookies)

	var idTokenFound bool
	for _, c := range cookies {
		if c.Name == "id_token" {
			idTokenFound = true
			require.NotEmpty(t, c.Value)
			break
		}
	}
	require.True(t, idTokenFound, "id_token cookie should be set")
}
