package cmd

import (
	"Spellbook/internal/Constants"
	"fmt"
	"log"
	"net/http"

	"github.com/blevesearch/bleve/v2"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	description    string
	contents       string
	tags           string
	language       string
	limit          int
	server_port    string
	server_address string
	bleve_index    string

	updateSpellCmd = &cobra.Command{
		Use: "update",
		Run: updateSpellCli,
	}

	createSpellCmd = &cobra.Command{
		Use: "create",
		Run: createSpellCli,
	}

	serverCmd = &cobra.Command{
		Use: "server",
	}

	serverStartCmd = &cobra.Command{
		Use: "start",
		Run: startServer,
	}

	initCmd = &cobra.Command{
		Use: "init",
		Run: initBleveIndex,
	}

	getSpellCmd = &cobra.Command{
		Use: "find",
		Run: getSpellCli,
	}

	getAllSpellsCmd = &cobra.Command{
		Use:  "all",
		Args: cobra.ExactArgs(1),
		Run:  getAllSpellsCli,
	}

	getSpellsByDescriptionCmd = &cobra.Command{
		Use:  "by-description",
		Args: cobra.ExactArgs(1),
		Run:  getSpellsByDescriptionCli,
	}

	getSpellByTagCmd = &cobra.Command{
		Use:  "by-tag",
		Args: cobra.ExactArgs(1),
		Run:  getAllSpellsCli,
	}
)

func init() {
	createSpellCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the spell")
	createSpellCmd.Flags().StringVarP(&contents, "contents", "c", "", "Contents of the spell")
	createSpellCmd.Flags().StringVarP(&tags, "tags", "t", "", "Tags for the spell")
	createSpellCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the spell")
	createSpellCmd.MarkFlagRequired("description")
	createSpellCmd.MarkFlagRequired("contents")
	createSpellCmd.MarkFlagRequired("tags")
	getSpellCmd.Flags().IntVarP(&limit, "limit", "l", Constants.SearchResultLimit, "Limit of search results")
	rootCmd.Flags().StringVarP(&bleve_index, "index", "i", "", "Bleve index file name")

	getSpellCmd.AddCommand(getAllSpellsCmd)
	getSpellCmd.AddCommand(getSpellsByDescriptionCmd)
	getSpellCmd.AddCommand(getSpellByTagCmd)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(getSpellCmd)
	rootCmd.AddCommand(createSpellCmd)
	rootCmd.AddCommand(updateSpellCmd)

	serverStartCmd.Flags().StringVarP(&server_address, "address", "a", "", "Server address")
	serverStartCmd.Flags().StringVarP(&server_port, "port", "p", "", "Server port")

	viper.BindPFlag("SERVER_PORT", serverStartCmd.Flags().Lookup("port"))
	viper.BindPFlag("SERVER_ADDRESS", serverStartCmd.Flags().Lookup("address"))
	viper.BindPFlag("BLEVE_INDEX", rootCmd.Flags().Lookup("index"))

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

	router := echo.New()

	api := router.Group("/api")
	{
		api.GET("/spell/:id", getSpellApi)
		api.GET("/spells", getSpellsApi)
		api.POST("/spell", createSpellApi)
		api.PATCH("/spell/:id", updateSpellApi)
		// api.POST("/spellbook", AddRemoteSpellbookApi)
		// api.GET("/spellbooks", GetAllSpellbooks)
	}
	router.GET("/ping", ping)
	// configFile := GetServerConfig()
	// port := Utils.GetKVFromConfig(configFile, "http_port", "server")

	router.Logger.Fatal(router.Start(fmt.Sprintf("%s:%s", viper.GetString("SERVER_ADDRESS"), viper.GetString("SERVER_PORT"))))
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong!")
}
