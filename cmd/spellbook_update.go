package cmd

import (
	"Spellbook/internal/Spellbook"
	"Spellbook/internal/Utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

func updateSpellApi(c echo.Context) error {
	spellToUpdate := new(Spellbook.Spell)
	spellID := c.Param("id")
	if err := c.Bind(spellToUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request Body")
	}

	sugar.Debugf("From POST CreateSpell: %s, %s, %s, %s", spellToUpdate.Language, spellToUpdate.Contents, spellToUpdate.Description, spellToUpdate.Tags)
	spell, err := spellbook.UpdateSpell(spellID, *spellToUpdate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, spell)
}

func updateSpellCli(c *cobra.Command, args []string) {
	tbl := Utils.GenerateTableHeader()
	contents, err := c.Flags().GetString("contents")
	if err != nil {
		sugar.Debugf("Could not get contents. %v", err)
	}
	description, err := c.Flags().GetString("description")
	if err != nil {
		sugar.Debugf("Could not get description. %v", err)
	}
	language, err := c.Flags().GetString("language")
	if err != nil {
		sugar.Debugf("Could not get language. %v", err)
	}
	tags, err := c.Flags().GetString("tags")
	if err != nil {
		sugar.Debugf("Could not get tags. %v", err)
	}
	author, err := c.Flags().GetString("author")
	if err != nil {
		sugar.Debugf("Could not get author. %v", err)
	}

	id, err := c.Flags().GetString("id")
	if err != nil {
		sugar.Debugf("Could not get id. %v", err)
	}
	result, err := updateSpell(id, contents, description, language, tags, author)
	if err != nil {
		sugar.Debugf("Could not get id. %v", err)
	}
	if result != (&Spellbook.Spell{}) {
		tbl.AddRow(result.ID, result.Description, result.Contents, result.Language, result.Tags)
	}
	tbl.Print()
}

func updateSpell(spell_id string, contents string, description string, language string, tags string, author string) (*Spellbook.Spell, error) {
	spellToUpdate := Spellbook.Spell{
		Description: description, Language: language, Contents: contents, Tags: tags,
	}
	result, err := spellbook.UpdateSpell(spell_id, spellToUpdate)
	if err != nil {
		return result, fmt.Errorf("error creating spell: %#v", err)
	}

	return result, nil
}
