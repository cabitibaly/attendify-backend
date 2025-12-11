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

func (r *GeorepRepository) FindAll(recherche string, page, limit int) ([]models.Georeperage, bool, int64, error) {
	var georeps []models.Georeperage
	var total int64

	db := r.db.Model(&models.Georeperage{})

	if recherche != "" {
		recherche = "%" + recherche + "%"
		db = db.Where("nom_site LIKE ?", recherche)
	}

	db.Count(&total)

	offset := (page - 1) * limit
	err := db.Offset(offset).Limit(limit).Order("id_georeperage DESC").Find(&georeps).Error

	hasNextPage := int64(limit*page) <= total

	return georeps, hasNextPage, total, err
}

func (r *GeorepRepository) Update(id uint, data map[string]any) error {
	return r.db.Model(&models.Georeperage{}).Where("id_georeperage = ?", id).Updates(data).Error
}

func (r *GeorepRepository) Delete(id uint) error {
	return r.db.Delete(&models.Georeperage{}, id).Error
}
