package repositories

import (
	"fmt"
	"time"

	"github.com/cabitibaly/internal/models"
	"gorm.io/gorm"
)

type PointageRepository struct {
	db *gorm.DB
}

func NewPointageRepository(db *gorm.DB) *PointageRepository {
	return &PointageRepository{db: db}
}

func (r *PointageRepository) Create(pointage *models.Pointage) error {
	return r.db.Create(&pointage).Error
}

func (r *PointageRepository) FindByID(id uint) (*models.Pointage, error) {
	var pointage models.Pointage

	err := r.db.First(&pointage, id).Error

	if err != nil {
		return nil, err
	}

	return &pointage, nil
}

func (r *PointageRepository) FindByUtilisateurID(utilisateurID uint) (*models.Pointage, error) {
	var pointage models.Pointage

	err := r.db.Where("utilisateur_id = ?", utilisateurID).Order("id_pointage DESC").First(&pointage).Error

	if err != nil {
		return nil, err
	}

	return &pointage, nil
}

func (r *PointageRepository) FindAll(utilisateurID uint, aujoudhui bool, date time.Time, page, limit int) ([]models.Pointage, bool, int64, error) {
	var pointages []models.Pointage
	var total int64

	db := r.db.Model(&models.Pointage{}).Preload("Utilisateur").Order("id_pointage DESC")

	location, errLoc := time.LoadLocation("Africa/Ouagadougou")

	if errLoc != nil {
		return nil, false, 0, fmt.Errorf("erreur de timezone: %w", errLoc)
	}

	if !date.IsZero() {
		debutJournee := time.Date(
			date.Year(),
			date.Month(),
			date.Day(),
			0, 0, 0, 0,
			location,
		)

		finJournee := debutJournee.Add(24 * time.Hour)

		db.Where("arrive >= ? AND arrive < ? ", debutJournee, finJournee)
		aujoudhui = false
	}

	if aujoudhui {
		debutJournee := time.Date(
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			0,
			0,
			0,
			0,
			location,
		)

		finJournee := debutJournee.Add(24 * time.Hour)

		db.Where("arrive >= ? AND arrive < ? ", debutJournee, finJournee)
	}

	if utilisateurID != 0 {
		db.Where("utilisateur_id = ?", utilisateurID)
	}

	db.Count(&total)

	offset := (page - 1) * limit

	err := db.Offset(offset).Limit(limit).Find(&pointages).Error

	hasNextPage := int64(limit*page) <= total

	return pointages, hasNextPage, total, err
}

func (r *PointageRepository) Update(id uint, data map[string]any) error {
	return r.db.Model(&models.Pointage{}).Where("id_pointage = ?", id).Updates(data).Error
}

func (r *PointageRepository) Delete(id uint) error {
	return r.db.Delete(&models.Pointage{}, id).Error
}

func (r *PointageRepository) GetTotalPresentAndRetard() (int64, int64, error) {
	location, errLoc := time.LoadLocation("Africa/Ouagadougou")

	if errLoc != nil {
		return 0, 0, fmt.Errorf("erreur de timezone: %w", errLoc)
	}

	aujoudhui := time.Now().In(location)
	debutJournee := time.Date(
		aujoudhui.Year(),
		aujoudhui.Month(),
		aujoudhui.Day(),
		0, 0, 0, 0,
		location,
	)

	finJournee := debutJournee.Add(24 * time.Hour)

	var totalPresent, totalRetard int64
	r.db.Model(&models.Pointage{}).Where("en_retard = true AND arrive >= ? AND arrive < ?", debutJournee, finJournee).Count(&totalRetard)
	r.db.Model(&models.Pointage{}).Where("est_present = true AND arrive >= ? AND arrive < ?", debutJournee, finJournee).Count(&totalPresent)

	return totalPresent, totalRetard, nil
}
