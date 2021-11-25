package Utils

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDatabaseConnection() (db *gorm.DB) {
	spellbookDB := GetSpellbookDBLocation()
	db, err := gorm.Open(sqlite.Open(spellbookDB), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect database")
	}
	return db
}

func GetSpellbookDBLocation() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configFile := fmt.Sprintf("%s/.config/Spellbook/Spellbook.ini", userHome)

	cfg, err := ini.Load(configFile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	spellbook := cfg.Section("").Key("spellbookdb").String()
	return spellbook
}
