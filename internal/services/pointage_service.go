package services

import (
	"errors"
	"fmt"
	"time"

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
		return fmt.Errorf("erreur de timezone: %w", errLoc)
	}

	maintenant := time.Now().In(location)
	aujourdhui := time.Date(maintenant.Year(), maintenant.Month(), maintenant.Day(), 0, 0, 0, 0, location)
	heureDebut := aujourdhui.Add(8 * time.Hour)

	dernierPointage, err := s.pointageRepo.FindByUtilisateurID(uint(employe.IDUtilisateur))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if dernierPointage != nil {
		dateArriveDernierPointage := dernierPointage.Arrive.In(location)
		if dateArriveDernierPointage.Year() == maintenant.Year() &&
			dateArriveDernierPointage.Month() == maintenant.Month() &&
			dateArriveDernierPointage.Day() == maintenant.Day() {
			return fmt.Errorf("vous avez déjà un pointage pour cette journée")
		}
	}

	newPointage := models.Pointage{
		Arrive:        maintenant,
		EstPresent:    true,
		EnRetard:      maintenant.After(heureDebut),
		UtilisateurID: int(utilisateurID),
	}

	return s.pointageRepo.Create(&newPointage)
}

func (s *PointageService) PointageDepart(utilisateurID uint, empLatitude, empLongitude float64) error {
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
		return fmt.Errorf("erreur de timezone: %w", errLoc)
	}

	maintenant := time.Now().In(location)
	aujourdhui := time.Date(maintenant.Year(), maintenant.Month(), maintenant.Day(), 0, 0, 0, 0, location)
	heureFin := aujourdhui.Add(16 * time.Hour)

	dernierPointage, err := s.pointageRepo.FindByUtilisateurID(utilisateurID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("vous n'avez pas de pointage pour cette journée")
		}

		return err
	}

	dateArrivee := dernierPointage.Arrive.In(location)
	if dateArrivee.Year() != maintenant.Year() ||
		dateArrivee.Month() != maintenant.Month() ||
		dateArrivee.Day() != maintenant.Day() {
		return fmt.Errorf("vous n'avez pas de pointage pour cette journée")
	}

	if dernierPointage.Depart != nil && !dernierPointage.Depart.IsZero() {
		return fmt.Errorf("vous avez déjà un pointage de départ")
	}

	return s.pointageRepo.Update(uint(dernierPointage.IDPointage), map[string]any{
		"depart":            maintenant,
		"departAnticipe":    maintenant.Before(heureFin),
		"heuresTravaillees": maintenant.Sub(dernierPointage.Arrive).Hours(),
	})
}

func (s *PointageService) TousLesPointages(utilisateurID uint, aujourdhui bool, date time.Time, page, limit int) ([]models.Pointage, bool, int64, error) {
	return s.pointageRepo.FindAll(utilisateurID, aujourdhui, date, page, limit)
}

func (s *PointageService) LireUnPointage(utilisateurID uint, date time.Time) (*models.Pointage, error) {
	return s.pointageRepo.FindOneByDate(utilisateurID, date)
}

func (s *PointageService) SupprimerPointage(pointageID uint) error {
	return s.pointageRepo.Delete(pointageID)
}

func (s *PointageService) Stats() (int64, int64, int64, error) {
	totalPresent, totalRetard, err := s.pointageRepo.GetTotalPresentAndRetard()
	return s.utilisateurRepo.GetTotalEmploye(), totalPresent, totalRetard, err
}
