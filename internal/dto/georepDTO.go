package dto

import (
	"time"

	"github.com/cabitibaly/internal/models"
)

type GeorepDTO struct {
	Id         int       `json:"id"`
	NomSite    string    `json:"site"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Rayon      float64   `json:"rayon"`
	HeureDebut time.Time `json:"heureDebut"`
	HeureFin   time.Time `json:"heureFin"`
}

type GeorepResponseDTO struct {
	Id         int     `json:"id"`
	NomSite    string  `json:"site"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Rayon      float64 `json:"rayon"`
	HeureDebut string  `json:"heureDebut"`
	HeureFin   string  `json:"heureFin"`
}

func ToGeorepDTO(georep *models.Georeperage) GeorepResponseDTO {
	return GeorepResponseDTO{
		Id:         georep.IDGeoreperage,
		NomSite:    georep.NomSite,
		Latitude:   georep.Latitude,
		Longitude:  georep.Longitude,
		Rayon:      georep.Rayon,
		HeureDebut: georep.HeureDebut.Format("15:04"),
		HeureFin:   georep.HeureFin.Format("15:04"),
	}
}

func ToGeorepDTOList(georeps []models.Georeperage) []GeorepResponseDTO {
	dtos := make([]GeorepResponseDTO, len(georeps))

	for i, georep := range georeps {
		dtos[i] = ToGeorepDTO(&georep)
	}

	return dtos
}
