package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "postgresql://postgres.bwlnihqelepcrqmdkleh:Aadity7531@@aws-0-ap-south-1.pooler.supabase.com:6543/postgres"
	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("unable to connect to the Database", err)
	}

}
