package loginfirebase_test

import (
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

	di := testhelpers.NewMockDeps(logger)

	r := router.NewRouter(di, http.FS(public.PublicFS))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/loginfirebase", nil)
	require.NoError(t, err)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	cookies := w.Result().Cookies()
	require.NotEmpty(t, cookies)

	var idTokenFound bool
	for _, c := range cookie {
		if c.Name == "id_token" {
			idTokenFound = true
			require.NotEmpty(t, c.Value)
			break
		}
	}
	require.True(t, idTokenFound, "id_token cookie should be set")
}
