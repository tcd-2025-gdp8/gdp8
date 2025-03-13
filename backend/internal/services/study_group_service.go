package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type AdminMemberOperationCommand string
type SelfMemberOperationCommand string

const (
	InviteMemberToStudyGroupCommand      AdminMemberOperationCommand = "inviteMemberToStudyGroup"
	AcceptRequestToJoinStudyGroupCommand AdminMemberOperationCommand = "acceptRequestToJoinStudyGroup"
	RejectRequestToJoinStudyGroupCommand AdminMemberOperationCommand = "rejectRequestToJoinStudyGroup"
	RemoveMemberFromStudyGroupCommand    AdminMemberOperationCommand = "removeMemberFromStudyGroup"

	AcceptStudyGroupInviteCommand  SelfMemberOperationCommand = "acceptInviteToStudyGroup"
	RejectStudyGroupInviteCommand  SelfMemberOperationCommand = "rejectInviteToStudyGroup"
	RequestToJoinStudyGroupCommand SelfMemberOperationCommand = "requestToJoinStudyGroup"
	LeaveStudyGroupCommand         SelfMemberOperationCommand = "leaveStudyGroup"
)

type StudyGroupService interface {
	GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroupView, error)
	GetAllStudyGroups() ([]models.StudyGroupView, error)
	GetAllRelevantStudyGroups(userID models.UserID) ([]models.StudyGroupView, error)

	CreateStudyGroup(
		studyGroupDetails models.StudyGroupDetails,
		creatorID models.UserID) (*models.StudyGroupView, error)

	UpdateStudyGroupDetails(
		id models.StudyGroupID,
		details models.StudyGroupDetails,
		requesterID models.UserID) (*models.StudyGroupView, error)

	DeleteStudyGroup(id models.StudyGroupID, requesterID models.UserID) error

	HandleAdminMemberOperation(
		command AdminMemberOperationCommand,
		studyGroupID models.StudyGroupID,
		targetUserID models.UserID,
		adminID models.UserID) error

	HandleSelfMemberOperation(
		command SelfMemberOperationCommand,
		studyGroupID models.StudyGroupID,
		memberID models.UserID) error

	RetrieveGroupRole(studyGroupID models.StudyGroupID, userID models.UserID) (*models.StudyGroupRole, error)
}

var ErrStudyGroupNotFound = errors.New("study group not found")
var ErrInvalidMemberOperation = errors.New("invalid study group member operation")
var ErrUnauthorizedMemberOperation = errors.New("unauthorized member operation")
var ErrStudyGroupFull = errors.New("study group is full")

type studyGroupServiceImpl struct {
	txMgr               persistence.TransactionManager
	studyGroupRepo      repositories.StudyGroupRepository
	notificationService NotificationService
}

func NewStudyGroupService(
	txMgr persistence.TransactionManager,
	studyGroupRepo repositories.StudyGroupRepository,
	notificationService NotificationService) StudyGroupService {

	return &studyGroupServiceImpl{
		txMgr:               txMgr,
		studyGroupRepo:      studyGroupRepo,
		notificationService: notificationService,
	}
}

func (s *studyGroupServiceImpl) GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroupView, error) {
	grp, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudyGroupView, error) {
		return s.studyGroupRepo.GetStudyGroupByID(tx, id)
	})

	err = resolveError(err, "fetching study group")
	if err != nil {
		return nil, err
	}

	return grp, nil
}

func (s *studyGroupServiceImpl) GetAllStudyGroups() ([]models.StudyGroupView, error) {
	grp, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.StudyGroupView, error) {
		return s.studyGroupRepo.GetAllStudyGroups(tx)
	})

	err = resolveError(err, "fetching all study groups")
	if err != nil {
		return nil, err
	}

	return grp, nil
}

// GetAllRelevantStudyGroups retrieves a list of study groups relevant to the specified user based on module choices.
func (s *studyGroupServiceImpl) GetAllRelevantStudyGroups(userID models.UserID) ([]models.StudyGroupView, error) {
	grp, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) ([]models.StudyGroupView, error) {
		return s.studyGroupRepo.GetAllRelevantStudyGroups(tx, userID)
	})

	err = resolveError(err, "fetching all relevant study groups")
	if err != nil {
		return nil, err
	}

	return grp, nil
}

