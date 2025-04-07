package services

import (
	"context"
	"log"
	"os"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarService struct {
	svc *calendar.Service
}

func NewCalendarService(credentialsPath string) (*CalendarService, error) {
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		log.Println("⚠️ Calendar service credentials not found:", credentialsPath)
		return nil, nil // Safe fallback for CI
	}

	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}

	return &CalendarService{svc: srv}, nil
}

func (cs *CalendarService) GetService() *calendar.Service {
	return cs.svc
}
