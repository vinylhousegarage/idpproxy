package github_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/config"
	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/user"
)

func TestGitHubUserRoute_Success(t *testing.T) {
	t.Parallel()

	mockGitHub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/user" {
					http.NotFound(w, r)
					return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":12345,"login":"octocat","name":"Mona","email":"mona@github.com"}`))
	}))
	t.Cleanup(mockGitHub.Close)

	apiCfg := &config.GitHubAPIConfig{
		APIVersion: "2022-11-28",
		BaseURL:    mockGitHub.URL,
		UserAgent:  "idpproxy-test",
	}
	deps := deps.NewGitHubAPIDeps(apiCfg, &http.Client{Timeout: time.Second}, zap.NewNop())

	r := gin.New()
	r.GET("/user", user.NewGitHubUserHandler(deps))

	sut := httptest.NewServer(r)
	t.Cleanup(sut.Close)

	req, _ := http.NewRequest(http.MethodGet, sut.URL+"/user", nil)
	req.Header.Set("Authorization", "Bearer dummy_token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	require.Contains(t, string(body), `"login":"octocat"`)
}
