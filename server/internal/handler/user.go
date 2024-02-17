package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/models"
	"github.com/morf1lo/blog-app/internal/utils"
	"github.com/morf1lo/blog-app/internal/utils/auth"
)

// Create and add token to cookie
func createSendToken(c *gin.Context, token models.Token) error {
	jwt, err := auth.GenerateToken(token.ID)
	if err != nil {
		return err
	}

	c.SetCookie("jwt", jwt, int(time.Now().Add(time.Hour * 24).Unix()), "/", "localhost", true, true)
	return nil
}

func (h *Handler) signup(c *gin.Context) {
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

	token, err := h.services.User.CreateUser(user)
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

func (h *Handler) signin(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.services.User.SignIn(user)
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

func (h *Handler) deleteUser(c *gin.Context) {
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

	if err := h.services.User.DeleteUser(user, confirmPassword.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", "", -1, "/", "localhost", true, true)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) getUser(c *gin.Context) {
	username := c.Param("uname")

	user, err := h.services.User.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

func (h *Handler) setAvatar(c *gin.Context) {
	user := utils.GetUser(c)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.services.User.SetAvatar(c, file, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) follow(c *gin.Context) {
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

	if err := h.services.User.Follow(user, uint64(following)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) getUserFollowers(c *gin.Context) {
	user := utils.GetUser(c)

	followers, err := h.services.User.GetUserFollowers(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": followers})
}

func (h *Handler) getUserFollows(c *gin.Context) {
	user := utils.GetUser(c)

	followers, err := h.services.User.GetUserFollows(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": followers})
}

