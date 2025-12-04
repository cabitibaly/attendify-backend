package handlers

import (
	"net/http"
	"strconv"

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

	utilisateur, errGetInfo := h.service.MesInformations(utilisateurID.(uint))

	if errGetInfo != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errGetInfo.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"utilisateur": utilisateur,
		"status":      http.StatusOK,
	})
}

func (h *UtilisateurHandler) TousLesEmployesHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 20
	}

	utilsateurs, hasNextPage, total, err := h.service.TousLesEmployes(page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"utilisateurs": utilsateurs,
		"hasNextPage":  hasNextPage,
		"total":        total,
		"status":       http.StatusOK,
	})
}

func (h *UtilisateurHandler) LireUnEmployeHandler(c *gin.Context) {
	utilisateurID, _ := strconv.Atoi(c.Param("id"))

	utilsateur, err := h.service.LireUnEmploye(uint(utilisateurID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"utilisateur": utilsateur,
		"status":      http.StatusOK,
	})
}

func (h *UtilisateurHandler) ModifierSonCompteHandler(c *gin.Context) {
	var data map[string]any

	utilisateurID, err := c.Get("utilisateurID")
	token, errToken := c.Get("token")

	if !err || !errToken {
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

	nouveauToken, errService := h.service.ModifierSonCompte(utilisateurID.(uint), token.(string), data)

	if errService != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  errService.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	if nouveauToken != "" {
		c.SetCookie(
			"jwt",
			nouveauToken,
			3600*24*3,
			"/",
			"",
			false,
			true,
		)
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Compte supprimé avec succès",
		"status":  http.StatusOK,
	})
}
