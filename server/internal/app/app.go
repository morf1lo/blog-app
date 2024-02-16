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
	config.Init()

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	services := service.NewService(db)
	handler := handler.NewHandler(services)

	router := gin.New()

	router.Static("/public", "./public")

	router.SetTrustedProxies(nil)

	handler.SetupRoutes(router)

	router.Run(":8080")
}
