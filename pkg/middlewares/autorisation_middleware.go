package middlewares

import "github.com/gin-gonic/gin"

func AutorisationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, err := c.Get("roleID")

		if !err {
			c.JSON(401, gin.H{
				"error":  "vous n'avez pas les droits pour effectuer cette action",
				"status": 401,
			})
			c.Abort()
			return
		}

		if roleID.(uint) != 1 {
			c.JSON(401, gin.H{
				"error":  "vous n'avez pas les droits pour effectuer cette action",
				"status": 401,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
