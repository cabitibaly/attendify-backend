package services

import (
	"fmt"

	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
)

type NofificationService struct {
	repo *repositories.NotificationRepository
}

func NewNoficationService(repo *repositories.NotificationRepository) *NofificationService {
	return &NofificationService{repo: repo}
}

func (s *NofificationService) ToutesLesNotifications(utilisateurID uint, page, limit int) ([]models.Notification, bool, error) {
	return s.repo.FindAll(utilisateurID, page, limit)
}

func (s *NofificationService) ModifierUneNotification(id uint, data map[string]any) error {
	notifiationExiste, _ := s.repo.FindByID(id)

	if notifiationExiste == nil {
		return fmt.Errorf("une erreur est survenue")
	}

	return s.repo.Update(id, data)
}

func (s *NofificationService) SupprimerUneNotification(id uint) error {
	return s.repo.Delele(id)
}
