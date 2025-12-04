package dto

type GeorepDTO struct {
	NomSite   string  `json:"site"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Rayon     float64 `json:"rayon"`
}
