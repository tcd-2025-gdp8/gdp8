package models

import "time"

type StudySessionID int64

type StudySession struct {
	StudySessionID StudySessionID
	StudyGroupID   StudyGroupID
	CreatorID      UserID
	StudySessionDetails
	EndTime time.Time
}

type StudySessionDetails struct {
	Title           string
	StartTime       time.Time
	DurationMinutes int64
}
