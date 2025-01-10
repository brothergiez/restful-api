package main

import (
	"log"
	"os"

	"github.com/brothergiez/restful-api/handlers"
	"github.com/brothergiez/restful-api/middlewares"
	"github.com/brothergiez/restful-api/repositories"
	"github.com/brothergiez/restful-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	repo := repositories.NewArticleRepository()
	handler := handlers.NewArticleHandler(repo)

	router := gin.Default()

	router.Use(middlewares.LoggingMiddleware())

	routes.RegisterArticleRoutes(router, handler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
