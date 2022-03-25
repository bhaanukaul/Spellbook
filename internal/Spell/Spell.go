package Spell

import (
	"Spellbook/internal/Utils"
	"log"

	"github.com/blevesearch/bleve/v2"
	"github.com/google/uuid"
)

var tableName = "spells"

type Spell struct {
	ID          string `json:"id,omitempty"`
	Language    string `json:"language,omitempty"`
	Contents    string `json:"contents,omitempty"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
	Tags        string `json:"tags,omitempty"`
	Remote      string `json:"remote,omitempty"` // This will be the alias to the remote DB server. Can find actual url in the config dir.
}

func CreateSpell(spell Spell, index bleve.Index) (Spell, error) {
	// db := Utils.GetDatabaseConnection()
	// var newSpell SpellDBModel
	var id string
	if spell.ID == "" {
		uuid := uuid.New()
		id = uuid.String()
	} else {
		id = spell.ID
	}
	newSpell := Spell{
		ID: id, Language: spell.Language, Contents: spell.Contents, Description: spell.Description,
		Tags: spell.Tags,
	}
	// db.Table(tableName).Create(&newSpell)
	index.Index(newSpell.ID, newSpell)
	return newSpell, nil
}

func GetAllSpells(index bleve.Index) (*bleve.SearchResult, error) {

	log.Print("Getting spells")
	query := bleve.NewMatchAllQuery()
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	results, _ := index.Search(search)

	return results, nil
}

func FindSpellsByTag(tag string, index bleve.Index) (*bleve.SearchResult, error) {
	query := bleve.NewPhraseQuery([]string{tag}, "tags")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	results, err := index.Search(search)
	if err != nil {
		Utils.Error("Error searching index.", err)
	}

	return results, nil
}

func GetSpellByID(spell_id string, index bleve.Index) (*bleve.SearchResult, error) {
	query := bleve.NewDocIDQuery([]string{spell_id})
	search := bleve.NewSearchRequest(query)
	// search.Fields = []string{"tags"}
	search.Fields = []string{"*"}
	results, err := index.Search(search)
	if err != nil {
		Utils.Error("Error searching index.", err)
	}

	return results, nil
}

func UpdateSpell(spell_id string, updatedSpell Spell, index bleve.Index) (Spell, error) {

	err := index.Index(spell_id, updatedSpell)

	if err != nil {
		Utils.Error("Error getting document from index.", err)
	}
	return updatedSpell, nil
}
