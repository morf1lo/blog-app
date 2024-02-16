package service

import (
	"database/sql"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/morf1lo/blog-app/internal/models"
	"github.com/morf1lo/blog-app/internal/utils/auth"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(user models.User) (models.Token, error) {
	insertedUser, err := s.db.Exec("INSERT INTO users(username, email, password) VALUES(?, ?, ?)", user.Username, user.Email, user.Password)
	if err != nil {
		return models.Token{}, err
	}

	id, err := insertedUser.LastInsertId()
	if err != nil {
		return models.Token{}, err
	}

	return models.Token{ID: id}, nil
}

func (s *UserService) Login(user models.User) (models.Token, error) {
	var existingUser models.User
	err := s.db.QueryRow("SELECT id, username, password, avatar FROM users WHERE username = ? OR email = ?", user.Username, user.Email).Scan(&existingUser.ID, &existingUser.Username, &existingUser.Password, &existingUser.Avatar)
	switch {
	case err == sql.ErrNoRows:
		return models.Token{}, errInvalidCredenials
	case err != nil:
		return models.Token{}, err
	}

	matchPassword := auth.VerifyPassword([]byte(existingUser.Password), []byte(user.Password))
	if !matchPassword {
		return models.Token{}, errInvalidCredenials
	}

	return models.Token{ID: existingUser.ID}, nil
}

func (s *UserService) DeleteUser(user models.Token, confirmPassword string) error {
	var password string
	err := s.db.QueryRow("SELECT password FROM users WHERE id = ?", user.ID).Scan(&password)
	if err != nil {
		return err
	}

	if matchPassword := auth.VerifyPassword([]byte(password), []byte(confirmPassword)); !matchPassword {
		return errInvalidPassword
	}

	if err := deleteUserData(s.db, user.ID); err != nil {
		return errInternalServer
	}

	return nil
}

func deleteUserData(db *sql.DB, userID int64) error {
	// Delete user profile picture
	path := "public/avatars"

	fileNamePattern := strconv.FormatInt(userID, 10) + ".*"

	files, err := filepath.Glob(filepath.Join(path, fileNamePattern))
	if err != nil {
		return err
	}

	for _, filePath := range files {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	// Delete user data from Database
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := [4]string{
		"DELETE FROM users WHERE id = ?",
		"DELETE FROM posts WHERE author = ?",
		"DELETE FROM comments WHERE author = ?",
		"DELETE FROM likes WHERE user = ?",
	}

	for _, query := range queries {
		_, err := tx.Exec(query, userID)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec("DELETE FROM followers WHERE user = ? OR following = ?", userID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *UserService) GetUserByUsername(username string) (interface{}, error) {
	var user struct {
		ID				uint64 `json:"id"`
		Username	string `json:"username"`
		Email			string `json:"email"`
		Avatar		string `json:"avatar"`
		CreatedAt	string `json:"created_at"`
	}

	err := s.db.QueryRow("SELECT id, username, email, avatar, created_at FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.Avatar, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) SetAvatar(c *gin.Context, file *multipart.FileHeader, user *models.Token) error {
	uploadPath := "public/avatars"

	recentFileNamePattern := strconv.FormatInt(user.ID, 10) + ".*"

	recentFiles, err := filepath.Glob(filepath.Join(uploadPath, recentFileNamePattern))
	if err != nil {
		return err
	}

	for _, recentFilePath := range recentFiles {
		if err := os.Remove(recentFilePath); err != nil {
			return err
		}
	}

	fileExt := filepath.Ext(file.Filename)
	fileName := strconv.FormatInt(int64(user.ID), 10) + fileExt
	filePath := filepath.Join(uploadPath, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return err
	}

	avatar := "http://localhost:8080/public/avatars/" + fileName
	_, err = s.db.Exec("UPDATE users SET avatar = ? WHERE id = ?", avatar, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Follow(user models.Token, following uint64) error {
	// Checking user existence
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", following).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errUserNotFound
	}

	var alreadyFollowed bool
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM followers WHERE user = ? AND following = ?)", user.ID, following).Scan(&alreadyFollowed)
	if err != nil {
		return err
	}

	if alreadyFollowed {
		_, err = s.db.Exec("DELETE FROM followers WHERE user = ? AND following = ?", user.ID, following)
		if err != nil {
			return err
		}
	} else {
		_, err = s.db.Exec("INSERT INTO followers(user, following) VALUES(?, ?)", user.ID, following)
		if err != nil {
			return err
		}
	}

	return nil
}
