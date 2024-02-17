package utils

import (
	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/models"
)

func GetUser(c *gin.Context) models.Token {
	claims, ok := c.Get("token")
	if !ok {
		return models.Token{}
	}

	user, ok := claims.(models.Token)
	if !ok {
		return models.Token{}
	}

	return user
}
