package services

import (
	"fmt"
	"log"
	"time"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
	"github.com/cabitibaly/pkg/utils"
)

type AuthService struct {
	refreshTokenRepo *repositories.RefreshTokenRepository
	utilisateurRepo  *repositories.UtilisateurRepository
	siteRepo         *repositories.SiteRepository
}

func NewAuthservice(
	refreshTokenRepo *repositories.RefreshTokenRepository,
	utilisateurRepo *repositories.UtilisateurRepository,
	siteRepo *repositories.SiteRepository,
) *AuthService {
	return &AuthService{
		refreshTokenRepo: refreshTokenRepo,
		utilisateurRepo:  utilisateurRepo,
		siteRepo:         siteRepo,
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

	if utilisateurDTO.SiteID == 0 {
		return nil, fmt.Errorf("vous devez sélectionner un site")
	}

	if _, err := s.siteRepo.FindByID(uint(utilisateurDTO.SiteID)); err != nil {
		return nil, fmt.Errorf("le site sélectionné n'existe pas")
	}

	if _, err := s.utilisateurRepo.FindByEmail(utilisateurDTO.Email); err == nil {
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
		SiteID:                   utilisateurDTO.SiteID,
	}

	err = s.utilisateurRepo.Create(&utilisateur)

	if err != nil {
		return nil, fmt.Errorf("une erreur est survenue lors de la création de l'utilisateur")
	}

	return &utilisateur, nil
}

func (s *AuthService) Connexion(username, motDePasse string, expectedRoleID int) (*models.Utilisateur, string, string, error) {

	if username == "" {
		return nil, "", "", fmt.Errorf("email ou telephone est obligatoire")
	}

	if motDePasse == "" {
		return nil, "", "", fmt.Errorf("mot de passe est obligatoire")
	}

	utilisateurExist, _ := s.utilisateurRepo.FindByEmailOrTelephone(username)

	if utilisateurExist == nil {
		return nil, "", "", fmt.Errorf("vos identifiants sont incorrects")
	}

	if utilisateurExist.RoleID != expectedRoleID {
		return nil, "", "", fmt.Errorf("vous n'avez pas les droits d'utiliser cette ressource")
	}

	if !utils.ComparePassword(motDePasse, utilisateurExist.MotDePasse) {
		return nil, "", "", fmt.Errorf("vos identifiants sont incorrects")
	}

	accessToken, errAT := utils.GenerateAccessToken(
		uint(utilisateurExist.IDUtilisateur),
		uint(utilisateurExist.RoleID),
		utilisateurExist.Email,
		utilisateurExist.Telephone,
	)
	if errAT != nil {
		log.Println("Une erreur est survenue lors de la génération du token d'accès :", errAT)
		return nil, "", "", fmt.Errorf("une erreur est survenue lors de la génération du token d'accès")
	}

	refreshToken, expireTime, errRT := utils.GenerateRefreshToken(
		uint(utilisateurExist.IDUtilisateur),
		uint(utilisateurExist.RoleID),
		utilisateurExist.Email,
		utilisateurExist.Telephone,
	)
	if errRT != nil {
		log.Println("Une erreur est survenue lors de la génération du refresh token :", errRT)
		return nil, "", "", fmt.Errorf("une erreur est survenue lors de la génération du refresh token")
	}

	jwt := &models.RefreshToken{
		Token:         refreshToken,
		UtilisateurID: uint(utilisateurExist.IDUtilisateur),
		ExpireAt:      *expireTime,
	}

	erreur := s.refreshTokenRepo.Create(jwt)

	if erreur != nil {
		return nil, "", "", fmt.Errorf("une erreur est survenue, veuillez réessayer")
	}

	return utilisateurExist, refreshToken, accessToken, nil
}

func (s *AuthService) Deconnexion(token string) error {
	return s.refreshTokenRepo.DeleteByToken(token)
}

func (s *AuthService) NouveauRefreshToken(refreshToken string) (string, string, error) {
	ancienRT, err := s.refreshTokenRepo.FindByToken(refreshToken)

	if err != nil {
		return "", "", err
	}

	location, _ := time.LoadLocation("Africa/Ouagadougou")
	maintenat := time.Now().In(location)
	if maintenat.After(ancienRT.ExpireAt) || ancienRT.RevokedAt != nil {
		return "", "", fmt.Errorf("le refresh token a expiré ou a été révoqué")
	}

	accessToken, errAT := utils.GenerateAccessToken(
		ancienRT.UtilisateurID,
		uint(ancienRT.Utilisateur.RoleID),
		ancienRT.Utilisateur.Email,
		ancienRT.Utilisateur.Telephone,
	)
	if errAT != nil {
		log.Println("Une erreur est survenue lors de la génération du token d'accès :", errAT)
		return "", "", fmt.Errorf("une erreur est survenue lors de la génération du token d'accès")
	}

	newRefreshToken, expireTime, errRT := utils.GenerateRefreshToken(
		ancienRT.UtilisateurID,
		uint(ancienRT.Utilisateur.RoleID),
		ancienRT.Utilisateur.Email,
		ancienRT.Utilisateur.Telephone,
	)
	if errRT != nil {
		log.Println("Une erreur est survenue lors de la génération du refresh token :", errRT)
		return "", "", fmt.Errorf("une erreur est survenue lors de la génération du refresh token")
	}

	s.refreshTokenRepo.Delete(uint(ancienRT.ID))

	token := &models.RefreshToken{
		Token:         newRefreshToken,
		ExpireAt:      *expireTime,
		UtilisateurID: ancienRT.UtilisateurID,
	}

	err = s.refreshTokenRepo.Create(token)

	if err != nil {
		return "", "", err
	}

	return newRefreshToken, accessToken, nil
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
