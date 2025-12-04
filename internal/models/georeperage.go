package models

import "time"

type Georeperage struct {
	IDGeoreperage          int       `gorm:"primaryKey" json:"id"`
	Position               string    `gorm:"type:varchar(255)" json:"position"`
	Latitude               float64   `gorm:"type:double" json:"latitude"`
	Longitude              float64   `gorm:"type:double" json:"longitude"`
	Rayon                  float64   `gorm:"type:double" json:"rayon"`
	DateCreationGeorep     time.Time `gorm:"autoCreateTime" json:"dateCreation"`
	DateModificationGeorep time.Time `gorm:"autoUpdateTime" json:"dateModification"`
}

func (g *Georeperage) TableName() string {
	return "georeperage"
}
