package routes

import (
	"github.com/cabitibaly/internal/handlers"
	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func NotificationRoutes(
	router *gin.Engine,
	notificationHandler *handlers.NotificationHandler,
	authService *services.AuthService,
) {
	notif := router.Group("/notification")
	notif.Use(middlewares.AuthMiddleware(authService))
	{
		notif.GET("/toutes-les-notifications", notificationHandler.ToutesLesNotificationsHandler)
		notif.PATCH("/modifier/:id", notificationHandler.ModifierUneNoticationHandler)
		notif.DELETE("/supprimer/:id", notificationHandler.SupprimerUneNotificationHandler)
	}
}
