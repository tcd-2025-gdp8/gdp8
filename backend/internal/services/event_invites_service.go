package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"

	"gdp8-backend/internal/models"
)

type EventDetails struct {
	Summary     string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

type EventInvitesService interface {
	SendEventInvites(eventDetails EventDetails, invitees []models.UserID) (string, error)
	CancelEvent(eventID string) error
}

type eventInvitesServiceImpl struct {
	googleCalendarService *calendar.Service
	userService           UserService
}

func NewEventInvitesService(googleCalendarService *calendar.Service, userService UserService) EventInvitesService {
	return &eventInvitesServiceImpl{
		googleCalendarService: googleCalendarService,
		userService:           userService,
	}
}

func (s *eventInvitesServiceImpl) SendEventInvites(
	eventDetails EventDetails,
	invitees []models.UserID,
) (string, error) {
	attendees := make([]*calendar.EventAttendee, 0, len(invitees))
	for _, userID := range invitees {
		user, err := s.userService.GetUser(userID)
		if err != nil {
			return "", fmt.Errorf("failed to get user details for %s: %w", userID, err)
		}
		attendees = append(attendees, &calendar.EventAttendee{Email: user.Email})
	}

	event := &calendar.Event{
		Summary:     eventDetails.Summary,
		Description: eventDetails.Description,
		Start: &calendar.EventDateTime{
			DateTime: eventDetails.StartTime.Format(time.RFC3339),
			TimeZone: "Europe/Dublin",
		},
		End: &calendar.EventDateTime{
			DateTime: eventDetails.EndTime.Format(time.RFC3339),
			TimeZone: "Europe/Dublin",
		},
		Attendees: attendees,
		Reminders: &calendar.EventReminders{
			UseDefault: false,
			Overrides: []*calendar.EventReminder{
				{Method: "email", Minutes: 24 * 60},
				{Method: "popup", Minutes: 10},
			},
			ForceSendFields: []string{"UseDefault"},
		},
	}

	ctx := context.Background()
	createdEvent, err := s.googleCalendarService.Events.Insert("primary", event).SendUpdates("all").Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to insert calendar event: %w", err)
	}

	log.Printf("Google Calendar event created with ID: %s", createdEvent.Id)
	return createdEvent.Id, nil
}

func (s *eventInvitesServiceImpl) CancelEvent(eventID string) error {
	ctx := context.Background()
	err := s.googleCalendarService.Events.Delete("primary", eventID).SendUpdates("all").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to cancel Google Calendar event with ID %s: %w", eventID, err)
	}
	log.Printf(" Successfully cancelled Google Calendar event ID: %s", eventID)
	return nil
}

type NoOpEventInvitesService struct{}

func (s *NoOpEventInvitesService) SendEventInvites(_ EventDetails, _ []models.UserID) (string, error) {
	log.Println("📭 NoOpEventInvitesService: Pretending to send invite.")
	return "no-op-event-id", nil
}

func (s *NoOpEventInvitesService) CancelEvent(_ string) error {
	log.Println(" NoOpEventInvitesService: Pretending to cancel event.")
	return nil
}
