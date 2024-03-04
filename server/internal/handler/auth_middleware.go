package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	tokenCookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not authorized"})
		c.Abort()
		return
	}

	parsedToken, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read authentication token"})
		c.Abort()
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication token is not valid"})
		c.Abort()
		return
	}

	id := int64(claims["uid"].(float64))

	user, err := h.services.User.FindUserById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set("user", *user)
	c.Next()
}
