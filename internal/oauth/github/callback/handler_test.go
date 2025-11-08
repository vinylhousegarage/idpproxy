package callback

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGitHubCallbackHandler_Serve(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tokenJSON := loadTestDataJSON(t, "testdata/token_success.json")
	userJSON := loadTestDataJSON(t, "testdata/user_success.json")
}
