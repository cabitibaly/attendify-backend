package dto

type PushTokenDTO struct {
	PushToken     string `json:"push_token"`
	DeviceType    string `json:"device_type"`
	DeviceName    string `json:"device_name"`
	Platform      string `json:"platform"`
	EstActive     bool   `json:"est_active"`
	UtilisateurID int    `json:"utilisateurID"`
}
