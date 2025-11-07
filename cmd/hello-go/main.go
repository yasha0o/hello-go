package main

import (
	"context"
	"log"
	"os"

	"hello-go/internal/controller"
	"hello-go/internal/mapper"
	"hello-go/internal/repository"
	"hello-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

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

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/archive_db?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("unexpeced error while run db pool")
		return
	}
	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
		return
	}

	archiveMapper := mapper.NewArchiveMapper()
	archiveRepository := repository.NewArchiveRepository(pool)

	archiveService := service.NewArchiveService(archiveRepository, archiveMapper)

	archiveController := controller.NewArchiveController(archiveService)
	archiveController.Init(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("unexpeced error while run controllers")
	}
}
