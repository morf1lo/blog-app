package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/handlers"
	"github.com/morf1lo/blog-app/internal/middlewares"
	"github.com/morf1lo/blog-app/internal/services"
)

func SetupPostRoutes(router *gin.Engine, postService services.PostService) {
	router.POST("/api/posts/create", middlewares.AuthMiddleware(), handlers.CreatePost(postService))
	router.GET("/api/posts/:id", middlewares.AuthMiddleware(), handlers.GetAuthorPosts(postService))
	router.PATCH("/api/posts/:id", middlewares.AuthMiddleware(), handlers.UpdatePost(postService))
	router.POST("/api/posts/like/:id", middlewares.AuthMiddleware(), handlers.LikePost(postService))
	router.DELETE("/api/posts/:id", middlewares.AuthMiddleware(), handlers.DeletePost(postService))
}