func (s *studyGroupServiceImpl) CreateStudyGroup(studyGroupDetails models.StudyGroupDetails,
	creatorID models.UserID) (*models.StudyGroupView, error) {
	grp, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudyGroupView, error) {
		// TODO check if creator exists in the users repo
		return s.studyGroupRepo.CreateStudyGroup(tx, studyGroupDetails, creatorID)
	})

	err = resolveError(err, "creating study group")
	if err != nil {
		return nil, err
	}

	return grp, nil
}

func (s *studyGroupServiceImpl) UpdateStudyGroupDetails(id models.StudyGroupID,
	details models.StudyGroupDetails, requesterID models.UserID) (*models.StudyGroupView, error) {
	grp, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudyGroupView, error) {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, id)
		if err != nil {
			return nil, err
		}

		if !hasRole(requesterID, models.RoleAdmin, studyGroup.Members) {
			return nil, ErrUnauthorizedMemberOperation
		}

		return s.studyGroupRepo.UpdateStudyGroupDetails(tx, id, details)
	})

	err = resolveError(err, "updating study group details")
	if err != nil {
		return nil, err
	}

	return grp, nil
}

func (s *studyGroupServiceImpl) DeleteStudyGroup(id models.StudyGroupID, requesterID models.UserID) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx *sql.Tx) error {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, id)
		if err != nil {
			return err
		}

		if !hasRole(requesterID, models.RoleAdmin, studyGroup.Members) {
			return ErrUnauthorizedMemberOperation
		}

		return s.studyGroupRepo.DeleteStudyGroup(tx, id)
	})

	err = resolveError(err, "deleting study group")

	return err
}

func (s *studyGroupServiceImpl) HandleAdminMemberOperation(command AdminMemberOperationCommand,
	studyGroupID models.StudyGroupID, targetUserID models.UserID, adminID models.UserID) error {
	studyGroup, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudyGroupView, error) {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, studyGroupID)
		if err != nil {
			return nil, err
		}

		if !hasRole(adminID, models.RoleAdmin, studyGroup.Members) {
			return nil, ErrUnauthorizedMemberOperation
		}

		switch command {
		case InviteMemberToStudyGroupCommand:
			err = s.inviteMember(tx, studyGroup, targetUserID)
		case AcceptRequestToJoinStudyGroupCommand:
			err = s.acceptRequestToJoin(tx, studyGroup, targetUserID)
		case RejectRequestToJoinStudyGroupCommand:
			err = s.rejectRequestToJoin(tx, studyGroup, targetUserID)
		case RemoveMemberFromStudyGroupCommand:
			err = s.removeMemberFromStudyGroup(tx, studyGroup, targetUserID, adminID)
		default:
			log.Printf("[ERROR] invalid admin member operation command: %s\n", command)
			return nil, errors.New("invalid admin member operation command")
		}

		if err != nil {
			return nil, err
		}

		return studyGroup, nil
	})

	err = resolveError(err, fmt.Sprintf("executing admin member operation %s", command))

	if err == nil {
		var notificationType models.NotificationType
		switch command {
		case InviteMemberToStudyGroupCommand:
			notificationType = models.NotificationTypeStudyGroupInvited
		case AcceptRequestToJoinStudyGroupCommand:
			notificationType = models.NotificationTypeStudyGroupAcceptedJoinRequest
		case RejectRequestToJoinStudyGroupCommand:
			notificationType = models.NotificationTypeStudyGroupRejectedJoinRequest
		case RemoveMemberFromStudyGroupCommand:
			notificationType = models.NotificationTypeStudyGroupRemovedMember
		}
		go func() {
			notificationErr := s.notificationService.AddStudyGroupEventNotification(
				notificationType, adminID, &targetUserID, studyGroupID, studyGroup.Members)
			if notificationErr != nil {
				log.Printf("Error sending notification: %v\n", notificationErr)
			}
		}()
	}

	return err
}

