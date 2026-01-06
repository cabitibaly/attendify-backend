package services

import (
	"fmt"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
)

type SiteService struct {
	siteRepo *repositories.SiteRepository
}

func NewSiteService(siteRepo *repositories.SiteRepository) *SiteService {
	return &SiteService{siteRepo: siteRepo}
}

func (s *SiteService) CreerUnSite(SiteDTO dto.SiteDTO) error {
	if SiteDTO.NomSite == "" ||
		SiteDTO.Latitude == 0 ||
		SiteDTO.Longitude == 0 ||
		SiteDTO.Rayon == 0 {
		return fmt.Errorf("tous les champs sont obligatoires")
	}

	if _, err := s.siteRepo.FindByNomSite(SiteDTO.NomSite); err == nil {
		return fmt.Errorf("ce site existe déjà")
	}

	georep := models.Site{
		NomSite:    SiteDTO.NomSite,
		Latitude:   SiteDTO.Latitude,
		Longitude:  SiteDTO.Longitude,
		Rayon:      SiteDTO.Rayon,
		HeureDebut: SiteDTO.HeureDebut,
		HeureFin:   SiteDTO.HeureFin,
	}

	return s.siteRepo.Create(&georep)
}

func (s *SiteService) LireUnSite(id uint) (*models.Site, error) {
	return s.siteRepo.FindByID(id)
}

func (s *SiteService) TousLesSites(recherche string, page, limit int) ([]models.Site, bool, int64, error) {
	return s.siteRepo.FindAll(recherche, page, limit)
}

func (s *SiteService) ModifierUnSite(id uint, data map[string]any) error {
	return s.siteRepo.Update(id, data)
}

func (s *SiteService) SupprimerUnSite(id uint) error {
	return s.siteRepo.Delete(id)
}
