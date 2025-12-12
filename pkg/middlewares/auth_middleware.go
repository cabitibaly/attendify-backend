package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Aucun token d'accès n'a été fourni",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Le token d'accès n'est pas valide",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		jwtClaims, err := utils.ValidateToken(parts[1])

		if err != nil {

			if strings.Contains(err.Error(), "invalid") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":  "Le token d'accès n'est pas valide",
					"status": http.StatusUnauthorized,
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  err.Error(),
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		if jwtClaims.ExpiresAt.Unix() < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "votre session a expiré",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		c.Set("utilisateurID", jwtClaims.UtilisateurID)
		c.Set("email", jwtClaims.Email)
		c.Set("telephone", jwtClaims.Telephone)
		c.Set("roleID", jwtClaims.RoleID)

		c.Next()
	}
}
