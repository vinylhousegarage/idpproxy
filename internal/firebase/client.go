package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

func NewFirebaseApp(
	ctx context.Context,
	opts ...option.ClientOption,
) (*firebase.App, error) {
	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func NewAuthClient(
	ctx context.Context,
	app *firebase.App,
	logger *zap.Logger,
) (*auth.Client, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		logger.Error("failed to initialize Firebase Auth Client", zap.Error(err))
		return nil, err
	}

	return client, nil
}

func NewFirestoreClient(
	ctx context.Context,
	app *firebase.App,
	logger *zap.Logger,
) (*firestore.Client, error) {
	client, err := app.Firestore(ctx)
	if err != nil {
		logger.Error("failed to initialize Firestore Client", zap.Error(err))
		return nil, err
	}

	return client, nil
}
