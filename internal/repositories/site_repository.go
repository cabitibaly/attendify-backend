package repositories

import (
	"github.com/cabitibaly/internal/models"
	"gorm.io/gorm"
)

type SiteRepository struct {
	db *gorm.DB
}

func NewSiteRepository(db *gorm.DB) *SiteRepository {
	return &SiteRepository{db: db}
}

func (r *SiteRepository) Create(site *models.Site) error {
	return r.db.Create(&site).Error
}

func (r *SiteRepository) FindByID(id uint) (*models.Site, error) {
	var site models.Site

	err := r.db.First(&site, id).Error

	if err != nil {
		return nil, err
	}

	return &site, nil
}

func (r *SiteRepository) FindByNomSite(nomSite string) (*models.Site, error) {
	var site models.Site

	err := r.db.Where("site = ?", nomSite).First(&site).Error

	if err != nil {
		return nil, err
	}

	return &site, nil
}

func (r *SiteRepository) FindAll(recherche string, page, limit int) ([]models.Site, bool, int64, error) {
	var sites []models.Site
	var total int64

	db := r.db.Model(&models.Site{})

	if recherche != "" {
		recherche = "%" + recherche + "%"
		db = db.Where("nom_site LIKE ?", recherche)
	}

	db.Count(&total)

	offset := (page - 1) * limit
	err := db.Offset(offset).Limit(limit).Order("id_site DESC").Find(&sites).Error

	hasNextPage := int64(limit*page) <= total

	return sites, hasNextPage, total, err
}

func (r *SiteRepository) Update(id uint, data map[string]any) error {
	return r.db.Model(&models.Site{}).Where("id_site = ?", id).Updates(data).Error
}

func (r *SiteRepository) Delete(id uint) error {
	return r.db.Delete(&models.Site{}, id).Error
}
