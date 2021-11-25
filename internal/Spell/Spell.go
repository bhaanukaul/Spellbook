package Spell

import (
	"Spellbook/internal/Utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

var tableName = "spells"

type Spell struct {
	ID          int    `json:"id,omitempty"`
	Language    string `json:"language"`
	Contents    string `json:"contents"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
}

func GetSpell(c *gin.Context) {
	return
}

func CreateSpell(language string, contents string, description string, tags string) (int, error) {
	db := Utils.GetDatabaseConnection()
	// var newSpell SpellDBModel
	newSpell := Spell{
		Language: language, Contents: contents, Description: description,
		Tags: tags,
	}
	db.Table(tableName).Create(&newSpell)
	return newSpell.ID, nil
}

func GetAllSpells() ([]Spell, error) {
	db := Utils.GetDatabaseConnection()
	var results []Spell
	db.Table(tableName).Find(&results)
	return results, nil
}

func FindSpellsByTag(tag string) ([]Spell, error) {
	db := Utils.GetDatabaseConnection()
	var results []Spell
	fmt.Printf("args in FindSpellsByTags %s\n", tag)

	// db.Table("Spells").Where("tags LIKE ?", "%?%", tags).Find(&dbResults)
	db.Table(tableName).Where("tags LIKE ?", "%"+tag+"%").Find(&results)

	return results, nil
}

func GetSpellByID(spell_id int) (Spell, error) {
	db := Utils.GetDatabaseConnection()
	fmt.Printf("ID in GetSpellbyID: %d\n", spell_id)
	var result Spell
	// db.Table("Spells").Where("tags LIKE ?", "%?%", tags).Find(&dbResults)
	db.Table(tableName).Find(&result, spell_id)
	return result, nil
}

func UpdateSpell(c *gin.Context) {
	return
}

func DeleteSpell(c *gin.Context) {
	return
}
