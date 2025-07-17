package login

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestGetGoogleLoginURL(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"authorization_endpoint": "https://example.com/oauth2/authorize"}`))
		}))
		defer ts.Close()

		endpoint, err := GetGoogleLoginURL(ts.URL, http.DefaultClient, logger)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/oauth2/authorize", endpoint)
	})

	t.Run("RequestCreationError", func(t *testing.T) {
		t.Parallel()

		_, err := GetGoogleLoginURL(":", http.DefaultClient, logger)
		assert.ErrorIs(t, err, ErrFailedToCreateRequest)
	})

	t.Run("HTTPClientError", func(t *testing.T) {
		t.Parallel()

		_, err := GetGoogleLoginURL("http://invalid.host.local", http.DefaultClient, logger)
		assert.ErrorIs(t, err, ErrFailedToFetchMetadata)
	})

	t.Run("UnexpectedStatusCode", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer ts.Close()

		_, err := GetGoogleLoginURL(ts.URL, http.DefaultClient, logger)
		assert.ErrorIs(t, err, ErrUnexpectedMetadataStatusCode)
	})

	t.Run("JSONDecodeError", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer ts.Close()

		_, err := GetGoogleLoginURL(ts.URL, http.DefaultClient, logger)
		assert.ErrorIs(t, err, ErrFailedToDecodeMetadata)
	})

	t.Run("MissingAuthorizationEndpoint", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"authorization_endpoint": ""}`))
		}))
		defer ts.Close()

		_, err := GetGoogleLoginURL(ts.URL, http.DefaultClient, logger)
		assert.ErrorIs(t, err, ErrMissingAuthorizationEndpoint)
	})
}
