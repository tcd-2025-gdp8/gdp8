package models

import (
	"time"
)

type NotificationID int64
type NotificationType string

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
	NotificationTypeStudyGroupChatMessage         NotificationType = "study-group-chat-message"
)

type NotificationDetails struct {
	Type             NotificationType
	TriggeringUserID UserID
	TargetUserID     *UserID
	StudyGroupID     StudyGroupID
	MessageID        *int64 // TODO no MessageID type exists yet, to be changed
}

type Notification struct {
	ID NotificationID
	NotificationDetails
	CreatedAt time.Time
}

type NotificationView struct {
	ID                 NotificationID
	Type               NotificationType
	TriggeringUserID   UserID
	TriggeringUserName string
	TargetUserID       *UserID
	TargetUserName     *string
	StudyGroupID       StudyGroupID
	StudyGroupName     string
	MessageID          *int64 // TODO no MessageID type exists yet, to be changed
	CreatedAt          time.Time
}
