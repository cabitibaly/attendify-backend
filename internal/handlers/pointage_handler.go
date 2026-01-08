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

type PointageHandler struct {
	service *services.PointageService
}

func NewPointageHandler(service *services.PointageService) *PointageHandler {
	return &PointageHandler{service: service}
}

func (h *PointageHandler) PointageArriveeHandler(c *gin.Context) {
	empLatitude, _ := strconv.ParseFloat(c.Query("latitude"), 64)
	empLongitute, _ := strconv.ParseFloat(c.Query("longitude"), 64)
	utilisateurID, _ := c.Get("utilisateurID")

	if empLatitude == 0 || empLongitute == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Une erreur est survenue veuillez reessayer!",
			"status": http.StatusBadRequest,
		})
		return
	}

	err := h.service.PointageArrivee(utilisateurID.(uint), empLatitude, empLongitute)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  err.Error(),
				"status": http.StatusNotFound,
			})
			return
		}

		if strings.Contains(err.Error(), "sur site") {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  err.Error(),
				"status": http.StatusForbidden,
			})
			return
		}

		if strings.Contains(err.Error(), "vous avez déjà") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  err.Error(),
				"status": http.StatusBadRequest,
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
		"message": "pointage enregistré avec succès",
		"status":  http.StatusOK,
	})
}

func (h *PointageHandler) PointageDepartHandler(c *gin.Context) {
	empLatitude, _ := strconv.ParseFloat(c.Query("latitude"), 64)
	empLongitute, _ := strconv.ParseFloat(c.Query("longitude"), 64)
	utilisateurID, _ := c.Get("utilisateurID")

	err := h.service.PointageDepart(utilisateurID.(uint), empLatitude, empLongitute)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  err.Error(),
				"status": http.StatusNotFound,
			})
			return
		}

		if strings.Contains(err.Error(), "sur site") {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  err.Error(),
				"status": http.StatusForbidden,
			})
			return
		}

		if strings.Contains(err.Error(), "vous n'avez pas") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  err.Error(),
				"status": http.StatusBadRequest,
			})
			return
		}

		if strings.Contains(err.Error(), "vous avez déjà") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  err.Error(),
				"status": http.StatusBadRequest,
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
		"message": "Départ enregistré avec succès",
		"status":  http.StatusOK,
	})
}

func (h *PointageHandler) TousLesPointagesHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	aujoudhui, _ := strconv.ParseBool(c.Query("aujourdhui"))
	date, _ := time.Parse("2006-01-02T15:04:05.000Z", c.Query("date"))
	utilisateurID, _ := strconv.Atoi(c.Query("utilisateurID"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	cacheKey := "pointages:All:" + strconv.Itoa(utilisateurID) + ":" + strconv.FormatBool(aujoudhui) + ":" + date.Format("2006-01-02") + ":" + strconv.Itoa(page) + ":" + strconv.Itoa(limit)
	if cached, err := configs.GetCache(cacheKey); err == nil {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
		return
	}

	pointages, hasNextPage, total, err := h.service.TousLesPointages(uint(utilisateurID), aujoudhui, date, page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	pointagesFormated := dto.ToPointageResponseDTOList(pointages)

	cacheValue := dto.PointagesResponse{
		Pointages: pointagesFormated,
		Pagination: dto.Pagination{
			HasNextPage: hasNextPage,
			Total:       total,
			Status:      http.StatusOK,
		},
	}

	jsonData, _ := json.Marshal(cacheValue)
	_ = configs.SetCache(cacheKey, jsonData, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"pointages":   pointagesFormated,
		"hasNextPage": hasNextPage,
		"total":       total,
		"status":      http.StatusOK,
	})
}

func (h *PointageHandler) TousLesPointagesEmployeHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	aujoudhui, _ := strconv.ParseBool(c.Query("aujourdhui"))
	date, _ := time.Parse("2006-01-02T15:04:05.00Z", c.Query("date"))

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

	cacheKey := "pointages:AllEmploye:" + fmt.Sprintf("%d", utilisateurID) + ":" + strconv.FormatBool(aujoudhui) + ":" + date.Format("2006-01-02") + ":" + strconv.Itoa(page) + ":" + strconv.Itoa(limit)
	if cached, err := configs.GetCache(cacheKey); err == nil {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
		return
	}

	pointages, hasNextPage, total, errService := h.service.TousLesPointages(utilisateurID.(uint), aujoudhui, date, page, limit)

	if errService != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	pointagesFormated := dto.ToPointageResponseDTOList(pointages)

	cacheValue := dto.PointagesResponse{
		Pointages: pointagesFormated,
		Pagination: dto.Pagination{
			HasNextPage: hasNextPage,
			Total:       total,
			Status:      http.StatusOK,
		},
	}

	jsonData, _ := json.Marshal(cacheValue)
	_ = configs.SetCache(cacheKey, jsonData, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"pointages":   pointagesFormated,
		"hasNextPage": hasNextPage,
		"total":       total,
		"status":      http.StatusOK,
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

func (h *PointageHandler) StatsHandler(c *gin.Context) {
	totalEmp, totalPresent, totalRetard, err := h.service.Stats()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalEmploye": totalEmp,
		"present":      totalPresent,
		"retard":       totalRetard,
		"absent":       totalEmp - totalPresent,
		"tauxPresence": float64(totalPresent) / float64(totalEmp) * 100,
		"status":       http.StatusOK,
	})
}
