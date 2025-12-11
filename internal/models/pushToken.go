package models

import "time"

type PushToken struct {
	IDPushToken               int       `gorm:"primary_key" json:"id"`
	PushToken                 string    `json:"push_token"`
	DeviceType                string    `gorm:"type:varchar(20)" json:"device_type"`
	EstActive                 bool      `gorm:"default:true" json:"est_active"`
	UtilisateurID             int       `gorm:"type:int" json:"utilisateurID"`
	DateCreationPushToken     time.Time `gorm:"autoCreateTime" json:"date_creation"`
	DateModificationPushToken time.Time `gorm:"autoUpdateTime" json:"date_modification"`
}

func (p *PushToken) TableName() string {
	return "push_token"
}
