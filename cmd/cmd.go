package cmd

import (
	"Spellbook/internal/Spellbook"
	"Spellbook/internal/Utils"
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var (
	spellbook Spellbook.Spellbook

	rootCmd = &cobra.Command{
		Use:              "spellbook",
		Short:            "A collection of personal code snippets",
		PersistentPreRun: spellbookPreRun,
	}
	sugar *zap.SugaredLogger
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer spellbook.Index.Close()
}

func init() {
	cobra.OnInitialize(initConfig)

}

func initConfig() {
	viper.SetConfigName("spellbook")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	sugar = logger.Sugar()

	if err := viper.ReadInConfig(); err != nil {
		sugar.Fatalf("fatal error config file: %w", err)
	}
}

func spellbookPreRun(cmd *cobra.Command, args []string) {
	index_file := viper.GetString("BLEVE_INDEX")
	sugar.Info("bleve index: %s", index_file)
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
