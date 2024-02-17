package service

import (
	"database/sql"
	"encoding/json"

	"github.com/morf1lo/blog-app/internal/models"
)

type PostService struct {
	db *sql.DB
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) CreatePost(post models.Post) error {
	_, err := s.db.Exec("INSERT INTO posts(author, title, text) VALUES(?, ?, ?)", post.Author, post.Title, post.Text)
	if err != nil {
		return errInternalServer
	}
	return nil
}

func (s *PostService) GetAuthorPosts(authorID int64) ([]models.Post, error) {
	rows, err := s.db.Query("SELECT * FROM posts WHERE author = ?", authorID)
	if err != nil {
		return []models.Post{}, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var authorJSON string
		if err := rows.Scan(&post.ID, &authorJSON, &post.Title, &post.Text, &post.Likes); err != nil {
			return []models.Post{}, errInternalServer
		}

		if err := json.Unmarshal([]byte(authorJSON), &post.Author); err != nil {
			return []models.Post{}, errInternalServer
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return []models.Post{}, err
	}

	return posts, nil
}

func (s *PostService) UpdatePost(updateQuery string, values []interface{}) error {
	_, err := s.db.Exec(updateQuery, values...)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) LikePost(postID int64, userID int64) error {
	// Checking post existence
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errPostNotFound
	}

	// Checking if user has already likes a post
	var alreadyLiked bool
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user = ? AND post = ?)", userID, postID).Scan(&alreadyLiked)
	if err != nil {
		return err
	}

	if alreadyLiked {
		_, err := s.db.Exec("UPDATE posts SET likes = likes - 1 WHERE id = ?", postID)
		if err != nil {
			return errInternalServer
		}

		_, err = s.db.Exec("DELETE FROM likes WHERE user = ? AND post = ?", userID, postID)
		if err != nil {
			return errInternalServer
		}
	} else {
		_, err := s.db.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
		if err != nil {
			return errInternalServer
		}

		_, err = s.db.Exec("INSERT INTO likes(user, post) VALUES(?, ?)", userID, postID)
		if err != nil {
			return errInternalServer
		}
	}

	return nil
}

func (s *PostService) DeletePost(postID int64, userID int64) error {
	if err := deletePostDataFromDB(s.db, postID, userID); err != nil {
		return err
	}
	return nil
}

func deletePostDataFromDB(db *sql.DB, postID int64, userID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := [2]string{
		"DELETE FROM posts WHERE id = ? AND author = ?",
		"DELETE FROM comments WHERE JSON_EXTRACT(post, '$.id') = ? AND JSON_EXTRACT(post, '$.author') = ?",
	}

	for i, query := range queries {
		if i == 0 {
			_, err := tx.Exec(query, postID, userID)
			if err != nil {
				return err
			}
		} else {
			_, err := tx.Exec(query, postID, userID)
			if err != nil {
				return err
			}
		}
	}

	_, err = tx.Exec("DELETE FROM likes WHERE post = ?", postID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostService) GetUserLikes(user models.Token) ([]models.Post, error) {
	var postIDs []int64
	rows, err := s.db.Query("SELECT post FROM likes WHERE user = ?", user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postID int64
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}
		postIDs = append(postIDs, postID)
	}

	var posts []models.Post
	for _, postID := range postIDs {
		var post models.Post
		err := s.db.QueryRow("SELECT id, author, title, likes FROM posts WHERE id = ?", postID).Scan(&post.ID, &post.Author, &post.Title, &post.Likes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
