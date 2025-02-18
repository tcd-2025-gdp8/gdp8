package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
)

type MockStudyGroupService struct {
	mock.Mock
}

func (m *MockStudyGroupService) GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroup, error) {
	args := m.Called(id)
	return args.Get(0).(*models.StudyGroup), args.Error(1)
}

func (m *MockStudyGroupService) GetAllStudyGroups() ([]models.StudyGroup, error) {
	args := m.Called()
	return args.Get(0).([]models.StudyGroup), args.Error(1)
}

func (m *MockStudyGroupService) CreateStudyGroup(studyGroupDetails models.StudyGroupDetails, creatorID models.UserID) (*models.StudyGroup, error) {
	args := m.Called(studyGroupDetails, creatorID)
	return args.Get(0).(*models.StudyGroup), args.Error(1)
}

func (m *MockStudyGroupService) UpdateStudyGroupDetails(id models.StudyGroupID, details models.StudyGroupDetails, requesterID models.UserID) (*models.StudyGroup, error) {
	args := m.Called(id, details, requesterID)
	return args.Get(0).(*models.StudyGroup), args.Error(1)
}

func (m *MockStudyGroupService) DeleteStudyGroup(id models.StudyGroupID, requesterID models.UserID) error {
	args := m.Called(id, requesterID)
	return args.Error(0)
}

func (m *MockStudyGroupService) HandleAdminMemberOperation(command services.AdminMemberOperationCommand, studyGroupID models.StudyGroupID, memberID models.UserID, adminID models.UserID) error {
	args := m.Called(command, studyGroupID, memberID, adminID)
	return args.Error(0)
}

func (m *MockStudyGroupService) HandleSelfMemberOperation(command services.SelfMemberOperationCommand, studyGroupID models.StudyGroupID, memberID models.UserID) error {
	args := m.Called(command, studyGroupID, memberID)
	return args.Error(0)
}

func TestStudyGroupHandler_GetStudyGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		url          string
		mockSetup    func(service *MockStudyGroupService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "Valid ID and found",
			url:  "/study-groups/123",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetStudyGroupByID", models.StudyGroupID(123)).
					Return(&models.StudyGroup{
						ID: 123,
						StudyGroupDetails: models.StudyGroupDetails{
							Name:        "Test Group",
							Description: "Test Description",
							Type:        models.TypePublic,
							ModuleID:    42,
						},
						Members: []models.StudyGroupMember{
							{
								UserID: 3,
								Role:   models.RoleAdmin,
							},
							{
								UserID: 4,
								Role:   models.RoleMember,
							},
						},
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":123,"name":"Test Group","description":"Test Description","type":"public"}` + "\n",
		},
		{
			name:         "Invalid ID format",
			url:          "/study-groups/invalid",
			mockSetup:    func(_ *MockStudyGroupService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid study group ID\n",
		},
		{
			name: "Valid ID but not found",
			url:  "/study-groups/999",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetStudyGroupByID", models.StudyGroupID(999)).
					Return((*models.StudyGroup)(nil), services.ErrStudyGroupNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "Study group not found\n",
		},
		{
			name: "Valid ID with internal service error",
			url:  "/study-groups/500",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetStudyGroupByID", models.StudyGroupID(500)).
					Return((*models.StudyGroup)(nil), errors.New("unexpected error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Error fetching study group\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := &MockStudyGroupService{}
			tt.mockSetup(mockService)

			handler := NewStudyGroupHandler(mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /study-groups/{id}", handler.GetStudyGroup)

			req := httptest.NewRequest(http.MethodGet, tt.url, bytes.NewReader([]byte{}))
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestStudyGroupHandler_GetAllStudyGroups(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockSetup    func(service *MockStudyGroupService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "Successfully fetch multiple study groups",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetAllStudyGroups").
					Return([]models.StudyGroup{
						{
							ID: 1,
							StudyGroupDetails: models.StudyGroupDetails{
								Name:        "Group 1",
								Description: "Description 1",
								Type:        models.TypePublic,
								ModuleID:    42,
							},
							Members: []models.StudyGroupMember{
								{
									UserID: 3,
									Role:   models.RoleAdmin,
								},
								{
									UserID: 4,
									Role:   models.RoleMember,
								},
							},
						},
						{
							ID: 2,
							StudyGroupDetails: models.StudyGroupDetails{
								Name:        "Group 2",
								Description: "Description 2",
								Type:        models.TypeClosed,
								ModuleID:    1,
							},
							Members: []models.StudyGroupMember{
								{
									UserID: 1,
									Role:   models.RoleAdmin,
								},
							},
						},
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"id":1,"name":"Group 1","description":"Description 1","type":"public"},{"id":2,"name":"Group 2","description":"Description 2","type":"closed"}]` + "\n",
		},
		{
			name: "Successfully fetch empty list of study groups",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetAllStudyGroups").
					Return([]models.StudyGroup{}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `[]` + "\n",
		},
		{
			name: "Error fetching study groups",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetAllStudyGroups").
					Return(([]models.StudyGroup)(nil), errors.New("internal error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Error fetching study groups\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := &MockStudyGroupService{}
			tt.mockSetup(mockService)

			handler := NewStudyGroupHandler(mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /study-groups", handler.GetAllStudyGroups)

			req := httptest.NewRequest(http.MethodGet, "/study-groups", bytes.NewReader([]byte{}))
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
