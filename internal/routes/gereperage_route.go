package routes

import (
	"github.com/cabitibaly/internal/handlers"
	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func GeorepRoutes(
	r *gin.Engine,
	georepHandler *handlers.GeorepHandler,
	authservice *services.AuthService,
) {
	georep := r.Group("/site")
	georep.Use(middlewares.AuthMiddleware(authservice))
	{
		georep.GET("/tous-les-sites/:id", georepHandler.LireUnSiteHandler)
	}

	admin := georep.Group("/")
	admin.Use(middlewares.AutorisationMiddleware(1))
	{
		admin.POST("/ajouter", georepHandler.CreerUnSiteHandler)
		admin.GET("/tous-les-sites", georepHandler.TousLesSitesHandler)
		admin.PATCH("/modifier/:id", georepHandler.ModifierUnSiteHandler)
		admin.DELETE("/supprimer/:id", georepHandler.SupprimerUnSiteHandler)
	}
}
