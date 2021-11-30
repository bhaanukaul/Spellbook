package Spell

import (
	"Spellbook/internal/Utils"

	"github.com/gin-gonic/gin"
)

var tableName = "spells"

type Spell struct {
	ID          int    `json:"id,omitempty"`
	Language    string `json:"language,omitempty"`
	Contents    string `json:"contents,omitempty"`
	Description string `json:"description,omitempty"`
	Tags        string `json:"tags,omitempty"`
}

func GetSpell(c *gin.Context) {
	return
}

func CreateSpell(language string, contents string, description string, tags string) (Spell, error) {
	db := Utils.GetDatabaseConnection()
	// var newSpell SpellDBModel
	newSpell := Spell{
		Language: language, Contents: contents, Description: description,
		Tags: tags,
	}
	db.Table(tableName).Create(&newSpell)
	return newSpell, nil
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
	// fmt.Printf("args in FindSpellsByTags %s\n", tag)

	// db.Table("Spells").Where("tags LIKE ?", "%?%", tags).Find(&dbResults)
	db.Table(tableName).Where("tags LIKE ?", "%"+tag+"%").Find(&results)

	return results, nil
}

func GetSpellByID(spell_id int) (Spell, error) {
	db := Utils.GetDatabaseConnection()
	// fmt.Printf("ID in GetSpellbyID: %d\n", spell_id)
	var result Spell
	// db.Table("Spells").Where("tags LIKE ?", "%?%", tags).Find(&dbResults)
	db.Table(tableName).Find(&result, spell_id)
	// fmt.Printf("Got spell by id: %#v", result)
	return result, nil
}

func UpdateSpell(spell_id int, spell Spell) (Spell, error) {
	db := Utils.GetDatabaseConnection()
	var spellToUpdate Spell
	db.Table(tableName).Model(&spellToUpdate).Where("id = ?", spell_id).Updates(spell)
	return spellToUpdate, nil
}
