package cmd

import (
	"Spellbook/internal/Constants"
	"Spellbook/internal/Spellbook"
	"Spellbook/internal/Utils"
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
	result, err := spellbook.GetSpellByID(id)

	if err != nil {
		panic(err)
	}
	return result
}

func getSpellApi(c echo.Context) error {
	spell_id := c.Param("id")
	result := getSpell(spell_id)

	if result == (&Spellbook.Spell{}) {
		return c.JSON(http.StatusNotFound, "Spell Not Found")
	}
	return c.JSON(http.StatusOK, result)
}

func getSpellCli(c *cobra.Command, args []string) {
	sugar.Debugf("Getting id: %s", args[0])
	spell_id := args[0]
	result := getSpell(spell_id)
	tbl := Utils.GenerateTableHeader()
	if result != (&Spellbook.Spell{}) {
		tbl.AddRow(result.ID, result.Description, result.Contents, result.Language, result.Tags)
	}
	tbl.Print()

}

/*
Get Alls spells, no search parameters
*/

func getAllSpellsCli(c *cobra.Command, args []string) {
	sugar.Debug("Getting all spells")
	tbl := Utils.GenerateTableHeader()
	results, err := spellbook.GetAllSpells()
	for _, result := range results {
		tbl.AddRow(result.ID, result.Description, result.Contents, result.Language, result.Tags)
	}
	if err != nil {
		panic(err)
	}
	tbl.Print()
}

func getAllSpells(limit int) ([]Spellbook.Spell, error) {
	results, err := spellbook.GetAllSpells()
	if err != nil {
		return results, err
	}

	return results, nil
}

/*
Get spells by tag
*/

func getSpellsByTag(tag string, limit int) ([]Spellbook.Spell, error) {
	sugar.Debugf("search by tag: %s", tag)

	results, err := spellbook.FindSpellsByTag(tag, limit)
	if err != nil {
		return results, err
	}
	return results, nil
}

func getSpellsByTagCli(c *cobra.Command, args []string) {
	tbl := Utils.GenerateTableHeader()
	limit, _ := c.Flags().GetInt("limit")
	sugar.Debugf("cli search by tag: %s with limit: %d", args[0], limit)
	results, err := getSpellsByTag(args[0], limit)
	if err != nil {
		log.Fatalf("Failed to get spells: %#v", err)
	}
	for _, result := range results {
		tbl.AddRow(result.ID, result.Description, result.Contents, result.Language, result.Tags)
	}
	tbl.Print()
}

/*
Get spells by description
*/

func getSpellsByDescription(description string, limit int) ([]Spellbook.Spell, error) {
	results, err := spellbook.FindSpellsByDescription(description, limit)
	if err != nil {
		return results, err
	}
	return results, nil
}

func getSpellsByDescriptionCli(c *cobra.Command, args []string) {
	tbl := Utils.GenerateTableHeader()
	limit, _ := c.Flags().GetInt("limit")
	results, err := getSpellsByDescription(args[0], limit)
	if err != nil {
		log.Fatalf("Failed to get spells: %#v", err)
	}
	for _, result := range results {
		tbl.AddRow(result.ID, result.Description, result.Contents, result.Language, result.Tags)
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
	sugar.Debugf("Query params: %s, %s", qf, qfv)
	var results []Spellbook.Spell
	var err error
	if qf != "" && qfv != "" {
		sugar.Debugf("Using query parameters %s, %s", qf, qfv)
		if qf == "tags" {
			results, err = getSpellsByTag(qfv, limit)
		} else if qf == "description" {
			results, err = getSpellsByDescription(qfv, limit)
		}
		if err != nil {
			sugar.Errorf("Error searching %s for %s", qf, qfv, err)
			c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, results)
	}
	results, err = getAllSpells(limit)
	if err != nil {
		sugar.Errorf("Error searching %s for %s", qf, qfv, err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, results)

}
