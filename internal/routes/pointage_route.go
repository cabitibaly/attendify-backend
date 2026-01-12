package routes

import (
	"github.com/cabitibaly/internal/handlers"
	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func PointageRoutes(
	r *gin.Engine,
	pointageHandler *handlers.PointageHandler,
	authService *services.AuthService,
) {

	pointage := r.Group("/pointage")
	pointage.Use(middlewares.AuthMiddleware(authService), middlewares.AutorisationMiddleware(2))
	{
		pointage.POST("/arrive", pointageHandler.PointageArriveeHandler)
		pointage.PATCH("/depart", pointageHandler.PointageDepartHandler)
		pointage.GET("/tous-mes-pointages", pointageHandler.TousLesPointagesEmployeHandler)
	}

	admin := r.Group("/pointage")
	admin.Use(middlewares.AuthMiddleware(authService), middlewares.AutorisationMiddleware(1))
	{
		admin.GET("/tous-les-pointages", pointageHandler.TousLesPointagesHandler)
		admin.GET("/stats", pointageHandler.StatsHandler)
		admin.GET("/export", pointageHandler.ExportHandler)
		admin.DELETE("/supprimer/:id", pointageHandler.SupprimerUnPointageHandler)
	}
}
