package dto

import (
	"time"

	"github.com/cabitibaly/internal/models"
)

type PointageUtilisateurDTO struct {
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`
	Image  string `json:"image"`
}

type PointageResponseDTO struct {
	ID                int                    `json:"id"`
	Arrive            time.Time              `json:"arrive"`
	Depart            *time.Time             `json:"depart"`
	HeuresTravaillees *float64               `json:"heuresTravaillees"`
	EstPresent        bool                   `json:"estPresent"`
	EnRetard          bool                   `json:"enRetard"`
	DepartAnticipe    bool                   `json:"departAnticipe"`
	Utilisateur       PointageUtilisateurDTO `json:"utilisateur"`
}

type PointagesResponse struct {
	Pointages []PointageResponseDTO `json:"pointages"`
	Pagination
}

func ToPointageResponseDTO(pointage *models.Pointage) PointageResponseDTO {
	dto := PointageResponseDTO{
		ID:                pointage.IDPointage,
		Arrive:            pointage.Arrive,
		Depart:            pointage.Depart,
		HeuresTravaillees: pointage.HeuresTravaillees,
		EstPresent:        pointage.EstPresent,
		EnRetard:          pointage.EnRetard,
		DepartAnticipe:    pointage.DepartAnticipe,
	}

	if pointage.Utilisateur != nil {
		dto.Utilisateur = PointageUtilisateurDTO{
			Nom:    pointage.Utilisateur.Nom,
			Prenom: pointage.Utilisateur.Prenom,
			Image:  pointage.Utilisateur.Image,
		}
	}

	return dto
}

func ToPointageResponseDTOList(pointages []models.Pointage) []PointageResponseDTO {
	dtos := make([]PointageResponseDTO, len(pointages))
	for i, pointage := range pointages {
		dtos[i] = ToPointageResponseDTO(&pointage)
	}
	return dtos
}
