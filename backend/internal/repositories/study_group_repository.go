package repositories

import (
	"errors"
	"sync"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
)

type StudyGroupRepository interface {
	GetStudyGroupByID(tx persistence.Transaction, id models.StudyGroupID) (*models.StudyGroup, error)
	GetAllStudyGroups(tx persistence.Transaction) ([]models.StudyGroup, error)
	CreateStudyGroup(tx persistence.Transaction, studyGroupDetails models.StudyGroupDetails,
		adminUserID models.UserID) (*models.StudyGroup, error)
	UpdateStudyGroup(tx persistence.Transaction, studyGroup *models.StudyGroup) (*models.StudyGroup, error)
	DeleteStudyGroup(tx persistence.Transaction, id models.StudyGroupID) error
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
	id models.StudyGroupID) (*models.StudyGroup, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	studyGroup, exists := r.studyGroups[id]
	if !exists {
		return nil, ErrStudyGroupNotFound
	}
	return &studyGroup, nil
}

func (r *MockStudyGroupRepository) GetAllStudyGroups(_ persistence.Transaction) ([]models.StudyGroup, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	studyGroupsList := make([]models.StudyGroup, 0, len(r.studyGroups))
	for _, studyGroup := range r.studyGroups {
		studyGroupsList = append(studyGroupsList, studyGroup)
	}
	return studyGroupsList, nil
}

func (r *MockStudyGroupRepository) CreateStudyGroup(_ persistence.Transaction,
	studyGroupDetails models.StudyGroupDetails, adminUserID models.UserID) (*models.StudyGroup, error) {
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

	return &studyGroup, nil
}

func (r *MockStudyGroupRepository) UpdateStudyGroup(_ persistence.Transaction,
	studyGroup *models.StudyGroup) (*models.StudyGroup, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := studyGroup.ID

	_, exists := r.studyGroups[id]
	if !exists {
		return nil, ErrStudyGroupNotFound
	}

	r.studyGroups[id] = *studyGroup

	return studyGroup, nil
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
