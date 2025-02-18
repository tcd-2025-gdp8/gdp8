package services

import (
	"errors"
	"fmt"
	"slices"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
)

type StudyGroupService interface {
	GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroup, error)
	GetAllStudyGroups() ([]models.StudyGroup, error)
	CreateStudyGroup(studyGroupDetails models.StudyGroupDetails, creatorID models.UserID) (*models.StudyGroup, error)
	UpdateStudyGroupDetails(id models.StudyGroupID, details models.StudyGroupDetails, requesterID models.UserID) (*models.StudyGroup, error)
	DeleteStudyGroup(id models.StudyGroupID, requesterID models.UserID) error

	InviteMemberToStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID, requesterID models.UserID) error
	AcceptInviteToStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID) error
	RejectInviteToStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID) error
	RequestToJoinStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID) error
	AcceptRequestToJoinStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID, adminID models.UserID) error
	RejectRequestToJoinStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID, adminID models.UserID) error
	RemoveMemberFromStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID, requesterID models.UserID) error
	LeaveStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID) error
}

var ErrStudyGroupNotFound = errors.New("study group not found")
var ErrInvalidMemberOperation = errors.New("invalid study group member operation")
var ErrUnauthorizedMemberOperation = errors.New("unauthorized member operation")

type studyGroupServiceImpl struct {
	txMgr          persistence.TransactionManager
	studyGroupRepo repositories.StudyGroupRepository
}

func NewStudyGroupService(
	txMgr persistence.TransactionManager,
	studyGroupRepo repositories.StudyGroupRepository) StudyGroupService {

	return &studyGroupServiceImpl{
		txMgr:          txMgr,
		studyGroupRepo: studyGroupRepo,
	}
}

func (s *studyGroupServiceImpl) GetStudyGroupByID(id models.StudyGroupID) (*models.StudyGroup, error) {
	studyGroup, err := persistence.WithTransaction(s.txMgr, func(tx persistence.Transaction) (*models.StudyGroup, error) {
		return s.studyGroupRepo.GetStudyGroupByID(tx, id)
	})

	err = resolveError(err, "fetching study group")
	if err != nil {
		return nil, err
	}

	return studyGroup, nil
}

func (s *studyGroupServiceImpl) GetAllStudyGroups() ([]models.StudyGroup, error) {
	studyGroups, err := persistence.WithTransaction(s.txMgr, func(tx persistence.Transaction) ([]models.StudyGroup, error) {
		return s.studyGroupRepo.GetAllStudyGroups(tx)
	})

	err = resolveError(err, "fetching all study groups")
	if err != nil {
		return nil, err
	}

	return studyGroups, nil
}

func (s *studyGroupServiceImpl) CreateStudyGroup(studyGroupDetails models.StudyGroupDetails, creatorID models.UserID) (*models.StudyGroup, error) {
	studyGroup, err := persistence.WithTransaction(s.txMgr, func(tx persistence.Transaction) (*models.StudyGroup, error) {
		// TODO check if creator exists in the users repo
		return s.studyGroupRepo.CreateStudyGroup(tx, studyGroupDetails, creatorID)
	})

	err = resolveError(err, "creating study group")
	if err != nil {
		return nil, err
	}

	return studyGroup, nil
}

func (s *studyGroupServiceImpl) UpdateStudyGroupDetails(id models.StudyGroupID, details models.StudyGroupDetails, requesterID models.UserID) (*models.StudyGroup, error) {
	studyGroup, err := persistence.WithTransaction(s.txMgr, func(tx persistence.Transaction) (*models.StudyGroup, error) {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, id)
		if err != nil {
			return nil, err
		}

		if !hasRole(requesterID, models.RoleAdmin, studyGroup.Members) {
			return nil, ErrUnauthorizedMemberOperation
		}

		studyGroup.StudyGroupDetails = details

		return s.studyGroupRepo.UpdateStudyGroup(tx, studyGroup)
	})

	err = resolveError(err, "updating study group details")
	if err != nil {
		return nil, err
	}

	return studyGroup, nil
}

