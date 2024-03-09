package service

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/morf1lo/blog-app/internal/models"
	"github.com/morf1lo/blog-app/internal/utils/auth"
)

type AuthService struct {
	db *sql.DB

	mail Mail
}

func NewAuthService(db *sql.DB, mail Mail) *AuthService {
	return &AuthService{
		db: db,
		mail: mail,
	}
}

func (s *AuthService) CreateUser(user models.User) (int64, error) {
	activationLink := uuid.New()

	insertedUser, err := s.db.Exec("INSERT INTO users(username, email, password, activation_link) VALUES(?, ?, ?, ?)", user.Username, user.Email, user.Password, activationLink)
	if err != nil {
		return 0, err
	}

	id, err := insertedUser.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := s.mail.sendActivationLink([]string{user.Email}, "http://localhost:8080/api/auth/activate/" + activationLink.String()); err != nil {
		return 0, nil
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

	_, err := s.db.Exec("UPDATE users SET activated = true WHERE activation_link = ?", activationLink)
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
