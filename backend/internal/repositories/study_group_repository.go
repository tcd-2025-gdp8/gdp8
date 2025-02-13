package repositories

import (
	"errors"
	"sync"

	"gdp8-backend/internal/models"
)

type StudyGroupRepository interface {
	GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroup, error)
	GetAllStudyGroups() ([]models.StudyGroup, error)
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
				ID:          1,
				Name:        "Math Study Group",
				Description: "A group for studying mathematics.",
				Type:        models.TypePublic,
				ModuleID:    1,
				Members:     []models.UserID{1, 2},
			},
			2: {
				ID:          2,
				Name:        "Agile Study Group",
				Description: "A group for studying agile methods.",
				Type:        models.TypeClosed,
				ModuleID:    2,
				Members:     []models.UserID{2},
			},
			3: {
				ID:          3,
				Name:        "Elite Study Group",
				Description: "Very elite invite-only study group.",
				Type:        models.TypeInviteOnly,
				ModuleID:    1,
				Members:     []models.UserID{1},
			},
		},
		counter: 4,
	}
}

func (r *MockStudyGroupRepository) GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroup, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	studyGroup, exists := r.studyGroups[id]
	if !exists {
		return nil, ErrStudyGroupNotFound
	}
	return &studyGroup, nil
}

func (r *MockStudyGroupRepository) GetAllStudyGroups() ([]models.StudyGroup, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	studyGroupsList := make([]models.StudyGroup, 0, len(r.studyGroups))
	for _, studyGroup := range r.studyGroups {
		studyGroupsList = append(studyGroupsList, studyGroup)
	}
	return studyGroupsList, nil
}
