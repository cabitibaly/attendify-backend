package dto

import (
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
	SiteID                   int    `json:"siteID"`
}

type UtilisateurResponseDTO struct {
	Id                       int     `json:"id"`
	Image                    string  `json:"image"`
	Nom                      string  `json:"nom"`
	Prenom                   string  `json:"prenom"`
	Email                    string  `json:"email"`
	Poste                    string  `json:"poste"`
	Telephone                string  `json:"telephone"`
	NombreConge              int     `json:"nombreConge"`
	SoldeConge               int     `json:"soldeConge"`
	MotdepasseAReinitialiser bool    `json:"motdepasseAReinitialiser"`
	Site                     SiteDTO `json:"site"`
}

type UtilisateurResponse struct {
	Utilisateur UtilisateurResponseDTO `json:"utilisateur"`
	Status      int                    `json:"status"`
}

type UtilisateursResponse struct {
	Utilisateurs []UtilisateurResponseDTO `json:"utilisateurs"`
	Pagination
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

	if utilisateur.Site != nil {
		dto.Site = SiteDTO{
			Id:        utilisateur.Site.IDSite,
			NomSite:   utilisateur.Site.NomSite,
			Latitude:  utilisateur.Site.Latitude,
			Longitude: utilisateur.Site.Longitude,
			Rayon:     utilisateur.Site.Rayon,
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
