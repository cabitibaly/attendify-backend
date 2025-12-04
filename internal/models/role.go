package models

import "time"

type Role struct {
	IDRole               int       `gorm:"primaryKey" json:"id"`
	Libelle              string    `gorm:"type:varchar(15)" json:"libelle"`
	DateCreationRole     time.Time `gorm:"autoCreateTime" json:"dateCreation"`
	DateModificationRole time.Time `gorm:"autoUpdateTime" json:"dateModification"`
}

func (r *Role) TableName() string {
	return "role"
}
