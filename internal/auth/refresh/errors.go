package refresh

import "errors"

var (
	ErrEmptyUserID  = errors.New("userID empty")
	ErrInvalidTTL   = errors.New("ttl must be > 0")
	ErrInvalidPurge = errors.New("purgeAfter must be >= ttl")
	ErrRandFailure  = errors.New("rand failure")
)
