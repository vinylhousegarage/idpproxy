package store

import "time"

func isActive(rec *RefreshTokenRecord, now time.Time) bool {
	if rec == nil {
		return false
	}

	if !rec.RevokedAt.IsZero() {
		return false
	}

	if rec.ReplacedBy != "" {
		return false
	}

	if !rec.ExpiresAt.IsZero() && now.After(rec.ExpiresAt) {
		return false
	}

	return true
}
