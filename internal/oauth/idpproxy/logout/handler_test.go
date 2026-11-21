package logout

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

type fakeSessionInvalidator struct {
	called    bool
	sessionID string
	err       error
}

func (f *fakeSessionInvalidator) Invalidate(
	_ context.Context,
	sessionID string,
) (*authsession.Session, error) {
	f.called = true
	f.sessionID = sessionID

	if f.err != nil {
		return nil, f.err
	}

	return &authsession.Session{}, nil
}

func TestHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	t.Run("invalidates session and deletes session cookie", func(t *testing.T) {
		t.Parallel()

		invalidator := &fakeSessionInvalidator{}
		handler := NewHandler(invalidator)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  authsession.SessionCookieName,
			Value: "session-123",
			Path:  "/",
		})
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusNoContent, recorder.Code)
		require.True(t, invalidator.called)
		require.Equal(t, "session-123", invalidator.sessionID)

		cookie := findCookie(
			t,
			recorder.Result().Cookies(),
			authsession.SessionCookieName,
		)
		require.Equal(t, "/", cookie.Path)
		require.Equal(t, -1, cookie.MaxAge)
	})

	t.Run("returns no content when session cookie is missing", func(t *testing.T) {
		t.Parallel()

		invalidator := &fakeSessionInvalidator{}
		handler := NewHandler(invalidator)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusNoContent, recorder.Code)
		require.False(t, invalidator.called)
	})

	t.Run("returns internal server error when invalidation fails", func(t *testing.T) {
		t.Parallel()

		invalidator := &fakeSessionInvalidator{
			err: errors.New("firestore update failed"),
		}
		handler := NewHandler(invalidator)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  authsession.SessionCookieName,
			Value: "session-123",
			Path:  "/",
		})
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		require.Equal(
			t,
			http.StatusInternalServerError,
			recorder.Code,
		)
		require.True(t, invalidator.called)
		require.Equal(t, "session-123", invalidator.sessionID)
		require.Empty(t, recorder.Result().Cookies())
	})
}

func findCookie(
	t *testing.T,
	cookies []*http.Cookie,
	name string,
) *http.Cookie {
	t.Helper()

	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	t.Fatalf("cookie %q was not set", name)

	return nil
}
