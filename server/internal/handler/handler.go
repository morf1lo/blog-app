package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/middlewares"
	"github.com/morf1lo/blog-app/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	user := router.Group("/api/users")
	{
		user.POST("/signup", h.Signup)
		user.POST("/login", h.Login)
		user.DELETE("/delete", middlewares.AuthMiddleware(), h.DeleteUser)
		user.POST("/logout", middlewares.AuthMiddleware(), h.Logout)
		user.GET("/:uname", middlewares.AuthMiddleware(), h.GetUser)
		user.POST("/avatar", middlewares.AuthMiddleware(), h.SetAvatar)
		user.POST("/follow/:id", middlewares.AuthMiddleware(), h.Follow)
	}

	post := router.Group("/api/posts")
	{
		post.POST("/create", middlewares.AuthMiddleware(), h.createPost)
		post.GET("/:id", middlewares.AuthMiddleware(), h.getAuthorPosts)
		post.PATCH("/:id", middlewares.AuthMiddleware(), h.updatePost)
		post.POST("/like/:id", middlewares.AuthMiddleware(), h.likePost)
		post.DELETE("/:id", middlewares.AuthMiddleware(), h.deletePost)
	}

	comment := router.Group("/api/comments")
	{
		comment.POST("/add/:post", middlewares.AuthMiddleware(), h.addComment)
		comment.GET("/:post", middlewares.AuthMiddleware(), h.getAllPostComments)
		comment.DELETE("/:post/:comment", middlewares.AuthMiddleware(), h.deleteComment)
	}
}
