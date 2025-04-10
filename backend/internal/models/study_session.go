package models

import (
	"database/sql"
	"time"
)

type StudySessionID int64

type StudySession struct {
	ID           StudySessionID
	StudyGroupID StudyGroupID
	CreatorID    UserID
	StudySessionDetails
	EndTime         time.Time
	CalendarEventID sql.NullString
}

type StudySessionDetails struct {
	Title           string
	StartTime       time.Time
	DurationMinutes int64
}
