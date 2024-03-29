package service

import (
	"database/sql"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/models"
)

type Mail interface {
	SendActivationLink(to []string, link string) error
	SendResetPasswordLink(to []string, link string) error
}

type Authorization interface {
	CreateUser(user models.User, activationLink string) (int64, error)
	Activate(activationLink string) error
	SignIn(user models.User) (int64, error)
	SaveResetToken(email string, token string, tokenExpiry time.Time) error
	ResetPassword(token string, newPassword string) error
}

type User interface {
	DeleteUser(userID int64, confirmPassword string) error
	FindUserById(userID int64) (*models.User, error)
	FindUserByUsername(username string) (*models.User, error)
	SetAvatar(c *gin.Context, file *multipart.FileHeader, userID int64) error
	Follow(userID int64, followingID int64) error
	FindUserFollowers(userID int64) (*[]models.User, error)
	FindUserFollows(userID int64) (*[]models.User, error)
}

type Post interface {
	CreatePost(post models.Post) error
	FindPostById(postID int64) (*models.Post, error)
	FindAuthorPosts(authorID int64) (*[]models.Post, error)
	UpdatePost(updateOpts models.PostUpdateOptions, postID int64, userID int64) error
	LikePost(postID int64, userID int64) error
	DeletePost(postID int64, userID int64) error
	FindUserLikes(userID int64) (*[]models.Post, error)
	SearchPosts(query string) (*[]models.Post, error)
}

type Comment interface {
	AddComment(comment models.Comment, userID int64, postID int64) error
	FindAllPostComments(postID int64) (*[]models.Comment, error)
	DeleteComment(commentID int64, userID int64, postID int64) error
}

type Service struct {
	Mail
	Authorization
	User
	Post
	Comment
}

func NewService(db *sql.DB) *Service {
	return &Service{
		Mail: NewMailService(db),
		Authorization: NewAuthService(db),
		User: NewUserService(db),
		Post: NewPostService(db),
		Comment: NewCommentService(db),
	}
}
