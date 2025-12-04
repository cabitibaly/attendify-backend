package dto

type UtilisateurDTO struct {
	Nom                      string `json:"nom"`
	Prenom                   string `json:"prenom"`
	Email                    string `json:"email"`
	Poste                    string `json:"poste"`
	Telephone                string `json:"telephone"`
	MotDePasse               string `json:"motDePasse"`
	MotdepasseAReinitialiser bool   `json:"motdepasseAReinitialiser"`
	GeoreperageID            int    `json:"georeperageID"`
}
