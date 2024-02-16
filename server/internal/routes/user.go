package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/handlers"
	"github.com/morf1lo/blog-app/internal/middlewares"
	"github.com/morf1lo/blog-app/internal/services"
)

func SetupUserRoutes(router *gin.Engine, userService services.UserService) {
	router.POST("/api/users/signup", handlers.Signup(userService))
	router.POST("/api/users/login", handlers.Login(userService))
	router.DELETE("/api/users/delete", middlewares.AuthMiddleware(), handlers.DeleteUser(userService))
	router.POST("/api/users/logout", middlewares.AuthMiddleware(), handlers.Logout())
	router.GET("/api/users/:uname", middlewares.AuthMiddleware(), handlers.GetUser(userService))
	router.POST("/api/users/avatar", middlewares.AuthMiddleware(), handlers.SetAvatar(userService))
	router.POST("/api/users/follow/:id", middlewares.AuthMiddleware(), handlers.Follow(userService))
}
