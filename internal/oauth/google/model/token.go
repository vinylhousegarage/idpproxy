package model

import "time"

type GoogleToken struct {
	AccessToken  string    `firestore:"access_token"`
	RefreshToken string    `firestore:"refresh_token"`
	CreatedAt    time.Time `firestore:"created_at"`
	ExpiresIn    int       `firestore:"expires_in"`
	ExpiresAt    time.Time `firestore:"expires_at"`
	LastUsedAt   time.Time `firestore:"last_used_at"`
}
