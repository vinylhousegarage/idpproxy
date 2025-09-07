package kms

import (
	"context"
	"errors"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, opts ...option.ClientOption) (*cloudkms.KeyManagementClient, error) {
	client, err := cloudkms.NewKeyManagementClient(ctx, opts...)
	if err != nil {
		return nil, errors.Join(ErrInitFailed, err)
	}

	return client, nil
}
