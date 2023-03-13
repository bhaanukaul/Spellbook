package Utils

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"gopkg.in/ini.v1"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SpellbookPing struct {
	Version string
}

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

func Error(msg string, err error) {
	fmt.Printf("%s, %v", msg, err)
	// os.Exit(1)
}

func AddKVToConfig(cf string, key string, value string, section string) {
	cfg, err := ini.Load(cf)
	if err != nil {
		Error("Failed to read file", err)
	}

	cfg.Section(section).Key(key).SetValue(value)
	cfg.SaveTo(cf)
}

func DoesSectionExist(cf string, section string) bool {
	cfg, err := ini.Load(cf)
	if err != nil {
		Error("Failed to read file", err)
		return false
	}

	sectionExists, err := cfg.GetSection(section)
	if err != nil {
		Error("Failed to get section", err)
		return false
	}

	log.Printf("section: %#v", sectionExists)
	return true
}

func GetKVFromConfig(cf string, key string, section string) string {
	cfg, err := ini.Load(cf)
	if err != nil {
		Error("Failed to read file", err)
	}

	value := cfg.Section(section).Key(key).String()
	return value
}

func GetBleveIndex() {

}

func InsertJsonToDB(objs []interface{}) {

}

func HandleError(e error) {

	if e != nil {
		fmt.Println(e)
	}
}
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func GenerateTableHeader() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "Description", "Contents", "Language", "Tags")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
