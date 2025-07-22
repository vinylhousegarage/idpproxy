import (
	"context"

	"firebase.google.com/go/v4"
)

func InitFirebaseApp() (*firebase.App, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return app, nil
}
