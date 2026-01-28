package routes

import (
	"github.com/cabitibaly/internal/handlers"
	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func UtilisateurRoutes(
	r *gin.Engine,
	utilisateurHandler *handlers.UtilisateurHandler,
	authService *services.AuthService,
) {
	utilisateur := r.Group("/compte")
	utilisateur.Use(middlewares.AuthMiddleware(authService))
	{
		utilisateur.PATCH("/modifier-son-mot-de-passe", utilisateurHandler.ModifierSonMotDePasseHandler)
		utilisateur.GET("/mes-informations", utilisateurHandler.MesInformationsHandler)
		utilisateur.PATCH("/modifier-un-compte", utilisateurHandler.ModifierSonCompteHandler)
	}

	admin := utilisateur.Group("/")
	admin.Use(middlewares.AutorisationMiddleware(1))
	{
		admin.GET("/tous-les-employes", utilisateurHandler.TousLesEmployesHandler)
		admin.GET("/tous-les-employes/:id", utilisateurHandler.LireUnEmployeHandler)
		admin.PATCH("/changer-de-site/:id/:siteID", utilisateurHandler.ChangerDeSiteHandler)
		admin.DELETE("/supprimer-un-compte/:id", utilisateurHandler.SupprimerUnCompteHandler)
	}
}
