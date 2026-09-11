package me

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/auth/session"
	"github.com/vinylhousegarage/idpproxy/internal/oauth/github/response"
)

type fakeGitHubMeService struct {
	user *response.GitHubUserAPIResponse
	err  error

	called      bool
	firebaseUID string
}

func (s *fakeGitHubMeService) Get(
	_ context.Context,
	firebaseUID string,
) (*response.GitHubUserAPIResponse, error) {
	s.called = true
	s.firebaseUID = firebaseUID

	if s.err != nil {
		return nil, s.err
	}

	return s.user, nil
}

func TestGitHubMeHandler_Serve(t *testing.T) {
	t.Parallel()

	t.Run("returns_GitHub_user_for_authenticated_user", func(t *testing.T) {
		t.Parallel()

		service := &fakeGitHubMeService{
			user: &response.GitHubUserAPIResponse{},
		}
		handler := NewGitHubMeHandler(service)

		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/me", nil)
		ctx.Set(session.ContextUserIDKey, "firebase-uid-123")

		handler.Serve(ctx)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		if !service.called {
			t.Fatal("GitHubMeService.Get was not called")
		}

		if got, want := service.firebaseUID, "firebase-uid-123"; got != want {
			t.Errorf("Get firebaseUID = %q, want %q", got, want)
		}
	})

	t.Run("returns_unauthorized_when_user_ID_is_missing", func(t *testing.T) {
		t.Parallel()

		service := &fakeGitHubMeService{
			user: &response.GitHubUserAPIResponse{},
		}
		handler := NewGitHubMeHandler(service)

		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/me", nil)

		handler.Serve(ctx)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}

		if service.called {
			t.Fatal("GitHubMeService.Get must not be called without a user ID")
		}
	})

	t.Run("returns_internal_server_error_when_service_fails", func(t *testing.T) {
		t.Parallel()

		service := &fakeGitHubMeService{
			err: errors.New("GitHub API unavailable"),
		}
		handler := NewGitHubMeHandler(service)

		rr := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rr)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/me", nil)
		ctx.Set(session.ContextUserIDKey, "firebase-uid-123")

		handler.Serve(ctx)

		if rr.Code != http.StatusInternalServerError {
			t.Fatalf(
				"status = %d, want %d",
				rr.Code,
				http.StatusInternalServerError,
			)
		}

		if !service.called {
			t.Fatal("GitHubMeService.Get was not called")
		}
	})
}
