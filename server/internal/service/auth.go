package service

import (
	"database/sql"

	"github.com/morf1lo/blog-app/internal/models"
	"github.com/morf1lo/blog-app/internal/utils/auth"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) CreateUser(user models.User) (int64, error) {
	insertedUser, err := s.db.Exec("INSERT INTO users(username, email, password) VALUES(?, ?, ?)", user.Username, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := insertedUser.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
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
