package models

import "time"

type Pointage struct {
	IDPointage               int        `gorm:"primaryKey" json:"id"`
	Arrivee                  time.Time  `gorm:"autoCreateTime" json:"arrivee"`
	Depart                   *time.Time `gorm:"timestamp" json:"depart"`
	EstPresent               bool       `gorm:"default:false" json:"estPresent"`
	EnRetard                 bool       `gorm:"default:false" json:"EnRetard"`
	UtilisateurID            int        `gorm:"type:int" json:"utilisateurID"`
	DateModificationPointage time.Time  `gorm:"autoUpdateTime" json:"dateModification"`

	// Relations
	Utilisateur *Utilisateur `gorm:"foreignKey:UtilisateurID;references:IDUtilisateur;constraint:OnDelete:CASCADE" json:"-"`
}

func (p *Pointage) TableName() string {
	return "pointage"
}
