package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/signup", h.signUp)
		auth.POST("/signin", h.signIn)
	}

	user := router.Group("/api/users")
	{
		user.POST("/logout", h.authMiddleware, h.logout)
		user.GET("/id/:id", h.authMiddleware, h.getUserById)
		user.GET("/name/:uname", h.authMiddleware, h.getUserByUsername)
		user.POST("/avatar", h.authMiddleware, h.setAvatar)
		user.POST("/follow/:id", h.authMiddleware, h.follow)
		user.GET("/:id/followers", h.authMiddleware, h.getUserFollowers)
		user.GET("/:id/follows", h.authMiddleware, h.getUserFollows)
		user.DELETE("/delete", h.authMiddleware, h.deleteUser)
	}

	post := router.Group("/api/posts")
	{
		post.POST("/create", h.authMiddleware, h.createPost)
		post.GET("/:id", h.authMiddleware, h.getPostById)
		post.GET("/user/:id", h.authMiddleware, h.getAuthorPosts)
		post.PATCH("/:id", h.authMiddleware, h.updatePost)
		post.POST("/like/:id", h.authMiddleware, h.likePost)
		post.GET("/my/likes", h.authMiddleware, h.getUserLikes)
		post.DELETE("/:id", h.authMiddleware, h.deletePost)
	}

	comment := router.Group("/api/comments")
	{
		comment.POST("/add/:post", h.authMiddleware, h.addComment)
		comment.GET("/:post", h.authMiddleware, h.getAllPostComments)
		comment.DELETE("/:post/:comment", h.authMiddleware, h.deleteComment)
	}
}
