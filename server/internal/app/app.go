package app

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/config"
	"github.com/morf1lo/blog-app/internal/db"
	"github.com/morf1lo/blog-app/internal/routes"
	"github.com/morf1lo/blog-app/internal/services"
)

func Run() {
	config.Init()

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.New()

	userService := services.NewUserService(db)
	postService := services.NewPostService(db)
	commentService := services.NewCommentService(db)

	router.Static("/public", "./public")

	router.SetTrustedProxies(nil)

	routes.SetupUserRoutes(router, *userService)
	routes.SetupPostRoutes(router, *postService)
	routes.SetupCommentRoutes(router, *commentService)

	router.Run(":8080")
}
