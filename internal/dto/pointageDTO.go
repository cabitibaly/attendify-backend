package dto

import "time"

type PointageDTO struct {
	Arrive            time.Time  `json:"arrive"`
	Depart            *time.Time `json:"depart"`
	EstPresent        bool       `json:"estPresent"`
	EnRetard          bool       `json:"enRetard"`
	HeuresTravaillees *float64   ` json:"heuresTravaillees"`
	DepartAnticipe    bool       ` json:"departAnticipe"`
	UtilisateurID     int        `json:"utilisateurID"`
}
