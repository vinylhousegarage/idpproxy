package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vinylhousegarage/idpproxy/internal/deps"
)

func NewGitHubUserHandler(apiDeps *deps.GitHubAPIDependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := ExtractAuthHeaderToken(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		req, err := NewGitHubUserRequest(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build request"})
			return
		}

		resp, err := apiDeps.HTTPClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to call GitHub"})
			return
		}

		githubUser, err := DecodeGitHubUserResponse(resp)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, githubUser)
	}
}
