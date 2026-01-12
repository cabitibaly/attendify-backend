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

type CongeService struct {
	congeRepo        *repositories.CongeRepository
	utilisateurRepo  *repositories.UtilisateurRepository
	notifRepo        *repositories.NotificationRepository
	notifpushService *PushTokenService
}

func NewCongeService(
	congeRepo *repositories.CongeRepository,
	utilisateurRepo *repositories.UtilisateurRepository,
	notifRepo *repositories.NotificationRepository,
	notifpushService *PushTokenService,
) *CongeService {
	return &CongeService{
		congeRepo:        congeRepo,
		utilisateurRepo:  utilisateurRepo,
		notifRepo:        notifRepo,
		notifpushService: notifpushService,
	}
}

func (s *CongeService) FaireUneDemande(congeDTO dto.CongeDTO, utilisateurID uint) error {
	utilisateur, err := s.utilisateurRepo.FindByID(utilisateurID)
	if err != nil {
		return fmt.Errorf("cet utilisateur n'existe pas")
	}

	location, errLoc := time.LoadLocation("Africa/Ouagadougou")

	if errLoc != nil {
		return fmt.Errorf("erreur de timezone: %w", errLoc)
	}

	if congeDTO.Raison == "" {
		return fmt.Errorf("qu'elle est la raison de votre congé")
	}

	if congeDTO.TypeConge == "" {
		return fmt.Errorf("quel type de congé vous avez")
	}

	if congeDTO.DateDepart.In(location).IsZero() || congeDTO.DateRetour.In(location).IsZero() || congeDTO.Raison == "" {
		return fmt.Errorf("date de départ, date de retour et raison sont obligatoires")
	}

	dateDebut := congeDTO.DateDepart.In(location)
	dateFin := congeDTO.DateRetour.In(location)
	maintenant := time.Now().In(location)

	debutJour := time.Date(dateDebut.Year(), dateDebut.Month(), dateDebut.Day(), 0, 0, 0, 0, location)
	finJour := time.Date(dateFin.Year(), dateFin.Month(), dateFin.Day(), 0, 0, 0, 0, location)
	aujourdhui := time.Date(maintenant.Year(), maintenant.Month(), maintenant.Day(), 0, 0, 0, 0, location)

	if debutJour.Equal(aujourdhui) {
		return fmt.Errorf("la date de départ ne peut pas être aujourd'hui")
	}

	if finJour.Before(debutJour) {
		return fmt.Errorf("la date de retour ne peut pas être antérieure à la date de départ")
	}

	if debutJour.Before(aujourdhui) {
		return fmt.Errorf("la date de debut ne peut pas être dans le passé")
	}

	// On verifie qu'il n'y a pas de congé dans cette période
	congesExistants, err := s.congeRepo.FindByUtilisateurIDAndPeriode(utilisateurID, debutJour, finJour)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if len(congesExistants) > 0 {
		return fmt.Errorf("vous avez déjà un congé dans cette période")
	}

	congeRestant, err := s.utilisateurRepo.GetSoldeConge(utilisateurID)
	if err != nil {
		return err
	}

	nombreDeJours := utils.CalculerNombreJours(debutJour, finJour)

	if nombreDeJours > congeRestant {
		return fmt.Errorf("le nombre de jours ne peut pas dépasser le solde de vos congés")
	}

	conge := &models.Conge{
		DateDepart:        debutJour,
		DateRetour:        finJour,
		Raison:            congeDTO.Raison,
		TypeConge:         congeDTO.TypeConge,
		PieceJointe:       congeDTO.PieceJointe,
		NombreJours:       nombreDeJours,
		StatutCongeID:     1,
		UtilisateurID:     int(utilisateurID),
		DateCreationConge: maintenant,
	}

	errCreate := s.congeRepo.Create(conge)

	if errCreate != nil {
		return errCreate
	}

	_ = s.notifRepo.Create(&models.Notification{
		Titre:            "Nouvelle demande",
		Message:          fmt.Sprintf("Une nouvelle demande de congé a été créée par %s %s", utilisateur.Nom, utilisateur.Prenom),
		TypeNotification: "SUCCESS",
		UtilisateurID:    1,
	})

	_ = s.notifpushService.EnvoyerNotificationPushAUnUtilisateur(
		1,
		"Nouvelle demande",
		"Vous avez une nouvelle demande de congé",
	)

	return nil
}

func (s *CongeService) TousLesConges(utilsateurID uint, statutID uint, page, limit int) ([]models.Conge, bool, int64, error) {
	return s.congeRepo.FindAll(utilsateurID, statutID, page, limit)
}

func (s *CongeService) LireUnConge(id uint) (*models.Conge, error) {
	return s.congeRepo.FindByID(id)
}

func (s *CongeService) ModifierUnConge(id uint, data map[string]any) error {
	congeExistant, err := s.congeRepo.FindByID(id)
	if err != nil {
		return err
	}

	if congeExistant.StatutCongeID == 2 || congeExistant.StatutCongeID == 3 {
		return fmt.Errorf("ce congé a été traité, donc il ne peut pas être modifié")
	}

	return s.congeRepo.Update(id, data)
}

func (s *CongeService) ModifierStatutConge(id uint, statutID uint) error {
	congeExistant, err := s.congeRepo.FindByID(id)
	if err != nil {
		return err
	}

	if congeExistant.StatutCongeID == 2 || congeExistant.StatutCongeID == 3 {
		return fmt.Errorf("ce congé a été traité, donc il ne peut pas être modifié")
	}

	if statutID == 2 && congeExistant.StatutCongeID == 1 {
		congeResant, err := s.utilisateurRepo.GetSoldeConge(uint(congeExistant.UtilisateurID))
		if err != nil {
			return err
		}

		errEmpUpdate := s.utilisateurRepo.Update(uint(congeExistant.UtilisateurID), map[string]any{"soldeConge": congeResant - congeExistant.NombreJours})
		if errEmpUpdate != nil {
			return errEmpUpdate
		}
	}

	errCongeUpdate := s.congeRepo.Update(id, map[string]any{"StatutCongeID": statutID})
	if errCongeUpdate != nil {
		return errCongeUpdate
	}

	message := ""

	if statutID == 2 {
		message = "approuvée"
	} else {
		message = "rejetée"
	}

	_ = s.notifRepo.Create(&models.Notification{
		Titre:            "Validation du congé",
		Message:          fmt.Sprintf("Votre demande a été %s", message),
		TypeNotification: "SUCCESS",
		UtilisateurID:    congeExistant.UtilisateurID,
	})

	if statutID == 2 {
		_ = s.notifpushService.EnvoyerNotificationPushAUnUtilisateur(
			uint(congeExistant.UtilisateurID),
			"Validation du congé",
			"Votre demande a été approuvée",
		)
	} else {
		_ = s.notifpushService.EnvoyerNotificationPushAUnUtilisateur(
			uint(congeExistant.UtilisateurID),
			"Validation du congé",
			"Votre demande a été rejetée",
		)
	}

	return nil
}

func (s *CongeService) SupprimerUnConge(id uint) error {
	return s.congeRepo.Delete(id)
}
