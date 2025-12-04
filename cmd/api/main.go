package main

import (
	"github.com/cabitibaly/configs"
	"github.com/cabitibaly/internal/database"
	"github.com/cabitibaly/internal/handlers"
	"github.com/cabitibaly/internal/repositories"
	"github.com/cabitibaly/internal/routes"
	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.LoadConfig()
	utils.InitializeJWTSecret(cfg.JWTSecret)

	if err := database.Connect(cfg); err != nil {
		panic(err)
	}

	if err := database.Migration(); err != nil {
		panic(err)
	}

	database.IntDefaultRole()
	database.InitDefaultStatutConge()
	database.CreateGeoreperage()
	database.CreateAdmin()

	db := database.GetDB()

	authRepo := repositories.NewJWTRepository(db)
	utilisateurRepo := repositories.NewUtilisateurRepository(db)
	utilisateurService := services.NewUtilisateurService(utilisateurRepo, authRepo)
	utilisateurHandler := handlers.NewUtilisateurHandler(utilisateurService)

	authService := services.NewAuthservice(authRepo, utilisateurRepo)
	authHandler := handlers.NewAuthandler(authService)

	georepRepo := repositories.NewGeorepRepository(db)
	georepService := services.NewGeorepService(georepRepo)
	georepHandler := handlers.NewGeorepHandler(georepService)

	router := gin.Default()

	routes.AuthRoutes(
		router,
		authHandler,
		authService,
	)

	routes.UtilisateurRoutes(
		router,
		utilisateurHandler,
		authService,
	)

	routes.GeorepRoutes(
		router,
		georepHandler,
		authService,
	)

	router.Run()
}
