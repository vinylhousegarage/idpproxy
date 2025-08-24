package idp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vinylhousegarage/idpproxy/internal/httpclient"
)

const (
	testAPIKey   = "test-api-key"
	testReqURI   = "https://idpproxy.com/auth_cb"
	testAccToken = "ACCESS_TOKEN_X"
)

type rewriteRoundTripper struct {
	target *url.URL
	next   http.RoundTripper
}

func (r *rewriteRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host == "identitytoolkit.googleapis.com" {
		req.URL.Scheme = r.target.Scheme
		req.URL.Host = r.target.Host
	}
	return r.next.RoundTrip(req)
}

func newRewritingHTTPClient(t *testing.T, target *url.URL) httpclient.HTTPClient {
	t.Helper()
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &rewriteRoundTripper{
			target: target,
			next:   http.DefaultTransport,
		},
	}
}

func TestSignInWithIdpByAccessToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErr    bool
		errSubstrs []string
		wantResp   *SignInWithIdpResp
	}{
		{
			name: "Success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/v1/accounts:signInWithIdp", r.URL.Path)
				require.Contains(t, r.URL.RawQuery, "key="+testAPIKey)
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))

				var payload signInPayload
				defer r.Body.Close()
				require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))

				require.Equal(t, testReqURI, payload.RequestURI)
				require.True(t, payload.ReturnSecureToken)

				values, err := url.ParseQuery(payload.PostBody)
				require.NoError(t, err)
				require.Equal(t, testAccToken, values.Get("access_token"))
				require.Equal(t, "github.com", values.Get("providerId"))

				resp := map[string]any{
					"providerId":   "github.com",
					"localId":      "firebase_local_123",
					"idToken":      "ID_TOKEN_ABC",
					"refreshToken": "REFRESH_TOKEN_DEF",
					"expiresIn":    "3600",
					"isNewUser":    false,
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(resp)
			},
			wantResp: &SignInWithIdpResp{
				ProviderID:   "github.com",
				LocalID:      "firebase_local_123",
				IDToken:      "ID_TOKEN_ABC",
				RefreshToken: "REFRESH_TOKEN_DEF",
				ExpiresIn:    "3600",
				IsNewUser:    false,
			},
		},
		{
			name: "ErrorWithMessage",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]any{
						"code":    400,
						"message": "INVALID_IDP_RESPONSE",
					},
				})
			},
			wantErr:    true,
			errSubstrs: []string{"INVALID_IDP_RESPONSE", "status=400"},
		},
		{
			name: "ErrorWithoutMessage",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{ "not_firebase_error": true }`))
			},
			wantErr:    true,
			errSubstrs: []string{"unexpected status 500"},
		},
		{
			name: "Error_200ButInvalidJSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{invalid-json"))
			},
			wantErr:    true,
			errSubstrs: []string{"decode:"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(tt.handler)
			t.Cleanup(ts.Close)

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)
			hc := newRewritingHTTPClient(t, u)

			out, err := SignInWithIdpByAccessToken(
				context.Background(),
				hc,
				testAPIKey,
				testReqURI,
				testAccToken,
			)

			if tt.wantErr {
				require.Error(t, err)
				for _, s := range tt.errSubstrs {
					require.Contains(t, err.Error(), s)
				}
				require.Nil(t, out)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, out)
			}
		})
	}

	t.Run("Validation_apiKeyEmpty", func(t *testing.T) {
		t.Parallel()
		_, err := SignInWithIdpByAccessToken(context.Background(), nil, "", testReqURI, testAccToken)
		require.Error(t, err)
		require.Contains(t, err.Error(), "apiKey is empty")
	})

	t.Run("Validation_requestURIEmpty", func(t *testing.T) {
		t.Parallel()
		_, err := SignInWithIdpByAccessToken(context.Background(), nil, testAPIKey, "", testAccToken)
		require.Error(t, err)
		require.Contains(t, err.Error(), "requestURI is empty")
	})

	t.Run("Validation_accessTokenEmpty", func(t *testing.T) {
		t.Parallel()
		_, err := SignInWithIdpByAccessToken(context.Background(), nil, testAPIKey, testReqURI, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "accessToken is empty")
	})
}
