package token

import (
	"errors"
	"net/http"
	"testing"
)

func TestGitHubAccessTokenResponse(t *testing.T) {
	t.Parallel()

	err := errors.New("something went wrong")
	internalInfo := "failed to parse token body"

	got := GitHubAccessTokenResponse(err, internalInfo)

	if got.Code != ErrorCodeGitHubAccessTokenResponse {
		t.Errorf("expected code %s, got %s", ErrorCodeGitHubAccessTokenResponse, got.Code)
	}
	if got.HTTPStatus != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", got.HTTPStatus)
	}
	if !errors.Is(got.Err, err) {
		t.Errorf("expected original error to be wrapped")
	}
	if got.Internal != internalInfo {
		t.Errorf("expected internal info %s, got %s", internalInfo, got.Internal)
	}
}
