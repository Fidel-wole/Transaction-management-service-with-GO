package middleware

import (
	"net/http"

	"github.com/Fidel-wole/Transaction_Management_Service/utils"
	"github.com/gin-gonic/gin"
)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
			return
		}

		// Check if the token starts with "Bearer " and remove it if present
		const bearerPrefix = "Bearer "
		if len(token) > len(bearerPrefix) && token[:len(bearerPrefix)] == bearerPrefix {
			token = token[len(bearerPrefix):]
		}

		// Update: Call VerifyToken and handle the returned userID
		userId, err := utils.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
			return
		}

		// Set the userID in the context
		c.Set("userId", userId)

		c.Next()
	}
}