func (s *studyGroupServiceImpl) HandleSelfMemberOperation(command SelfMemberOperationCommand,
	studyGroupID models.StudyGroupID, memberID models.UserID) error {
	studyGroup, err := persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudyGroupView, error) {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, studyGroupID)
		if err != nil {
			return nil, err
		}

		switch command {
		case AcceptStudyGroupInviteCommand:
			err = s.acceptStudyGroupInvite(tx, studyGroup, memberID)
		case RejectStudyGroupInviteCommand:
			err = s.rejectStudyGroupInvite(tx, studyGroup, memberID)
		case RequestToJoinStudyGroupCommand:
			err = s.requestToJoinStudyGroup(tx, studyGroup, memberID)
		case LeaveStudyGroupCommand:
			err = s.leaveStudyGroup(tx, studyGroup, memberID)
		default:
			log.Printf("[ERROR] invalid member operation command: %s\n", command)
			return nil, errors.New("invalid member operation command")
		}

		if err != nil {
			return nil, err
		}

		return studyGroup, nil
	})

	err = resolveError(err, fmt.Sprintf("executing member operation %s", command))

	if err == nil {
		var notificationType models.NotificationType
		switch command {
		case AcceptStudyGroupInviteCommand:
			notificationType = models.NotificationTypeStudyGroupAcceptedInvite
		case RejectStudyGroupInviteCommand:
			notificationType = models.NotificationTypeStudyGroupRejectedInvite
		case RequestToJoinStudyGroupCommand:
			if studyGroup.Type == models.TypePublic {
				notificationType = models.NotificationTypeStudyGroupJoined
			} else {
				notificationType = models.NotificationTypeStudyGroupRequestedToJoin
			}
		case LeaveStudyGroupCommand:
			notificationType = models.NotificationTypeStudyGroupLeft
		}
		go func() {
			notificationErr := s.notificationService.AddStudyGroupEventNotification(
				notificationType, memberID, nil, studyGroupID, studyGroup.Members)
			if notificationErr != nil {
				log.Printf("Error sending notification: %v\n", notificationErr)
			}
		}()
	}

	return err
}

func (s *studyGroupServiceImpl) RetrieveGroupRole(studyGroupID models.StudyGroupID,
	userID models.UserID) (*models.StudyGroupRole, error) {

	return persistence.WithTransaction(s.txMgr, func(tx *sql.Tx) (*models.StudyGroupRole, error) {
		return s.studyGroupRepo.RetrieveGroupRole(tx, studyGroupID, userID)
	})
}

func resolveError(err error, operation string) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ErrStudyGroupFull):
		return err
	case errors.Is(err, ErrInvalidMemberOperation):
		return err
	case errors.Is(err, ErrStudyGroupNotFound):
		return err
	case errors.Is(err, ErrUnauthorizedMemberOperation):
		return err
	case errors.Is(err, repositories.ErrStudyGroupNotFound):
		return ErrStudyGroupNotFound
	default:
		return fmt.Errorf("error %s: %w", operation, err)
	}
}

func (s *studyGroupServiceImpl) inviteMember(tx *sql.Tx,
	studyGroup *models.StudyGroupView, memberID models.UserID) error {
	// TODO check if member exists in the users repo

	if studyGroup.StudyGroupDetails.MaxMembers > 0 &&
		len(studyGroup.Members) >= studyGroup.StudyGroupDetails.MaxMembers {
		return fmt.Errorf("%w: study group is full", ErrStudyGroupFull)
	}

	for _, member := range studyGroup.Members {
		if member.UserID == memberID {
			return fmt.Errorf("%w: member already exists in the study group", ErrInvalidMemberOperation)
		}
	}

	role := models.RoleInvitee
	return s.studyGroupRepo.UpdateStudyGroupMember(tx, studyGroup.ID, memberID, &role)
}

func (s *studyGroupServiceImpl) acceptRequestToJoin(tx *sql.Tx,
	studyGroup *models.StudyGroupView, memberID models.UserID) error {

	if studyGroup.StudyGroupDetails.MaxMembers > 0 &&
		len(studyGroup.Members) >= studyGroup.StudyGroupDetails.MaxMembers {
		return fmt.Errorf("%w: study group is full", ErrStudyGroupFull)
	}
	if !hasRole(memberID, models.RoleRequester, studyGroup.Members) {
		return fmt.Errorf("%w: member hasn't requested to join the study group", ErrInvalidMemberOperation)
	}

	role := models.RoleMember
	return s.studyGroupRepo.UpdateStudyGroupMember(tx, studyGroup.ID, memberID, &role)
}

