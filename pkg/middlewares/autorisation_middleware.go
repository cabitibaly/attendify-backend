package middlewares

import "github.com/gin-gonic/gin"

func AutorisationMiddleware(roleAutorise int) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, err := c.Get("roleID")

		if !err {
			c.JSON(401, gin.H{
				"error":  "vous n'êtes pas authentifié",
				"status": 401,
			})
			c.Abort()
			return
		}

		if roleID.(uint) != uint(roleAutorise) {
			c.JSON(403, gin.H{
				"error":  "vous n'avez pas les droits pour effectuer cette action",
				"status": 403,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
