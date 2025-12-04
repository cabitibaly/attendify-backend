package middlewares

import (
	"time"

	"github.com/cabitibaly/internal/services"
	"github.com/cabitibaly/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("jwt")

		if err != nil {
			c.JSON(401, gin.H{
				"error":  "impossible de vous authentifier",
				"status": 401,
			})
			c.Abort()
			return
		}

		jwtClaims, err := utils.ValidateToken(cookie)

		if err != nil {
			c.JSON(401, gin.H{
				"error":  "impossible de vous authentifier",
				"status": 401,
			})
			c.Abort()
			return
		}

		if jwtClaims.ExpiresAt.Unix() < time.Now().Unix() {
			c.JSON(401, gin.H{
				"error":  "votre session a expiré",
				"status": 401,
			})
			c.Abort()
			return
		}

		c.Set("utilisateurID", jwtClaims.UtilisateurID)
		c.Set("email", jwtClaims.Email)
		c.Set("telephone", jwtClaims.Telephone)
		c.Set("roleID", jwtClaims.RoleID)
		c.Set("token", cookie)

		c.Next()
	}
}
