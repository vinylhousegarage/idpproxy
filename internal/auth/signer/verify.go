package signer

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type VerifyOptions struct {
	Now    func() time.Time
	Leeway time.Duration

	ExpectKID  string
	RequireTyp bool
	ExpectTyp  string
}

type VerifyResult struct {
	Alg    string
	Typ    string
	KID    string
	Claims jwt.MapClaims
}

func (s *HMACSigner) parseAndVerifyHS256(token string, now func() time.Time, leeway time.Duration) (*jwt.Token, jwt.MapClaims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithLeeway(leeway),
		jwt.WithTimeFunc(now),
	)

	var claims jwt.MapClaims
	t, err := parser.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidAlg
		}

		return s.key, nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !t.Valid {
		return nil, nil, ErrInvalidToken
	}

	return t, claims, nil
}

func checkHeaderPolicy(t *jwt.Token, opt VerifyOptions) error {
	if opt.RequireTyp {
		typ, _ := t.Header["typ"].(string)
		if typ == "" || typ != opt.ExpectTyp {
			return ErrInvalidTyp
		}
	}
	if expKid := opt.ExpectKID; expKid != "" {
		kid, _ := t.Header["kid"].(string)
		if kid != expKid {
			return ErrUnexpectedKID
		}
	}

	return nil
}

func normalizeVerifyOptions(in *VerifyOptions) VerifyOptions {
	var o VerifyOptions
	if in != nil {
		o = *in
	}
	if o.RequireTyp && o.ExpectTyp == "" {
		o.ExpectTyp = "JWT"
	}
	return o
}

func (s *HMACSigner) Verify(ctx context.Context, token string, opt *VerifyOptions) (*VerifyResult, error) {
	_ = ctx
	if len(s.key) == 0 {
		return nil, ErrEmptyKey
	}
	if token == "" {
		return nil, ErrEmptyToken
	}

	o := normalizeVerifyOptions(opt)

	nowFn := s.Now
	if o.Now != nil {
		nowFn = o.Now
	}

	tok, claims, err := s.parseAndVerifyHS256(token, nowFn, o.Leeway)
	if err != nil {
		return nil, err
	}

	if err := checkHeaderPolicy(tok, o); err != nil {
		return nil, err
	}

	alg, _ := tok.Header["alg"].(string)
	typ, _ := tok.Header["typ"].(string)
	kid, _ := tok.Header["kid"].(string)

	return &VerifyResult{
		Alg:    alg,
		Typ:    typ,
		KID:    kid,
		Claims: claims,
	}, nil
}
