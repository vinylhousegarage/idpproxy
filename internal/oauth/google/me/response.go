package me

type MeResponse struct {
	Sub string `json:"sub"`
	Iss string `json:"iss"`
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
}
