package middleware

import (
	"assignment-23/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Membaca header Authorization dari request
		authHeader := c.GetHeader("Authorization")

		// Memeriksa apakah header Authorization ada dan tidak kosong
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Memeriksa apakah token diawali dengan "Bearer "
		parts := strings.Split(authHeader, " ")

		// Memeriksa indeks 0 dan panjang parts
		userId, err := utils.ValidateToken(parts[1])

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userId)
		c.Next()
	}
}
