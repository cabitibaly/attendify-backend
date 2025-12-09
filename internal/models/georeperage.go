package models

import "time"

type Georeperage struct {
	IDGeoreperage          int       `gorm:"primaryKey" json:"id"`
	NomSite                string    `gorm:"type:varchar(255);not null;uniqueIndex;" json:"site"`
	Latitude               float64   `gorm:"type:double;not null;" json:"latitude"`
	Longitude              float64   `gorm:"type:double;not null;" json:"longitude"`
	Rayon                  float64   `gorm:"type:double;not null;" json:"rayon"`
	HeureDebut             time.Time `gorm:"type:time;default:08:00:00;" json:"heureDebut"`
	HeureFin               time.Time `gorm:"type:time;default:16:00:00;" json:"heureFin"`
	DateCreationGeorep     time.Time `gorm:"autoCreateTime" json:"dateCreation"`
	DateModificationGeorep time.Time `gorm:"autoUpdateTime" json:"dateModification"`
}

func (g *Georeperage) TableName() string {
	return "georeperage"
}