func (s *studyGroupServiceImpl) DeleteStudyGroup(id models.StudyGroupID, requesterID models.UserID) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx persistence.Transaction) error {
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

func (s *studyGroupServiceImpl) InviteMemberToStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID, requesterID models.UserID) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx persistence.Transaction) error {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, studyGroupID)
		if err != nil {
			return err
		}

		if !hasRole(requesterID, models.RoleAdmin, studyGroup.Members) {
			return ErrUnauthorizedMemberOperation
		}

		// TODO check if member exists in the users repo

		for _, member := range studyGroup.Members {
			if member.UserID == memberID {
				return fmt.Errorf("%w: member already exists in the study group", ErrInvalidMemberOperation)
			}
		}

		studyGroup.Members = append(studyGroup.Members, models.StudyGroupMember{
			UserID: memberID,
			Role:   models.RoleInvitee,
		})

		_, err = s.studyGroupRepo.UpdateStudyGroup(tx, studyGroup)
		return err
	})

	err = resolveError(err, "inviting member to the study group")

	// TODO send a notification

	return err
}

func (s *studyGroupServiceImpl) RequestToJoinStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx persistence.Transaction) error {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, studyGroupID)
		if err != nil {
			return err
		}

		// TODO check if member exists in the users repo

		for _, member := range studyGroup.Members {
			if member.UserID == memberID {
				return fmt.Errorf("%w: member already exists in the study group", ErrInvalidMemberOperation)
			}
		}

		switch studyGroup.Type {
		case models.TypePublic:
			studyGroup.Members = append(studyGroup.Members, models.StudyGroupMember{
				UserID: memberID,
				Role:   models.RoleMember,
			})
		case models.TypeClosed:
			studyGroup.Members = append(studyGroup.Members, models.StudyGroupMember{
				UserID: memberID,
				Role:   models.RoleRequester,
			})
		case models.TypeInviteOnly:
			return fmt.Errorf("%w: the study group is invite-only", ErrInvalidMemberOperation)
		}

		_, err = s.studyGroupRepo.UpdateStudyGroup(tx, studyGroup)
		return err
	})

	err = resolveError(err, "requesting to join the study group")

	// TODO send a notification

	return err
}

func (s *studyGroupServiceImpl) AcceptInviteToStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx persistence.Transaction) error {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, studyGroupID)
		if err != nil {
			return err
		}

		if !hasRole(memberID, models.RoleInvitee, studyGroup.Members) {
			return fmt.Errorf("%w: member not invited to join the study group", ErrInvalidMemberOperation)
		}

		studyGroup.Members = setRole(memberID, models.RoleMember, studyGroup.Members)

		_, err = s.studyGroupRepo.UpdateStudyGroup(tx, studyGroup)
		return err
	})

	err = resolveError(err, "accepting invite to join the study group")

	// TODO send a notification (?)

	return err
}

func (s *studyGroupServiceImpl) RejectInviteToStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx persistence.Transaction) error {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, studyGroupID)
		if err != nil {
			return err
		}

		if !hasRole(memberID, models.RoleInvitee, studyGroup.Members) {
			return fmt.Errorf("%w: member not invited to join the study group", ErrInvalidMemberOperation)
		}

		studyGroup.Members = slices.DeleteFunc(studyGroup.Members, func(m models.StudyGroupMember) bool {
			return m.UserID == memberID
		})

		_, err = s.studyGroupRepo.UpdateStudyGroup(tx, studyGroup)
		return err
	})

	err = resolveError(err, "rejecting invite to join the study group")

	// TODO send a notification (?)

	return err
}

