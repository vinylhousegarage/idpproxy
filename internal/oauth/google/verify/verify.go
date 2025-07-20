package verify

import (
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
}

func ParseAndVerifyJWT(
	idToken string,
	pubKey *rsa.PublicKey,
	expectedIss, expectedAud string,
) (*MyClaims, error) {
	claims := &MyClaims{}

	_, err := jwt.ParseWithClaims(
		idToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, ErrInvalidSigningAlg
			}

			return pubKey, nil
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidSigningAlg):
			return nil, ErrInvalidSigningAlg
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		default:
			return nil, ErrJWTParseFailed
		}
	}

	if claims.Issuer != expectedIss {
		return nil, ErrInvalidIssuer
	}

	if len(claims.Audience) == 0 {
		return nil, ErrMissingAudience
	}

	validAud := false
	for _, aud := range claims.Audience {
		if aud == expectedAud {
			validAud = true
			break
		}
	}

	if !validAud {
		return nil, ErrInvalidAudience
	}

	if claims.Subject == "" {
		return nil, ErrMissingSubject
	}

	return claims, nil
}
