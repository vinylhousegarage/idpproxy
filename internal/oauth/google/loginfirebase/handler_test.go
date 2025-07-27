package loginfirebase

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vinylhousegarage/idpproxy/test/testhelpers"
)

func TestLoginFirebaseHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockVerifier := &testhelpers.MockVerifier{
			VerifyFunc: func(ctx context.Context, idToken string) (*firebaseauth.Token, error) {
				return &firebaseauth.Token{UID: "test-uid"}, nil
			},
		}

		body := []byte(`{"id_token":"dummy.token.value"}`)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := &LoginFirebaseHandler{
			Logger:   zap.NewNop(),
			Verifier: mockVerifier,
		}
		err := handler.LoginFirebaseHandler(rr, req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)

		cookies := rr.Result().Cookies()
		require.Len(t, cookies, 1)
		require.Equal(t, "id_token", cookies[0].Name)
		require.Equal(t, "dummy.token.value", cookies[0].Value)
	})

	t.Run("invalid request", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`invalid`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := &LoginFirebaseHandler{}
		err := handler.LoginFirebaseHandler(rr, req)

		require.ErrorIs(t, err, ErrInvalidRequest)
	})

	t.Run("invalid token", func(t *testing.T) {
		t.Parallel()

		mockVerifier := &testhelpers.MockVerifier{
			VerifyFunc: func(ctx context.Context, idToken string) (*firebaseauth.Token, error) {
				return nil, errors.New("invalid token")
			},
		}

		body := []byte(`{"id_token":"invalid.token"}`)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler := &LoginFirebaseHandler{
			Logger:   zap.NewNop(),
			Verifier: mockVerifier,
		}
		err := handler.LoginFirebaseHandler(rr, req)

		require.ErrorIs(t, err, ErrInvalidIDToken)
	})
}
