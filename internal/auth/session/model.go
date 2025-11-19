package session

import "time"

type Session struct {
	SessionID string     `firestore:"session_id"`
	UserID    string     `firestore:"user_id"`
	Status    string     `firestore:"status"`
	ExpiresAt time.Time  `firestore:"expires_at"`
	CreatedAt time.Time  `firestore:"created_at"`
	UpdatedAt *time.Time `firestore:"updated_at,omitempty"`
	LastUsed  *time.Time `firestore:"last_used,omitempty"`
}
