package idtoken

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"strings"
)

func computeAtHash(alg, accessToken string) (string, error) {
	if accessToken == "" {
		return "", errors.New("at_hash: empty access token")
	}

	var sum []byte
	switch strings.ToUpper(alg) {
	case "RS256", "PS256", "ES256", "HS256":
		s := sha256.Sum256([]byte(accessToken))
		sum = s[:]
	case "RS384", "PS384", "ES384", "HS384":
		s := sha512.Sum384([]byte(accessToken))
		sum = s[:]
	case "RS512", "PS512", "ES512", "HS512":
		s := sha512.Sum512([]byte(accessToken))
		sum = s[:]
	default:
		return "", errors.New("at_hash: unsupported alg: " + alg)
	}
	half := sum[:len(sum)/2]
	return base64.RawURLEncoding.EncodeToString(half), nil
}
