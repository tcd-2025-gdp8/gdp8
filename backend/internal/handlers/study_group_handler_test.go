package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/models"
	"gdp8-backend/internal/services"
)

type MockStudyGroupService struct {
	mock.Mock
}

func (m *MockStudyGroupService) GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroupView, error) {
	args := m.Called(id)
	return args.Get(0).(*models.StudyGroupView), args.Error(1)
}

func (m *MockStudyGroupService) GetAllStudyGroups() ([]models.StudyGroupView, error) {
	args := m.Called()
	return args.Get(0).([]models.StudyGroupView), args.Error(1)
}

func (m *MockStudyGroupService) GetAllRelevantStudyGroups(userID models.UserID) ([]models.StudyGroupView, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.StudyGroupView), args.Error(1)
}

func (m *MockStudyGroupService) CreateStudyGroup(studyGroupDetails models.StudyGroupDetails, creatorID models.UserID) (*models.StudyGroupView, error) {
	args := m.Called(studyGroupDetails, creatorID)
	return args.Get(0).(*models.StudyGroupView), args.Error(1)
}

func (m *MockStudyGroupService) UpdateStudyGroupDetails(id models.StudyGroupID, details models.StudyGroupDetails, requesterID models.UserID) (*models.StudyGroupView, error) {
	args := m.Called(id, details, requesterID)
	return args.Get(0).(*models.StudyGroupView), args.Error(1)
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

func (m *MockStudyGroupService) RetrieveGroupRole(studyGroupID models.StudyGroupID, userID models.UserID) (*models.StudyGroupRole, error) {
	args := m.Called(studyGroupID, userID)
	return args.Get(0).(*models.StudyGroupRole), args.Error(1)
}

func (m *MockStudyGroupService) IsGroupMember(studyGroupID models.StudyGroupID, userID models.UserID) (bool, error) {
	args := m.Called(studyGroupID, userID)
	return args.Get(0).(bool), args.Error(1)
}

func TestStudyGroupHandler_GetStudyGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		url          string
		mockSetup    func(service *MockStudyGroupService)
		ctxSetup     func(r *http.Request) *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name: "Valid ID, found, closed, and members are returned",
			url:  "/study-groups/123",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetStudyGroupByID", models.StudyGroupID(123)).
					Return(&models.StudyGroupView{
						ID: 123,
						StudyGroupDetails: models.StudyGroupDetails{
							Name:        "Test Group",
							Description: "Test Description",
							Type:        models.TypeClosed,
							ModuleID:    42,
							MaxMembers:  21,
						},
						Members: []models.StudyGroupMemberView{
							{
								UserID: "3",
								Name:   "Name 3",
								Role:   models.RoleAdmin,
							},
							{
								UserID: "4",
								Name:   "Name 4",
								Role:   models.RoleMember,
							},
						},
					}, nil)
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "4"))
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":123,"name":"Test Group","description":"Test Description","type":"closed","maxMembers":21,"moduleId":42,"members":[{"id":"3","name":"Name 3","role":"admin"},{"id":"4","name":"Name 4","role":"member"}]}` + "\n",
		},
		{
			name: "Valid ID, found, closed, and members are not returned",
			url:  "/study-groups/123",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetStudyGroupByID", models.StudyGroupID(123)).
					Return(&models.StudyGroupView{
						ID: 123,
						StudyGroupDetails: models.StudyGroupDetails{
							Name:        "Test Group",
							Description: "Test Description",
							Type:        models.TypeClosed,
							ModuleID:    42,
							MaxMembers:  7,
						},
						Members: []models.StudyGroupMemberView{
							{
								UserID: "3",
								Name:   "Name 3",
								Role:   models.RoleAdmin,
							},
							{
								UserID: "4",
								Name:   "Name 4",
								Role:   models.RoleMember,
							},
						},
					}, nil)
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "7"))
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":123,"name":"Test Group","description":"Test Description","type":"closed","maxMembers":7,"moduleId":42}` + "\n",
		},
		{
			name: "Valid ID, found, invite-only, not a member",
			url:  "/study-groups/123",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetStudyGroupByID", models.StudyGroupID(123)).
					Return(&models.StudyGroupView{
						ID: 123,
						StudyGroupDetails: models.StudyGroupDetails{
							Name:        "Test Group",
							Description: "Test Description",
							Type:        models.TypeInviteOnly,
							ModuleID:    42,
							MaxMembers:  5,
						},
						Members: []models.StudyGroupMemberView{
							{
								UserID: "3",
								Name:   "Name 3",
								Role:   models.RoleAdmin,
							},
							{
								UserID: "4",
								Name:   "Name 4",
								Role:   models.RoleMember,
							},
						},
					}, nil)
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "7"))
			},
			expectedCode: http.StatusForbidden,
			expectedBody: "Forbidden\n",
		},
		{
			name: "Valid ID, found, invite-only, member",
			url:  "/study-groups/123",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetStudyGroupByID", models.StudyGroupID(123)).
					Return(&models.StudyGroupView{
						ID: 123,
						StudyGroupDetails: models.StudyGroupDetails{
							Name:        "Test Group",
							Description: "Test Description",
							Type:        models.TypeInviteOnly,
							ModuleID:    42,
							MaxMembers:  9,
						},
						Members: []models.StudyGroupMemberView{
							{
								UserID: "3",
								Name:   "Name 3",
								Role:   models.RoleAdmin,
							},
							{
								UserID: "4",
								Name:   "Name 4",
								Role:   models.RoleMember,
							},
						},
					}, nil)
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "4"))
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":123,"name":"Test Group","description":"Test Description","type":"invite-only","maxMembers":9,"moduleId":42,"members":[{"id":"3","name":"Name 3","role":"admin"},{"id":"4","name":"Name 4","role":"member"}]}` + "\n",
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
					Return((*models.StudyGroupView)(nil), services.ErrStudyGroupNotFound)
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "1"))
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
					Return((*models.StudyGroupView)(nil), errors.New("unexpected error"))
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "1"))
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
			if tt.ctxSetup != nil {
				req = tt.ctxSetup(req)
			}
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestStudyGroupHandler_GetRelevantStudyGroups(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockSetup    func(service *MockStudyGroupService)
		ctxSetup     func(r *http.Request) *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name: "Successfully fetch multiple study groups",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetAllRelevantStudyGroups", models.UserID("2")).
					Return([]models.StudyGroupView{
						{
							ID: 1,
							StudyGroupDetails: models.StudyGroupDetails{
								Name:        "Group 1",
								Description: "Description 1",
								Type:        models.TypePublic,
								ModuleID:    42,
								MaxMembers:  5,
							},
							Members: []models.StudyGroupMemberView{
								{
									UserID: "3",
									Name:   "Name 3",
									Role:   models.RoleAdmin,
								},
								{
									UserID: "4",
									Name:   "Name 4",
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
								MaxMembers:  3,
							},
							Members: []models.StudyGroupMemberView{
								{
									UserID: "1",
									Name:   "Name 1",
									Role:   models.RoleAdmin,
								},
							},
						},
						{
							ID: 3,
							StudyGroupDetails: models.StudyGroupDetails{
								Name:        "Group 3",
								Description: "Description 3",
								Type:        models.TypeClosed,
								ModuleID:    1,
								MaxMembers:  4,
							},
							Members: []models.StudyGroupMemberView{
								{
									UserID: "2",
									Name:   "Name 2",
									Role:   models.RoleAdmin,
								},
							},
						},
						{
							ID: 4,
							StudyGroupDetails: models.StudyGroupDetails{
								Name:        "Group 4",
								Description: "Description 4",
								Type:        models.TypeInviteOnly,
								ModuleID:    1,
								MaxMembers:  10,
							},
							Members: []models.StudyGroupMemberView{
								{
									UserID: "1",
									Name:   "Name 1",
									Role:   models.RoleAdmin,
								},
							},
						},
						{
							ID: 5,
							StudyGroupDetails: models.StudyGroupDetails{
								Name:        "Group 5",
								Description: "Description 5",
								Type:        models.TypeInviteOnly,
								ModuleID:    1,
								MaxMembers:  15,
							},
							Members: []models.StudyGroupMemberView{
								{
									UserID: "2",
									Name:   "Name 2",
									Role:   models.RoleAdmin,
								},
							},
						},
					}, nil)
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "2"))
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"id":1,"name":"Group 1","description":"Description 1","type":"public","maxMembers":5,"moduleId":42,"members":[{"id":"3","name":"Name 3","role":"admin"},{"id":"4","name":"Name 4","role":"member"}]},{"id":2,"name":"Group 2","description":"Description 2","type":"closed","maxMembers":3,"moduleId":1},{"id":3,"name":"Group 3","description":"Description 3","type":"closed","maxMembers":4,"moduleId":1,"members":[{"id":"2","name":"Name 2","role":"admin"}]},{"id":5,"name":"Group 5","description":"Description 5","type":"invite-only","maxMembers":15,"moduleId":1,"members":[{"id":"2","name":"Name 2","role":"admin"}]}]` + "\n",
		},
		{
			name: "Successfully fetch empty list of study groups",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetAllRelevantStudyGroups", models.UserID("1")).
					Return([]models.StudyGroupView{}, nil)
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "1"))
			},
			expectedCode: http.StatusOK,
			expectedBody: `[]` + "\n",
		},
		{
			name: "Error fetching study groups",
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("GetAllRelevantStudyGroups", models.UserID("1")).
					Return(([]models.StudyGroupView)(nil), errors.New("internal error"))
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "1"))
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
			mux.HandleFunc("GET /study-groups", handler.GetRelevantStudyGroups)

			req := httptest.NewRequest(http.MethodGet, "/study-groups", bytes.NewReader([]byte{}))
			if tt.ctxSetup != nil {
				req = tt.ctxSetup(req)
			}
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestStudyGroupHandler_CreateStudyGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		body         string
		mockSetup    func(service *MockStudyGroupService)
		ctxSetup     func(r *http.Request) *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name: "Valid creation",
			body: `{"name":"Group A","description":"Desc A","type":"public", "moduleId":5, "maxMembers":7}`,
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("CreateStudyGroup", models.StudyGroupDetails{
						Name:        "Group A",
						Description: "Desc A",
						Type:        models.TypePublic,
						ModuleID:    5,
						MaxMembers:  7,
					}, models.UserID("123")).
					Return(&models.StudyGroupView{
						ID: 1,
						StudyGroupDetails: models.StudyGroupDetails{
							Name:        "Group A",
							Description: "Desc A",
							Type:        models.TypePublic,
							ModuleID:    5,
							MaxMembers:  7,
						},
					}, nil)
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "123"))
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"name":"Group A","description":"Desc A","type":"public","maxMembers":7,"moduleId":5}` + "\n",
		},
		{
			name:         "Missing user context",
			body:         `{"name":"Group B","description":"Desc B","type":"closed","maxMembers":7,"moduleId":5}`,
			mockSetup:    func(_ *MockStudyGroupService) {},
			ctxSetup:     func(r *http.Request) *http.Request { return r },
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
		{
			name:      "Invalid input payload",
			body:      `invalid-json`,
			mockSetup: func(_ *MockStudyGroupService) {},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "123"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request payload\n",
		},
		{
			name: "Service error during creation",
			body: `{"name":"Group C","description":"Desc C","type":"closed", "moduleId":5, "maxMembers":7}`,
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("CreateStudyGroup", models.StudyGroupDetails{
						Name:        "Group C",
						Description: "Desc C",
						Type:        models.TypeClosed,
						ModuleID:    5,
						MaxMembers:  7,
					}, models.UserID("123")).
					Return((*models.StudyGroupView)(nil), errors.New("service error"))
			},
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "123"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Error creating study group\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := &MockStudyGroupService{}
			tt.mockSetup(mockService)

			handler := NewStudyGroupHandler(mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /study-groups", handler.CreateStudyGroup)

			req := httptest.NewRequest(http.MethodPost, "/study-groups", bytes.NewReader([]byte(tt.body)))
			if tt.ctxSetup != nil {
				req = tt.ctxSetup(req)
			}
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestStudyGroupHandler_HandleStudyMemberOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		url          string
		command      string
		body         string
		ctxSetup     func(r *http.Request) *http.Request
		mockSetup    func(service *MockStudyGroupService)
		expectedCode int
		expectedBody string
	}{
		{
			name:    "Accept invite successfully",
			url:     "/study-groups/123/accept-invite",
			command: "accept-invite",
			body:    ``,
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "123"))
			},
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("HandleSelfMemberOperation", services.AcceptStudyGroupInviteCommand, models.StudyGroupID(123), models.UserID("123")).
					Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:    "Request to join study group unauthorized",
			url:     "/study-groups/123/request-to-join",
			command: "request-to-join",
			body:    ``,
			ctxSetup: func(r *http.Request) *http.Request {
				return r
			},
			mockSetup:    func(_ *MockStudyGroupService) {},
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Unauthorized\n",
		},
		{
			name:    "Invite member with missing payload",
			url:     "/study-groups/123/invite",
			command: "invite",
			body:    `{}`,
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "123"))
			},
			mockSetup:    func(_ *MockStudyGroupService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request payload\n",
		},
		{
			name:    "Invite member",
			url:     "/study-groups/123/invite",
			command: "invite",
			body:    `{"targetUserId":"456"}`,
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "123"))
			},
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("HandleAdminMemberOperation", services.InviteMemberToStudyGroupCommand, models.StudyGroupID(123), models.UserID("456"), models.UserID("123")).
					Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:    "Remove member study group not found",
			url:     "/study-groups/123/remove-member",
			command: "remove-member",
			body:    `{"targetUserId":"456"}`,
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "123"))
			},
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("HandleAdminMemberOperation", services.RemoveMemberFromStudyGroupCommand, models.StudyGroupID(123), models.UserID("456"), models.UserID("123")).
					Return(services.ErrStudyGroupNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "Study group not found\n",
		},

		{
			name:    "Reject request invalid operation",
			url:     "/study-groups/123/reject-request-to-join",
			command: "reject-request-to-join",
			body:    `{"targetUserId":"456"}`,
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "123"))
			},
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("HandleAdminMemberOperation", services.RejectRequestToJoinStudyGroupCommand, models.StudyGroupID(123), models.UserID("456"), models.UserID("123")).
					Return(services.ErrInvalidMemberOperation)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid study group operation\n",
		},
		{
			name:    "Reject request forbidden (unauthorized member operation)",
			url:     "/study-groups/123/reject-request-to-join",
			command: "reject-request-to-join",
			body:    `{"targetUserId":"456"}`,
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "123"))
			},
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("HandleAdminMemberOperation", services.RejectRequestToJoinStudyGroupCommand, models.StudyGroupID(123), models.UserID("456"), models.UserID("123")).
					Return(services.ErrUnauthorizedMemberOperation)
			},
			expectedCode: http.StatusForbidden,
			expectedBody: "Unauthorized study group operation\n",
		},
		{
			name:    "Leave study group internal error",
			url:     "/study-groups/123/leave",
			command: "leave",
			body:    ``,
			ctxSetup: func(r *http.Request) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), middleware.UIDCtxKey{}, "123"))
			},
			mockSetup: func(service *MockStudyGroupService) {
				service.
					On("HandleSelfMemberOperation", services.LeaveStudyGroupCommand, models.StudyGroupID(123), models.UserID("123")).
					Return(errors.New("unexpected error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Error processing study group operation\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := &MockStudyGroupService{}
			tt.mockSetup(mockService)

			handler := NewStudyGroupHandler(mockService)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /study-groups/{id}/{command}", handler.HandleStudyMemberOperation)

			req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewReader([]byte(tt.body)))
			if tt.ctxSetup != nil {
				req = tt.ctxSetup(req)
			}
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
