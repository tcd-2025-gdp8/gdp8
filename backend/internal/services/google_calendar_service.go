package services

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// GoogleCalendarService defines methods for interacting with Google Calendar.
type GoogleCalendarService interface {
	// CreateEvent creates an event on the specified calendar.
	CreateEvent(calendarID string, event *calendar.Event) (*calendar.Event, error)
}

type googleCalendarServiceImpl struct {
	calendarService *calendar.Service
}

// NewGoogleCalendarService creates a new instance of GoogleCalendarService.
// The subject parameter is the email address of the user to impersonate.
func NewGoogleCalendarService(
	ctx context.Context,
	credentialsFile string,
	subject string,
) (GoogleCalendarService, error) {
	data, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	config, err := google.JWTConfigFromJSON(data, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT config: %w", err)
	}
	config.Subject = subject

	svc, err := calendar.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}
	return &googleCalendarServiceImpl{
		calendarService: svc,
	}, nil
}

// CreateEvent creates an event on the specified calendar.
func (g *googleCalendarServiceImpl) CreateEvent(calendarID string, event *calendar.Event) (*calendar.Event, error) {
	createdEvent, err := g.calendarService.Events.Insert(calendarID, event).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	return createdEvent, nil
}
