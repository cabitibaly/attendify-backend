package handlers

import (
	"net/http"
	"strconv"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/services"
	"github.com/gin-gonic/gin"
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

	utilisateur, token, err := h.service.Connexion(data.Email, data.Telephone, data.MotDePasse, 1)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  err.Error(),
			"status": http.StatusUnauthorized,
		})
		return
	}

	c.SetCookie(
		"jwt",
		token,
		3600*24*3,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"utilisateur": utilisateur,
		"status":      http.StatusOK,
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

	utilisateur, token, err := h.service.Connexion(data.Email, data.Telephone, data.MotDePasse, 2)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  err.Error(),
			"status": http.StatusUnauthorized,
		})
		return
	}

	c.SetCookie(
		"jwt",
		token,
		3600*24*3,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"utilisateur": utilisateur,
		"status":      http.StatusOK,
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
