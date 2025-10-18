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
	router.POST("/shorten", routes.ShortenURL)
	// router.GET("/:shortenedURL", routes.GetURL)
	// router.POST("/addTag", routes.AddTag)
	// router.POST("/editUrl", routes.EditURL)
	// router.POST("/deleteUrl", routes.DeleteURL)
}
