package Spellbook

import (
	"fmt"
	"strconv"

	"github.com/blevesearch/bleve/v2"
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

func (s *Spellbook) CreateSpell(spell Spell) (*Spell, error) {
	var id string
	doc_count, err := s.Index.DocCount()
	if err != nil {
		return nil, fmt.Errorf("error searching index: %#v", err)
	}

	doc_count += 1
	id = strconv.Itoa(int(doc_count))
	newSpell := Spell{
		ID:       id,
		Language: spell.Language, Contents: spell.Contents, Description: spell.Description, Author: spell.Author,
		Tags: spell.Tags,
	}
	// db.Table(tableName).Create(&newSpell)
	s.Index.Index(newSpell.ID, newSpell)
	results, err := s.GetSpellByID(id)
	if err != nil {
		return nil, fmt.Errorf("error searching index: %#v", err)
	}

	return results, nil
}

func (s *Spellbook) GetAllSpells() ([]Spell, error) {
	query := bleve.NewMatchAllQuery()
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	results, err := s.Index.Search(search)
	if err != nil {
		return nil, fmt.Errorf("error searching index: %#v", err)
	}
	return bleveResultsToSpell(results), nil
}

func (s *Spellbook) GetSearchRange(from int, to int) ([]Spell, error) {
	query := bleve.NewMatchAllQuery()
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.From = from
	search.Size = to
	results, err := s.Index.Search(search)
	if err != nil {
		return nil, fmt.Errorf("error searching index: %#v", err)
	}
	return bleveResultsToSpell(results), nil
}

func (s *Spellbook) FindSpellsByDescription(description string, result_size int) ([]Spell, error) {
	query := bleve.NewPhraseQuery([]string{description}, "description")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.Size = result_size
	results, err := s.Index.Search(search)
	if err != nil {
		return nil, fmt.Errorf("error searching index: %#v", err)
	}

	return bleveResultsToSpell(results), nil
}

func (s *Spellbook) FindSpellsByTag(tag string, result_size int) ([]Spell, error) {
	query := bleve.NewPhraseQuery([]string{tag}, "tags")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"*"}
	search.Size = result_size
	results, err := s.Index.Search(search)
	if err != nil {
		return nil, fmt.Errorf("error searching index: %#v", err)
	}

	return bleveResultsToSpell(results), nil
}

func (s *Spellbook) GetSpellByID(spell_id string) (*Spell, error) {
	query := bleve.NewDocIDQuery([]string{spell_id})
	search := bleve.NewSearchRequest(query)
	// search.Fields = []string{"tags"}
	search.Fields = []string{"*"}
	results, err := s.Index.Search(search)
	if err != nil {
		return nil, fmt.Errorf("error searching index: %#v", err)
	}
	spell := bleveResultsToSpell(results)
	if len(spell) == 0 {
		return &Spell{}, nil
	}
	return &spell[0], nil
}

func (s *Spellbook) UpdateSpell(spell_id string, updatedSpell Spell) (*Spell, error) {

	err := s.Index.Index(spell_id, updatedSpell)

	if err != nil {
		return nil, fmt.Errorf("error searching index: %#v", err)
	}
	result, err := s.GetSpellByID(spell_id)

	if err != nil {
		return nil, fmt.Errorf("error searching index: %#v", err)
	}
	return result, nil
}

func bleveResultsToSpell(results *bleve.SearchResult) []Spell {
	var spells []Spell
	for _, doc := range results.Hits {
		spell := Spell{
			ID:          doc.Fields["id"].(string),
			Description: doc.Fields["description"].(string),
			Language:    doc.Fields["language"].(string),
			Contents:    doc.Fields["contents"].(string),
			Tags:        doc.Fields["tags"].(string),
			Author:      doc.Fields["author"].(string),
		}
		spells = append(spells, spell)
	}
	return spells
}
