package models

import "time"

type Georeperage struct {
	IDGeoreperage          int       `gorm:"primaryKey" json:"id"`
	NomSite                string    `gorm:"type:varchar(255);not null;uniqueIndex;" json:"position"`
	Latitude               float64   `gorm:"type:double;not null;" json:"latitude"`
	Longitude              float64   `gorm:"type:double;not null;" json:"longitude"`
	Rayon                  float64   `gorm:"type:double;not null;" json:"rayon"`
	DateCreationGeorep     time.Time `gorm:"autoCreateTime" json:"dateCreation"`
	DateModificationGeorep time.Time `gorm:"autoUpdateTime" json:"dateModification"`
}

func (g *Georeperage) TableName() string {
	return "georeperage"
}