func (s *studyGroupServiceImpl) rejectRequestToJoin(tx *sql.Tx,
	studyGroup *models.StudyGroupView, memberID models.UserID) error {
	if !hasRole(memberID, models.RoleRequester, studyGroup.Members) {
		return fmt.Errorf("%w: member hasn't requested to join the study group", ErrInvalidMemberOperation)
	}

	return s.studyGroupRepo.UpdateStudyGroupMember(tx, studyGroup.ID, memberID, nil)
}

func (s *studyGroupServiceImpl) removeMemberFromStudyGroup(tx *sql.Tx,
	studyGroup *models.StudyGroupView, memberID models.UserID, adminID models.UserID) error {
	if memberID == adminID {
		return fmt.Errorf("%w: cannot remove self from the study group", ErrInvalidMemberOperation)
	}

	return s.studyGroupRepo.UpdateStudyGroupMember(tx, studyGroup.ID, memberID, nil)
}

func (s *studyGroupServiceImpl) acceptStudyGroupInvite(tx *sql.Tx,
	studyGroup *models.StudyGroupView, memberID models.UserID) error {

	if studyGroup.StudyGroupDetails.MaxMembers > 0 &&
		len(studyGroup.Members) >= studyGroup.StudyGroupDetails.MaxMembers {
		return fmt.Errorf("%w: study group is full", ErrStudyGroupFull)
	}
	if !hasRole(memberID, models.RoleInvitee, studyGroup.Members) {
		return fmt.Errorf("%w: member not invited to join the study group", ErrInvalidMemberOperation)
	}

	role := models.RoleMember
	return s.studyGroupRepo.UpdateStudyGroupMember(tx, studyGroup.ID, memberID, &role)
}

func (s *studyGroupServiceImpl) rejectStudyGroupInvite(tx *sql.Tx,
	studyGroup *models.StudyGroupView, memberID models.UserID) error {
	if !hasRole(memberID, models.RoleInvitee, studyGroup.Members) {
		return fmt.Errorf("%w: member not invited to join the study group", ErrInvalidMemberOperation)
	}

	return s.studyGroupRepo.UpdateStudyGroupMember(tx, studyGroup.ID, memberID, nil)
}

func (s *studyGroupServiceImpl) requestToJoinStudyGroup(tx *sql.Tx,
	studyGroup *models.StudyGroupView, memberID models.UserID) error {
	// TODO check if member exists in the users repo

	if studyGroup.StudyGroupDetails.MaxMembers > 0 &&
		len(studyGroup.Members) >= studyGroup.StudyGroupDetails.MaxMembers {
		return fmt.Errorf("%w: study group is full", ErrStudyGroupFull)
	}
	for _, member := range studyGroup.Members {
		if member.UserID == memberID {
			return fmt.Errorf("%w: member already exists in the study group", ErrInvalidMemberOperation)
		}
	}

	switch studyGroup.Type {
	case models.TypePublic:
		role := models.RoleMember
		return s.studyGroupRepo.UpdateStudyGroupMember(tx, studyGroup.ID, memberID, &role)
	case models.TypeClosed:
		role := models.RoleRequester
		return s.studyGroupRepo.UpdateStudyGroupMember(tx, studyGroup.ID, memberID, &role)
	case models.TypeInviteOnly:
		return fmt.Errorf("%w: the study group is invite-only", ErrInvalidMemberOperation)
	default:
		log.Printf("[ERROR] invalid study group type: %s\n", studyGroup.Type)
		return errors.New("invalid study group type")
	}
}

func (s *studyGroupServiceImpl) leaveStudyGroup(tx *sql.Tx,
	studyGroup *models.StudyGroupView, memberID models.UserID) error {
	return s.studyGroupRepo.UpdateStudyGroupMember(tx, studyGroup.ID, memberID, nil)
}

func hasRole(userID models.UserID, role models.StudyGroupRole, members []models.StudyGroupMemberView) bool {
	return slices.ContainsFunc(members, func(m models.StudyGroupMemberView) bool {
		return m.UserID == userID && m.Role == role
	})
}
