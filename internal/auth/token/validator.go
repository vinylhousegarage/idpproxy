package token

import "time"

func validateParams(userID string, ttl, purgeAfter time.Duration) error {
	switch {
	case userID == "":
		return ErrEmptyUserID
	case ttl <= 0:
		return ErrInvalidTTL
	case purgeAfter < ttl:
		return ErrInvalidPurge
	default:
		return nil
	}
}
