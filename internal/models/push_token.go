package models

import "time"

type PushToken struct {
	IDPushToken               int       `gorm:"primary_key" json:"id"`
	PushToken                 string    `gorm:"unique" json:"push_token"`
	DeviceType                string    `gorm:"type:varchar(20)" json:"device_type"`
	DeviceName                string    `gorm:"type:varchar(100)" json:"device_name"`
	Platform                  string    `gorm:"type:varchar(20)" json:"platform"`
	EstActive                 bool      `gorm:"default:true" json:"est_active"`
	UtilisateurID             int       `gorm:"type:int" json:"utilisateur_id"`
	DateCreationPushToken     time.Time `gorm:"autoCreateTime" json:"date_creation"`
	DateModificationPushToken time.Time `gorm:"autoUpdateTime" json:"date_modification"`
}

func (p *PushToken) TableName() string {
	return "push_token"
}
