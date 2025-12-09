package handlers

import (
	"strconv"

	"github.com/cabitibaly/internal/services"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *services.NofificationService
}

func NewNoficationHandler(service *services.NofificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) ToutesLesNotificationsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	utilisateurID, err := c.Get("utilisateurID")
	if !err {
		c.JSON(400, gin.H{
			"error":  "Erreur dans l'authentification",
			"status": 400,
		})
		return
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	notifications, hasNextPage, errService := h.service.ToutesLesNotifications(utilisateurID.(uint), page, limit)
	if errService != nil {
		c.JSON(500, gin.H{
			"error":  errService.Error(),
			"status": 500,
		})
		return
	}

	c.JSON(200, gin.H{
		"notifications": notifications,
		"hasNextPage":   hasNextPage,
		"status":        200,
	})
}

func (h *NotificationHandler) ModifierUneNoticationHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error":  "Une erreur est survenue",
			"status": 400,
		})
		return
	}

	estLu, _ := strconv.ParseBool(c.Query("estLu"))

	errService := h.service.ModifierUneNotification(uint(id), map[string]any{"estLu": estLu})
	if errService != nil {
		c.JSON(500, gin.H{
			"error":  errService.Error(),
			"status": 500,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Notification modifiée avec succès",
		"status":  200,
	})
}

func (h *NotificationHandler) SupprimerUneNotificationHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error":  "Une erreur est survenue",
			"status": 400,
		})
		return
	}

	errService := h.service.SupprimerUneNotification(uint(id))
	if errService != nil {
		c.JSON(500, gin.H{
			"error":  errService.Error(),
			"status": 500,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Notification supprimée avec succès",
		"status":  200,
	})
}
