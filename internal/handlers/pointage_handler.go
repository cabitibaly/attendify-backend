package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PointageHandler struct {
	service *services.PointageService
}

func NewPointageHandler(service *services.PointageService) *PointageHandler {
	return &PointageHandler{service: service}
}

func (h *PointageHandler) EstSurSiteHandler(c *gin.Context) {
	empLatitude, _ := strconv.ParseFloat(c.Query("latitude"), 64)
	empLongitute, _ := strconv.ParseFloat(c.Query("longitude"), 64)
	utilisateurID, _ := c.Get("utilisateurID")

	estSurSite, err := h.service.EstSurSite(utilisateurID.(uint), empLatitude, empLongitute)

	if !estSurSite && err != nil {
		if strings.Contains(err.Error(), "n'existe pas") {
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
		"est_sur_site": estSurSite,
		"status":       http.StatusOK,
	})
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

	pointages, hasNextPage, total, err := h.service.TousLesPointages(uint(utilisateurID), aujoudhui, date, page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	pointagesFormated := dto.ToPointageResponseDTOList(pointages)

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

	pointages, hasNextPage, total, errService := h.service.TousLesPointages(utilisateurID.(uint), aujoudhui, date, page, limit)

	if errService != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	pointagesFormated := dto.ToPointageResponseDTOList(pointages)

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

func (h *PointageHandler) ExportHandler(c *gin.Context) {
	debut := c.Query("debut")
	fin := c.Query("fin")

	if debut == "" || fin == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Les dates de début et de fin sont obligatoires",
			"status": http.StatusBadRequest,
		})
		return
	}

	debutParsed, err := time.Parse("2006-01-02T15:04:05.000Z", debut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Format de date début invalide (attendu: JJ/MM/AAAA)",
			"status": http.StatusBadRequest,
		})
		return
	}

	finParsed, err := time.Parse("2006-01-02T15:04:05.000Z", fin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Format de date fin invalide (attendu: JJ/MM/AAAA)",
			"status": http.StatusBadRequest,
		})
		return
	}

	if finParsed.Before(debutParsed) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "La date de fin ne peut pas être antérieure à la date de début",
			"status": http.StatusBadRequest,
		})
		return
	}

	filBytes, err := h.service.Export(debutParsed, finParsed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Erreur lors de la génération du fichier: " + err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	filename := utils.GenerateFileName(debutParsed, finParsed)
	c.Header("Content-Description", "File Transfert")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", filBytes)
}
