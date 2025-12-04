package services

import (
	"fmt"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
)

type GeorepService struct {
	georepRepo *repositories.GeorepRepository
}

func NewGeorepService(georepRepo *repositories.GeorepRepository) *GeorepService {
	return &GeorepService{georepRepo: georepRepo}
}

func (s *GeorepService) CreerUnSite(georepDTO dto.GeorepDTO) error {
	if georepDTO.NomSite == "" ||
		georepDTO.Latitude == 0 ||
		georepDTO.Longitude == 0 ||
		georepDTO.Rayon == 0 {
		return fmt.Errorf("tous les champs sont obligatoires")
	}

	if _, err := s.georepRepo.FindByNomSite(georepDTO.NomSite); err == nil {
		return fmt.Errorf("ce site existe déjà")
	}

	georep := models.Georeperage{
		NomSite:   georepDTO.NomSite,
		Latitude:  georepDTO.Latitude,
		Longitude: georepDTO.Longitude,
		Rayon:     georepDTO.Rayon,
	}

	return s.georepRepo.Create(&georep)
}

func (s *GeorepService) LireUnSite(id uint) (*models.Georeperage, error) {
	return s.georepRepo.FindByID(id)
}

func (s *GeorepService) TousLesSites(page, limit int) ([]models.Georeperage, bool, int64, error) {
	return s.georepRepo.FindAll(page, limit)
}

func (s *GeorepService) ModifierUnSite(id uint, data map[string]any) error {
	return s.georepRepo.Update(id, data)
}

func (s *GeorepService) SupprimerUnSite(id uint) error {
	return s.georepRepo.Delete(id)
}
