package repositories

import (
	"github.com/cabitibaly/internal/models"
	"gorm.io/gorm"
)

type GeorepRepository struct {
	db *gorm.DB
}

func NewGeorepRepository(db *gorm.DB) *GeorepRepository {
	return &GeorepRepository{db: db}
}

func (r *GeorepRepository) Create(georep *models.Georeperage) error {
	return r.db.Create(&georep).Error
}

func (r *GeorepRepository) FindByID(id uint) (*models.Georeperage, error) {
	var georep models.Georeperage

	err := r.db.First(&georep, id).Error

	if err != nil {
		return nil, err
	}

	return &georep, nil
}

func (r *GeorepRepository) FindByNomSite(nomSite string) (*models.Georeperage, error) {
	var georep models.Georeperage

	err := r.db.Where("nomSite = ?", nomSite).First(&georep).Error

	if err != nil {
		return nil, err
	}

	return &georep, nil
}

func (r *GeorepRepository) FindAll(page, limit int) ([]models.Georeperage, bool, int64, error) {
	var georeps []models.Georeperage
	var total int64

	r.db.Model(&models.Georeperage{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Find(&georeps).Error

	hasNextPage := int64(limit*page) <= total

	return georeps, hasNextPage, total, err
}

func (r *GeorepRepository) Update(id uint, data map[string]any) error {
	return r.db.Model(&models.Georeperage{}).Updates(data).Error
}

func (r *GeorepRepository) Delete(id uint) error {
	return r.db.Delete(&models.Georeperage{}, id).Error
}
