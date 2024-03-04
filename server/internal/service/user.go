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

func (s *UserService) DeleteUser(userID int64, confirmPassword string) error {
	var password string
	err := s.db.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&password)
	if err != nil {
		return err
	}

	if matchPassword := auth.VerifyPassword([]byte(password), []byte(confirmPassword)); !matchPassword {
		return errInvalidPassword
	}

	if err := deleteUserData(s.db, userID); err != nil {
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
		"DELETE FROM posts WHERE author_id = ?",
		"DELETE FROM comments WHERE author_id = ?",
		"DELETE FROM likes WHERE user_id = ?",
	}

	for _, query := range queries {
		_, err := tx.Exec(query, userID)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec("DELETE FROM followers WHERE user_id = ? OR following_id = ?", userID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *UserService) FindUserById(userID int64) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, username, email, avatar, created_at FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email, &user.Avatar, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, username, email, avatar, created_at FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.Avatar, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) SetAvatar(c *gin.Context, file *multipart.FileHeader, userID int64) error {
	uploadPath := "public/avatars"

	recentFileNamePattern := strconv.FormatInt(userID, 10) + ".*"

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
	fileName := strconv.FormatInt(int64(userID), 10) + fileExt
	filePath := filepath.Join(uploadPath, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return err
	}

	avatar := "http://localhost:8080/public/avatars/" + fileName
	_, err = s.db.Exec("UPDATE users SET avatar = ? WHERE id = ?", avatar, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Follow(userID int64, followingID int64) error {
	// Checking user existence
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", followingID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errUserNotFound
	}

	var alreadyFollowed bool
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM followers WHERE user_id = ? AND following_id = ?)", userID, followingID).Scan(&alreadyFollowed)
	if err != nil {
		return err
	}

	if alreadyFollowed {
		_, err = s.db.Exec("DELETE FROM followers WHERE user_id = ? AND following_id = ?", userID, followingID)
		if err != nil {
			return err
		}
	} else {
		_, err = s.db.Exec("INSERT INTO followers(user_id, following_id) VALUES(?, ?)", userID, followingID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserService) FindUserFollowers(userID int64) (*[]models.User, error) {
	var followerIDs []int64
	rows, err := s.db.Query("SELECT user_id FROM followers WHERE following_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followerID int64
		if err := rows.Scan(&followerID); err != nil {
			return nil, err
		}
		followerIDs = append(followerIDs, followerID)
	}

	var followers []models.User
	for _, followingID := range followerIDs {
		var follower models.User
		err := s.db.QueryRow("SELECT id, username, avatar FROM users WHERE id = ?", followingID).Scan(&follower.ID, &follower.Username, &follower.Avatar)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	return &followers, nil
}

func (s *UserService) FindUserFollows(userID int64) (*[]models.User, error) {
	var followingIDs []int64
	rows, err := s.db.Query("SELECT following_id FROM followers WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followingID int64
		if err := rows.Scan(&followingID); err != nil {
			return nil, err
		}
		followingIDs = append(followingIDs, followingID)
	}

	var follows []models.User
	for _, followingID := range followingIDs {
		var follow models.User
		err := s.db.QueryRow("SELECT id, username, avatar FROM users WHERE id = ?", followingID).Scan(&follow.ID, &follow.Username, &follow.Avatar)
		if err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}

	return &follows, nil
}
