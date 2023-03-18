package Spellbook

import (
	"Spellbook/internal/Utils"

	"log"

	"github.com/blevesearch/bleve/v2"
	"github.com/google/uuid"
)

type Spell struct {
	ID          string `json:"id,omitempty"`
	Language    string `json:"language,omitempty"`
	Contents    string `json:"contents,omitempty"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
	Tags        string `json:"tags,omitempty"`
	// Remote      string `json:"remote,omitempty"` // This will be the alias to the remote DB server. Can find actual url in the config dir.
}

type Spellbook struct {
	Index bleve.Index
}

func (s *Spellbook) CreateSpell(spell Spell) (Spell, error) {
	var id string
	if spell.ID == "" {
		id = uuid.New().String()
	} else {
		id = spell.ID
	}
	newSpell := Spell{
		ID: id, Language: spell.Language, Contents: spell.Contents, Description: spell.Description, Author: spell.Author,
		Tags: spell.Tags,
	}
	// db.Table(tableName).Create(&newSpell)
	s.Index.Index(newSpell.ID, newSpell)
	return newSpell, nil
}

func (s *Spellbook) GetAllSpells() (*bleve.SearchResult, error) {

	log.Print("Getting spells")
	query := bleve.NewMatchAllQuery()
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	results, _ := s.Index.Search(search)
	return results, nil
}

func (s *Spellbook) GetSearchRange(from int, to int) (*bleve.SearchResult, error) {
	log.Printf("Getting %d spells", to)
	query := bleve.NewMatchAllQuery()
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.From = from
	search.Size = to
	results, _ := s.Index.Search(search)
	return results, nil
}

func (s *Spellbook) FindSpellsByDescription(description string, result_size int) (*bleve.SearchResult, error) {
	query := bleve.NewPhraseQuery([]string{description}, "description")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.Size = result_size
	results, err := s.Index.Search(search)
	if err != nil {
		Utils.Error("Error searching index.", err)
	}

	return results, nil
}

func (s *Spellbook) FindSpellsByTag(tag string, result_size int) (*bleve.SearchResult, error) {
	query := bleve.NewPhraseQuery([]string{tag}, "tags")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.Size = result_size
	results, err := s.Index.Search(search)
	if err != nil {
		Utils.Error("Error searching index.", err)
	}

	return results, nil
}

func (s *Spellbook) GetSpellByID(spell_id string) (*bleve.SearchResult, error) {
	query := bleve.NewDocIDQuery([]string{spell_id})
	search := bleve.NewSearchRequest(query)
	// search.Fields = []string{"tags"}
	search.Fields = []string{"*"}
	results, err := s.Index.Search(search)
	if err != nil {
		Utils.Error("Error searching index.", err)
	}

	return results, nil
}

func (s *Spellbook) UpdateSpell(spell_id string, updatedSpell Spell) (Spell, error) {

	err := s.Index.Index(spell_id, updatedSpell)

	if err != nil {
		Utils.Error("Error getting document from index.", err)
	}
	return updatedSpell, nil
}
