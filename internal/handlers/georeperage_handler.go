package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GeorepHandler struct {
	service *services.GeorepService
}

func NewGeorepHandler(service *services.GeorepService) *GeorepHandler {
	return &GeorepHandler{service: service}
}

func (h *GeorepHandler) CreerUnSiteHandler(c *gin.Context) {
	var data dto.GeorepDTO

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	err := h.service.CreerUnSite(data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Le site a été créé avec succès",
		"status":  http.StatusCreated,
	})
}

func (h *GeorepHandler) TousLesSitesHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	sites, hasNextPage, total, err := h.service.TousLesSites(page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sites":       sites,
		"hasNextPage": hasNextPage,
		"total":       total,
		"status":      http.StatusOK,
	})
}

func (h *GeorepHandler) LireUnSiteHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	site, err := h.service.LireUnSite(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Le site n'existe pas",
			"status": http.StatusNotFound,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"site":   site,
		"status": http.StatusOK,
	})

}

func (h *GeorepHandler) ModifierUnSiteHandler(c *gin.Context) {
	var data map[string]any
	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	err := h.service.ModifierUnSite(uint(id), data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Le site a été modifié avec succès",
		"status":  http.StatusOK,
	})
}

func (h *GeorepHandler) SupprimerUnSiteHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.SupprimerUnSite(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Le site a été supprimé avec succès",
		"status":  http.StatusOK,
	})
}
