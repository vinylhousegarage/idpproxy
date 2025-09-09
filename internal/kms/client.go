package kms

import (
	"context"
	"errors"
	"os"
	"regexp"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

var reSA = regexp.MustCompile(`^[a-zA-Z0-9\-]+@[a-z0-9\-]+\.iam\.gserviceaccount\.com$`)

func NewClient(ctx context.Context) (*cloudkms.KeyManagementClient, error) {
	if sa := os.Getenv("IMPERSONATE_SERVICE_ACCOUNT"); sa != "" {
		if !reSA.MatchString(sa) {
				return nil, fmt.Errorf("%w: invalid service account email %q", ErrInitFailed, sa)
		}
		ts, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
				TargetPrincipal: sa,
				Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
		})
		if err != nil {
				return nil, errors.Join(ErrInitFailed, err)
		}
		return cloudkms.NewKeyManagementClient(ctx, option.WithTokenSource(ts))
	}

	return cloudkms.NewKeyManagementClient(ctx)
}
