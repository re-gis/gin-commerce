package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/re-gis/gin-commerce/database"
	"github.com/re-gis/gin-commerce/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading the .env file")
	}

	r := gin.Default()

	// Database
	database.Connect()

	// ROutes setup
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run()
}
