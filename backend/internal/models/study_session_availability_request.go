package models

import "time"

type StudySessionAvailabilityRequestID int64

type StudySessionAvailabilityRequest struct {
	ID           StudySessionAvailabilityRequestID `json:"id" db:"id"`
	StudyGroupID StudyGroupID                      `json:"study_group_id" db:"study_group_id"`
	CreatorID    UserID                            `json:"creator_id" db:"creator_id"`
	StudySessionAvailabilityRequestDetails
	Entries []StudySessionAvailabilityEntry `json:"entries"`
}

type StudySessionAvailabilityRequestDetails struct {
	Title                   string    `json:"title" db:"title"`
	AvailabilityPeriodStart time.Time `json:"availability_period_start" db:"availability_period_start"`
	AvailabilityPeriodEnd   time.Time `json:"availability_period_end" db:"availability_period_end"`
}

type StudySessionAvailabilityEntry struct {
	UserID            UserID            `json:"user_id" db:"user_id"`
	AvailabilityEntry AvailabilityEntry `json:"availability_entry"`
}

type AvailabilityEntry struct {
	AvailabilityEntryStart time.Time `json:"availability_entry_start" db:"availability_entry_start"`
	AvailabilityEntryEnd   time.Time `json:"availability_entry_end" db:"availability_entry_end"`
}
