package kms

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"golang.org/x/oauth2"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

const scopeCloudPlatform = "https://www.googleapis.com/auth/cloud-platform"

var (
	reSA           = regexp.MustCompile(`^[a-z][a-z0-9-]*@[a-z][a-z0-9-]*\.iam\.gserviceaccount\.com$`)
	newTokenSource = func(ctx context.Context, cfg impersonate.CredentialsConfig) (oauth2.TokenSource, error) {
		return impersonate.CredentialsTokenSource(ctx, cfg)
	}
)

func NewClient(ctx context.Context, impersonateSA string) (*cloudkms.KeyManagementClient, error) {
	impersonateSA = strings.ToLower(strings.TrimSpace(impersonateSA))
	if impersonateSA != "" {
		if !reSA.MatchString(impersonateSA) {
			return nil, fmt.Errorf("%w: invalid service account email %q", ErrInitFailed, impersonateSA)
		}

		ts, err := newTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: impersonateSA,
			Scopes:          []string{scopeCloudPlatform},
		})
		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("%w: impersonate failed for %q", ErrInitFailed, impersonateSA),
				err,
			)
		}

		cli, err := cloudkms.NewKeyManagementClient(ctx, option.WithTokenSource(ts))
		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("%w: kms client init failed (impersonated %q)", ErrInitFailed, impersonateSA),
				err,
			)
		}
		return cli, nil
	}

	cli, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("%w: kms client init failed (using ADC)", ErrInitFailed),
			err,
		)
	}

	return cli, nil
}
