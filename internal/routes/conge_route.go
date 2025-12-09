package routes

import (
	"github.com/cabitibaly/internal/handlers"
	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func CongeRoute(
	r *gin.Engine,
	congeHandler *handlers.CongeHandler,
	authService *services.AuthService,
) {
	conge := r.Group("/conge")
	conge.Use(middlewares.AuthMiddleware(authService), middlewares.AutorisationMiddleware(2))
	{
		conge.POST("/faire-une-demande", congeHandler.FaireUneDemandeHandler)
		conge.GET("/tous-les-conges-employe", congeHandler.TousLesCongésEmployeHandler)
		conge.GET("/tous-les-conges-employe/:id", congeHandler.LireUnCongeHandler)
		conge.PATCH("/modifier/:id", congeHandler.ModifierUnCongeHandler)
	}

	admin := r.Group("/conge")
	admin.Use(middlewares.AuthMiddleware(authService), middlewares.AutorisationMiddleware(1))
	{
		admin.GET("/tous-les-conges", congeHandler.TousLesCongesHandler)
		admin.GET("/tous-les-conges/:id", congeHandler.LireUnCongeHandler)
		admin.PATCH("/modifier-statut/:id", congeHandler.ModifierStatutCongeHandler)
		admin.DELETE("/supprimer/:id", congeHandler.SupprimerUnCongeHandler)
	}
}