func (s *studyGroupServiceImpl) AcceptRequestToJoinStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID, adminID models.UserID) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx persistence.Transaction) error {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, studyGroupID)
		if err != nil {
			return err
		}

		if !hasRole(adminID, models.RoleAdmin, studyGroup.Members) {
			return ErrUnauthorizedMemberOperation
		}

		if !hasRole(memberID, models.RoleRequester, studyGroup.Members) {
			return fmt.Errorf("%w: member hasn't requested to join the study group", ErrInvalidMemberOperation)
		}

		studyGroup.Members = setRole(memberID, models.RoleMember, studyGroup.Members)

		_, err = s.studyGroupRepo.UpdateStudyGroup(tx, studyGroup)
		return err
	})

	err = resolveError(err, "accepting request to join the study group")

	// TODO send a notification (?)

	return err
}

func (s *studyGroupServiceImpl) RejectRequestToJoinStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID, adminID models.UserID) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx persistence.Transaction) error {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, studyGroupID)
		if err != nil {
			return err
		}

		if !hasRole(adminID, models.RoleAdmin, studyGroup.Members) {
			return ErrUnauthorizedMemberOperation
		}

		if !hasRole(memberID, models.RoleRequester, studyGroup.Members) {
			return fmt.Errorf("%w: member hasn't requested to join the study group", ErrInvalidMemberOperation)
		}

		studyGroup.Members = slices.DeleteFunc(studyGroup.Members, func(m models.StudyGroupMember) bool {
			return m.UserID == memberID
		})

		_, err = s.studyGroupRepo.UpdateStudyGroup(tx, studyGroup)
		return err
	})

	err = resolveError(err, "rejecting request to join the study group")

	// TODO send a notification (?)

	return err
}

func (s *studyGroupServiceImpl) RemoveMemberFromStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID, requesterID models.UserID) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx persistence.Transaction) error {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, studyGroupID)
		if err != nil {
			return err
		}

		if !hasRole(requesterID, models.RoleAdmin, studyGroup.Members) {
			return ErrUnauthorizedMemberOperation
		}

		if memberID == requesterID {
			return fmt.Errorf("%w: cannot remove self from the study group", ErrInvalidMemberOperation)
		}

		for i, member := range studyGroup.Members {
			if member.UserID == memberID {
				studyGroup.Members = append(studyGroup.Members[:i], studyGroup.Members[i+1:]...)
				_, err = s.studyGroupRepo.UpdateStudyGroup(tx, studyGroup)
				return err
			}
		}

		return fmt.Errorf("%w: member not found in the study group", ErrInvalidMemberOperation)
	})

	err = resolveError(err, "removing member from the study group")

	// TODO send a notification

	return err
}

func (s *studyGroupServiceImpl) LeaveStudyGroup(studyGroupID models.StudyGroupID, memberID models.UserID) error {
	err := persistence.WithTransactionNoReturnVal(s.txMgr, func(tx persistence.Transaction) error {
		studyGroup, err := s.studyGroupRepo.GetStudyGroupByID(tx, studyGroupID)
		if err != nil {
			return err
		}

		for i, member := range studyGroup.Members {
			if member.UserID == memberID {
				studyGroup.Members = append(studyGroup.Members[:i], studyGroup.Members[i+1:]...)
				_, err = s.studyGroupRepo.UpdateStudyGroup(tx, studyGroup)
				return err
			}
		}

		return fmt.Errorf("%w: not currently a member of the study group", ErrInvalidMemberOperation)
	})

	err = resolveError(err, "leaving the study group")

	return err
}

func resolveError(err error, operation string) error {
	switch {
	case err == nil:
		return nil
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

func hasRole(userID models.UserID, role models.StudyGroupRole, members []models.StudyGroupMember) bool {
	return slices.ContainsFunc(members, func(m models.StudyGroupMember) bool {
		return m.UserID == userID && m.Role == role
	})
}

func setRole(userID models.UserID, role models.StudyGroupRole, members []models.StudyGroupMember) []models.StudyGroupMember {
	for i, member := range members {
		if member.UserID == userID {
			members[i].Role = role
			return members
		}
	}

	members = append(members, models.StudyGroupMember{
		UserID: userID,
		Role:   role,
	})
	return members
}
