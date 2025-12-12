package routes

import (
	"github.com/cabitibaly/internal/handlers"
	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	authService *services.AuthService,
) {
	auth := r.Group("/auth")

	auth.POST("/connexion-admin", authHandler.ConnexionAdminHandler)
	auth.POST("/connexion-employe", authHandler.ConnexionEmployeHandler)
	auth.POST("/deconnexion", authHandler.DeconnexionHandler)
	auth.POST("/refresh-token", authHandler.RefreshTokenHandler)

	auth.Use(
		middlewares.AuthMiddleware(authService),
		middlewares.AutorisationMiddleware(1),
	)

	auth.POST("/nouveau-compte-employe", authHandler.CreerUnCompteHandler)
	auth.PATCH("/reinitialiser-mot-de-passe/:id", authHandler.ReinitialiserMotDePasseHandler)
}
