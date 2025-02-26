package routes

import (
	"net/http"

	"firebase.google.com/go/v4/auth"

	"gdp8-backend/internal/handlers"
	"gdp8-backend/internal/middleware"
	"gdp8-backend/internal/services"
)

func RegisterNotificationRoutes(firebaseAuth *auth.Client, notificationService services.NotificationService) {
	handler := handlers.NewNotificationHandler(notificationService)

	http.HandleFunc("GET /api/notifications", middleware.WithFirebaseAuth(
		firebaseAuth, handler.GetUserNotifications))
	http.HandleFunc("GET /api/notifications/", middleware.WithFirebaseAuth(
		firebaseAuth, handler.GetUserNotifications))
	http.HandleFunc("DELETE /api/notifications/{id}", middleware.WithFirebaseAuth(
		firebaseAuth, handler.MarkNotificationAsRead))
}
