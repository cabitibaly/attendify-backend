package routes

import (
	"github.com/cabitibaly/internal/handlers"
	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func SiteRoutes(
	r *gin.Engine,
	siteHandler *handlers.SiteHandler,
	authservice *services.AuthService,
) {
	site := r.Group("/site")
	site.Use(middlewares.AuthMiddleware(authservice))
	{
		site.GET("/tous-les-sites/:id", siteHandler.LireUnSiteHandler)
	}

	admin := site.Group("/")
	admin.Use(middlewares.AutorisationMiddleware(1))
	{
		admin.POST("/ajouter", siteHandler.CreerUnSiteHandler)
		admin.GET("/tous-les-sites", siteHandler.TousLesSitesHandler)
		admin.PATCH("/modifier/:id", siteHandler.ModifierUnSiteHandler)
		admin.DELETE("/supprimer/:id", siteHandler.SupprimerUnSiteHandler)
	}
}
