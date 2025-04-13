package main

import (
	"assignment-23/config"
	"assignment-23/models/entity"
	"assignment-23/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.ConnectDatabase()
	db.AutoMigrate(&entity.Product{}, &entity.Inventory{}, &entity.Order{}, &entity.OrderItem{})
	router := gin.Default()

	routes.SetupRoutes(router, db)
	err = router.Run(":8080")

	if err != nil {
		return
	}
}
