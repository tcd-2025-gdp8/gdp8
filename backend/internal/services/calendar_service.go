package services

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarService struct {
	srv *calendar.Service
}

func NewCalendarService(credentialsPath string) (*CalendarService, error) {
	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}
	return &CalendarService{srv: srv}, nil
}

type CalendarInviteInput struct {
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	StartTime   string   `json:"startTime"` // ISO 8601
	EndTime     string   `json:"endTime"`   // ISO 8601
	Attendees   []string `json:"attendees"`
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

	for _, email := range input.Attendees {
		event.Attendees = append(event.Attendees, &calendar.EventAttendee{Email: email})
	}

	createdEvent, err := c.srv.Events.Insert("primary", event).Do()
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	log.Printf("✅ Calendar event created: %s", createdEvent.HtmlLink)
	return nil
}
