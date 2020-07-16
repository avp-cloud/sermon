package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "sermon.db"
	}
	database, err := gorm.Open("sqlite3", dbPath)

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&DBService{})

	DB = database
}
