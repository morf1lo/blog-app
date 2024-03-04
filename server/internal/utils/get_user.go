package utils

import (
	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/models"
)

func GetUserFromRequest(c *gin.Context) *models.User {
	claims, _ := c.Get("user")

	user, ok := claims.(models.User)
	if !ok {
		return nil
	}

	return &user
}
