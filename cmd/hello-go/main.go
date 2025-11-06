package main

import (
	"fmt"

	"hello-go/internal/controller"
	"hello-go/internal/service"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "hello-go/docs"
)

// @title Archive API
// @version 1.0
// @description API для работы с архивами
// @host localhost:8080
// @BasePath /
func main() {
	router := gin.Default()

	archiveService := service.NewArchiveService()

	archiveController := controller.NewArchiveController(archiveService)
	archiveController.Init(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(":8080")
	if err != nil {
		fmt.Printf("unexpeced error while run controllers")
	}
}
