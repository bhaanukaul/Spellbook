package cmd

import (
	"Spellbook/internal/Spellbook"
	"Spellbook/internal/Utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func createSpellApi(c *gin.Context) {
	var newSpell Spellbook.Spell
	if err := c.BindJSON(&newSpell); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Request Body")
		return
	}

	log.Printf("From POST CreateSpell: %s, %s, %s, %s", newSpell.Language, newSpell.Contents, newSpell.Description, newSpell.Tags)
	spell, err := spellbook.CreateSpell(newSpell)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, spell)
}

func createSpellCli(c *cobra.Command, args []string) {
	contents, err := c.Flags().GetString("contents")
	if err != nil {
		log.Printf("Could not get contents. %v", err)
	}
	description, err := c.Flags().GetString("description")
	if err != nil {
		log.Printf("Could not get description. %v", err)
	}
	language, err := c.Flags().GetString("language")
	if err != nil {
		log.Printf("Could not get language. %v", err)
	}
	tags, err := c.Flags().GetString("tags")
	if err != nil {
		log.Printf("Could not get tags. %v", err)
	}
	author, err := c.Flags().GetString("author")
	if err != nil {
		log.Printf("Could not get author. %v", err)
	}
	spell, err := createSpell(contents, description, language, tags, author)
	tbl := Utils.GenerateTableHeader()
	if err == nil {
		if spell != nil {
			tbl.AddRow(spell.ID, spell.Description, spell.Contents, spell.Language, spell.Tags)
		}
		tbl.Print()
	}

}

func createSpell(contents string, description string, language string, tags string, author string) (*Spellbook.Spell, error) {
	spellToCreate := Spellbook.Spell{
		Description: description, Language: language, Contents: contents, Tags: tags,
	}
	spell, err := spellbook.CreateSpell(spellToCreate)
	if err != nil {
		return nil, fmt.Errorf("error creating spell: %#v", err)
	}
	return &spell, nil
}
