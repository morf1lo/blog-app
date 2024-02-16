package service

import (
	"database/sql"
	"encoding/json"

	"github.com/morf1lo/blog-app/internal/models"
)

type CommentService struct {
	db *sql.DB
}

func NewCommentService(db *sql.DB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) AddComment(comment models.Comment, userID int64, postID int64) error {
	// Checking post existence
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errPostNotFound
	}

	var postAuthorId int64
	err = s.db.QueryRow("SELECT author FROM posts WHERE id = ?", postID).Scan(&postAuthorId)
	if err != nil {
		return err
	}

	postData := models.CommentPost{
		ID: postID,
		Author: postAuthorId,
	}
	postDataJSON, err := json.Marshal(postData)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO comments (post, author, text) VALUES(?, ?, ?)", postDataJSON, userID, comment.Text)
	if err != nil {
		return err
	}

	return nil
}

func (s *CommentService) GetAllPostComments(postID int64) ([]models.Comment, error) {
	rows, err := s.db.Query("SELECT * FROM comments WHERE JSON_EXTRACT(post, '$.id') = ?", postID)
	if err != nil {
		return []models.Comment{}, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() { 
		var comment models.Comment
		var postDataJSON string
		if err := rows.Scan(&comment.ID, &postDataJSON, &comment.Author, &comment.Text); err != nil {
			return []models.Comment{}, err
		}

		if err := json.Unmarshal([]byte(postDataJSON), &comment.Post); err != nil {
			return []models.Comment{}, errInternalServer
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return []models.Comment{}, err
	}

	return comments, nil
}

func (s *CommentService) DeleteComment(commentID int64, userID int64, postID int64) error {
	var postAuthorId int64
	err := s.db.QueryRow("SELECT author FROM posts WHERE id = ?", postID).Scan(&postAuthorId)
	if err != nil {
		return err
	}

	var commentAuthorId int64
	err = s.db.QueryRow("SELECT author FROM comments WHERE id = ?", commentID).Scan(&commentAuthorId)
	if err != nil {
		return err
	}

	if userID == postAuthorId || userID == commentAuthorId {
		_, err = s.db.Exec("DELETE FROM comments WHERE JSON_EXTRACT(post, '$.id') = ? AND id = ?", postID, commentID)
		if err != nil {
			return errInternalServer
		}
	} else {
		return errNoAccess
	}

	return nil
}
