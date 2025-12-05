package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
	"github.com/cabitibaly/pkg/utils"
	"gorm.io/gorm"
)

type PointageService struct {
	pointageRepo    *repositories.PointageRepository
	utilisateurRepo *repositories.UtilisateurRepository
}

func NewPointageService(
	pointageRepo *repositories.PointageRepository,
	utilisateurRepo *repositories.UtilisateurRepository,
) *PointageService {
	return &PointageService{
		pointageRepo:    pointageRepo,
		utilisateurRepo: utilisateurRepo,
	}
}

func (s *PointageService) PointageArrivee(utilisateurID uint, empLatitude, empLongitude float64) error {
	employe, errEmp := s.utilisateurRepo.FindWithGeoreperage(utilisateurID)

	if errEmp != nil {
		return errEmp
	}

	if !utils.EstDansLaZone(
		empLatitude,
		empLongitude,
		employe.Georeperage.Latitude,
		employe.Georeperage.Longitude,
		employe.Georeperage.Rayon,
	) {
		return fmt.Errorf("vous n'êtes pas sur site")
	}

	location, errLoc := time.LoadLocation("Africa/Ouagadougou")

	if errLoc != nil {
		panic(errLoc)
	}

	heureDebut := time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		8, 0, 0, 0,
		location,
	)

	pointage, err := s.pointageRepo.FindByUtilisateurID(uint(employe.IDUtilisateur))

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || pointage == nil {

		pointage = &models.Pointage{
			EstPresent:    true,
			EnRetard:      time.Now().After(heureDebut),
			UtilisateurID: int(utilisateurID),
		}

		return s.pointageRepo.Create(pointage)
	}

	datePointage := pointage.Arrive.Truncate(24 * time.Hour)
	dateDTO := time.Now().Truncate(24 * time.Hour)

	if datePointage.Equal(dateDTO) {
		return fmt.Errorf("vous avez déjà un pointage pour cette journée")
	}

	newPointage := models.Pointage{
		EstPresent:    true,
		EnRetard:      time.Now().After(heureDebut),
		UtilisateurID: int(utilisateurID),
	}

	return s.pointageRepo.Create(&newPointage)
}

func (s *PointageService) PointageDepart(pointageDTO dto.PointageDTO, utilisateurID uint, empLatitude, empLongitude float64) error {
	employe, errEmp := s.utilisateurRepo.FindWithGeoreperage(utilisateurID)

	if errEmp != nil {
		return errEmp
	}

	if !utils.EstDansLaZone(
		empLatitude,
		empLongitude,
		employe.Georeperage.Latitude,
		employe.Georeperage.Longitude,
		employe.Georeperage.Rayon,
	) {
		return fmt.Errorf("vous n'êtes pas sur site")
	}

	location, errLoc := time.LoadLocation("Africa/Ouagadougou")

	if errLoc != nil {
		panic(errLoc)
	}

	heureFin := time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		16, 0, 0, 0,
		location,
	)

	pointage, err := s.pointageRepo.FindByUtilisateurID(utilisateurID)

	if err != nil {
		return err
	}

	datePointage := pointage.Arrive.Truncate(24 * time.Hour)
	dateDTO := pointageDTO.Depart.Truncate(24 * time.Hour)

	if !datePointage.Equal(dateDTO) {
		return fmt.Errorf("vous n'avez pas de pointage pour cette journée")
	}

	return s.pointageRepo.Update(uint(pointage.IDPointage), map[string]any{
		"depart":            pointageDTO.Depart,
		"departAnticipe":    pointageDTO.Depart.Before(heureFin),
		"heuresTravaillees": pointageDTO.Depart.Sub(pointage.Arrive).Hours(),
	})
}

func (s *PointageService) TousLesPointages(utilisateurID uint, aujourdhui bool, date time.Time, page, limit int) ([]models.Pointage, bool, int64, error) {
	return s.pointageRepo.FindAll(utilisateurID, aujourdhui, date, page, limit)
}

func (s *PointageService) LireUnPointage(utilisateurID uint, date time.Time) (models.Pointage, error) {
	return s.pointageRepo.FindOneByDate(utilisateurID, date)
}

func (s *PointageService) SupprimerPointage(pointageID uint) error {
	return s.pointageRepo.Delete(pointageID)
}
