package models

import "time"

type RefreshToken struct {
	ID                  int        `gorm:"primarykey" json:"id"`
	Token               string     `gorm:"size:255;not null" json:"token"`
	UtilisateurID       uint       `gorm:"not null" json:"utilisateurID"`
	ExpireAt            time.Time  `gorm:"not null" json:"expireA"`
	RevokedAt           *time.Time `gorm:"default:null" json:"revoked_at"`
	DateCreationJwt     time.Time  `gorm:"autoCreateTime" json:"dateCreation"`
	DateModificationJwt time.Time  `gorm:"autoUpdateTime" json:"dateModification"`

	Utilisateur Utilisateur `gorm:"foreignKey:UtilisateurID;references:IDUtilisateur;constraint:OnDelete:CASCADE" json:"-"`
}

func (j *RefreshToken) TableName() string {
	return "refresh_token"
}
