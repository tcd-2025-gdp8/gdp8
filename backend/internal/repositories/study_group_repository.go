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
					Name:        "CS Wizards",
					Description: "A group for computer science wizards who excel in coding and problem-solving.",
					Type:        models.TypePublic,
					ModuleID:    3,
				},
				Members: []models.StudyGroupMember{
					{UserID: "David", Role: models.RoleAdmin},
					{UserID: "Eve", Role: models.RoleMember},
					{UserID: "Frank", Role: models.RoleMember},
				},
			},
			3: {
				ID: 3,
				StudyGroupDetails: models.StudyGroupDetails{
					Name:        "The Elites",
					Description: "A group for elite students who aim for excellence in their academic pursuits.",
					Type:        models.TypePublic,
					ModuleID:    1,
				},
				Members: []models.StudyGroupMember{
					{UserID: "Grace", Role: models.RoleAdmin},
					{UserID: "Hannah", Role: models.RoleMember},
					{UserID: "Ian", Role: models.RoleMember},
				},
			},
			4: {
				ID: 4,
				StudyGroupDetails: models.StudyGroupDetails{
					Name:        "The Fun Group",
					Description: "A group for students who believe in having fun while learning and collaborating.",
					Type:        models.TypePublic,
					ModuleID:    2,
				},
				Members: []models.StudyGroupMember{
					{UserID: "Jack", Role: models.RoleAdmin},
					{UserID: "Kate", Role: models.RoleMember},
					{UserID: "Leo", Role: models.RoleMember},
					{UserID: "Blake", Role: models.RoleMember},
					{UserID: "Robert", Role: models.RoleMember},
					{UserID: "Marco", Role: models.RoleMember},
				},
			},
			5: {
				ID: 5,
				StudyGroupDetails: models.StudyGroupDetails{
					Name:        "The Prefects",
					Description: "A group for prefects who lead by example and strive for academic and personal growth.",
					Type:        models.TypePublic,
					ModuleID:    3,
				},
				Members: []models.StudyGroupMember{
					{UserID: "Mike", Role: models.RoleAdmin},
					{UserID: "Nina", Role: models.RoleMember},
					{UserID: "Oscar", Role: models.RoleMember},
					{UserID: "Alessandro", Role: models.RoleMember},
					{UserID: "Alice", Role: models.RoleMember},
					{UserID: "David", Role: models.RoleMember},
					{UserID: "Grace", Role: models.RoleMember},
					{UserID: "Ava", Role: models.RoleMember},
				},
			},
			6: {
				ID: 6,
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
		counter: 7,
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
