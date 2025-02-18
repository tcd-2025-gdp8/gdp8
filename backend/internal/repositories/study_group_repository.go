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
	CreateStudyGroup(tx persistence.Transaction, studyGroupDetails models.StudyGroupDetails, adminUserID models.UserID) (*models.StudyGroup, error)
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
					Name:        "Math Study Group",
					Description: "A group for studying mathematics.",
					Type:        models.TypePublic,
					ModuleID:    1,
				},
				Members: []models.StudyGroupMember{
					{
						UserID: 1,
						Role:   models.RoleAdmin,
					},
					{
						UserID: 2,
						Role:   models.RoleMember,
					},
				},
			},
			2: {
				ID: 2,
				StudyGroupDetails: models.StudyGroupDetails{
					Name:        "Agile Study Group",
					Description: "A group for studying agile methods.",
					Type:        models.TypeClosed,
					ModuleID:    2,
				},
				Members: []models.StudyGroupMember{
					{
						UserID: 2,
						Role:   models.RoleAdmin,
					},
				},
			},
			3: {
				ID: 3,
				StudyGroupDetails: models.StudyGroupDetails{
					Name:        "Elite Study Group",
					Description: "Very elite invite-only study group.",
					Type:        models.TypeInviteOnly,
					ModuleID:    1,
				},
				Members: []models.StudyGroupMember{
					{
						UserID: 1,
						Role:   models.RoleAdmin,
					},
				},
			},
		},
		counter: 4,
	}
}

func (r *MockStudyGroupRepository) GetStudyGroupByID(_ persistence.Transaction, id models.StudyGroupID) (*models.StudyGroup, error) {
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

func (r *MockStudyGroupRepository) CreateStudyGroup(tx persistence.Transaction, studyGroupDetails models.StudyGroupDetails, adminUserID models.UserID) (*models.StudyGroup, error) {
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

func (r *MockStudyGroupRepository) UpdateStudyGroup(tx persistence.Transaction, studyGroup *models.StudyGroup) (*models.StudyGroup, error) {
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

func (r *MockStudyGroupRepository) DeleteStudyGroup(tx persistence.Transaction, id models.StudyGroupID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.studyGroups[id]
	if !exists {
		return ErrStudyGroupNotFound
	}

	delete(r.studyGroups, id)

	return nil
}
