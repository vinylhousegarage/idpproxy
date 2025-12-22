package token

type TokenResponse struct {
	IDToken   string `json:"id_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
}
