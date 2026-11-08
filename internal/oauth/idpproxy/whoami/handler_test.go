package whoami

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	authsession "github.com/vinylhousegarage/idpproxy/internal/auth/session"
)

func TestWhoAmIHandler_Serve(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns_authenticated_user_id", func(t *testing.T) {
		handler := &WhoAmIHandler{}

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Set(authsession.ContextUserIDKey, "firebase-user-123")

		handler.Serve(ctx)

		require.Equal(t, http.StatusOK, recorder.Code)

		var response struct {
			UserID string `json:"user_id"`
		}
		require.NoError(
			t,
			json.NewDecoder(recorder.Body).Decode(&response),
		)
		require.Equal(t, "firebase-user-123", response.UserID)
	})
}
