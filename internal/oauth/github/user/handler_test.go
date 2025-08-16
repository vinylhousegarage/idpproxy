package user

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/response"
)

type mockHTTPClient struct {
	do func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) { return m.do(req) }

var _ httpclient.HTTPClient = (*mockHTTPClient)(nil)

func newRouterForTest(h gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/github/user", h)
	return r
}

func newJSONResponse(status int, v any) *http.Response {
	var body []byte
	if v != nil {
		body, _ = json.Marshal(v)
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     make(http.Header),
	}
}

func newTextResponse(status int, s string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewReader([]byte(s))),
		Header:     make(http.Header),
	}
}

func TestGitHubUserHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMock  func(t *testing.T) httpclient.HTTPClient
		setupReq   func() *http.Request
		wantStatus int
		assertBody func(t *testing.T, body string)
	}{
		{
			name: "Success",
			setupMock: func(t *testing.T) httpclient.HTTPClient {
				return &mockHTTPClient{
					do: func(req *http.Request) (*http.Response, error) {
						require.Equal(t, "Bearer ghp_success", req.Header.Get("Authorization"))
						return newJSONResponse(http.StatusOK, response.GitHubUserAPIResponse{
							ID:    12345,
							Login: "octocat",
							Email: "octo@example.com",
							Name:  "The Octo",
						}), nil
					},
				}
			},
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/github/user", nil)
				req.Header.Set("Authorization", "Bearer ghp_success")
				return req
			},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, body string) {
				var got response.GitHubUserAPIResponse
				require.NoError(t, json.Unmarshal([]byte(body), &got))
				require.Equal(t, int64(12345), got.ID)
				require.Equal(t, "octocat", got.Login)
				require.Equal(t, "octo@example.com", got.Email)
				require.Equal(t, "The Octo", got.Name)
			},
		},
		{
			name: "MissingAuthorizationHeader",
			setupMock: func(t *testing.T) httpclient.HTTPClient {
				return &mockHTTPClient{
					do: func(req *http.Request) (*http.Response, error) {
						t.Fatalf("HTTPClient.Do should not be called when header is missing")
						return nil, nil
					},
				}
			},
			setupReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/github/user", nil) // Authorization なし
			},
			wantStatus: http.StatusUnauthorized,
			assertBody: func(t *testing.T, body string) {
				require.Contains(t, body, ErrMissingAuthorizationHeader.Error())
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupMock: func(t *testing.T) httpclient.HTTPClient {
				return &mockHTTPClient{
					do: func(req *http.Request) (*http.Response, error) {
						t.Fatalf("HTTPClient.Do should not be called on invalid header format")
						return nil, nil
					},
				}
			},
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/github/user", nil)
				req.Header.Set("Authorization", "Token ghp_wrong")
				return req
			},
			wantStatus: http.StatusUnauthorized,
			assertBody: func(t *testing.T, body string) {
				require.Contains(t, body, ErrInvalidAuthorizationHeaderFormat.Error())
			},
		},
		{
			name: "UpstreamNon2xx",
			setupMock: func(t *testing.T) httpclient.HTTPClient {
				return &mockHTTPClient{
					do: func(req *http.Request) (*http.Response, error) {
						return newTextResponse(http.StatusForbidden, "forbidden-xxxxx"), nil
					},
				}
			},
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/github/user", nil)
				req.Header.Set("Authorization", "Bearer ghp_forbidden")
				return req
			},
			wantStatus: http.StatusBadGateway,
			assertBody: func(t *testing.T, body string) {
				require.Contains(t, body, "non-2xx status")
				require.Contains(t, body, "403")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			apiDeps := &deps.GitHubAPIDependencies{
				HTTPClient: tt.setupMock(t),
				Logger:     zap.NewNop(),
			}
			h := NewGitHubUserHandler(apiDeps)
			r := newRouterForTest(h)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, tt.setupReq())

			require.Equal(t, tt.wantStatus, w.Code)
			if tt.assertBody != nil {
				tt.assertBody(t, w.Body.String())
			}
		})
	}
}
