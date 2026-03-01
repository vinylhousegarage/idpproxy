package authcode

import "time"

type ProxyCode struct {
	Code      string
	UserID    string
	ClientID  string
	ExpiresAt time.Time
}
