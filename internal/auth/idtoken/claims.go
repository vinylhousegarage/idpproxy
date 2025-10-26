package idtoken

type IDTokenClaims struct {
	Iss      string   `json:"iss"`
	Sub      string   `json:"sub"`
	Aud      string   `json:"aud"`
	Exp      int64    `json:"exp"`
	Iat      int64    `json:"iat"`
	Nonce    string   `json:"nonce,omitempty"`
	AuthTime int64    `json:"auth_time,omitempty"`
	AMR      []string `json:"amr,omitempty"`
	AtHash   string   `json:"at_hash,omitempty"`
	Azp      string   `json:"azp,omitempty"`
}

func (c *IDTokenClaims) Validate() error {
	switch {
	case c.Iss == "":
		return ErrInvalidIssuer
	case c.Sub == "":
		return ErrInvalidSubject
	case c.Aud == "":
		return ErrInvalidAudience
	case c.Iat <= 0:
		return ErrInvalidIat
	case c.Exp <= c.Iat:
		return ErrInvalidExp
	default:
		return nil
	}
}
