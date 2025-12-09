package services

import (
	"fmt"

	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
)

type NofificationService struct {
	notifRepo       *repositories.NotificationRepository
	utilisateurRepo *repositories.UtilisateurRepository
}

func NewNoficationService(
	notifRepo *repositories.NotificationRepository,
	utilisateurRepo *repositories.UtilisateurRepository,
) *NofificationService {
	return &NofificationService{
		notifRepo:       notifRepo,
		utilisateurRepo: utilisateurRepo,
	}
}

func (s *NofificationService) ModifierUneNotification(id uint, data map[string]any) error {
	notifiationExiste, _ := s.notifRepo.FindByID(id)

	if notifiationExiste == nil {
		return fmt.Errorf("Une erreur est survenue")
	}

	return s.notifRepo.Update(id, data)
}

func (s *NofificationService) ToutesLesNotifications(utilisateurID uint, page, limit int) ([]models.Notification, bool, error) {
	return s.notifRepo.FindAll(utilisateurID, page, limit)
}

func (s *NofificationService) SupprimerUneNotification(id uint) error {
	return s.notifRepo.Delele(id)
}
