package services

import (
	"fmt"
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
	SendEventInvites(eventDetails EventDetails, invitees []models.UserID) error
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

func (s *eventInvitesServiceImpl) SendEventInvites(eventDetails EventDetails, invitees []models.UserID) error {
	var attendees []*calendar.EventAttendee
	for _, userID := range invitees {
		user, err := s.userService.GetUser(userID)
		if err != nil {
			return fmt.Errorf("failed to get user details for %s: %w", userID, err)
		}
		attendees = append(attendees, &calendar.EventAttendee{
			Email: user.Email,
		})
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
			Overrides: []*calendar.EventReminder{
				{Method: "email", Minutes: 24 * 60},
				{Method: "popup", Minutes: 10},
			},
			UseDefault:      false,
			ForceSendFields: []string{"UseDefault"},
		},
	}

	event, err := s.googleCalendarService.Events.Insert("primary", event).SendUpdates("all").Do()
	return err
}

type NoOpEventInvitesService struct{}

func (s *NoOpEventInvitesService) SendEventInvites(_ EventDetails, _ []models.UserID) error {
	return nil
}
