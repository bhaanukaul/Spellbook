package Utils

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDatabaseConnection() (db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("./Spellbook.db"), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect database")
	}
	return db
}
