package models

import "time"

type Site struct {
	IDSite               int       `gorm:"primaryKey" json:"id"`
	NomSite              string    `gorm:"type:varchar(255);not null;uniqueIndex;" json:"site"`
	Latitude             float64   `gorm:"type:double;not null;" json:"latitude"`
	Longitude            float64   `gorm:"type:double;not null;" json:"longitude"`
	Rayon                float64   `gorm:"type:double;not null;" json:"rayon"`
	HeureDebut           time.Time `gorm:"type:time;default:08:00:00;" json:"heureDebut"`
	HeureFin             time.Time `gorm:"type:time;default:16:00:00;" json:"heureFin"`
	DateCreationSite     time.Time `gorm:"autoCreateTime" json:"dateCreation"`
	DateModificationSite time.Time `gorm:"autoUpdateTime" json:"dateModification"`
}

func (g *Site) TableName() string {
	return "site"
}
