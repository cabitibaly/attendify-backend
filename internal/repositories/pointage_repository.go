package repositories

import (
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

	db := r.db.Model(&models.Pointage{}).Preload("Utilisateur")

	location, errLoc := time.LoadLocation("Africa/Ouagadougou")

	if errLoc != nil {
		panic(errLoc)
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

func (r *PointageRepository) FindOneByDate(utilisateurID uint, date time.Time) (models.Pointage, error) {
	var pointage models.Pointage

	location, errLoc := time.LoadLocation("Africa/Ouagadougou")

	if errLoc != nil {
		panic(errLoc)
	}

	debutJournee := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		location,
	)

	finJournee := debutJournee.Add(24 * time.Hour)

	err := r.db.Where("arrive >= ? AND arrive < ? AND utilisateur_id = ?", debutJournee, finJournee, utilisateurID).First(&pointage).Error

	if err != nil {
		return pointage, err
	}

	return pointage, nil

}

func (r *PointageRepository) Update(id uint, data map[string]any) error {
	return r.db.Model(&models.Pointage{}).Where("id_pointage = ?", id).Updates(data).Error
}

func (r *PointageRepository) Delete(id uint) error {
	return r.db.Delete(&models.Pointage{}, id).Error
}
