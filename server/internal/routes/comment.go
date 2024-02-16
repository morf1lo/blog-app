package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/handlers"
	"github.com/morf1lo/blog-app/internal/middlewares"
	"github.com/morf1lo/blog-app/internal/services"
)

func SetupCommentRoutes(router *gin.Engine, commentService services.CommentService) {
	router.POST("/api/comments/add/:post", middlewares.AuthMiddleware(), handlers.AddComment(commentService))
	router.GET("/api/comments/:post", middlewares.AuthMiddleware(), handlers.GetAllPostComments(commentService))
	router.DELETE("/api/comments/:post/:comment", middlewares.AuthMiddleware(), handlers.DeleteComment(commentService))
}
