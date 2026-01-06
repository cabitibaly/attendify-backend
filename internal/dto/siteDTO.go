package dto

import (
	"time"

	"github.com/cabitibaly/internal/models"
)

type SiteDTO struct {
	Id         int       `json:"id"`
	NomSite    string    `json:"site"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Rayon      float64   `json:"rayon"`
	HeureDebut time.Time `json:"heureDebut"`
	HeureFin   time.Time `json:"heureFin"`
}

type SiteResponseDTO struct {
	Id         int     `json:"id"`
	NomSite    string  `json:"site"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Rayon      float64 `json:"rayon"`
	HeureDebut string  `json:"heureDebut"`
	HeureFin   string  `json:"heureFin"`
}

type SiteResponse struct {
	Site   SiteResponseDTO `json:"site"`
	Status int             `json:"status"`
}

type SitesResponse struct {
	Sites []SiteResponseDTO `json:"sites"`
	Pagination
}

func ToSiteDTO(site *models.Site) SiteResponseDTO {
	return SiteResponseDTO{
		Id:         site.IDSite,
		NomSite:    site.NomSite,
		Latitude:   site.Latitude,
		Longitude:  site.Longitude,
		Rayon:      site.Rayon,
		HeureDebut: site.HeureDebut.Format("15:04"),
		HeureFin:   site.HeureFin.Format("15:04"),
	}
}

func ToSiteDTOList(georeps []models.Site) []SiteResponseDTO {
	dtos := make([]SiteResponseDTO, len(georeps))

	for i, site := range georeps {
		dtos[i] = ToSiteDTO(&site)
	}

	return dtos
}
