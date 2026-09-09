package session

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeSessionValidator struct {
	session *Session
	err     error

	called    bool
	sessionID string
}

func (f *fakeSessionValidator) Validate(
	_ context.Context,
	sessionID string,
) (*Session, error) {
	f.called = true
	f.sessionID = sessionID

	if f.err != nil {
		return nil, f.err
	}

	return f.session, nil
}

func TestRequireSession(t *testing.T) {
	t.Parallel()

	t.Run("sets_user_ID_and_calls_next_handler_for_valid_session", func(t *testing.T) {
		t.Parallel()

		validator := &fakeSessionValidator{
			session: &Session{
				SessionID: "session-123",
				UserID:    "user-123",
				Status:    "active",
			},
		}

		router := gin.New()
		router.Use(RequireSession(validator))
		router.GET("/me", func(c *gin.Context) {
			userID, ok := c.Get(contextUserIDKey)
			if !ok {
				t.Error("user ID was not set in Gin context")
				c.Status(http.StatusInternalServerError)

				return
			}

			if got, want := userID, "user-123"; got != want {
				t.Errorf("context user ID = %v, want %q", got, want)
			}

			c.Status(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.AddCookie(&http.Cookie{
			Name:  sessionCookieName,
			Value: "session-123",
			Path:  "/",
		})
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
		}

		if !validator.called {
			t.Fatal("Validate was not called")
		}

		if got, want := validator.sessionID, "session-123"; got != want {
			t.Errorf("Validate sessionID = %q, want %q", got, want)
		}
	})

	t.Run("returns_unauthorized_when_session_cookie_is_missing", func(t *testing.T) {
		t.Parallel()

		validator := &fakeSessionValidator{
			session: &Session{UserID: "user-123"},
		}

		nextCalled := false

		router := gin.New()
		router.Use(RequireSession(validator))
		router.GET("/me", func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}

		if validator.called {
			t.Fatal("Validate must not be called without a session cookie")
		}

		if nextCalled {
			t.Fatal("next handler must not be called without a session cookie")
		}
	})

	t.Run("returns_unauthorized_when_session_validation_fails", func(t *testing.T) {
		t.Parallel()

		validator := &fakeSessionValidator{
			err: errors.New("expired session"),
		}

		nextCalled := false

		router := gin.New()
		router.Use(RequireSession(validator))
		router.GET("/me", func(c *gin.Context) {
			nextCalled = true
			c.Status(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.AddCookie(&http.Cookie{
			Name:  sessionCookieName,
			Value: "expired-session-123",
			Path:  "/",
		})
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}

		if !validator.called {
			t.Fatal("Validate was not called")
		}

		if got, want := validator.sessionID, "expired-session-123"; got != want {
			t.Errorf("Validate sessionID = %q, want %q", got, want)
		}

		if nextCalled {
			t.Fatal("next handler must not be called when validation fails")
		}
	})
}
