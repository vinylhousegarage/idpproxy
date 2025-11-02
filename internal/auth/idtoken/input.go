package idtoken

import "time"

type IDTokenInput struct {
	UserID   string
	ClientID string
	Now      time.Time
	TTL      time.Duration
	AuthTime *time.Time
	AMR      []string

	AccessToken string
	SignAlg     string

	Nonce string
	Azp   string
}
