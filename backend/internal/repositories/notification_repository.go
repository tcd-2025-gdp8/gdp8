package repositories

import (
	"database/sql"

	"gdp8-backend/internal/models"
)

type NotificationRepository interface {
	GetUserNotifications(tx *sql.Tx, id models.UserID) ([]models.Notification, error)
	DeleteUserNotification(tx *sql.Tx, id models.NotificationID, userID models.UserID) error
	AddNotification(tx *sql.Tx, notification models.NotificationDetails, usersToBeNotified []models.UserID) error
}

type SQLNotificationRepository struct {
}

func (s *SQLNotificationRepository) GetUserNotifications(tx *sql.Tx,
	id models.UserID) ([]models.Notification, error) {

	query := `
        SELECT 
            n.id,
            n.type,
            n.payload,
            n.created_at
        FROM notifications n
        INNER JOIN user_notifications un ON n.id = un.notification_id
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
	return err
}

func (s *SQLNotificationRepository) AddNotification(tx *sql.Tx,
	notification models.NotificationDetails, usersToBeNotified []models.UserID) error {

	queryInsertNotification := `
        INSERT INTO notifications (
            type,
            payload
        ) VALUES (?, ?)
    `

	res, err := tx.Exec(
		queryInsertNotification,
		notification.Type,
		notification.Payload,
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

func readNotification(s scanner) (*models.Notification, error) {
	var notification models.Notification

	if err := s.Scan(
		&notification.ID,
		&notification.Type,
		&notification.Payload,
		&notification.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &notification, nil
}

func readNotifications(rows *sql.Rows) ([]models.Notification, error) {
	var notifications []models.Notification

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
