package dto

type GeorepDTO struct {
	NomSite   string  `json:"nomSite"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Rayon     float64 `json:"rayon"`
}
