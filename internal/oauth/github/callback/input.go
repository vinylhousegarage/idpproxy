package github

import (
	"time"
)

type SaveTokenInput struct {
	UserID      string
	AccessToken string
	Scopes      []string
	TokenType   string
	ExpiresAt   *time.Time
	Now         time.Time
}
