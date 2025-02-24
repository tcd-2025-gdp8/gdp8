package repositories

import (
	"errors"
	"sync"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
)

type StudyGroupRepository interface {
	GetStudyGroupByID(tx persistence.Transaction, id models.StudyGroupID) (*models.StudyGroupView, error)
	GetAllStudyGroups(tx persistence.Transaction) ([]models.StudyGroupView, error)
	CreateStudyGroup(tx persistence.Transaction, studyGroupDetails models.StudyGroupDetails,
		adminUserID models.UserID) (*models.StudyGroupView, error)
	UpdateStudyGroupDetails(tx persistence.Transaction, id models.StudyGroupID,
		details models.StudyGroupDetails) (*models.StudyGroupView, error)
	DeleteStudyGroup(tx persistence.Transaction, id models.StudyGroupID) error
	UpdateStudyGroupMember(tx persistence.Transaction, id models.StudyGroupID,
		userID models.UserID, role *models.StudyGroupRole) error
}

var ErrStudyGroupNotFound = errors.New("study group not found")

type MockStudyGroupRepository struct {
	studyGroups map[models.StudyGroupID]models.StudyGroup
	counter     int
	mu          sync.Mutex
}

func NewMockStudyGroupRepository() StudyGroupRepository {
	return &MockStudyGroupRepository{
		studyGroups: map[models.StudyGroupID]models.StudyGroup{
			1: {
				ID: 1,
				StudyGroupDetails: models.StudyGroupDetails{
					Name:        "Tech Nerds",
					Description: "A group for tech enthusiasts who love to explore new technologies and innovations.",
					Type:        models.TypePublic,
					ModuleID:    1,
				},
				Members: []models.StudyGroupMember{
					{UserID: "Alice", Role: models.RoleAdmin},
					{UserID: "Bob", Role: models.RoleMember},
					{UserID: "Charlie", Role: models.RoleMember},
					{UserID: "Maria", Role: models.RoleMember},
					{UserID: "Catriona", Role: models.RoleMember},
				},
			},
			2: {
				ID: 2,
				StudyGroupDetails: models.StudyGroupDetails{
					Name:        "The Elites",
					Description: "A group for elite students who aim for excellence in their academic pursuits.",
					Type:        models.TypeClosed,
					ModuleID:    1,
				},
				Members: []models.StudyGroupMember{
					{UserID: "Grace", Role: models.RoleAdmin},
					{UserID: "Alessandro", Role: models.RoleMember},
					{UserID: "Ian", Role: models.RoleMember},
				},
			},
			3: {
				ID: 3,
				StudyGroupDetails: models.StudyGroupDetails{
					Name:        "Trinners for Winners",
					Description: "A group for final year project students who are dedicated to achieving outstanding results.",
					Type:        models.TypePublic,
					ModuleID:    6,
				},
				Members: []models.StudyGroupMember{
					{UserID: "Paul", Role: models.RoleAdmin},
					{UserID: "Quinn", Role: models.RoleMember},
					{UserID: "Rachel", Role: models.RoleMember},
					{UserID: "Jade", Role: models.RoleMember},
					{UserID: "Robert", Role: models.RoleMember},
					{UserID: "Bob", Role: models.RoleMember},
					{UserID: "Hannah", Role: models.RoleMember},
					{UserID: "Bianca", Role: models.RoleMember},
					{UserID: "Oscar", Role: models.RoleMember},
					{UserID: "Ava", Role: models.RoleMember},
				},
			},
		},
		counter: 4,
	}
}

func (r *MockStudyGroupRepository) GetStudyGroupByID(_ persistence.Transaction,
	id models.StudyGroupID) (*models.StudyGroupView, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	studyGroup, exists := r.studyGroups[id]
	if !exists {
		return nil, ErrStudyGroupNotFound
	}

	return convertToView(&studyGroup), nil
}

func (r *MockStudyGroupRepository) GetAllStudyGroups(
	_ persistence.Transaction) ([]models.StudyGroupView, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	studyGroupsView := make([]models.StudyGroupView, 0, len(r.studyGroups))
	for _, studyGroup := range r.studyGroups {
		studyGroupsView = append(studyGroupsView, *convertToView(&studyGroup))
	}
	return studyGroupsView, nil
}

func (r *MockStudyGroupRepository) CreateStudyGroup(_ persistence.Transaction,
	studyGroupDetails models.StudyGroupDetails, adminUserID models.UserID) (*models.StudyGroupView, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := models.StudyGroupID(r.counter)
	r.counter++

	studyGroup := models.StudyGroup{
		ID:                id,
		StudyGroupDetails: studyGroupDetails,
		Members: []models.StudyGroupMember{
			{
				UserID: adminUserID,
				Role:   models.RoleAdmin,
			},
		},
	}
	r.studyGroups[id] = studyGroup

	return convertToView(&studyGroup), nil
}

func (r *MockStudyGroupRepository) UpdateStudyGroupDetails(_ persistence.Transaction,
	id models.StudyGroupID, details models.StudyGroupDetails) (*models.StudyGroupView, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	studyGroup, exists := r.studyGroups[id]
	if !exists {
		return nil, ErrStudyGroupNotFound
	}

	studyGroup.StudyGroupDetails = details
	r.studyGroups[id] = studyGroup

	return convertToView(&studyGroup), nil
}

// UpdateStudyGroupMember updates the role of a member in a study group or removes the member if the role is nil.
// The operation is idempotent - it will not return an error if the user already has the requested role.
// Returns an error if the study group does not exist.
func (r *MockStudyGroupRepository) UpdateStudyGroupMember(_ persistence.Transaction,
	id models.StudyGroupID, userID models.UserID, role *models.StudyGroupRole) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	studyGroup, exists := r.studyGroups[id]
	if !exists {
		return ErrStudyGroupNotFound
	}

	for i, member := range studyGroup.Members {
		if member.UserID == userID {
			if role == nil {
				studyGroup.Members = append(studyGroup.Members[:i], studyGroup.Members[i+1:]...)
			} else {
				studyGroup.Members[i].Role = *role
			}
			r.studyGroups[id] = studyGroup
			return nil
		}
	}

	studyGroup.Members = append(studyGroup.Members, models.StudyGroupMember{
		UserID: userID,
		Role:   *role,
	})
	r.studyGroups[id] = studyGroup
	return nil
}

func (r *MockStudyGroupRepository) DeleteStudyGroup(_ persistence.Transaction, id models.StudyGroupID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.studyGroups[id]
	if !exists {
		return ErrStudyGroupNotFound
	}

	delete(r.studyGroups, id)

	return nil
}

func convertToView(studyGroup *models.StudyGroup) *models.StudyGroupView {
	membersView := make([]models.StudyGroupMemberView, 0, len(studyGroup.Members))
	for _, member := range studyGroup.Members {
		membersView = append(membersView, models.StudyGroupMemberView{
			UserID: member.UserID,
			Name:   string(member.UserID),
			Role:   member.Role,
		})
	}

	return &models.StudyGroupView{
		ID:                studyGroup.ID,
		StudyGroupDetails: studyGroup.StudyGroupDetails,
		Members:           membersView,
	}
}
