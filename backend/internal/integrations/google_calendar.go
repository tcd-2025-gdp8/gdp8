package integrations

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetCalendarService() (*calendar.Service, error) {
	ctx := context.Background()

	clientID := os.Getenv("GOOGLE_CALENDAR_HOST_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CALENDAR_HOST_CLIENT_SECRET")
	refreshToken := os.Getenv("GOOGLE_CALENDAR_HOST_REFRESH_TOKEN")

	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return nil, fmt.Errorf("required environment variables are not set: GOOGLE_CALENDAR_HOST_CLIENT_ID, " +
			"GOOGLE_CALENDAR_HOST_CLIENT_SECRET, GOOGLE_CALENDAR_HOST_REFRESH_TOKEN")
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		Scopes: []string{"https://www.googleapis.com/auth/calendar"},
	}

	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	client := conf.Client(ctx, token)

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create google calendar service: %w", err)
	}
	return srv, nil
}
