package service

import (
	"database/sql"
	"time"

	"github.com/morf1lo/blog-app/internal/models"
	"github.com/morf1lo/blog-app/internal/utils/auth"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{
		db: db,
	}
}

func (s *AuthService) CreateUser(user models.User, activationLink string) (int64, error) {
	insertedUser, err := s.db.Exec("INSERT INTO users(username, email, password, activation_link) VALUES(?, ?, ?, ?)", user.Username, user.Email, user.Password, activationLink)
	if err != nil {
		return 0, err
	}

	id, err := insertedUser.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *AuthService) Activate(activationLink string) error {
	var exists bool
	if err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE activation_link = ?)", activationLink).Scan(&exists); err != nil {
		return errInternalServer
	}
	if !exists {
		return errUserNotFound
	}

	_, err := s.db.Exec("UPDATE users SET activated = true, activation_link = null WHERE activation_link = ?", activationLink)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SignIn(user models.User) (int64, error) {
	var existingUser models.User
	err := s.db.QueryRow("SELECT id, username, password, avatar FROM users WHERE username = ? OR email = ?", user.Username, user.Email).Scan(&existingUser.ID, &existingUser.Username, &existingUser.Password, &existingUser.Avatar)
	if err != nil {
		return 0, errInvalidCredenials
	}

	matchPassword := auth.VerifyPassword([]byte(existingUser.Password), []byte(user.Password))
	if !matchPassword {
		return 0, errInvalidCredenials
	}

	return existingUser.ID, nil
}

func (s *AuthService) SaveResetToken(email string, token string, tokenExpiry time.Time) error {
	_, err := s.db.Exec("UPDATE users SET reset_token = ?, reset_token_expiry = ? WHERE email = ?", token, tokenExpiry, email)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ResetPassword(token string, newPassword string) error {
	var user models.User
	var resetTokenExpiryStr string
	if err := s.db.QueryRow("SELECT id, reset_token, reset_token_expiry FROM users WHERE reset_token = ?", token).Scan(&user.ID, &user.ResetToken, &resetTokenExpiryStr); err != nil {
		return err
	}

	resetTokenExpiry, err := time.Parse("2006-01-02 15:04:05", resetTokenExpiryStr)
	if err != nil {
		return err
	}
	user.ResetTokenExpiry = resetTokenExpiry

	if time.Now().After(user.ResetTokenExpiry) {
		_, err := s.db.Exec("UPDATE users SET reset_token = null, reset_token_expiry = null WHERE reset_token = ?", token)
		if err != nil {
			return err
		}
		return errTokenHasExpired
	}

	hash, err := auth.HashPassword([]byte(newPassword))
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE users SET password = ?, reset_token = null, reset_token_expiry = null WHERE id = ?", hash, user.ID)
	if err != nil {
		return err
	}

	return nil
}
