package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/services"
	"github.com/gin-gonic/gin"
)

type PointageHandler struct {
	service *services.PointageService
}

func NewPointageHandler(service *services.PointageService) *PointageHandler {
	return &PointageHandler{service: service}
}

func (h *PointageHandler) PointageArriveeHandler(c *gin.Context) {
	var data dto.PointageDTO
	empLatitude, _ := strconv.ParseFloat(c.Query("latitude"), 64)
	empLongitute, _ := strconv.ParseFloat(c.Query("longitude"), 64)

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	err := h.service.PointageArrivee(data, empLatitude, empLongitute)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pointage enregistré avec succès",
		"status":  http.StatusOK,
	})
}

func (h *PointageHandler) PointageDepartHandler(c *gin.Context) {
	var data dto.PointageDTO

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	err := h.service.PointageDepart(data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pointage supprimé avec succès",
		"status":  http.StatusOK,
	})
}

func (h *PointageHandler) TousLesPointagesHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	aujoudhui, _ := strconv.ParseBool(c.Query("aujourdhui"))
	utilisateurID, _ := strconv.Atoi(c.Query("utilisateurID"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	pointages, hasNextPage, total, err := h.service.TousLesPointages(uint(utilisateurID), aujoudhui, page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pointages":   pointages,
		"hasNextPage": hasNextPage,
		"total":       total,
		"status":      http.StatusOK,
	})
}

func (h *PointageHandler) TousLesPointagesParDateHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	date, err := time.Parse("2006-01-02T15:04:05.00Z", c.Query("date"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	pointages, hasNextPage, total, err := h.service.TousLesPointagesParDate(date, page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pointages":   pointages,
		"hasNextPage": hasNextPage,
		"total":       total,
		"status":      http.StatusOK,
	})
}

func (h *PointageHandler) LireUnPointageHandler(c *gin.Context) {
	date, err := time.Parse("2006-01-02T15:04:05.00Z", c.Query("date"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	pointage, err := h.service.LireUnPointage(date)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pointage": pointage,
		"status":   http.StatusOK,
	})
}

func (h *PointageHandler) SupprimerUnPointageHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.SupprimerPointage(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pointage supprimé avec succès",
		"status":  http.StatusOK,
	})
}
