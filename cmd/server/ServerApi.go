package main

import (
	"Spellbook/internal/Spell"
	"Spellbook/internal/Utils"
	"fmt"
	"log"
	"net/http"

	"github.com/blevesearch/bleve/v2"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

/*
API Functions
*/

type RemoteSection struct {
	Alias       string `json:"alias"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func GetAllSpellbooks(c *gin.Context) {
	cd := GetServerConfig()

	cfg, err := ini.Load(cd)
	if err != nil {
		Utils.Error("Cannot load config.", err)
	}

	names := cfg.ChildSections("remotes")
	var returnSb []RemoteSection
	for _, n := range names {
		toAdd := RemoteSection{
			Alias: n.Key("alias").String(), Description: n.Key("description").String(), Url: n.Key("url").String(),
		}
		returnSb = append(returnSb, toAdd)
	}
	c.JSON(http.StatusOK, returnSb)
}

func AddRemoteSpellbookApi(c *gin.Context) {
	var remoteSection RemoteSection
	if err := c.BindJSON(&remoteSection); err != nil {
		Utils.Error("", err)
		return
	}

	AddRemoteServer(remoteSection.Alias, remoteSection.Url, remoteSection.Description)
	c.JSON(http.StatusCreated, remoteSection)
}

func UpdateSpell(c *gin.Context) {
	var spellToUpdate Spell.Spell
	if err := c.BindJSON(&spellToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Request Body")
		return
	}
	spell_id_param := c.Param("id")
	spellToUpdate.ID = spell_id_param
	index := GetBleveIndex()
	_, err := Spell.CreateSpell(spellToUpdate, index)
	if err != nil {

		c.JSON(http.StatusInternalServerError, err)
	}
	// Spell.UpdateSpell(spell_id_param, spellToUpdate, index)
	index.Close()
}

func CreateSpell(c *gin.Context) {

	var newSpell Spell.Spell
	if err := c.BindJSON(&newSpell); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Request Body")
		return
	}

	log.Printf("From POST CreateSpell: %s, %s, %s, %s", newSpell.Language, newSpell.Contents, newSpell.Description, newSpell.Tags)
	index := GetBleveIndex()
	spell, err := Spell.CreateSpell(newSpell, index)
	index.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, spell)

}

func GetSpell(c *gin.Context) {
	spell_id_param := c.Param("id")
	index := GetBleveIndex()
	searchResult, err := Spell.GetSpellByID(spell_id_param, index)
	index.Close()

	if err != nil {
		panic(err)
	}

	if searchResult.Hits.Len() == 0 {
		c.JSON(http.StatusNotFound, "Spell Not Found")
	} else {
		var spell Spell.Spell
		result := searchResult.Hits[0].Fields
		spell.ID = result["id"].(string)
		spell.Description = result["description"].(string)
		spell.Language = result["language"].(string)
		spell.Contents = result["contents"].(string)
		spell.Tags = result["tags"].(string)
		spell.Author = result["author"].(string)

		// c.BindJSON(&result)
		c.JSON(http.StatusOK, result)
	}
}

func GetAllSpells(c *gin.Context) {
	qf := c.Query("field")
	qfv := c.Query("value")
	log.Printf("Query params: %s, %s", qf, qfv)
	index := GetBleveIndex()

	if qf != "" && qfv != "" {
		log.Printf("Using query parameters %s, %s", qf, qfv)
		var searchResult *bleve.SearchResult
		var err error
		var spells []Spell.Spell
		if qf == "tags" {
			searchResult, err = Spell.FindSpellsByTag(qfv, index)
		} else if qf == "description" {
			searchResult, err = Spell.FindSpellsByDescription(qfv, index)
		}

		if err != nil {
			Utils.Error(fmt.Sprintf("Error searching %s for %s", qf, qfv), err)
			c.JSON(http.StatusInternalServerError, err)
		}

		for _, doc := range searchResult.Hits {
			spell := Spell.Spell{
				ID:          doc.Fields["id"].(string),
				Description: doc.Fields["description"].(string),
				Language:    doc.Fields["language"].(string),
				Contents:    doc.Fields["contents"].(string),
				Tags:        doc.Fields["tags"].(string),
				Author:      doc.Fields["author"].(string),
			}
			spells = append(spells, spell)
		}
		index.Close()

		c.JSON(http.StatusOK, spells)
		return
	}

	log.Print("Getting all spells")

	var spells []Spell.Spell
	searchResult, err := Spell.GetAllSpells(index)
	index.Close()
	for _, doc := range searchResult.Hits {
		spell := Spell.Spell{
			ID:          doc.Fields["id"].(string),
			Description: doc.Fields["description"].(string),
			Language:    doc.Fields["language"].(string),
			Contents:    doc.Fields["contents"].(string),
			Tags:        doc.Fields["tags"].(string),
			Author:      doc.Fields["author"].(string),
		}
		spells = append(spells, spell)
	}
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, spells)
}

func Ping(c *gin.Context) {
	ping := Utils.SpellbookPing{
		Version: "1.0",
	}

	c.JSON(http.StatusOK, ping)
}
