package cmd

import (
	"Spellbook/internal/Constants"
	"Spellbook/internal/Utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	description    string
	contents       string
	tags           string
	language       string
	id             string
	limit          int
	server_port    string
	server_address string
	bleve_index    string

	updateSpellCmd = &cobra.Command{
		Use:              "update",
		Run:              updateSpellCli,
		PersistentPreRun: spellbookPreRun,
	}

	createSpellCmd = &cobra.Command{
		Use:              "create",
		Run:              createSpellCli,
		PersistentPreRun: spellbookPreRun,
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
		Use:              "find",
		PersistentPreRun: spellbookPreRun,
	}

	getAllSpellsCmd = &cobra.Command{
		Use: "all",
		Run: getAllSpellsCli,
	}

	getSpellsByDescriptionCmd = &cobra.Command{
		Use:  "by-description",
		Args: cobra.ExactArgs(1),
		Run:  getSpellsByDescriptionCli,
	}

	getSpellByTagCmd = &cobra.Command{
		Use:  "by-tag",
		Args: cobra.ExactArgs(1),
		Run:  getSpellsByTagCli,
	}

	getSpellByIdCmd = &cobra.Command{
		Use:  "id",
		Args: cobra.ExactArgs(1),
		Run:  getSpellCli,
	}

	completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: fmt.Sprintf(`To load completions:
	
	Bash:
	
	  $ source <(%[1]s completion bash)
	
	  # To load completions for each session, execute once:
	  # Linux:
	  $ %[1]s completion bash > /etc/bash_completion.d/%[1]s
	  # macOS:
	  $ %[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s
	
	Zsh:
	
	  # If shell completion is not already enabled in your environment,
	  # you will need to enable it.  You can execute the following once:
	
	  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
	
	  # To load completions for each session, execute once:
	  $ %[1]s completion zsh > "${fpath[1]}/_%[1]s"
	
	  # You will need to start a new shell for this setup to take effect.
	
	fish:
	
	  $ %[1]s completion fish | source
	
	  # To load completions for each session, execute once:
	  $ %[1]s completion fish > ~/.config/fish/completions/%[1]s.fish
	
	PowerShell:
	
	  PS> %[1]s completion powershell | Out-String | Invoke-Expression
	
	  # To load completions for every new session, run:
	  PS> %[1]s completion powershell > %[1]s.ps1
	  # and source this file from your PowerShell profile.
	`, rootCmd.Name()),
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				rootCmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				rootCmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				rootCmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				rootCmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
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

	getSpellCmd.PersistentFlags().IntVarP(&limit, "limit", "l", Constants.SearchResultLimit, "Limit of search results")

	updateSpellCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the spell")
	updateSpellCmd.Flags().StringVarP(&contents, "contents", "c", "", "Contents of the spell")
	updateSpellCmd.Flags().StringVarP(&tags, "tags", "t", "", "Tags for the spell")
	updateSpellCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the spell")
	updateSpellCmd.Flags().StringVar(&id, "id", "", "Id of the spell")
	updateSpellCmd.MarkFlagRequired("id")

	rootCmd.Flags().StringVarP(&bleve_index, "index", "i", "", "Bleve index file name")

	getSpellCmd.AddCommand(getAllSpellsCmd)
	getSpellCmd.AddCommand(getSpellsByDescriptionCmd)
	getSpellCmd.AddCommand(getSpellByTagCmd)
	getSpellCmd.AddCommand(getSpellByIdCmd)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(getSpellCmd)
	rootCmd.AddCommand(createSpellCmd)
	rootCmd.AddCommand(updateSpellCmd)
	rootCmd.AddCommand(completionCmd)

	serverStartCmd.Flags().StringVarP(&server_address, "address", "a", "", "Server address")
	serverStartCmd.Flags().StringVarP(&server_port, "port", "p", "", "Server port")

	viper.BindPFlag("SERVER_PORT", serverStartCmd.Flags().Lookup("port"))
	viper.BindPFlag("SERVER_ADDRESS", serverStartCmd.Flags().Lookup("address"))
	viper.BindPFlag("BLEVE_INDEX", rootCmd.Flags().Lookup("index"))

	serverCmd.AddCommand(serverStartCmd)

}

func initBleveIndex(c *cobra.Command, args []string) {
	mapping := bleve.NewIndexMapping()
	index_name := viper.GetString("BLEVE_INDEX")
	_, err := bleve.New(index_name, mapping)
	if err != nil {
		log.Fatalf("Error creating bleve index: %#v", err)
	}
}

func startServer(c *cobra.Command, args []string) {

	router := echo.New()

	keyAuth := middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == viper.GetString("SERVER_SECRET"), nil
	})

	api := router.Group("/api")
	{
		api.GET("/spell/:id", getSpellApi)
		api.GET("/spells", getSpellsApi)
		api.POST("/spell", createSpellApi, keyAuth)
		api.PATCH("/spell/:id", updateSpellApi, keyAuth)
	}
	router.GET("/ping", ping)

	router.Logger.Fatal(router.Start(fmt.Sprintf("%s:%s", viper.GetString("SERVER_ADDRESS"), viper.GetString("SERVER_PORT"))))
	defer spellbook.Index.Close()

}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong!")
}

func spellbookPreRun(cmd *cobra.Command, args []string) {
	index_file := viper.GetString("BLEVE_INDEX")
	sugar.Infof("bleve index: %s", index_file)
	if !Utils.FileExists(index_file) {
		fmt.Println("No bleve index, create one using: spellbook init")
		os.Exit(1)
	} else {
		index, err := bleve.Open(index_file)
		if err != nil {
			sugar.Fatalf("Error opening bleve index: %#v", err)
		}
		sugar.Info("Bleve index exists.")
		spellbook.Index = index
	}
}
