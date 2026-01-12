package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
	"github.com/cabitibaly/pkg/utils"
	"github.com/xuri/excelize/v2"
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

func (s *PointageService) WritePointageRow(f *excelize.File, sheetName string, row int, pointage models.Pointage, style int) {
	nomUtilisateur := "N/A"
	if pointage.Utilisateur != nil {
		nomUtilisateur = pointage.Utilisateur.Nom + " " + pointage.Utilisateur.Prenom
	}

	depart := "N/A"
	if pointage.Depart != nil {
		depart = pointage.Depart.Format("02/01/2006 15:04")
	}

	heureTravaillees := "N/A"
	if pointage.HeuresTravaillees != nil {
		heureTravaillees = utils.FormatHeure(*pointage.HeuresTravaillees)
	}

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), pointage.IDPointage)
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), nomUtilisateur)
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), pointage.Arrive.Format("02/01/2006 15:04"))
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), depart)
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), heureTravaillees)
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), utils.BoolToOuiNon(pointage.EstPresent))
	f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), utils.BoolToOuiNon(pointage.EnRetard))
	f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), utils.BoolToOuiNon(pointage.DepartAnticipe))

	for col := 'A'; col <= 'H'; col++ {
		cell := fmt.Sprintf("%s%d", string(col), row)
		f.SetCellStyle(sheetName, cell, cell, style)
	}
}

func (s *PointageService) CreateExcel(debut, fin time.Time) (*excelize.File, error) {
	pointages, err := s.pointageRepo.FindByDateRange(debut, fin)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération pointages: %v", err)
	}

	f := excelize.NewFile()
	sheetName := "Pointages"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("erreur de création de la feuille: %v", err)
	}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    utils.CreateBorder(),
	})

	dataStyle, _ := f.NewStyle(&excelize.Style{
		Border:    utils.CreateBorder(),
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	headers := []string{
		"ID", "Utilisateur", "Arrivée", "Départ", "Heures Travaillées",
		"Présent", "En Retard", "Départ Anticipé",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	for i, pointage := range pointages {
		row := i + 2
		s.WritePointageRow(f, sheetName, row, pointage, dataStyle)
	}

	columnWidths := map[string]float64{
		"A": 10, "B": 25, "C": 20, "D": 20, "E": 18,
		"F": 12, "G": 12, "H": 15, "I": 20,
	}

	for col, width := range columnWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return f, nil
}

func (s *PointageService) Export(debut, fin time.Time) ([]byte, error) {

	f, err := s.CreateExcel(debut, fin)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'écriture du buffer: %v", err)
	}

	return buffer.Bytes(), nil
}
