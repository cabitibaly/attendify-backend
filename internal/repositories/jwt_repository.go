package repositories

import (
	"github.com/cabitibaly/internal/models"
	"gorm.io/gorm"
)

type JWTRepository struct {
	db *gorm.DB
}

func NewJWTRepository(db *gorm.DB) *JWTRepository {
	return &JWTRepository{db: db}
}

func (r *JWTRepository) Create(token *models.Jwt) error {
	return r.db.Create(&token).Error
}

func (r *JWTRepository) FindByToken(token string) (*models.Jwt, error) {
	var jwt models.Jwt

	err := r.db.Where("token = ?", token).First(&jwt).Error

	if err != nil {
		return nil, err
	}

	return &jwt, nil
}

func (r *JWTRepository) FindByTokenAndUtilisateurID(token string, utilisateurID uint) (*models.Jwt, error) {
	var jwt models.Jwt

	err := r.db.Where("token = ? AND utilisateur_id = ?", token, utilisateurID).First(&jwt).Error

	if err != nil {
		return nil, err
	}

	return &jwt, nil
}

func (r *JWTRepository) Delete(id uint) error {
	return r.db.Delete(&models.Jwt{}, id).Error
}
