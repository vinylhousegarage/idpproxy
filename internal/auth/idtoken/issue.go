package idtoken

import (
	"context"
	"errors"
	"time"
)

type IssueIDTokenUsecase struct {
	Issuer string
	Signer Signer
}

func (uc *IssueIDTokenUsecase) Issue(ctx context.Context, in *IDTokenInput) (token string, kid string, err error) {
	if uc == nil || uc.Signer == nil || uc.Issuer == "" {
		return "", "", errors.New("idtoken: invalid usecase configuration")
	}
	if in == nil {
		return "", "", errors.New("idtoken: nil input")
	}
	if in.TTL <= 0 {
		return "", "", errors.New("idtoken: TTL must be > 0")
	}

	now := in.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	exp := now.Add(in.TTL)

	claims := &IDTokenClaims{
		Iss: uc.Issuer,
		Sub: in.UserID,
		Aud: in.ClientID,
		Iat: now.Unix(),
		Exp: exp.Unix(),
		AMR: in.AMR,
	}
	if in.AuthTime != nil {
		claims.AuthTime = in.AuthTime.UTC().Unix()
	}

	if err := claims.Validate(); err != nil {
		return "", "", err
	}

	payload := map[string]any{
		"iss": claims.Iss,
		"sub": claims.Sub,
		"aud": claims.Aud,
		"iat": claims.Iat,
		"exp": claims.Exp,
	}
	if claims.AuthTime != 0 {
		payload["auth_time"] = claims.AuthTime
	}
	if len(claims.AMR) > 0 {
		payload["amr"] = claims.AMR
	}
	if claims.Nonce != "" {
		payload["nonce"] = claims.Nonce
	}
	if claims.AtHash != "" {
		payload["at_hash"] = claims.AtHash
	}
	if claims.Azp != "" {
		payload["azp"] = claims.Azp
	}

	return uc.Signer.SignJWT(ctx, payload)
}
