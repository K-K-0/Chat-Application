package controllers

import (
	"Chat/database"
	"Chat/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateRoom(c *gin.Context) {
	var Room models.Room

	if err := c.ShouldBindJSON(&Room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&Room).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error while creating Room": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Room)

}

func JoinRoom(c *gin.Context) {
	userID := c.GetInt64("user_id")
	roomID := c.Param("id")

	var Room models.Room
	if err := database.DB.First(&Room, "id = ?", roomID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room not found"})
		return
	}

	var existing models.RoomMember
	if err := database.DB.First(&existing, "room_id = ? AND user_id = ?", Room.Id, userID).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "already Joined", "chair": existing.ChairPosition})
		return
	}

	var count int64
	database.DB.Model(&models.RoomMember{}).Where("room_id = ?", roomID).Count(&count)
	if count >= int64(Room.MaxSeat) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Room is full"})
		return
	}

	var taken []int

	database.DB.Model(&models.RoomMember{}).Where("room_id = ?", Room.Id).Pluck("chair position", &taken)

	chairMap := make(map[int]bool)
	for _, c := range taken {
		chairMap[c] = true
	}

	var chairPos int

	for i := 1; i <= Room.MaxSeat; i++ {
		if !chairMap[i] {
			chairPos = i
			break
		}
	}

	member := models.RoomMember{
		UserID:        uint(userID),
		RoomID:        Room.Id,
		ChairPosition: uint(chairPos),
		LightOn:       true,
		JoinedAt:      time.Now(),
	}

	if err := database.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to join room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "joined", "chair": chairPos})
}

func LeaveRoom(c *gin.Context) {
	userID := c.GetInt("user_id")
	roomID := c.Param("id")

	if err := database.DB.Where("room_id = ? AND user_id = ?", roomID, userID).
		Delete(&models.RoomMember{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Feailed to leave Room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left"})

}
