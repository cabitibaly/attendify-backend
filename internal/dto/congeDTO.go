package dto

import (
	"time"

	"github.com/cabitibaly/internal/models"
)

type CongeDTO struct {
	Id             int       `json:"id"`
	DateDepart     time.Time `json:"dateDepart"`
	DateRetour     time.Time `json:"dateRetour"`
	Raison         string    `json:"raison"`
	TypeConge      string    `json:"typeConge"`
	PieceJointe    string    `json:"pieceJointe"`
	PieceJointeURL string    `json:"pieceJointeURL"`
	StatutConge    string    `json:"statutConge"`
}

type CongeUtilisateurDTO struct {
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`
	Image  string `json:"image"`
	Poste  string `json:"poste"`
}

type CongeResponseAdminDTO struct {
	CongeDTO
	Utilisateur CongeUtilisateurDTO `json:"utilisateur"`
}

type CongeResponseEmp struct {
	Conge  any `json:"conge"`
	Status int `json:"status"`
}

type CongesResponseEmp struct {
	Conges []CongeDTO `json:"conges"`
	Pagination
}

type CongesResponseAdmin struct {
	Conges []CongeResponseAdminDTO `json:"conges"`
	Pagination
}

type CongeResponseAdmin struct {
	Conge  any `json:"conge"`
	Status int `json:"status"`
}

func ToCongeResponseAdminDTO(conge *models.Conge) CongeResponseAdminDTO {
	dto := CongeResponseAdminDTO{
		CongeDTO: CongeDTO{
			Id:             conge.IDConge,
			DateDepart:     conge.DateDepart,
			DateRetour:     conge.DateRetour,
			Raison:         conge.Raison,
			TypeConge:      conge.TypeConge,
			PieceJointe:    conge.PieceJointe,
			PieceJointeURL: conge.PieceJointeURL,
		},
	}

	if conge.StatutConge != nil {
		dto.StatutConge = conge.StatutConge.LibelleStatut
	}

	if conge.Utilisateur != nil {
		dto.Utilisateur = CongeUtilisateurDTO{
			Nom:    conge.Utilisateur.Nom,
			Prenom: conge.Utilisateur.Prenom,
			Image:  conge.Utilisateur.Image,
			Poste:  conge.Utilisateur.Poste,
		}
	}

	return dto
}

func ToCongeResponseEmpDTO(conge *models.Conge) CongeDTO {
	dto := CongeDTO{
		Id:             conge.IDConge,
		DateDepart:     conge.DateDepart,
		DateRetour:     conge.DateRetour,
		Raison:         conge.Raison,
		TypeConge:      conge.TypeConge,
		PieceJointe:    conge.PieceJointe,
		PieceJointeURL: conge.PieceJointeURL,
	}

	if conge.StatutConge != nil {
		dto.StatutConge = conge.StatutConge.LibelleStatut
	}

	return dto
}

func ToCongeResponseAdminDTOList(conges []models.Conge) []CongeResponseAdminDTO {
	dtos := make([]CongeResponseAdminDTO, len(conges))

	for i, conge := range conges {
		dtos[i] = ToCongeResponseAdminDTO(&conge)
	}

	return dtos
}

func ToCongeResponseEmpDTOList(conges []models.Conge) []CongeDTO {
	dtos := make([]CongeDTO, len(conges))

	for i, conge := range conges {
		dtos[i] = ToCongeResponseEmpDTO(&conge)
	}

	return dtos
}
