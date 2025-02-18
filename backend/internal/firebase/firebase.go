package firebase

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	fb "firebase.google.com/go/v4"
)

// InitializeFirebase initializes the Firebase Admin SDK and returns an *auth.Client.
func InitializeFirebase(credentialsPath string) (*auth.Client, error) {
	// credentialsPath should point to your serviceAccountKey.json
	opt := option.WithCredentialsFile(credentialsPath)

	// Create the Firebase App
	app, err := fb.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create new Firebase app: %w", err)
	}

	// Get the Auth client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get auth client: %w", err)
	}

	return authClient, nil
}
