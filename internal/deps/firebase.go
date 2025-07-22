package deps

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func InitFirebaseApp(jsonPath string) (*firebase.App, error) {
	opt := option.WithCredentialsFile(jsonPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}
