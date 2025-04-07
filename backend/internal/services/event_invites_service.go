package services

import (
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
}

func NewEventInvitesService(googleCalendarService *calendar.Service) EventInvitesService {
	return &eventInvitesServiceImpl{
		googleCalendarService: googleCalendarService,
	}
}

func (s *eventInvitesServiceImpl) SendEventInvites(eventDetails EventDetails, _ []models.UserID) error {
	// TODO get invitees emails

	attendees := []*calendar.EventAttendee{{Email: "bronickm@tcd.ie"}}

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
