package secret

import (
	"context"
	"errors"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, opts ...option.ClientOption) (*secretmanager.Client, error) {
	c, err := secretmanager.NewClient(ctx, opts...)
	if err != nil {
		return nil, errors.Join(ErrInitFailed, err)
	}

	return c, nil
}
