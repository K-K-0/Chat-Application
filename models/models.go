package models

import "time"

type User struct {
	Id        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex:not null"`
	Password  string `gorm:"not null"`
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time `gorm:"default:null"`
}

type Room struct {
	Id        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	CreatedBy uint   `gorm:"not null"`
	Creator   User   `gorm:"foreignKey:CreatedBy"`
	MaxSeat   int    `gorm:"default: 10"`
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type RoomMember struct {
	Id            uint `gorm:"primaryKey"`
	UserID        uint `gorm:"not null"`
	User          User `gorm:"foreignKey:UserID"`
	RoomID        uint `gorm:"not null"`
	Room          Room `gorm:"foreignKey:RoomID"`
	ChairPosition uint
	LightOn       bool
	JoinedAt      time.Time
}

type Message struct {
	Id        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	RoomID    uint   `gorm:"not null"`
	Room      Room   `gorm:"foreignKey:RoomID"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
}
