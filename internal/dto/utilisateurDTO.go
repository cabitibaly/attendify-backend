package dto

import (
	"log"

	"github.com/cabitibaly/internal/models"
)

type UtilisateurDTO struct {
	Nom                      string `json:"nom"`
	Image                    string `json:"image"`
	Prenom                   string `json:"prenom"`
	Email                    string `json:"email"`
	Poste                    string `json:"poste"`
	NombreConge              int    `json:"nombreConge"`
	Telephone                string `json:"telephone"`
	MotDePasse               string `json:"motDePasse"`
	MotdepasseAReinitialiser bool   `json:"motdepasseAReinitialiser"`
	GeoreperageID            int    `json:"georeperageID"`
}

type UtilisateurResponseDTO struct {
	Id                       int       `json:"id"`
	Image                    string    `json:"image"`
	Nom                      string    `json:"nom"`
	Prenom                   string    `json:"prenom"`
	Email                    string    `json:"email"`
	Poste                    string    `json:"poste"`
	Telephone                string    `json:"telephone"`
	NombreConge              int       `json:"nombreConge"`
	SoldeConge               int       `json:"soldeConge"`
	MotdepasseAReinitialiser bool      `json:"motdepasseAReinitialiser"`
	Georeperage              GeorepDTO `json:"georeperage"`
}

func ToUtilisateurResponseDTO(utilisateur *models.Utilisateur) UtilisateurResponseDTO {
	dto := UtilisateurResponseDTO{
		Id:                       utilisateur.IDUtilisateur,
		Image:                    utilisateur.Image,
		Nom:                      utilisateur.Nom,
		Prenom:                   utilisateur.Prenom,
		Email:                    utilisateur.Email,
		Poste:                    utilisateur.Poste,
		Telephone:                utilisateur.Telephone,
		NombreConge:              utilisateur.NombreConge,
		SoldeConge:               utilisateur.SoldeConge,
		MotdepasseAReinitialiser: utilisateur.MotdepasseAReinitialiser,
	}

	log.Println("georep:", utilisateur.Georeperage)

	if utilisateur.Georeperage != nil {
		dto.Georeperage = GeorepDTO{
			NomSite:   utilisateur.Georeperage.NomSite,
			Latitude:  utilisateur.Georeperage.Latitude,
			Longitude: utilisateur.Georeperage.Longitude,
			Rayon:     utilisateur.Georeperage.Rayon,
		}
	}

	return dto
}

func ToUtilisateurResponseDTOList(utilisateurs []models.Utilisateur) []UtilisateurResponseDTO {
	dtos := make([]UtilisateurResponseDTO, len(utilisateurs))
	for i, utilisateur := range utilisateurs {
		dtos[i] = ToUtilisateurResponseDTO(&utilisateur)
	}
	return dtos
}
