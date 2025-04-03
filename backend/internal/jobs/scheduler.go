package jobs

import (
	"database/sql"
	"log"
	"time"

	"gdp8-backend/internal/models"
	"gdp8-backend/internal/persistence"
	"gdp8-backend/internal/repositories"
	"gdp8-backend/internal/services"
)

// ProcessUpcomingSessions is a helper that processes sessions based on the current time.
func ProcessUpcomingSessions(
	now time.Time,
	studySessionRepo repositories.StudySessionRepository,
	studyGroupRepo repositories.StudyGroupRepository,
	notificationService services.NotificationService,
	txMgr persistence.TransactionManager,
) {
	reminderOffset := 15 * time.Minute
	windowStart := now.Add(reminderOffset)
	windowEnd := windowStart.Add(1 * time.Minute)

	sessions, err := persistence.WithTransaction(txMgr, func(tx *sql.Tx) ([]models.StudySession, error) {
		return studySessionRepo.GetUpcomingSessions(tx, windowStart, windowEnd)
	})
	if err != nil {
		log.Printf("Error fetching upcoming sessions: %v", err)
		return
	}

	for _, session := range sessions {
		members, err := persistence.WithTransaction(txMgr, func(tx *sql.Tx) ([]models.StudyGroupMemberView, error) {
			return studyGroupRepo.GetMembers(tx, session.StudyGroupID)
		})
		if err != nil {
			log.Printf("Error fetching members for study group %d: %v", session.StudyGroupID, err)
			continue
		}

		if err := notificationService.AddStudySessionReminderNotification(&session, members); err != nil {
			log.Printf("Error sending reminder for session %d: %v", session.ID, err)
		} else {
			log.Printf("Reminder sent for session %d", session.ID)
		}
	}
}

// StartStudySessionReminderScheduler starts a background job that calls ProcessUpcomingSessions periodically.
func StartStudySessionReminderScheduler(
	studySessionRepo repositories.StudySessionRepository,
	studyGroupRepo repositories.StudyGroupRepository,
	notificationService services.NotificationService,
	txMgr persistence.TransactionManager,
) {
	ticker := time.NewTicker(1 * time.Minute)
	log.Println("Study session reminder scheduler started")

	go func() {
		for now := range ticker.C {
			ProcessUpcomingSessions(now, studySessionRepo, studyGroupRepo, notificationService, txMgr)
		}
	}()
}
