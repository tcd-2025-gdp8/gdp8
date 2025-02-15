package models

type StudyGroupID int64
type StudyGroupType string
type StudyGroupRole string

const (
	TypePublic     StudyGroupType = "public"
	TypeClosed     StudyGroupType = "closed"
	TypeInviteOnly StudyGroupType = "invite-only"
	RoleAdmin      StudyGroupRole = "admin"
	RoleMember     StudyGroupRole = "member"
	RoleInvitee    StudyGroupRole = "invitee"
	RoleRequester  StudyGroupRole = "requester"
)

type StudyGroupMember struct {
	UserID UserID
	Role   StudyGroupRole
}

type StudyGroup struct {
	ID          StudyGroupID
	Name        string
	Description string
	Type        StudyGroupType
	ModuleID    ModuleID
	Members     []StudyGroupMember
}
