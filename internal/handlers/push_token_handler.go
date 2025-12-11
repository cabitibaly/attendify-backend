package handlers

import (
	"net/http"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/services"
	"github.com/gin-gonic/gin"
)

type PushTokenHandler struct {
	service *services.PushTokenService
}

func NewPushTokenHandler(service *services.PushTokenService) *PushTokenHandler {
	return &PushTokenHandler{service: service}
}

func (h *PushTokenHandler) EnregistrerOuModifierPushTokenHandler(c *gin.Context) {
	var data dto.PushTokenDTO

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Une erreur est survenue",
			"status": http.StatusBadRequest,
		})
		return
	}

	utilisateurID, err := c.Get("utilisateurID")
	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Vous n'êtes pas autorisé à accéder à cette ressource",
			"status": http.StatusUnauthorized,
		})
		return
	}

	data.UtilisateurID = int(utilisateurID.(uint))

	if err := h.service.EnregistrerOuModifierPushToken(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Une erreur est survenue",
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Le token a bien été enregistré",
		"status":  http.StatusOK,
	})
}
