package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/cabitibaly/configs"
	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) CreerUnCompteHandler(c *gin.Context) {
	var data dto.UtilisateurDTO

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	_, err := h.service.CreerUnCompte(data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Compte créé avec succès",
		"status":  http.StatusCreated,
	})
}

func (h *AuthHandler) ConnexionAdminHandler(c *gin.Context) {
	var data dto.ConnexionDTO

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	utilisateur, refreshToken, accessToken, err := h.service.Connexion(data.Email, data.Telephone, data.MotDePasse, 1)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  err.Error(),
			"status": http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"utilisateur":   dto.ToUtilisateurResponseDTO(utilisateur),
		"refresh_token": refreshToken,
		"access_token":  accessToken,
		"status":        http.StatusOK,
	})
}

func (h *AuthHandler) ConnexionEmployeHandler(c *gin.Context) {
	var data dto.ConnexionDTO

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	utilisateur, refreshToken, accessToken, err := h.service.Connexion(data.Email, data.Telephone, data.MotDePasse, 2)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  err.Error(),
			"status": http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"utilisateur":   dto.ToUtilisateurResponseDTO(utilisateur),
		"refresh_token": refreshToken,
		"access_token":  accessToken,
		"status":        http.StatusOK,
	})
}

func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	var data struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON((&data)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	refreshToken, accessToken, err := h.service.NouveauRefreshToken(data.Token)
	if err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  err.Error(),
				"status": http.StatusNotFound,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"refresh_token": refreshToken,
		"access_token":  accessToken,
		"status":        http.StatusOK,
	})

}

func (h *AuthHandler) DeconnexionHandler(c *gin.Context) {
	jti, errExiste := c.Get("jti")

	if !errExiste {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Vous n'êtes pas connecté(e)s",
			"status": http.StatusUnauthorized,
		})
		return
	}

	cacheKey := "access_token:" + jti.(string)

	var data struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON((&data)); err != nil {
		log.Println("une erreur est survenue:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Vous n'êtes pas connecté(e)s",
			"status": http.StatusUnauthorized,
		})
		return
	}

	err := h.service.Deconnexion(data.Token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	_ = configs.DeleteCache(cacheKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "Vous êtes déconnecté(e)",
		"status":  http.StatusOK,
	})
}

func (h *AuthHandler) ReinitialiserMotDePasseHandler(c *gin.Context) {

	utilisateurID, _ := strconv.Atoi(c.Param("id"))

	motDePasse, err := h.service.ReinitialiserMotDePasse(uint(utilisateurID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"motDePasse": motDePasse,
		"status":     http.StatusOK,
	})
}
