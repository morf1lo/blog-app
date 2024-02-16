package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/morf1lo/blog-app/internal/models"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(401, gin.H{"error": "User is not authorized"})
			c.Abort()
			return
		}

		parsedToken, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "Cannot read authentication token"})
			c.Abort()
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			c.JSON(401, gin.H{"error": "Authentication token is not valid"})
			c.Abort()
			return
		}

		token := models.Token{
			ID: int64(claims["id"].(float64)),
		}

		c.Set("token", token)
		c.Next()
	}
}
