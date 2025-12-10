package services

import (
	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
)

type PushTokenService struct {
	repo *repositories.PushTokenRepository
}

func NewPushTokenService(repo *repositories.PushTokenRepository) *PushTokenService {
	return &PushTokenService{repo: repo}
}

func (s *PushTokenService) EnregistrerOuModifierPushToken(pushTokenDTO dto.PushTokenDTO) error {
	return s.repo.CreateOrUpdate(&models.PushToken{
		PushToken:     pushTokenDTO.PushToken,
		DeviceType:    pushTokenDTO.DeviceType,
		UtilisateurID: pushTokenDTO.UtilisateurID,
	})
}
