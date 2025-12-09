package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CongeHandler struct {
	service *services.CongeService
}

func NewCongeHandler(service *services.CongeService) *CongeHandler {
	return &CongeHandler{service: service}
}

func (h *CongeHandler) FaireUneDemandeHandler(c *gin.Context) {
	var data dto.CongeDTO
	utilisateurID, err := c.Get("utilisateurID")

	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "erreur de connexion",
			"status": http.StatusUnauthorized,
		})
		return
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	errService := h.service.FaireUneDemande(data, utilisateurID.(uint))
	if errService != nil {
		if strings.Contains(errService.Error(), "date") ||
			strings.Contains(errService.Error(), "congé") ||
			strings.Contains(errService.Error(), "congés") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  errService.Error(),
				"status": http.StatusBadRequest,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "congé créé avec succès",
		"status":  http.StatusOK,
	})
}

func (h *CongeHandler) TousLesCongésEmployeHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	statutID, _ := strconv.Atoi(c.Query("statutID"))
	utilisateurID, err := c.Get("utilisateurID")

	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Erreur de connexion",
			"status": http.StatusBadRequest,
		})
		return
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	conges, hasNextPage, total, errService := h.service.TousLesConges(utilisateurID.(uint), uint(statutID), page, limit)
	if errService != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conges":      conges,
		"hasNextPage": hasNextPage,
		"total":       total,
		"status":      http.StatusOK,
	})
}

func (h *CongeHandler) TousLesCongesHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	statutID, _ := strconv.Atoi(c.Query("statutID"))
	utilisateurID, _ := strconv.Atoi(c.Query("utilisateurID"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	conges, hasNextPage, total, errService := h.service.TousLesConges(uint(utilisateurID), uint(statutID), page, limit)
	if errService != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conges":      conges,
		"hasNextPage": hasNextPage,
		"total":       total,
		"status":      http.StatusOK,
	})
}

func (h *CongeHandler) LireUnCongeHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Erreur de paramètre",
			"status": http.StatusBadRequest,
		})
		return
	}

	conge, errService := h.service.LireUnConge(uint(id))
	if errService != nil {
		if errors.Is(errService, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  "congé non trouvé",
				"status": http.StatusNotFound,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conge":  conge,
		"status": http.StatusOK,
	})
}

func (h *CongeHandler) ModifierUnCongeHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Erreur de paramètre",
			"status": http.StatusBadRequest,
		})
		return
	}

	var data map[string]any
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Verifier les données",
			"status": http.StatusBadRequest,
		})
		return
	}

	errService := h.service.ModifierUnConge(uint(id), data)
	if errService != nil {
		if errors.Is(errService, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  errService.Error(),
				"status": http.StatusNotFound,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "congé modifié avec succès",
		"status":  http.StatusOK,
	})
}

func (h *CongeHandler) ModifierStatutCongeHandler(c *gin.Context) {
	statutID, _ := strconv.Atoi(c.Query("statutID"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Erreur de paramètre",
			"status": http.StatusBadRequest,
		})
		return
	}

	if statutID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Vous devez choisir un statut",
			"status": http.StatusBadRequest,
		})
		return
	}

	errService := h.service.ModifierStatutConge(uint(id), uint(statutID))
	if errService != nil {
		if errors.Is(errService, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  errService.Error(),
				"status": http.StatusNotFound,
			})
			return
		}

		if strings.Contains(errService.Error(), "congé") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  errService.Error(),
				"status": http.StatusBadRequest,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "statut de congé modifié avec succès",
		"status":  http.StatusOK,
	})
}

func (h *CongeHandler) SupprimerUnCongeHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Erreur de paramètre",
			"status": http.StatusBadRequest,
		})
		return
	}

	errService := h.service.SupprimerUnConge(uint(id))
	if errService != nil {
		if errors.Is(errService, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  errService.Error(),
				"status": http.StatusNotFound,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "congé supprimé avec succès",
		"status":  http.StatusOK,
	})
}
