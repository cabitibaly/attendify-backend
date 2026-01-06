package database

import (
	"fmt"
	"log"
	"time"

	"github.com/cabitibaly/configs"
	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *configs.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("impossible de se connecter à la base donnée")
	}

	sqlDB, erreur := DB.DB()

	if erreur != nil {
		return erreur
	}

	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Connexion à la base de donnée réussi ........ ✅")
	return nil

}

func GetDB() *gorm.DB {
	return DB
}

func Migration() error {
	err := DB.AutoMigrate(
		models.Utilisateur{},
		models.Role{},
		models.Site{},
		models.Pointage{},
		models.Conge{},
		models.StatutConge{},
		models.RefreshToken{},
		models.Notification{},
		models.PushToken{},
	)

	if err != nil {
		return fmt.Errorf("une erreur est survenue lors de la migration ........ ❌")
	}

	log.Println("Migration réussie ........ ✅")
	return nil
}

func IntDefaultRole() {
	var count int64
	DB.Model(&models.Role{}).Count(&count)

	if count == 0 {
		roles := []models.Role{
			{IDRole: 1, Libelle: "ADMIN"},
			{IDRole: 2, Libelle: "EMPLOYEE"},
		}
		DB.Create(&roles)
	}
}

func InitDefaultStatutConge() {
	var count int64
	DB.Model(&models.StatutConge{}).Count(&count)

	if count == 0 {
		statutConges := []models.StatutConge{
			{IDStatutConge: 1, LibelleStatut: "EN_ATTENTE"},
			{IDStatutConge: 2, LibelleStatut: "APPROUVEE"},
			{IDStatutConge: 3, LibelleStatut: "REJETEE"},
		}
		DB.Create(&statutConges)
	}
}

func CreateSite() {
	var site models.Site
	err := DB.Where("nom_site = ?", "POWERTECH BOBO site 22").First(&site).Error

	if err != nil {
		site.NomSite = "POWERTECH BOBO site 22"
		site.Latitude = 11.1895678
		site.Longitude = -4.3152317
		site.Rayon = 100

		DB.Create(&site)
	}
}

func CreateAdmin() {
	var utilisateur models.Utilisateur
	err := DB.Where("email = ?", "admin@powertechbf.com").First(&utilisateur).Error

	if err != nil {
		hashedPassword, err := utils.HashPassword("admin")

		if err != nil {
			panic(err)
		}

		var site models.Site

		erreur := DB.Where("nom_site = ?", "POWERTECH BOBO site 22").First(&site).Error

		if erreur != nil {
			panic(erreur)
		}

		utilisateur.Nom = "admin"
		utilisateur.Prenom = "admin"
		utilisateur.Email = "admin@powertechbf.com"
		utilisateur.Poste = "Directeur"
		utilisateur.Telephone = "64141525"
		utilisateur.RoleID = 1
		utilisateur.SiteID = site.IDSite
		utilisateur.MotDePasse = hashedPassword

		DB.Create(&utilisateur)
	}
}
