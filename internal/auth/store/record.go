package store

import "time"

type RefreshTokenRecord struct {
	RefreshID string `firestore:"refresh_id"`
	UserID    string `firestore:"user_id"`
	DigestB64 string `firestore:"digest_b64"`
	KeyID     string `firestore:"key_id"`

	FamilyID     string    `firestore:"family_id"`
	ReplacedBy   string    `firestore:"replaced_by"`
	RevokedAt    time.Time `firestore:"revoked_at"`
	RevokeReason string    `firestore:"revoke_reason"`

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
