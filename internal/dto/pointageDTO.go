package dto

import "time"

type PointageDTO struct {
	Arrivee           time.Time  `json:"arrivee"`
	Depart            *time.Time `json:"depart"`
	EstPresent        bool       `json:"estPresent"`
	EnRetard          bool       `json:"enRetard"`
	HeuresTravaillees *float64   ` json:"heuresTravaillees"`
	DepartAnticipe    bool       ` json:"departAnticipe"`
	UtilisateurID     int        `json:"utilisateurID"`
}
