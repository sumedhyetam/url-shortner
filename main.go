package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sumedhyetam/url-shortner/api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	router := gin.Default()

	setupRoutes(router)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Fatal(router.Run(":" + port))
}

func setupRoutes(router *gin.Engine) {
	router.POST("/api/v1", routes.ShortenURL)
	router.GET("/api/v1/:shortID", routes.GetByShortID)
	router.POST("/api/v1/addTag", routes.AddTag)
	router.PUT("/api/v1/:shortID", routes.EditURL)
	router.DELETE("/api/v1/:shortID", routes.DeleteURL)
}
