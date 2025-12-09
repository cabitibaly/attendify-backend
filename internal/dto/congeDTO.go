package dto

import "time"

type CongeDTO struct {
	DateDepart  time.Time `json:"dateDepart"`
	DateRetour  time.Time `json:"dateRetour"`
	Raison      string    `json:"raison"`
	TypeConge   string    `json:"typeConge"`
	PieceJointe string    `json:"pieceJointe"`
	StatutConge string    `json:"statutConge"`
}
