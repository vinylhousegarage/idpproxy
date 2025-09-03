package kms

import (
	"context"
	"fmt"

	"cloud.google.com/go/kms/apiv1"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, logger *zap.Logger, opts ...option.ClientOption) (*kms.KeyManagementClient, error) {
	client, err := kms.NewKeyManagementClient(ctx, opts...)
	if err != nil {
		if logger != nil {
			logger.Error("failed to initialize KMS client", zap.Error(err))
		}
		return nil, fmt.Errorf("kms.NewClient: %w", err)
	}
	return client, nil
}
