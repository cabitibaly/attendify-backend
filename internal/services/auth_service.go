package services

import (
	"fmt"
	"time"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
	"github.com/cabitibaly/pkg/utils"
)

type AuthService struct {
	jwtRepo         *repositories.JWTRepository
	utilisateurRepo *repositories.UtilisateurRepository
}

func NewAuthservice(jwtRepo *repositories.JWTRepository, utilisateurRepo *repositories.UtilisateurRepository) *AuthService {
	return &AuthService{
		jwtRepo:         jwtRepo,
		utilisateurRepo: utilisateurRepo,
	}
}

func (s *AuthService) CreerUnCompte(utilisateurDTO dto.UtilisateurDTO) (*models.Utilisateur, error) {
	if utilisateurDTO.Nom == "" ||
		utilisateurDTO.Email == "" ||
		utilisateurDTO.Poste == "" ||
		utilisateurDTO.Telephone == "" ||
		utilisateurDTO.MotDePasse == "" {
		return nil, fmt.Errorf("nom, email, poste, telephone et mot de passe sont obligatoires")
	}

	utilisateur_exist, _ := s.utilisateurRepo.FindByEmail(utilisateurDTO.Email)

	if utilisateur_exist != nil {
		return nil, fmt.Errorf("cet email est déjà utilisé")
	}

	if _, err := s.utilisateurRepo.FindByTelephone(utilisateurDTO.Telephone); err == nil {
		return nil, fmt.Errorf("cet numero de telephone est déjà utilisé")
	}

	hashedPassword, err := utils.HashPassword(utilisateurDTO.MotDePasse)

	if err != nil {
		return nil, fmt.Errorf("une erreur est survenue")
	}

	utilisateur := models.Utilisateur{
		Nom:                      utilisateurDTO.Nom,
		Prenom:                   utilisateurDTO.Prenom,
		Email:                    utilisateurDTO.Email,
		Poste:                    utilisateurDTO.Poste,
		Telephone:                utilisateurDTO.Telephone,
		MotdepasseAReinitialiser: true,
		MotDePasse:               hashedPassword,
		RoleID:                   2,
		GeoreperageID:            utilisateurDTO.GeoreperageID,
	}

	err = s.utilisateurRepo.Create(&utilisateur)

	if err != nil {
		return nil, fmt.Errorf("une erreur est survenue lors de la création de l'utilisateur")
	}

	return &utilisateur, nil
}

func (s *AuthService) Connexion(email, telephone, motDePasse string, expectedRoleID int) (*models.Utilisateur, string, error) {

	if email == "" && telephone == "" {
		return nil, "", fmt.Errorf("email ou telephone est obligatoire")
	}

	if motDePasse == "" {
		return nil, "", fmt.Errorf("mot de passe est obligatoire")
	}

	utilisateurExist, _ := s.utilisateurRepo.FindByEmailOrTelephone(email, telephone)

	if utilisateurExist == nil {
		return nil, "", fmt.Errorf("vos identifiants sont incorrects")
	}

	if utilisateurExist.RoleID != expectedRoleID {
		return nil, "", fmt.Errorf("vous n'avez pas les droits d'utiliser cette ressource")
	}

	if !utils.ComparePassword(motDePasse, utilisateurExist.MotDePasse) {
		return nil, "", fmt.Errorf("vos identifiants sont incorrects")
	}

	tokenString, err := utils.GenerateToken(
		uint(utilisateurExist.IDUtilisateur),
		uint(utilisateurExist.RoleID),
		utilisateurExist.Email,
		utilisateurExist.Telephone,
	)

	if err != nil {
		return nil, "", fmt.Errorf("une erreur est survenue, veuillez réessayer")
	}

	jwt := &models.Jwt{
		Token:         tokenString,
		UtilisateurID: uint(utilisateurExist.IDUtilisateur),
		ExpireAt:      time.Now().Add(time.Hour * 72),
	}

	erreur := s.jwtRepo.Create(jwt)

	if erreur != nil {
		return nil, "", fmt.Errorf("une erreur est survenue, veuillez réessayer")
	}

	return utilisateurExist, tokenString, nil
}

func (s *AuthService) ReinitialiserMotDePasse(utilisateurID uint) (string, error) {
	motDePasse, err := utils.GeneratePassword(12)

	if err != nil {
		return "", fmt.Errorf("une erreur est survenue, veuillez réessayer")
	}

	hashedPassword, errPassword := utils.HashPassword(motDePasse)

	if errPassword != nil {
		return "", fmt.Errorf("une erreur est survenue, veuillez réessayer")
	}

	errUpdate := s.utilisateurRepo.Update(utilisateurID, map[string]any{
		"motDePasse":               hashedPassword,
		"motdepasseAReinitialiser": true,
	})

	if errUpdate != nil {
		return "", fmt.Errorf("une erreur est survenue, veuillez réessayer")
	}

	return motDePasse, nil
}
