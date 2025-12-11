package services

import (
	"fmt"
	"time"

	"github.com/cabitibaly/internal/models"
	"github.com/cabitibaly/internal/repositories"
	"github.com/cabitibaly/pkg/utils"
)

type UtilisateurService struct {
	utilisateurRepo *repositories.UtilisateurRepository
	jwtRepo         *repositories.JWTRepository
}

func NewUtilisateurService(utilisateurRepo *repositories.UtilisateurRepository, jwtRepo *repositories.JWTRepository) *UtilisateurService {
	return &UtilisateurService{
		utilisateurRepo: utilisateurRepo,
		jwtRepo:         jwtRepo,
	}
}

func (s *UtilisateurService) MesInformations(utilisateurID uint) (*models.Utilisateur, error) {
	return s.utilisateurRepo.FindByID(utilisateurID)
}

func (s *UtilisateurService) TousLesEmployes(recherche string, page, limit int) ([]models.Utilisateur, bool, int64, error) {
	return s.utilisateurRepo.FindAll(recherche, page, limit)
}

func (s *UtilisateurService) LireUnEmploye(utilisateurID uint) (*models.Utilisateur, error) {
	return s.utilisateurRepo.FindByID(utilisateurID)
}

func (s *UtilisateurService) ModifierSonCompte(utilisateurID uint, token string, data map[string]any) (string, error) {
	utilisateurExist, err := s.utilisateurRepo.FindByID(utilisateurID)

	if err != nil {
		return "", err
	}

	emailChanged := data["email"] != nil && data["email"] != utilisateurExist.Email
	telephoneChanged := data["telephone"] != nil && data["telephone"] != utilisateurExist.Telephone

	tokenString := ""

	if emailChanged || telephoneChanged {

		nouvelEmail := ""
		if data["email"] != nil {
			nouvelEmail = data["email"].(string)
		}

		if _, err := s.utilisateurRepo.FindByEmail(nouvelEmail); err == nil && emailChanged {
			return "", fmt.Errorf("cet email est déjà utilisé")
		}

		nouveauTelephone := ""
		if data["telephone"] != nil {
			nouveauTelephone = data["telephone"].(string)
		}

		if _, err := s.utilisateurRepo.FindByTelephone(nouveauTelephone); err == nil && telephoneChanged {
			return "", fmt.Errorf("cet numero de telephone est déjà utilisé")
		}

		tokenDB, err := s.jwtRepo.FindByTokenAndUtilisateurID(token, utilisateurID)

		if err != nil {
			return "", err
		}

		if tokenDB != nil {
			err := s.jwtRepo.Delete(uint(tokenDB.ID))

			if err != nil {
				return "", err
			}
		}

		tokenGenerer, err := utils.GenerateToken(
			uint(utilisateurExist.IDUtilisateur),
			uint(utilisateurExist.RoleID),
			nouvelEmail,
			nouveauTelephone,
		)

		if err != nil {
			return "", err
		}

		jwt := &models.Jwt{
			Token:         tokenGenerer,
			UtilisateurID: uint(utilisateurExist.IDUtilisateur),
			ExpireAt:      time.Now().Add(time.Hour * 72),
		}

		err = s.jwtRepo.Create(jwt)

		if err != nil {
			return "", err
		}

		tokenString = tokenGenerer
	}

	return tokenString, s.utilisateurRepo.Update(utilisateurID, data)
}

func (s *UtilisateurService) ModifierSonMotDePasse(utilisateurID uint, ancien, nouveau string) error {
	utilisateurExist, err := s.utilisateurRepo.FindByID(utilisateurID)

	if err != nil {
		return err
	}

	if !utilisateurExist.MotdepasseAReinitialiser {
		if !utils.ComparePassword(ancien, utilisateurExist.MotDePasse) {
			return fmt.Errorf("acient mot de passe incorrect")
		}
	}

	hashedPassword, errorHashPwd := utils.HashPassword(nouveau)

	if errorHashPwd != nil {
		return errorHashPwd
	}

	errUpdate := s.utilisateurRepo.Update(utilisateurID, map[string]any{
		"motDePasse":               hashedPassword,
		"motdepasseAReinitialiser": false,
	})

	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (s *UtilisateurService) SupprimerUnCompte(utilisateurID uint) error {
	return s.utilisateurRepo.Delete(utilisateurID)
}
