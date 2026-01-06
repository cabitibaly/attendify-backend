package routes

import (
	"github.com/cabitibaly/internal/handlers"
	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func PushTokenRoute(
	r *gin.Engine,
	pushTokenHandler *handlers.PushTokenHandler,
	authService *services.AuthService,
) {
	pushToken := r.Group("/notification-push")
	pushToken.Use(middlewares.AuthMiddleware(authService))
	{
		pushToken.POST("", pushTokenHandler.EnregistrerOuModifierPushTokenHandler)
		pushToken.DELETE("/:token", pushTokenHandler.SupprimerPushTokenHandler)
	}
}
