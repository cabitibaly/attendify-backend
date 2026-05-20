package repositories

import (
	"strings"

	"github.com/cabitibaly/internal/models"
	"gorm.io/gorm"
)

type UtilisateurRepository struct {
	db *gorm.DB
}

func NewUtilisateurRepository(db *gorm.DB) *UtilisateurRepository {
	return &UtilisateurRepository{db: db}
}

func (r *UtilisateurRepository) Create(utilsateur *models.Utilisateur) error {
	return r.db.Create(&utilsateur).Error
}

func (r *UtilisateurRepository) FindByID(id uint) (*models.Utilisateur, error) {
	var utilisateur models.Utilisateur

	err := r.db.Preload("Site").First(&utilisateur, id).Error

	if err != nil {
		return nil, err
	}

	return &utilisateur, nil
}

func (r *UtilisateurRepository) FindByEmail(email string) (*models.Utilisateur, error) {
	var utilisateur models.Utilisateur

	err := r.db.Where("email = ?", email).First(&utilisateur).Error

	if err != nil {
		return nil, err
	}

	return &utilisateur, nil
}

func (r *UtilisateurRepository) FindByTelephone(telephone string) (*models.Utilisateur, error) {
	var utilisateur models.Utilisateur

	err := r.db.Where("telephone = ?", telephone).First(&utilisateur).Error

	if err != nil {
		return nil, err
	}

	return &utilisateur, nil
}

func (r *UtilisateurRepository) FindByEmailOrTelephone(username string) (*models.Utilisateur, error) {
	var utilisateur models.Utilisateur

	err := r.db.Where("email = ? OR telephone = ?", username, username).First(&utilisateur).Error

	if err != nil {
		return nil, err
	}

	return &utilisateur, nil
}

func (r *UtilisateurRepository) FindAll(recherche string, page, limt int) ([]models.Utilisateur, bool, int64, error) {
	var utilisateurs []models.Utilisateur
	var total int64

	db := r.db.Model(&models.Utilisateur{}).Where("role_id = ?", 2).Order("id_utilisateur DESC")

	if recherche != "" {
		mots := strings.Fields(recherche)

		for _, mot := range mots {
			like := "%" + mot + "%"

			db = db.Where("nom LIKE ? OR prenom LIKE ?", like, like)
		}
	}

	db.Count(&total)

	offset := (page - 1) * limt
	err := db.Offset(offset).Limit(limt).Preload("Site").Find(&utilisateurs).Error

	hasNextPage := int64(limt*page) <= total

	return utilisateurs, hasNextPage, total, err
}

func (r *UtilisateurRepository) Update(id uint, data map[string]any) error {
	return r.db.Model(&models.Utilisateur{}).Where("id_utilisateur = ?", id).Updates(data).Error
}

func (r *UtilisateurRepository) Delete(id uint) error {
	return r.db.Delete(&models.Utilisateur{}, id).Error
}

func (r *UtilisateurRepository) GetSoldeConge(id uint) (int, error) {
	var utilisateur models.Utilisateur

	err := r.db.Where("id_utilisateur = ?", id).First(&utilisateur).Error

	if err != nil {
		return 0, err
	}

	return utilisateur.SoldeConge, nil
}

func (r *UtilisateurRepository) GetTotalEmploye() int64 {
	var total int64
	r.db.Model(&models.Utilisateur{}).Where("role_id = ?", 2).Count(&total)
	return total
}
