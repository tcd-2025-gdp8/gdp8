package models

import "time"

type StudySessionAvailabilityRequestID int64

type StudySessionAvailabilityRequest struct {
	ID           StudySessionAvailabilityRequestID
	StudyGroupID StudyGroupID
	StudySessionAvailabilityRequestDetails
	Entries []StudySessionAvailabilityEntry
}

type StudySessionAvailabilityRequestDetails struct {
	Title                   string
	AvailabilityPeriodStart time.Time
	AvailabilityPeriodEnd   time.Time
}

type StudySessionAvailabilityEntry struct {
	UserID UserID
	AvailabilityEntry
}

type AvailabilityEntry struct {
	AvailabilityEntryStart time.Time
	AvailabilityEntryEnd   time.Time
}
