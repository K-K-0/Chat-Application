package controllers

import (
	"Chat/database"
	"Chat/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Registration(c *gin.Context) {
	var User models.User

	if err := c.ShouldBindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while Creating User"})
		return
	}

	c.JSON(http.StatusOK, User)
}
