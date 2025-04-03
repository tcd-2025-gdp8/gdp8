package models

import (
	"encoding/json"
	"time"
)

type NotificationID int64
type NotificationType string
type NotificationPayloadType json.RawMessage

const (
	NotificationTypeStudyGroupJoined              NotificationType = "study-group-joined"
	NotificationTypeStudyGroupRequestedToJoin     NotificationType = "study-group-requested-to-join"
	NotificationTypeStudyGroupLeft                NotificationType = "study-group-left"
	NotificationTypeStudyGroupAcceptedInvite      NotificationType = "study-group-accepted-invite"
	NotificationTypeStudyGroupRejectedInvite      NotificationType = "study-group-rejected-invite"
	NotificationTypeStudyGroupInvited             NotificationType = "study-group-invited"
	NotificationTypeStudyGroupAcceptedJoinRequest NotificationType = "study-group-accepted-join-request"
	NotificationTypeStudyGroupRejectedJoinRequest NotificationType = "study-group-rejected-join-request"
	NotificationTypeStudyGroupRemovedMember       NotificationType = "study-group-removed-member"

	NotificationTypeStudyGroupChatMessage NotificationType = "study-group-chat-message"

	NotificationTypeStudySessionScheduled NotificationType = "study-session-scheduled"
	NotificationTypeStudySessionUpdated   NotificationType = "study-session-updated"
	NotificationTypeStudySessionCancelled NotificationType = "study-session-cancelled"
	NotificationTypeStudySessionReminder  NotificationType = "study-session-reminder"

	NotificationTypeStudySessionAvailabilityRequestCreated NotificationType = "study-session-availability-request-created"
	NotificationTypeStudySessionAvailabilityEntriesUpdated NotificationType = "study-session-availability-request-updated-entries"
)

type NotificationDetails struct {
	Type    NotificationType
	Payload NotificationPayloadType
}

type Notification struct {
	ID NotificationID
	NotificationDetails
	CreatedAt time.Time
}
