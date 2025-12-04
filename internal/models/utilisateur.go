package models

import "time"

type Utilisateur struct {
	IDUtilisateur               int       `gorm:"primaryKey" json:"id"`
	Nom                         string    `gorm:"type:varchar(20);not null" json:"nom"`
	Prenom                      string    `gorm:"type:varchar(50)" json:"prenom"`
	Telephone                   string    `gorm:"type:varchar(20);not null;uniqueIndex" json:"telephone"`
	Email                       string    `gorm:"type:varchar(150);not null;uniqueIndex" json:"email"`
	Poste                       string    `gorm:"type:varchar(50);not null" json:"poste"`
	MotdepasseAReinitialiser    bool      `gorm:"default:true" json:"motdepasseAReinitialiser"`
	MotDePasse                  string    `gorm:"type:varchar(100);not null" json:"-"`
	RoleID                      int       `gorm:"type:int" json:"roleID"`
	GeoreperageID               int       `gorm:"type:int" json:"georeperageID"`
	DateCreationUtilisateur     time.Time `gorm:"autoCreateTime" json:"dateCreation"`
	DateModificationUtilisateur time.Time `gorm:"autoUpdateTime" json:"dateModification"`

	// Relations
	Role        *Role        `gorm:"foreignKey:RoleID;references:IDRole" json:"-"`
	Georeperage *Georeperage `gorm:"foreignKey:GeoreperageID;references:IDGeoreperage;constraint:OnDelete:SET NULL;" json:"-"`
}

func (u *Utilisateur) TableName() string {
	return "utilisateur"
}
