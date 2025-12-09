package models

import "time"

type Notification struct {
	IDNoficiation          int       `gorm:"primaryKey" json:"id"`
	EstLu                  bool      `gorm:"type:bool" json:"estLu"`
	TypeNoficcation        string    `gorm:"type:varchar(50)" json:"typeNoficcation"`
	Titre                  string    `gorm:"type:varchar(255)" json:"titre"`
	Message                string    `gorm:"type:text" json:"message"`
	DateCreationNofication time.Time `gorm:"autoCreateTime" json:"dateCreation"`
	UtilisateurID          int       `gorm:"type:int" json:"utilisateurID"`

	// Relations
	Utilisateur *Utilisateur `gorm:"foreignKey:UtilisateurID;references:IDUtilisateur;constraint:OnDelete:CASCADE" json:"utilisateur"`
}
