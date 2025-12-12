package repositories

import (
	"github.com/cabitibaly/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(&token).Error
}

func (r *RefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var jwt models.RefreshToken

	err := r.db.Preload("Utilisateur").Where("token = ?", token).First(&jwt).Error

	if err != nil {
		return nil, err
	}

	return &jwt, nil
}

func (r *RefreshTokenRepository) FindByTokenAndUtilisateurID(token string, utilisateurID uint) (*models.RefreshToken, error) {
	var jwt models.RefreshToken

	err := r.db.Where("token = ? AND utilisateur_id = ?", token, utilisateurID).First(&jwt).Error

	if err != nil {
		return nil, err
	}

	return &jwt, nil
}

func (r *RefreshTokenRepository) Delete(id uint) error {
	return r.db.Delete(&models.RefreshToken{}, id).Error
}

func (r *RefreshTokenRepository) DeleteByToken(token string) error {
	var jwt models.RefreshToken

	err := r.db.Where("token = ?", token).First(&jwt).Error

	if err != nil {
		return err
	}

	return r.db.Delete(&jwt).Error
}
