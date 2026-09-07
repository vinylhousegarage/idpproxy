package me

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	githubstore "github.com/vinylhousegarage/idpproxy/internal/oauth/github/store"
)

func TestService_Get(t *testing.T) {
	t.Parallel()

	const firebaseUID = "firebase-uid-123"
	const accessToken = "github-access-token"

	t.Run("gets_GitHub_user_and_touches_last_used", func(t *testing.T) {
		t.Parallel()

		tokenRepo := &fakeTokenRepository{
			record: &githubstore.GitHubTokenRecord{
				FirebaseUID: firebaseUID,
				AccessToken: accessToken,
			},
		}
		httpClient := &fakeHTTPClient{
			response: okJSONResponse(`{}`),
		}
		service := NewService(tokenRepo, httpClient)

		got, err := service.Get(context.Background(), firebaseUID)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		if got == nil {
			t.Fatal("Get() user = nil, want non-nil")
		}

		if tokenRepo.getFirebaseUID != firebaseUID {
			t.Fatalf(
				"GetByFirebaseUID() firebaseUID = %q, want %q",
				tokenRepo.getFirebaseUID,
				firebaseUID,
			)
		}

		if httpClient.request == nil {
			t.Fatal("GitHub API request was not sent")
		}

		if got := httpClient.request.Header.Get("Authorization"); got != "Bearer "+accessToken {
			t.Fatalf(
				"Authorization header = %q, want %q",
				got,
				"Bearer "+accessToken,
			)
		}

		if !tokenRepo.touchCalled {
			t.Fatal("TouchLastUsed() was not called")
		}

		if tokenRepo.touchFirebaseUID != firebaseUID {
			t.Fatalf(
				"TouchLastUsed() firebaseUID = %q, want %q",
				tokenRepo.touchFirebaseUID,
				firebaseUID,
			)
		}
	})

	t.Run("returns_error_when_token_is_not_found", func(t *testing.T) {
		t.Parallel()

		tokenRepo := &fakeTokenRepository{
			getErr: githubstore.ErrGitHubTokenNotFound,
		}
		httpClient := &fakeHTTPClient{}
		service := NewService(tokenRepo, httpClient)

		_, err := service.Get(context.Background(), firebaseUID)
		if !errors.Is(err, githubstore.ErrGitHubTokenNotFound) {
			t.Fatalf(
				"Get() error = %v, want %v",
				err,
				githubstore.ErrGitHubTokenNotFound,
			)
		}

		if httpClient.request != nil {
			t.Fatal("GitHub API request must not be sent")
		}

		if tokenRepo.touchCalled {
			t.Fatal("TouchLastUsed() must not be called")
		}
	})

	t.Run("returns_error_when_GitHub_API_call_fails", func(t *testing.T) {
		t.Parallel()

		githubAPIError := errors.New("GitHub API unavailable")

		tokenRepo := &fakeTokenRepository{
			record: &githubstore.GitHubTokenRecord{
				FirebaseUID: firebaseUID,
				AccessToken: accessToken,
			},
		}
		httpClient := &fakeHTTPClient{
			err: githubAPIError,
		}
		service := NewService(tokenRepo, httpClient)

		_, err := service.Get(context.Background(), firebaseUID)
		if !errors.Is(err, githubAPIError) {
			t.Fatalf("Get() error = %v, want %v", err, githubAPIError)
		}

		if tokenRepo.touchCalled {
			t.Fatal("TouchLastUsed() must not be called when GitHub API call fails")
		}
	})

	t.Run("returns_error_when_touch_last_used_fails", func(t *testing.T) {
		t.Parallel()

		touchError := errors.New("update Firestore last used time")

		tokenRepo := &fakeTokenRepository{
			record: &githubstore.GitHubTokenRecord{
				FirebaseUID: firebaseUID,
				AccessToken: accessToken,
			},
			touchErr: touchError,
		}
		httpClient := &fakeHTTPClient{
			response: okJSONResponse(`{}`),
		}
		service := NewService(tokenRepo, httpClient)

		_, err := service.Get(context.Background(), firebaseUID)
		if !errors.Is(err, touchError) {
			t.Fatalf("Get() error = %v, want %v", err, touchError)
		}

		if !tokenRepo.touchCalled {
			t.Fatal("TouchLastUsed() was not called")
		}
	})
}

type fakeTokenRepository struct {
	record           *githubstore.GitHubTokenRecord
	getErr           error
	touchErr         error
	getFirebaseUID   string
	touchFirebaseUID string
	touchCalled      bool
}

func (r *fakeTokenRepository) GetByFirebaseUID(
	_ context.Context,
	firebaseUID string,
) (*githubstore.GitHubTokenRecord, error) {
	r.getFirebaseUID = firebaseUID

	if r.getErr != nil {
		return nil, r.getErr
	}

	return r.record, nil
}

func (r *fakeTokenRepository) TouchLastUsed(
	_ context.Context,
	firebaseUID string,
) error {
	r.touchCalled = true
	r.touchFirebaseUID = firebaseUID

	return r.touchErr
}

type fakeHTTPClient struct {
	response *http.Response
	err      error
	request  *http.Request
}

func (c *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.request = req

	if c.err != nil {
		return nil, c.err
	}

	return c.response, nil
}

func okJSONResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
}
