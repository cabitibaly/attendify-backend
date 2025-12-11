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
	utils.InitializePushURL(cfg.PushURL)

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

	notifRepo := repositories.NewNotificationRepository(db)
	authRepo := repositories.NewJWTRepository(db)
	georepRepo := repositories.NewGeorepRepository(db)
	utilisateurRepo := repositories.NewUtilisateurRepository(db)
	utilisateurService := services.NewUtilisateurService(utilisateurRepo, authRepo)
	utilisateurHandler := handlers.NewUtilisateurHandler(utilisateurService)

	authService := services.NewAuthservice(authRepo, utilisateurRepo, georepRepo)
	authHandler := handlers.NewAuthandler(authService)

	georepService := services.NewGeorepService(georepRepo)
	georepHandler := handlers.NewGeorepHandler(georepService)

	pushTokenRepo := repositories.NewPushTokenRepository(db)
	pushTokenService := services.NewPushTokenService(pushTokenRepo)
	pushTokenHandler := handlers.NewPushTokenHandler(pushTokenService)

	pointageRepo := repositories.NewPointageRepository(db)
	pointageService := services.NewPointageService(pointageRepo, utilisateurRepo, pushTokenService)
	pointageHandler := handlers.NewPointageHandler(pointageService)

	congeRepo := repositories.NewCongeRepository(db)
	congeService := services.NewCongeService(congeRepo, utilisateurRepo, notifRepo, pushTokenService)
	congeHandler := handlers.NewCongeHandler(congeService)

	notifService := services.NewNoficationService(notifRepo)
	notifHandler := handlers.NewNoficationHandler(notifService)

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

	routes.PointageRoutes(
		router,
		pointageHandler,
		authService,
	)

	routes.CongeRoute(
		router,
		congeHandler,
		authService,
	)

	routes.NotificationRoutes(
		router,
		notifHandler,
		authService,
	)

	routes.PushTokenRoute(
		router,
		pushTokenHandler,
		authService,
	)

	router.Run()
}
