package models

import "time"

type ChatMessageID int64

type ChatMessageView struct {
	ID        ChatMessageID
	UserID    UserID
	UserName  string
	Text      string
	Timestamp time.Time
}

type ChatMessage struct {
	ID           ChatMessageID
	StudyGroupID StudyGroupID
	ChatMessageDetails
}

type ChatMessageDetails struct {
	UserID    UserID
	Text      string
	Timestamp time.Time
}
