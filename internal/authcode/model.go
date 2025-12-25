package authcode

import "time"

type AuthCode struct {
	Code      string
	UserID    string
	ClientID  string
	ExpiresAt time.Time
}
