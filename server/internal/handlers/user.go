package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/models"
	"github.com/morf1lo/blog-app/internal/services"
	"github.com/morf1lo/blog-app/internal/utils"
	"github.com/morf1lo/blog-app/internal/utils/auth"
)

// Create and add token to cookie
func createSendToken(c *gin.Context, token models.Token) error {
	jwt, err := auth.GenerateToken(token.ID, token.Username, token.Avatar)
	if err != nil {
		return err
	}

	c.SetCookie("jwt", jwt, int(time.Now().Add(time.Hour * 24).Unix()), "/", "localhost", true, true)
	return nil
}

func Signup(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if contains := strings.Contains(user.Password, " "); contains {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
			return
		}

		if err := user.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := auth.HashPassword([]byte(user.Password))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		user.Password = hash
		user.Username = strings.TrimSpace(user.Username)

		token, err := userService.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := createSendToken(c, token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func Login(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := userService.Login(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := createSendToken(c, token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func DeleteUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetUser(c)

		var reqBody map[string]interface{}
		if err := c.ShouldBind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		confirmPassword, exists := reqBody["confirm_password"]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide confirm password"})
			return
		}

		if err := userService.DeleteUser(user, confirmPassword.(string)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.SetCookie("jwt", "", -1, "/", "localhost", true, true)

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("jwt", "", -1, "/", "localhost", true, true)
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func GetUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("uname")

		user, err := userService.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
	}
}

func SetAvatar(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetUser(c)

		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := userService.SetAvatar(c, file, &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := createSendToken(c, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func Follow(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetUser(c)

		followingParam := c.Param("id")
		following, err := strconv.Atoi(followingParam)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.ID == int64(following) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot follow yourself"})
			return
		}

		if err := userService.Follow(user, uint64(following)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
