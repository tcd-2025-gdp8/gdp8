package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarService struct {
	srv *calendar.Service
}

func NewCalendarService(credentialsPath string) (*CalendarService, error) {
	ctx := context.Background()

	// Read credentials file
	_, err := os.Stat(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("credentials file not found: %w", err)
	}

	srv, err := calendar.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize calendar service: %w", err)
	}

	return &CalendarService{srv: srv}, nil
}

type CalendarInviteInput struct {
	Summary     string
	Description string
	Location    string
	StartTime   string
	EndTime     string
	Attendees   []string
	Organizer   string
}

func (c *CalendarService) SendInvite(input CalendarInviteInput) error {
	event := &calendar.Event{
		Summary:     input.Summary,
		Location:    input.Location,
		Description: input.Description,
		Start: &calendar.EventDateTime{
			DateTime: input.StartTime,
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: input.EndTime,
			TimeZone: "UTC",
		},
	}

	calendarID := "primary"
	createdEvent, err := c.srv.Events.Insert(calendarID, event).Do()
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	log.Printf("✅ Calendar event created: %s\n", createdEvent.HtmlLink)
	return nil
}
