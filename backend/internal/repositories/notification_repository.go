package repositories

import (
	"database/sql"

	"gdp8-backend/internal/models"
)

type NotificationRepository interface {
	GetUserNotifications(tx *sql.Tx, id models.UserID) ([]models.NotificationView, error)
	DeleteUserNotification(tx *sql.Tx, id models.NotificationID, userID models.UserID) error
	AddNotification(tx *sql.Tx, notification models.NotificationDetails, usersToBeNotified []models.UserID) error
}

type SQLNotificationRepository struct {
}

func (s *SQLNotificationRepository) GetUserNotifications(tx *sql.Tx,
	id models.UserID) ([]models.NotificationView, error) {

	query := `
		SELECT 
			n.id,
			n.type,
			n.triggering_user_id,
			triggering_user.name as triggering_user_name,
			n.target_user_id,
			target_user.name as target_user_name,
			n.study_group_id,
			sg.name as study_group_name,
			n.message_id,
			n.created_at
		FROM notifications n
		INNER JOIN user_notifications un ON n.id = un.notification_id
		INNER JOIN users triggering_user ON n.triggering_user_id = triggering_user.id
		LEFT JOIN users target_user ON n.target_user_id = target_user.id
		INNER JOIN study_groups sg ON n.study_group_id = sg.id
		WHERE un.user_id = ?`

	rows, err := tx.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	notifications, err := readNotifications(rows)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *SQLNotificationRepository) DeleteUserNotification(tx *sql.Tx,
	id models.NotificationID, userID models.UserID) error {

	query := `
		DELETE FROM user_notifications
		WHERE user_id = ? AND notification_id = ?
	`

	_, err := tx.Exec(query, userID, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLNotificationRepository) AddNotification(tx *sql.Tx,
	notification models.NotificationDetails, usersToBeNotified []models.UserID) error {

	queryInsertNotification := `
		INSERT INTO notifications (
			type,
			triggering_user_id,
			target_user_id,
			study_group_id,
			message_id
		) VALUES (?, ?, ?, ?, ?)
	`

	res, err := tx.Exec(
		queryInsertNotification,
		notification.Type,
		notification.TriggeringUserID,
		notification.TargetUserID,
		notification.StudyGroupID,
		notification.MessageID,
	)
	if err != nil {
		return err
	}

	notificationID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	queryInsertUserNotification := `
		INSERT INTO user_notifications (
			user_id,
			notification_id
		) VALUES (?, ?)
	`

	for _, userID := range usersToBeNotified {
		_, err := tx.Exec(queryInsertUserNotification, userID, notificationID)
		if err != nil {
			return err
		}
	}

	return nil
}

func readNotification(s scanner) (*models.NotificationView, error) {
	var notification models.NotificationView
	var targetUserID sql.NullString
	var targetUserName sql.NullString
	var messageID sql.NullInt64

	if err := s.Scan(
		&notification.ID,
		&notification.Type,
		&notification.TriggeringUserID,
		&notification.TriggeringUserName,
		&targetUserID,
		&targetUserName,
		&notification.StudyGroupID,
		&notification.StudyGroupName,
		&messageID,
		&notification.CreatedAt,
	); err != nil {
		return nil, err
	}

	if targetUserID.Valid {
		id := models.UserID(targetUserID.String)
		notification.TargetUserID = &id
	}

	if targetUserName.Valid {
		name := targetUserName.String
		notification.TargetUserName = &name
	}

	if messageID.Valid {
		notification.MessageID = &messageID.Int64
	}

	return &notification, nil
}

func readNotifications(rows *sql.Rows) ([]models.NotificationView, error) {
	var notifications []models.NotificationView

	for rows.Next() {
		notification, err := readNotification(rows)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, *notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}
