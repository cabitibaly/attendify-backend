package services

import (
	"fmt"
	"time"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
	"github.com/cabitibaly/pkg/utils"
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

func (s *PointageService) PointageArrivee(pointageDTO dto.PointageDTO, empLatitude, empLongitude float64) error {
	employe, errEmp := s.utilisateurRepo.FindByID(uint(pointageDTO.UtilisateurID))

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

	heureDebut := time.Date(
		pointageDTO.Arrivee.Year(),
		pointageDTO.Arrivee.Month(),
		pointageDTO.Arrivee.Day(),
		8, 0, 0, 0,
		pointageDTO.Arrivee.Location(),
	)

	pointage, err := s.pointageRepo.FindByUtilisateurID(uint(pointageDTO.UtilisateurID))

	if err != nil {
		return err
	}

	datePointage := pointage.Arrivee.Truncate(24 * time.Hour)
	dateDTO := pointageDTO.Arrivee.Truncate(24 * time.Hour)

	if datePointage.Equal(dateDTO) {
		return fmt.Errorf("vous avez déjà un pointage pour cette journée")
	}

	newPointage := models.Pointage{
		// Arrivee:       pointageDTO.Arrivee,
		EstPresent:    true,
		EnRetard:      pointageDTO.Arrivee.After(heureDebut),
		UtilisateurID: pointageDTO.UtilisateurID,
	}

	return s.pointageRepo.Create(&newPointage)
}

func (s *PointageService) PointageDepart(pointageDTO dto.PointageDTO) error {

	heureFin := time.Date(
		pointageDTO.Arrivee.Year(),
		pointageDTO.Arrivee.Month(),
		pointageDTO.Arrivee.Day(),
		16, 0, 0, 0,
		pointageDTO.Arrivee.Location(),
	)

	pointage, err := s.pointageRepo.FindByUtilisateurID(uint(pointageDTO.UtilisateurID))

	if err != nil {
		return err
	}

	datePointage := pointage.Arrivee.Truncate(24 * time.Hour)
	dateDTO := pointageDTO.Depart.Truncate(24 * time.Hour)

	if !datePointage.Equal(dateDTO) {
		return fmt.Errorf("vous n'avez pas de pointage pour cette journée")
	}

	return s.pointageRepo.Update(uint(pointage.IDPointage), map[string]any{
		"depart":            pointageDTO.Depart,
		"departAnticipe":    pointageDTO.Depart.Before(heureFin),
		"heuresTravaillees": pointageDTO.Depart.Sub(pointage.Arrivee).Hours(),
	})
}

func (s *PointageService) TousLesPointages(utilisateurID uint, aujourdhui bool, page, limit int) ([]models.Pointage, bool, int64, error) {
	return s.pointageRepo.FindAll(aujourdhui, utilisateurID, page, limit)
}

func (s *PointageService) LireUnPointage(date time.Time) (models.Pointage, error) {
	return s.pointageRepo.FindOneByDate(date)
}

func (s *PointageService) TousLesPointagesParDate(date time.Time, page, limit int) ([]models.Pointage, bool, int64, error) {
	return s.pointageRepo.FindAllByDate(date, page, limit)
}

func (s *PointageService) SupprimerPointage(pointageID uint) error {
	return s.pointageRepo.Delete(pointageID)
}
