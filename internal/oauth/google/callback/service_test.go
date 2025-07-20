package callback

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestValidateCallbackRequest_ValidInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc123&state=xyz789", nil)

	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "xyz789",
	})

	code, err := ValidateCallbackRequest(req)

	assert.NoError(t, err)
	assert.Equal(t, "abc123", code)
}

func TestValidateCallbackRequest_MissingCode(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	_, err := ValidateCallbackRequest(req)

	assert.ErrorIs(t, err, ErrMissingCode)
}

func TestValidateCallbackRequest_MissingState(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	_, err := ValidateCallbackRequest(req)

	assert.ErrorIs(t, err, ErrMissingState)
}

func TestValidateCallbackRequest_MissingCookie(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc&state=xyz", nil)

	_, err := ValidateCallbackRequest(req)

	assert.ErrorIs(t, err, ErrMissingStateCookie)
}

func TestValidateCallbackRequest_StateMismatch(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc&state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "wrong"})

	_, err := ValidateCallbackRequest(req)

	assert.ErrorIs(t, err, ErrInvalidState)
}

func TestGetCallbackURL_Success(t *testing.T) {
	t.Parallel()

	want := "https://example.com/token"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"token_endpoint": "` + want + `"}`))
	}))
	defer ts.Close()

	endpoint, err := GetCallbackURL(ts.URL, http.DefaultClient, zap.NewNop())

	assert.NoError(t, err)
	assert.Equal(t, want, endpoint)
}

func TestGetCallbackURL_CreateRequestError(t *testing.T) {
	t.Parallel()

	badURL := "http://[::1]:namedport"

	_, err := GetCallbackURL(badURL, http.DefaultClient, zap.NewNop())

	assert.ErrorIs(t, err, ErrFailedToCreateRequest)
}

type errorClient struct{}

func (e *errorClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("network failure")
}

func TestGetCallbackURL_FetchError(t *testing.T) {
	t.Parallel()

	_, err := GetCallbackURL("http://localhost", &errorClient{}, zap.NewNop())

	assert.ErrorIs(t, err, ErrFailedToFetchMetadata)
}

func TestGetCallbackURL_UnexpectedStatusCode(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer ts.Close()

	_, err := GetCallbackURL(ts.URL, http.DefaultClient, zap.NewNop())

	assert.ErrorIs(t, err, ErrUnexpectedMetadataStatusCode)
}

func TestGetCallbackURL_DecodeError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"token_endpoint":`))
	}))
	defer ts.Close()

	_, err := GetCallbackURL(ts.URL, http.DefaultClient, zap.NewNop())

	assert.ErrorIs(t, err, ErrFailedToDecodeMetadata)
}

func TestGetCallbackURL_EmptyTokenEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"token_endpoint": ""}`))
	}))
	defer ts.Close()

	_, err := GetCallbackURL(ts.URL, http.DefaultClient, zap.NewNop())

	assert.ErrorIs(t, err, ErrMissingTokenEndpoint)
}
