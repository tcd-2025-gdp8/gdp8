package firebase

import (
	"context"
	"fmt"
	"os"

	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	fb "firebase.google.com/go/v4"
)

func InitializeFirebase(credentialsPath string) (*auth.Client, error) {
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("credentials file does not exist: %s", credentialsPath)
	}

	opt := option.WithCredentialsFile(credentialsPath)

	app, err := fb.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create new Firebase app: %w", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get auth client: %w", err)
	}

	return authClient, nil
}
