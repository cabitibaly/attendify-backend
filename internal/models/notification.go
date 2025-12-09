package models

import "time"

type Notification struct {
	IDNotification           int       `gorm:"primaryKey" json:"id"`
	EstLu                    bool      `gorm:"type:bool;default:false" json:"estLu"`
	TypeNotification         string    `gorm:"type:varchar(50)" json:"typeNotification"`
	Titre                    string    `gorm:"type:varchar(255)" json:"titre"`
	Message                  string    `gorm:"type:text" json:"message"`
	DateCreationNotification time.Time `gorm:"autoCreateTime" json:"dateCreation"`
	UtilisateurID            int       `gorm:"type:int" json:"utilisateurID"`

	// Relations
	Utilisateur *Utilisateur `gorm:"foreignKey:UtilisateurID;references:IDUtilisateur;constraint:OnDelete:CASCADE" json:"utilisateur"`
}
