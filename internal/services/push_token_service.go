package services

import (
	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
	"github.com/cabitibaly/pkg/utils"
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
		DeviceName:    pushTokenDTO.DeviceName,
		Platform:      pushTokenDTO.Platform,
		UtilisateurID: pushTokenDTO.UtilisateurID,
	})
}

func (s *PushTokenService) EnvoyerNotificationPushAUnUtilisateur(utilisateurID uint, title, message string) error {
	tokens, err := s.repo.FindActivePushTokenByUtilisateurID(utilisateurID)
	if err != nil || len(tokens) == 0 {
		return err
	}

	for _, token := range tokens {
		if err := utils.EnvoyerNotificationPush(token.PushToken, title, message); err != nil {
			s.repo.DesactivePushToken(token.PushToken)
		}
	}
	return nil
}

func (s *PushTokenService) SupprimerUnPushToken(pushToken string) error {
	return s.repo.Delete(pushToken)
}
