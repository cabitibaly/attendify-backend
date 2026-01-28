package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cabitibaly/configs"
	"github.com/cabitibaly/internal/dto"
	"github.com/cabitibaly/internal/services"
	"github.com/gin-gonic/gin"
)

type UtilisateurHandler struct {
	service *services.UtilisateurService
}

func NewUtilisateurHandler(service *services.UtilisateurService) *UtilisateurHandler {
	return &UtilisateurHandler{service: service}
}

func (h *UtilisateurHandler) MesInformationsHandler(c *gin.Context) {
	utilisateurID, err := c.Get("utilisateurID")

	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Vous n'êtes pas connecté(e)s",
			"status": http.StatusUnauthorized,
		})
		return
	}

	cacheKey := "utilisateur:" + fmt.Sprintf("%d", utilisateurID)
	if cached, err := configs.GetCache(cacheKey); err == nil {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
		return
	}

	utilisateur, errGetInfo := h.service.MesInformations(utilisateurID.(uint))

	if errGetInfo != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errGetInfo.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	utilisateurFormated := dto.ToUtilisateurResponseDTO(utilisateur)

	cacheValue := dto.UtilisateurResponse{
		Utilisateur: utilisateurFormated,
		Status:      http.StatusOK,
	}

	jsonData, _ := json.Marshal(cacheValue)
	_ = configs.SetCache(cacheKey, jsonData, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"utilisateur": utilisateurFormated,
		"status":      http.StatusOK,
	})
}

func (h *UtilisateurHandler) TousLesEmployesHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	recherche := c.Query("recherche")

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 20
	}

	utilsateurs, hasNextPage, total, err := h.service.TousLesEmployes(recherche, page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	utilisateursFormated := dto.ToUtilisateurResponseDTOList(utilsateurs)

	c.JSON(http.StatusOK, gin.H{
		"utilisateurs": utilisateursFormated,
		"hasNextPage":  hasNextPage,
		"total":        total,
		"status":       http.StatusOK,
	})
}

func (h *UtilisateurHandler) LireUnEmployeHandler(c *gin.Context) {
	utilisateurID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Erreur de paramètre",
			"status": http.StatusBadRequest,
		})
		return
	}

	cacheKey := "utilisateur:" + strconv.Itoa(utilisateurID)
	if cached, err := configs.GetCache(cacheKey); err == nil {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
		return
	}

	utilisateur, err := h.service.LireUnEmploye(uint(utilisateurID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	utilisateurFormated := dto.ToUtilisateurResponseDTO(utilisateur)

	cacheValue := dto.UtilisateurResponse{
		Utilisateur: utilisateurFormated,
		Status:      http.StatusOK,
	}

	jsonData, _ := json.Marshal(cacheValue)
	_ = configs.SetCache(cacheKey, jsonData, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"utilisateur": utilisateurFormated,
		"status":      http.StatusOK,
	})
}

func (h *UtilisateurHandler) ChangerDeSiteHandler(c *gin.Context) {
	utilisateurID, _ := strconv.Atoi(c.Param("id"))
	siteID, _ := strconv.Atoi(c.Param("siteID"))

	if utilisateurID < 1 || siteID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Erreur de paramètre",
			"status": http.StatusBadRequest,
		})
		return
	}

	err := h.service.ChangerDeSite(uint(utilisateurID), uint(siteID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	cacheKey := "utilisateur:" + fmt.Sprintf("%d", utilisateurID)
	_ = configs.DeleteCache(cacheKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "Changement de site effectué avec succès",
		"status":  http.StatusOK,
	})
}

func (h *UtilisateurHandler) ModifierSonCompteHandler(c *gin.Context) {
	var data map[string]any

	utilisateurID, err := c.Get("utilisateurID")

	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Vous n'êtes pas connecté(e)s",
			"status": http.StatusUnauthorized,
		})
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}

	if data["refresh_token"] == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Un parametre manquant",
			"status": http.StatusUnauthorized,
		})
		return
	}

	refreshToken := data["refresh_token"].(string)
	delete(data, "refresh_token")

	nouveauRefreshToken, errService := h.service.ModifierSonCompte(utilisateurID.(uint), refreshToken, data)

	if errService != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	if nouveauRefreshToken != "" {
		c.JSON(http.StatusOK, gin.H{
			"message":       "Compte modifié avec succès",
			"refresh_token": nouveauRefreshToken,
			"status":        http.StatusOK,
		})
		return
	}

	cacheKey := "utilisateur:" + fmt.Sprintf("%d", utilisateurID)
	_ = configs.DeleteCache(cacheKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "Compte modifié avec succès",
		"status":  http.StatusOK,
	})
}

func (h *UtilisateurHandler) ModifierSonMotDePasseHandler(c *gin.Context) {
	var data struct {
		Ancien  string `json:"ancien"`
		Nouveau string `json:"nouveau"`
	}

	utilisateurID, err := c.Get("utilisateurID")

	if !err {
		c.JSON(401, gin.H{
			"error":  "vous n'êtes pas autorisé à accéder à cette ressource",
			"status": 401,
		})
		return
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"error":  err.Error(),
			"status": 400,
		})
		return
	}

	errService := h.service.ModifierSonMotDePasse(utilisateurID.(uint), data.Ancien, data.Nouveau)

	if errService != nil {
		c.JSON(400, gin.H{
			"error":  errService.Error(),
			"status": 400,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "mot de passe modifié avec succès",
		"status":  200,
	})

}

func (h *UtilisateurHandler) SupprimerUnCompteHandler(c *gin.Context) {
	utilisateurID, _ := strconv.Atoi(c.Param("id"))

	err := h.service.SupprimerUnCompte(uint(utilisateurID))

	if err != nil {
		if strings.Contains(err.Error(), "cet utilisateur n'existe pas") {
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

	cacheKey := "utilisateur:" + strconv.Itoa(utilisateurID)
	_ = configs.DeleteCache(cacheKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "Compte supprimé avec succès",
		"status":  http.StatusOK,
	})
}
