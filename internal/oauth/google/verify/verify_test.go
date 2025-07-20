package verify

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateParseAndVerifyJWTTestToken(t *testing.T, claims jwt.RegisteredClaims, privateKey *rsa.PrivateKey) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	assert.NoError(t, err)
	return signed
}

func TestParseAndVerifyJWT(t *testing.T) {
	t.Parallel()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		claims := jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			Subject:   "user-123",
			Audience:  []string{"test-aud"},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		}
		token := generateParseAndVerifyJWTTestToken(t, claims, privKey)
		result, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "test-aud")
		assert.NoError(t, err)
		assert.Equal(t, "user-123", result.Subject)
	})

	t.Run("invalid-alg", func(t *testing.T) {
		t.Parallel()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			Subject:   "user-123",
			Audience:  []string{"test-aud"},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		})
		signed, err := token.SignedString([]byte("secret"))
		assert.NoError(t, err)

		_, err = ParseAndVerifyJWT(signed, &privKey.PublicKey, "test-issuer", "test-aud")
		assert.ErrorIs(t, err, ErrInvalidSigningAlg)
	})

	t.Run("invalid-token", func(t *testing.T) {
		t.Parallel()

		brokenToken := "this.is.not.a.valid.jwt"

		_, err := ParseAndVerifyJWT(brokenToken, &privKey.PublicKey, "test-issuer", "test-aud")
		assert.ErrorIs(t, err, ErrJWTParseFailed)
	})

	t.Run("expired", func(t *testing.T) {
		t.Parallel()

		claims := jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			Subject:   "user-123",
			Audience:  []string{"test-aud"},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-1 * time.Hour)),
		}

		token := generateParseAndVerifyJWTTestToken(t, claims, privKey)

		_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "test-aud")
		assert.ErrorIs(t, err, ErrTokenExpired)
	})

	t.Run("invalid-issuer", func(t *testing.T) {
		t.Parallel()

		claims := jwt.RegisteredClaims{
			Issuer:    "wrong-issuer",
			Subject:   "user-123",
			Audience:  []string{"test-aud"},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		}

		token := generateParseAndVerifyJWTTestToken(t, claims, privKey)

		_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "expected-issuer", "test-aud")
		assert.ErrorIs(t, err, ErrInvalidIssuer)
	})

	t.Run("invalid-audience", func(t *testing.T) {
		t.Parallel()

		claims := jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			Subject:   "user-123",
			Audience:  []string{"wrong-aud"},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		}

		token := generateParseAndVerifyJWTTestToken(t, claims, privKey)

		_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "expected-aud")
		assert.ErrorIs(t, err, ErrInvalidAudience)
	})

	t.Run("missing-sub", func(t *testing.T) {
		t.Parallel()

		claims := jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			Audience:  []string{"test-aud"},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		}

		token := generateParseAndVerifyJWTTestToken(t, claims, privKey)

		_, err := ParseAndVerifyJWT(token, &privKey.PublicKey, "test-issuer", "test-aud")
		assert.ErrorIs(t, err, ErrMissingSubject)
	})
}
