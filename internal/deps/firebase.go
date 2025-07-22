package deps

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"go.uber.org/zap"
)

func NewFirebaseApp() (*firebase.App, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func NewFirestoreClient(app *firebase.App, logger *zap.Logger) (*firestore.Client, error) {
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		logger.Error("failed to initialize Firestore via Firebase App", zap.Error(err))
		return nil, err
	}

	return client, nil
}
