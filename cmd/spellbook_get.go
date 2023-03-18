package cmd

import (
	"Spellbook/internal/Constants"
	"Spellbook/internal/Spellbook"
	"Spellbook/internal/Utils"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

/*

Get individual spell

*/

func getSpell(id string) *Spellbook.Spell {
	searchResult, err := spellbook.GetSpellByID(id)

	if err != nil {
		panic(err)
	}

	if searchResult.Hits.Len() == 0 {
		return nil
	} else {
		var spell Spellbook.Spell
		result := searchResult.Hits[0].Fields
		spell.ID = result["id"].(string)
		spell.Description = result["description"].(string)
		spell.Language = result["language"].(string)
		spell.Contents = result["contents"].(string)
		spell.Tags = result["tags"].(string)
		spell.Author = result["author"].(string)

		// c.BindJSON(&result)
		return &spell
	}
}

func getSpellApi(c echo.Context) error {
	spell_id := c.Param("id")
	searchResult, err := spellbook.GetSpellByID(spell_id)

	if err != nil {
		panic(err)
	}

	if searchResult.Hits.Len() == 0 {
		return c.JSON(http.StatusNotFound, "Spell Not Found")
	} else {
		var spell Spellbook.Spell
		result := searchResult.Hits[0].Fields
		spell.ID = result["id"].(string)
		spell.Description = result["description"].(string)
		spell.Language = result["language"].(string)
		spell.Contents = result["contents"].(string)
		spell.Tags = result["tags"].(string)
		spell.Author = result["author"].(string)

		// c.BindJSON(&result)
		return c.JSON(http.StatusOK, result)
	}
}

func getSpellCli(c *cobra.Command, args []string) {
	log.Printf("Getting id: %s", args[0])
	spell_id := args[0]
	spell := getSpell(spell_id)
	tbl := Utils.GenerateTableHeader()
	if spell != nil {
		tbl.AddRow(spell.ID, spell.Description, spell.Contents, spell.Language, spell.Tags)
	}
	tbl.Print()
}

/*
Get Alls spells, no search parameters
*/

func getAllSpellsCli(c *cobra.Command, args []string) {
	log.Print("Getting all spells")
	tbl := Utils.GenerateTableHeader()
	tbl.Print()
	searchResult, err := spellbook.GetAllSpells()
	for _, doc := range searchResult.Hits {
		spell := Spellbook.Spell{
			ID:          doc.Fields["id"].(string),
			Description: doc.Fields["description"].(string),
			Language:    doc.Fields["language"].(string),
			Contents:    doc.Fields["contents"].(string),
			Tags:        doc.Fields["tags"].(string),
			Author:      doc.Fields["author"].(string),
		}
		tbl.AddRow(spell.ID, spell.Description, spell.Contents, spell.Language, spell.Tags)
	}
	if err != nil {
		panic(err)
	}
	tbl.Print()
}

func getAllSpells(limit int) ([]Spellbook.Spell, error) {
	var spells []Spellbook.Spell
	searchResult, err := spellbook.GetAllSpells()
	for _, doc := range searchResult.Hits {
		spell := Spellbook.Spell{
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
		return spells, err
	}
	return spells, nil
}

/*
Get spells by tag
*/

func getSpellsByTag(tag string, limit int) ([]Spellbook.Spell, error) {
	var spells []Spellbook.Spell
	searchResult, err := spellbook.FindSpellsByTag(tag, limit)
	for _, doc := range searchResult.Hits {
		spell := Spellbook.Spell{
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
		return spells, err
	}
	return spells, nil
}

/*
Get spells by description
*/

func getSpellsByDescription(description string, limit int) ([]Spellbook.Spell, error) {
	var spells []Spellbook.Spell
	searchResult, err := spellbook.FindSpellsByDescription(description, limit)
	for _, doc := range searchResult.Hits {
		spell := Spellbook.Spell{
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
		return spells, err
	}
	return spells, nil
}

func getSpellsByDescriptionCli(c *cobra.Command, args []string) {
	tbl := Utils.GenerateTableHeader()
	limit, _ := c.Flags().GetInt("limit")
	spells, err := getSpellsByDescription(args[0], limit)
	if err != nil {
		log.Fatalf("Failed to get spells: %#v", err)
	}
	for _, spell := range spells {
		tbl.AddRow(spell.ID, spell.Description, spell.Contents, spell.Language, spell.Tags)
	}
	tbl.Print()
}

/*
Get all spells, or search by tag or description.
*/

func getSpellsApi(c echo.Context) error {
	qf := c.QueryParam("field")
	qfv := c.QueryParam("value")
	var limit int
	if c.QueryParam("limit") != "" {
		limit, _ = strconv.Atoi(c.QueryParam("limit"))
	} else {
		limit = Constants.SearchResultLimit
	}
	log.Printf("Query params: %s, %s", qf, qfv)
	var spells []Spellbook.Spell
	var err error
	if qf != "" && qfv != "" {
		log.Printf("Using query parameters %s, %s", qf, qfv)
		if qf == "tags" {
			spells, err = getSpellsByTag(qfv, limit)
		} else if qf == "description" {
			spells, err = getSpellsByDescription(qfv, limit)
		}
		if err != nil {
			Utils.Error(fmt.Sprintf("Error searching %s for %s", qf, qfv), err)
			c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, spells)
	}
	spells, err = getAllSpells(limit)
	if err != nil {
		Utils.Error(fmt.Sprintf("Error searching %s for %s", qf, qfv), err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, spells)

}
