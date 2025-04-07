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

// StudyGroupMemberView represents a member with their user details
type StudyGroupMemberView struct {
	UserID UserID
	Name   string
	Role   StudyGroupRole
}

type StudyGroupDetails struct {
	Name        string
	Description string
	Type        StudyGroupType
	ModuleID    ModuleID
	MaxMembers  int
}

type StudyGroup struct {
	ID StudyGroupID
	StudyGroupDetails
	Members []StudyGroupMember
}

// StudyGroupView represents a study group with additional member details
type StudyGroupView struct {
	ID StudyGroupID
	StudyGroupDetails
	Members []StudyGroupMemberView
}

type StudyGroupStatsDTO struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Members     int     `json:"members"`
	TotalHours  int64   `json:"totalHours"`
	WeeklyHours []int64 `json:"weeklyHours"`
}

type StudyGroupsMap map[string][]StudyGroupStatsDTO
