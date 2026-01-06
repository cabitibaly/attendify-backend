package repositories

import (
	"time"

	"github.com/cabitibaly/internal/models"
	"gorm.io/gorm"
)

type PushTokenRepository struct {
	db *gorm.DB
}

func NewPushTokenRepository(db *gorm.DB) *PushTokenRepository {
	return &PushTokenRepository{db: db}
}

func (r *PushTokenRepository) CreateOrUpdate(pushToken *models.PushToken) error {
	var pushTokenExiste models.PushToken

	err := r.db.Where("utilisateur_id = ? AND push_token = ?", pushToken.UtilisateurID, pushToken.PushToken).First(&pushTokenExiste).Error

	if err == nil {
		pushTokenExiste.PushToken = pushToken.PushToken
		pushTokenExiste.DeviceType = pushToken.DeviceType
		pushTokenExiste.DeviceName = pushToken.DeviceName
		pushTokenExiste.Platform = pushToken.Platform
		pushTokenExiste.EstActive = true
		pushTokenExiste.DateModificationPushToken = time.Now()

		return r.db.Save(&pushTokenExiste).Error
	}

	return r.db.Create(&pushToken).Error
}

func (r *PushTokenRepository) FindActivePushTokenByUtilisateurID(utilisateurID uint) ([]models.PushToken, error) {
	var pushTokens []models.PushToken

	err := r.db.Where("utilisateurID = ? AND est_active = ?", utilisateurID, 1).Find(pushTokens).Error
	return pushTokens, err
}

func (r *PushTokenRepository) DesactivePushToken(token string) error {
	return r.db.Model(&models.PushToken{}).Where("push_token = ?", token).Update("est_active", false).Error
}

func (r *PushTokenRepository) Delete(token string) error {
	return r.db.Where("push_token = ?", token).Delete(&models.PushToken{}).Error
}
