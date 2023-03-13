package cmd

import (
	"Spellbook/internal/Spellbook"
	"Spellbook/internal/Utils"
	"fmt"
	"log"

	"github.com/blevesearch/bleve/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use: "server",
}

var serverStartCmd = &cobra.Command{
	Use: "start",
	Run: startServer,
}

var initCmd = &cobra.Command{
	Use: "init",
	Run: initBleveIndex,
}

var getSpellCmd = &cobra.Command{
	Use:  "find",
	Args: cobra.ExactArgs(1),
	Run:  getSpellCli,
}

var createSpellCmd = &cobra.Command{
	Use: "create",
	Run: createSpellCli,
}

var (
	description string
	contents    string
	tags        string
	language    string
	author      string
)

func init() {
	createSpellCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the spell")
	createSpellCmd.Flags().StringVarP(&contents, "contents", "c", "", "Contents of the spell")
	createSpellCmd.Flags().StringVarP(&tags, "tags", "t", "", "Tags for the spell")
	createSpellCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the spell")
	createSpellCmd.MarkFlagRequired("description")
	createSpellCmd.MarkFlagRequired("contents")
	createSpellCmd.MarkFlagRequired("language")

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(getSpellCmd)
	rootCmd.AddCommand(createSpellCmd)

	serverCmd.AddCommand(serverStartCmd)

}

func initBleveIndex(c *cobra.Command, args []string) {
	mapping := bleve.NewIndexMapping()
	log.Printf("fasdfsffsafvs")
	index_name := viper.GetString("BLEVE_INDEX")
	_, err := bleve.New(index_name, mapping)
	if err != nil {
		log.Fatalf("Error creating bleve index: %#v", err)
	}
}

func startServer(c *cobra.Command, args []string) {

	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/spell/:id", getSpellApi)
		api.GET("/spells", getAllSpells)
		api.POST("/spell", createSpellApi)
		api.PATCH("/spell/:id", updateSpell)
		// api.POST("/spellbook", AddRemoteSpellbookApi)
		// api.GET("/spellbooks", GetAllSpellbooks)
	}
	router.GET("/ping", ping)
	// configFile := GetServerConfig()
	// port := Utils.GetKVFromConfig(configFile, "http_port", "server")
	router.Run(fmt.Sprintf("%s:%s", viper.GetString("SERVER_ADDRESS"), viper.GetString("SERVER_PORT")))
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

func getAllSpells(c *gin.Context) {

}

func createSpellApi(c *gin.Context) {

}

func createSpellCli(c *cobra.Command, args []string) {
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

func updateSpell(c *gin.Context) {

}

func ping(c *gin.Context) {

}
