package models

type StudyGroupID int64
type StudyGroupType string

const (
	TypePublic     StudyGroupType = "public"
	TypeClosed     StudyGroupType = "closed"
	TypeInviteOnly StudyGroupType = "invite-only"
)

type StudyGroup struct {
	ID          StudyGroupID
	Name        string
	Description string
	Type        StudyGroupType
	ModuleID    ModuleID
	Members     []UserID
}
