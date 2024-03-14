package auth

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": id,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	jwt, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func CreateSendToken(c *gin.Context, userID int64) error {
	jwt, err := generateToken(userID)
	if err != nil {
		return err
	}

	c.SetCookie("jwt", jwt, int(time.Now().Add(time.Hour * 24).Unix()), "/", "localhost", true, true)
	return nil
}

func GenerateResetToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(base64.URLEncoding.EncodeToString(tokenBytes), "%", ""), nil
}
