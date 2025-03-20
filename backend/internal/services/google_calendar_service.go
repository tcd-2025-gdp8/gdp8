package services

import (
	"context"
	"fmt"

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

func NewGoogleCalendarService(ctx context.Context, credentialsFile string) (GoogleCalendarService, error) {
	svc, err := calendar.NewService(ctx, option.WithCredentialsFile(credentialsFile))
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
