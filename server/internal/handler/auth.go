package handler

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/morf1lo/blog-app/internal/models"
	"github.com/morf1lo/blog-app/internal/utils/auth"
)

func (h *Handler) signUp(c *gin.Context) {
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

	activationLink := uuid.New()

	userID, err := h.services.Authorization.CreateUser(user, activationLink.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.services.Mail.SendActivationLink([]string{user.Email}, os.Getenv("SERVER_URL") + "/api/auth/activate/" + activationLink.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := auth.CreateSendToken(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) activate(c *gin.Context) {
	activationLink := c.Param("link")

	if err := h.services.Authorization.Activate(activationLink); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) signIn(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.services.Authorization.SignIn(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := auth.CreateSendToken(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) requestToResetPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"email,required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resetTokenExpiry := time.Now().Add(time.Hour * 12)
	if err := h.services.Authorization.SaveResetToken(request.Email, token, resetTokenExpiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resetLink := os.Getenv("CLIENT_URL") + "/resetpass/" + token

	if err := h.services.Mail.SendResetPasswordLink([]string{request.Email}, resetLink); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) resetPassword(c *gin.Context) {
	token := c.Param("token")

	var request struct {
		NewPassword string `json:"new_password" binding:"min=8,max=32,required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.services.Authorization.ResetPassword(token, request.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
