package token

import (
	"errors"
	"net/http"
	"testing"
)

func TestGitHubAccessTokenRequestError(t *testing.T) {
	t.Parallel()

	err := errors.New("something went wrong")
	internalInfo := "failed to parse token body"

	got := GitHubAccessTokenRequestError(err, internalInfo)

	if got.Code != ErrorCodeGitHubAccessTokenRequest {
		t.Errorf("expected code %s, got %s", ErrorCodeGitHubAccessTokenRequest, got.Code)
	}
	if got.HTTPStatus != http.StatusBadGateway {
		t.Errorf("expected status 502, got %d", got.HTTPStatus)
	}
	if !errors.Is(got.Err, err) {
		t.Errorf("expected original error to be wrapped")
	}
	if got.Internal != internalInfo {
		t.Errorf("expected internal info %s, got %s", internalInfo, got.Internal)
	}
}
