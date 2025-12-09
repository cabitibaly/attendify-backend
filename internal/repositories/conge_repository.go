package repositories

import (
	"time"

	"github.com/cabitibaly/internal/models"
	"gorm.io/gorm"
)

type CongeRepository struct {
	db *gorm.DB
}

func NewCongeRepository(db *gorm.DB) *CongeRepository {
	return &CongeRepository{db: db}
}

func (r *CongeRepository) Create(conge *models.Conge) error {
	return r.db.Create(&conge).Error
}

func (r *CongeRepository) FindByID(id uint) (*models.Conge, error) {
	var conge models.Conge

	err := r.db.Preload("Utilisateur").Preload("StatutConge").First(&conge, id).Error
	if err != nil {
		return nil, err
	}

	return &conge, nil
}

func (r *CongeRepository) FindAll(utilisateurID uint, statutID uint, page, limit int) ([]models.Conge, bool, int64, error) {
	var conges []models.Conge
	var total int64

	db := r.db.Model(&models.Conge{}).Preload("Utilisateur").Preload("StatutConge").Order("id_conge DESC")

	if utilisateurID != 0 {
		db.Where("utilisateur_id = ?", utilisateurID)
	}

	if statutID != 0 {
		db.Where("statut_conge_id = ?", statutID)
	}

	db.Model(&models.Conge{}).Count(&total)

	offset := (page - 1) * limit
	err := db.Offset(offset).Limit(limit).Find(&conges).Error
	if err != nil {
		return nil, false, 0, err
	}

	hasNextPage := int64(limit*page) <= total
	return conges, hasNextPage, total, nil
}

func (r *CongeRepository) FindByUtilisateurIDAndPeriode(utilisateurID uint, debutJour, finJour time.Time) ([]models.Conge, error) {
	var conges []models.Conge

	err := r.db.Where("utilisateur_id = ? AND date_depart <= ? AND date_retour >= ?",
		utilisateurID,
		finJour,
		debutJour,
	).Find(&conges).Error

	if err != nil {
		return nil, err
	}

	return conges, nil
}

func (r *CongeRepository) FindByUtilisateurAndAnnee(utilisateurID uint, location *time.Location) ([]models.Conge, error) {
	var conges []models.Conge

	maintenant := time.Now().In(location)
	debutAnnee := time.Date(maintenant.Year(), time.January, 1, 0, 0, 0, 0, location)
	finAnnee := time.Date(maintenant.Year()+1, time.January, 1, 0, 0, 0, 0, location)

	err := r.db.Where(
		"utilisateur_id = ? AND date_depart < ? AND date_retour >= ?",
		utilisateurID,
		finAnnee,
		debutAnnee,
	).Find(&conges).Error

	if err != nil {
		return nil, err
	}

	return conges, nil
}

func (r *CongeRepository) Update(id uint, data map[string]any) error {
	return r.db.Model(&models.Conge{}).Where("id_conge = ?", id).Updates(data).Error
}

func (r *CongeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Conge{}, id).Error
}
