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
	pointageRepo     *repositories.PointageRepository
	utilisateurRepo  *repositories.UtilisateurRepository
	notifpushService *PushTokenService
}

func NewPointageService(
	pointageRepo *repositories.PointageRepository,
	utilisateurRepo *repositories.UtilisateurRepository,
	notifpushService *PushTokenService,
) *PointageService {
	return &PointageService{
		pointageRepo:     pointageRepo,
		utilisateurRepo:  utilisateurRepo,
		notifpushService: notifpushService,
	}
}

func (s *PointageService) PointageArrivee(utilisateurID uint, empLatitude, empLongitude float64) error {
	employe, errEmp := s.utilisateurRepo.FindByID(utilisateurID)

	if errEmp != nil {
		return errEmp
	}

	if !utils.EstDansLaZone(
		empLatitude,
		empLongitude,
		employe.Site.Latitude,
		employe.Site.Longitude,
		employe.Site.Rayon,
	) {
		return fmt.Errorf("vous n'êtes pas sur site")
	}

	location, errLoc := time.LoadLocation("Africa/Ouagadougou")

	if errLoc != nil {
		return fmt.Errorf("erreur de timezone: %w", errLoc)
	}

	heureDebutGeorep := employe.Site.HeureDebut.Hour()
	minuteDebutGeorep := employe.Site.HeureDebut.Minute()
	secondeDebutGeorep := employe.Site.HeureDebut.Second()

	maintenant := time.Now().In(location)
	heureDebut := time.Date(maintenant.Year(), maintenant.Month(), maintenant.Day(), heureDebutGeorep, minuteDebutGeorep, secondeDebutGeorep, 0, location)

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

	errCreate := s.pointageRepo.Create(&newPointage)
	if errCreate != nil {
		return errCreate
	}

	_ = s.notifpushService.EnvoyerNotificationPushAUnUtilisateur(
		uint(utilisateurID),
		"Pointage d'arrivée",
		"Bonjour, nous sommes ravie de vous revoir.\n Bonne journée !",
	)

	_ = s.notifpushService.EnvoyerNotificationPushAUnUtilisateur(
		1,
		"Pointage d'arrivée",
		fmt.Sprintf("%s %s est arrivé(e) a %s", employe.Nom, employe.Prenom, maintenant.Format("08:00:00")),
	)

	return nil
}

func (s *PointageService) PointageDepart(utilisateurID uint, empLatitude, empLongitude float64) error {
	employe, errEmp := s.utilisateurRepo.FindByID(utilisateurID)

	if errEmp != nil {
		return errEmp
	}

	if !utils.EstDansLaZone(
		empLatitude,
		empLongitude,
		employe.Site.Latitude,
		employe.Site.Longitude,
		employe.Site.Rayon,
	) {
		return fmt.Errorf("vous n'êtes pas sur site")
	}

	location, errLoc := time.LoadLocation("Africa/Ouagadougou")

	if errLoc != nil {
		return fmt.Errorf("erreur de timezone: %w", errLoc)
	}

	heureFinGeorep := employe.Site.HeureFin.Hour()
	minuteFinGeorep := employe.Site.HeureFin.Minute()
	secondeFinGeorep := employe.Site.HeureFin.Second()

	maintenant := time.Now().In(location)
	heureFin := time.Date(maintenant.Year(), maintenant.Month(), maintenant.Day(), heureFinGeorep, minuteFinGeorep, secondeFinGeorep, 0, location)

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

	errUpdate := s.pointageRepo.Update(uint(dernierPointage.IDPointage), map[string]any{
		"depart":            maintenant,
		"departAnticipe":    maintenant.Before(heureFin),
		"heuresTravaillees": maintenant.Sub(dernierPointage.Arrive).Hours(),
	})

	if errUpdate != nil {
		return errUpdate
	}

	_ = s.notifpushService.EnvoyerNotificationPushAUnUtilisateur(
		uint(dernierPointage.UtilisateurID),
		"Pointage de départ",
		"Vous avez terminé votre journée, nous avons hate de vous revoir",
	)

	_ = s.notifpushService.EnvoyerNotificationPushAUnUtilisateur(
		1,
		"Pointage de depart",
		fmt.Sprintf("%s %s est parti(e) a %s", employe.Nom, employe.Prenom, maintenant.Format("16:00:00")),
	)

	return nil
}

func (s *PointageService) TousLesPointages(utilisateurID uint, aujourdhui bool, date time.Time, page, limit int) ([]models.Pointage, bool, int64, error) {
	return s.pointageRepo.FindAll(utilisateurID, aujourdhui, date, page, limit)
}

func (s *PointageService) SupprimerPointage(pointageID uint) error {
	return s.pointageRepo.Delete(pointageID)
}

func (s *PointageService) Stats() (int64, int64, int64, error) {
	totalPresent, totalRetard, err := s.pointageRepo.GetTotalPresentAndRetard()
	return s.utilisateurRepo.GetTotalEmploye(), totalPresent, totalRetard, err
}
