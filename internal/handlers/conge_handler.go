package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cabitibaly/configs"
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

		if strings.Contains(errService.Error(), "n'existe pas") {
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

	c.JSON(http.StatusCreated, gin.H{
		"message": "congé créé avec succès",
		"status":  http.StatusCreated,
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

	cacheKey := "conges:AllEmploye:" + fmt.Sprintf("%d", utilisateurID) + ":" + strconv.Itoa(statutID) + ":" + strconv.Itoa(page) + ":" + strconv.Itoa(limit)
	if cached, err := configs.GetCache(cacheKey); err == nil {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
		return
	}

	conges, hasNextPage, total, errService := h.service.TousLesConges(utilisateurID.(uint), uint(statutID), page, limit)
	if errService != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	congesFormated := dto.ToCongeResponseEmpDTOList(conges)

	cacheValue := dto.CongesResponseEmp{
		Conges: congesFormated,
		Pagination: dto.Pagination{
			HasNextPage: hasNextPage,
			Total:       total,
			Status:      http.StatusOK,
		},
	}

	jsonData, _ := json.Marshal(cacheValue)
	_ = configs.SetCache(cacheKey, jsonData, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"conges":      congesFormated,
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

	congesFormated := dto.ToCongeResponseAdminDTOList(conges)

	c.JSON(http.StatusOK, gin.H{
		"conges":      congesFormated,
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

	roleID, errRole := c.Get("roleID")
	if !errRole {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Erreur de connexion",
			"status": http.StatusUnauthorized,
		})
		return
	}

	cacheKey := "conge:" + strconv.Itoa(id) + ":roleID:" + fmt.Sprintf("%d", roleID)
	if cached, err := configs.GetCache(cacheKey); err == nil {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
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

	var response any
	if roleID.(uint) == 2 {
		response = dto.ToCongeResponseEmpDTO(conge)
		cacheValue := dto.CongeResponseEmp{
			Conge:  response,
			Status: http.StatusOK,
		}

		jsonData, _ := json.Marshal(cacheValue)
		_ = configs.SetCache(cacheKey, jsonData, 5*time.Minute)
	} else {
		response = dto.ToCongeResponseAdminDTO(conge)
		cacheValue := dto.CongeResponseAdmin{
			Conge:  response,
			Status: http.StatusOK,
		}

		jsonData, _ := json.Marshal(cacheValue)
		_ = configs.SetCache(cacheKey, jsonData, 5*time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{
		"conge":  response,
		"status": http.StatusOK,
	})
}

func (h *CongeHandler) ModifierUnCongeHandler(c *gin.Context) {
	roleID, _ := c.Get("roleID")
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

	cacheKey := "conge:" + strconv.Itoa(id) + ":roleID:" + fmt.Sprintf("%d", roleID)
	_ = configs.DeleteCache(cacheKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "congé modifié avec succès",
		"status":  http.StatusOK,
	})
}

func (h *CongeHandler) ModifierStatutCongeHandler(c *gin.Context) {
	roleID, _ := c.Get("roleID")
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

	cacheKey := "conge:" + strconv.Itoa(id) + ":roleID:" + fmt.Sprintf("%d", roleID)
	_ = configs.DeleteCache(cacheKey)

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
