package models

import "time"

type Conge struct {
	IDConge               int       `gorm:"primaryKey" json:"id"`
	DateDepart            time.Time `gorm:"type:timestamp" json:"dateDepart"`
	DateRetour            time.Time `gorm:"type:timestamp" json:"dateRetour"`
	Raison                string    `gorm:"type:text" json:"raison"`
	TypeConge             string    `gorm:"type:varchar(50)" json:"typeConge"`
	PieceJointe           string    `gorm:"type:varchar(255)" json:"pieceJointe"`
	UtilisateurID         int       `gorm:"type:int" json:"utilisateurID"`
	StatutCongeID         int       `gorm:"type:int" json:"statutCongeID"`
	DateCreationConge     time.Time `gorm:"autoCreateTime" json:"dateCreation"`
	DateModificationConge time.Time `gorm:"autoUpdateTime" json:"dateModification"`

	// Relations
	Utilisateur *Utilisateur `gorm:"foreignKey:UtilisateurID;references:IDUtilisateur;constraint:OnDelete:CASCADE" json:"-"`
	StatutConge *StatutConge `gorm:"foreignKey:StatutCongeID;references:IDStatutConge" json:"statutConge"`
}

func (c *Conge) TableName() string {
	return "conge"
}
