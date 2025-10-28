package idtoken

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeSigner struct {
	got map[string]any
	err error
}

func (f *fakeSigner) SignJWT(_ context.Context, payload map[string]any) (string, string, error) {
	f.got = payload
	if f.err != nil {
		return "", "", f.err
	}
	return "jwt.mock", "kid-1", nil
}

func fixed(t time.Time) time.Time { return t.UTC() }

func TestIssue_MinimalSuccess(t *testing.T) {
	t.Parallel()

	s := &fakeSigner{}
	uc := &IssueIDTokenUsecase{Issuer: "https://idpproxy.com", Signer: s}

	now := time.Unix(1_800_000_000, 0).UTC()
	in := &IDTokenInput{
		UserID:   "user-123",
		ClientID: "client-abc",
		Now:      now,
		TTL:      time.Hour,
	}

	jwt, kid, err := uc.Issue(context.Background(), in)
	require.NoError(t, err)
	require.Equal(t, "jwt.mock", jwt)
	require.Equal(t, "kid-1", kid)

	require.Equal(t, "https://idpproxy.com", s.got["iss"])
	require.Equal(t, "user-123", s.got["sub"])
	require.Equal(t, "client-abc", s.got["aud"])
	require.EqualValues(t, now.Unix(), s.got["iat"])
	require.EqualValues(t, now.Add(time.Hour).Unix(), s.got["exp"])

	require.NotContains(t, s.got, "auth_time")
	require.NotContains(t, s.got, "amr")
	require.NotContains(t, s.got, "nonce")
	require.NotContains(t, s.got, "at_hash")
	require.NotContains(t, s.got, "azp")
}

func TestIssue_WithAMRAndAuthTime(t *testing.T) {
	t.Parallel()

	s := &fakeSigner{}
	uc := &IssueIDTokenUsecase{Issuer: "https://idpproxy.com", Signer: s}

	now := time.Unix(1_900_000_000, 0).UTC()
	auth := now.Add(-5 * time.Minute)

	in := &IDTokenInput{
		UserID:   "u",
		ClientID: "c",
		Now:      now,
		TTL:      30 * time.Minute,
		AuthTime: &auth,
		AMR:      []string{"pwd", "mfa"},
	}

	_, _, err := uc.Issue(context.Background(), in)
	require.NoError(t, err)

	require.EqualValues(t, auth.Unix(), s.got["auth_time"])
	require.ElementsMatch(t, []string{"pwd", "mfa"}, s.got["amr"])
}

func TestIssue_ErrorCases(t *testing.T) {
	t.Parallel()

	t.Run("TTL<=0", func(t *testing.T) {
		t.Parallel()
		uc := &IssueIDTokenUsecase{Issuer: "x", Signer: &fakeSigner{}}
		_, _, err := uc.Issue(context.Background(), &IDTokenInput{
			UserID: "u", ClientID: "c", Now: time.Now(), TTL: 0,
		})
		require.Error(t, err)
	})

	t.Run("empty user -> ErrInvalidSubject", func(t *testing.T) {
		t.Parallel()
		uc := &IssueIDTokenUsecase{Issuer: "x", Signer: &fakeSigner{}}
		_, _, err := uc.Issue(context.Background(), &IDTokenInput{
			UserID: "", ClientID: "c", Now: time.Now(), TTL: time.Minute,
		})
		require.ErrorIs(t, err, ErrInvalidSubject)
	})

	t.Run("empty client -> ErrInvalidAudience", func(t *testing.T) {
		t.Parallel()
		uc := &IssueIDTokenUsecase{Issuer: "x", Signer: &fakeSigner{}}
		_, _, err := uc.Issue(context.Background(), &IDTokenInput{
			UserID: "u", ClientID: "", Now: time.Now(), TTL: time.Minute,
		})
		require.ErrorIs(t, err, ErrInvalidAudience)
	})
}

func TestIssue_SignerErrorPropagates(t *testing.T) {
	t.Parallel()

	want := errors.New("signer boom")
	s := &fakeSigner{err: want}
	uc := &IssueIDTokenUsecase{Issuer: "x", Signer: s}

	_, _, err := uc.Issue(context.Background(), &IDTokenInput{
		UserID: "u", ClientID: "c", Now: time.Now(), TTL: time.Minute,
	})
	require.ErrorIs(t, err, want)
}
