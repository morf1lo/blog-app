package service

import (
	"database/sql"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/morf1lo/blog-app/internal/models"
)

type User interface {
	CreateUser(user models.User) (models.Token, error)
	SignIn(user models.User) (models.Token, error)
	DeleteUser(user models.Token, confirmPassword string) error
	GetUserByUsername(username string) (interface{}, error)
	SetAvatar(c *gin.Context, file *multipart.FileHeader, user *models.Token) error
	Follow(user models.Token, following uint64) error
	GetUserFollowers(user models.Token) ([]models.User, error)
	GetUserFollows(user models.Token) ([]models.User, error)
}

type Post interface {
	CreatePost(post models.Post) error
	GetAuthorPosts(authorID int64) ([]models.Post, error)
	UpdatePost(updateQuery string, values []interface{}) error
	LikePost(postID int64, userID int64) error
	DeletePost(postID int64, userID int64) error
	GetUserLikes(user models.Token) ([]models.Post, error)
}

type Comment interface {
	AddComment(comment models.Comment, userID int64, postID int64) error
	GetAllPostComments(postID int64) ([]models.Comment, error)
	DeleteComment(commentID int64, userID int64, postID int64) error
}

type Service struct {
	User
	Post
	Comment
}

func NewService(db *sql.DB) *Service {
	return &Service{
		User: NewUserService(db),
		Post: NewPostService(db),
		Comment: NewCommentService(db),
	}
}
