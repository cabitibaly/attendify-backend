package dto

type PushTokenDTO struct {
	PushToken     string `json:"push_token"`
	DeviceType    string `json:"device_type"`
	EstActive     bool   `json:"est_active"`
	UtilisateurID int    `json:"utilisateurID"`
}
