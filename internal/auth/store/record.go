package store

import "time"

type RefreshTokenRecord struct {
	UserID     string    `firestore:"user_id"`
	CreatedAt  time.Time `firestore:"created_at"`
	LastUsedAt time.Time `firestore:"last_used_at"`
	ExpiresAt  time.Time `firestore:"expires_at"`
	DeleteAt   time.Time `firestore:"delete_at"`
}

type AccessGenerationRecord struct {
	UserID    string    `firestore:"user_id"`
	Gen       int       `firestore:"gen"`
	UpdatedAt time.Time `firestore:"updated_at"`
}
