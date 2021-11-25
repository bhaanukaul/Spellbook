package Spell

import (
	"Spellbook/internal/Utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type SpellDBModel struct {
	ID          int    `json:"id,omitempty"`
	Language    string `json:"language"`
	Contents    string `json:"contents"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
}

type Spell struct {
	ID          int      `json:"id,omitempty"`
	Language    string   `json:"language"`
	Contents    string   `json:"contents"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func GetSpell(c *gin.Context) {
	return
}

func CreateSpell(language string, contents string, description string, tags string) (int, error) {
	db := Utils.GetDatabaseConnection()
	// var newSpell SpellDBModel
	newSpell := SpellDBModel{
		Language: language, Contents: contents, Description: description,
		Tags: tags,
	}
	db.Table("Spells").Create(&newSpell)
	return newSpell.ID, nil
}

func GetAllSpells() ([]Spell, error) {
	db := Utils.GetDatabaseConnection()
	var dbResults []SpellDBModel
	db.Table("Spells").Find(&dbResults)
	var results []Spell
	for _, spell := range dbResults {
		spellToAppend := Spell{
			ID:       spell.ID,
			Language: spell.Language, Contents: spell.Contents, Description: spell.Description,
			Tags: strings.Split(spell.Tags, " "),
		}
		results = append(results, spellToAppend)
	}
	return results, nil
}

func FindSpellsByTag(tag string) ([]Spell, error) {
	db := Utils.GetDatabaseConnection()
	var dbResults []SpellDBModel
	var results []Spell
	fmt.Printf("args in FindSpellsByTags %s\n", tag)

	// db.Table("Spells").Where("tags LIKE ?", "%?%", tags).Find(&dbResults)
	db.Table("Spells").Where("tags LIKE ?", "%"+tag+"%").Find(&dbResults)
	for _, spell := range dbResults {
		spellToAppend := Spell{
			ID:       spell.ID,
			Language: spell.Language, Contents: spell.Contents, Description: spell.Description,
			Tags: strings.Split(spell.Tags, " "),
		}
		results = append(results, spellToAppend)
	}
	return results, nil
}

func GetSpellByID(spell_id int) (Spell, error) {
	db := Utils.GetDatabaseConnection()
	var dbresult []SpellDBModel
	fmt.Printf("ID in GetSpellbyID: %d\n", spell_id)
	var result Spell
	// db.Table("Spells").Where("tags LIKE ?", "%?%", tags).Find(&dbResults)
	db.Table("Spells").Find(&dbresult, spell_id)
	if len(dbresult) == 0 {
		return result, nil
	}
	result = Spell{
		ID:       dbresult[0].ID,
		Language: dbresult[0].Language, Contents: dbresult[0].Contents, Description: dbresult[0].Description,
		Tags: strings.Split(dbresult[0].Tags, " "),
	}

	return result, nil
}

func UpdateSpell(c *gin.Context) {
	return
}

func DeleteSpell(c *gin.Context) {
	return
}
