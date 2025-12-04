package models

import "time"

type StatutConge struct {
	IDStatutConge          int       `gorm:"primaryKey" json:"id"`
	LibelleStatut          string    `gorm:"type:varchar(15)" json:"libelle"`
	DateCreationSatatut    time.Time `gorm:"autoCreateTime" json:"dateCreation"`
	DateModificationStatut time.Time `gorm:"autoUpdateTime" json:"dateModification"`
}
