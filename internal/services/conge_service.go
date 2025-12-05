package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
	"gorm.io/gorm"
)

type CongeService struct {
	repo *repositories.CongeRepository
}

func NewCongeService(repo *repositories.CongeRepository) *CongeService {
	return &CongeService{repo: repo}
}

func (s *CongeService) FaireUneDemande(congeDTO dto.CongeDTO, utilisateurID uint) error {
	location, errLoc := time.LoadLocation("Africa/Ouagadougou")

	if errLoc != nil {
		return fmt.Errorf("erreur de timezone: %w", errLoc)
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

	if finJour.Before(debutJour) {
		return fmt.Errorf("la date de retour ne peut pas être antérieure à la date de départ")
	}

	if debutJour.Before(aujourdhui) {
		return fmt.Errorf("la date de debut ne peut pas être dans le passé")
	}

	_, err := s.repo.FindByUtilisateurIDAndPeriode(utilisateurID, debutJour, finJour)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return err
	}

	// Le nombre de congés d'un utilisateur: implementation (ref: claude)

	conge := &models.Conge{
		DateDepart:        debutJour,
		DateRetour:        finJour,
		Raison:            congeDTO.Raison,
		TypeConge:         congeDTO.TypeConge,
		PieceJointe:       congeDTO.PieceJointe,
		StatutCongeID:     1,
		UtilisateurID:     int(utilisateurID),
		DateCreationConge: maintenant,
	}

	return s.repo.Create(conge)
}

func (s *CongeService) TousLesPointages(utilsateurID uint, statutID uint, page, limit int) ([]models.Conge, bool, int64, error) {
	return s.repo.FindAll(utilsateurID, statutID, page, limit)
}

func (s *CongeService) LireUnConge(id uint) (*models.Conge, error) {
	return s.repo.FindByID(id)
}

func (s *CongeService) ModifierUnConge(id uint, data map[string]any) error {
	return s.repo.Update(id, data)
}

func (s *CongeService) SupprimerUnConge(id uint) error {
	return s.repo.Delete(id)
}
