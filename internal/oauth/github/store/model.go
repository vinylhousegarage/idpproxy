package store

import "time"

type GitHubTokenRecord struct {
	GitHubID    string    `firestore:"github_id"`
	Provider    string    `firestore:"provider"`
	FirebaseUID string    `firestore:"firebase_uid"`
	Login       string    `firestore:"login"`
	Scopes      []string  `firestore:"scopes"`
	TokenType   string    `firestore:"token_type"`
	AccessToken string    `firestore:"access_token"`
	ExpiresAt   time.Time `firestore:"expires_at"`
	LastUsedAt  time.Time `firestore:"last_used_at"`
	CreatedAt   time.Time `firestore:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at"`
	DeleteAt    time.Time `firestore:"delete_at"`
}
