package repositories

import (
	"github.com/cabitibaly/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(&notification).Error
}

func (r *NotificationRepository) FindAll(utilisateurID uint, page, limit int) ([]models.Notification, bool, error) {
	var notifications []models.Notification
	var total int64
	r.db.Model(&models.Notification{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Where("utilisateur_id = ?", utilisateurID).Find(&notifications).Error
	if err != nil {
		return nil, false, err
	}

	hasNextPage := int64(page*limit) <= total

	return notifications, hasNextPage, nil
}

func (r *NotificationRepository) FindByID(id uint) (*models.Notification, error) {
	var notif models.Notification

	err := r.db.Where(&notif, id).Error
	if err != nil {
		return nil, err
	}

	return &notif, nil
}

func (r *NotificationRepository) Update(id uint, data map[string]any) error {
	return r.db.Model(&models.Notification{}).Where("id_notification = ?", id).Updates(data).Error
}

func (r *NotificationRepository) Delele(id uint) error {
	return r.db.Delete(&models.Notification{}, id).Error
}
