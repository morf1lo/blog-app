package app

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/morf1lo/blog-app/internal/config"
	"github.com/morf1lo/blog-app/internal/db"
	"github.com/morf1lo/blog-app/internal/handler"
	"github.com/morf1lo/blog-app/internal/service"
)

func Run() {
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	services := service.NewService(db)
	handlers := handler.NewHandler(services)

	router := gin.New()

	router.Static("/public", "./public")

	router.SetTrustedProxies(nil)

	handlers.SetupRoutes(router)

	router.Run(":8080")
}
