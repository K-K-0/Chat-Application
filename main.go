package main

import (
	"Chat/controllers"
	"Chat/database"
	"Chat/middlewares"
	"net/http"

	"Chat/models"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Print("Env is not loading")
	}
}

func main() {

	database.Connect()
	err := database.DB.AutoMigrate(&models.User{}, &models.Room{}, &models.RoomMember{}, &models.Message{})

	if err != nil {
		log.Fatalf("AutoMigrate failed %v", err)
	}

	router := gin.Default()

	router.POST("/signup", controllers.Registration)

	router.Use(middlewares.AuthMiddleware())
	router.GET("/Me", func(c *gin.Context) {
		userID := c.MustGet("UserId").(string)
		c.JSON(http.StatusOK, gin.H{"message": "Authenticated", "user_id": userID})
	})

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "No auth required"})
	})

	port := os.Getenv("PORT")
	fmt.Print(port)

	router.Run(port)
}
