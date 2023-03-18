package cmd

import (
	"Spellbook/internal/Spellbook"
	"Spellbook/internal/Utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

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

func getSpellApi(c *gin.Context) {

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

func getAllSpellsCli(c *cobra.Command, args []string) {

}

func getAllSpellsApi(c *gin.Context) {

}

func getAllSpells() {

}
