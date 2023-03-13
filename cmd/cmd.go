package cmd

import (
	"Spellbook/internal/Spellbook"
	"Spellbook/internal/Utils"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	spellbook Spellbook.Spellbook

	rootCmd = &cobra.Command{
		Use:   "spellbook",
		Short: "A collection of personal code snippets",
	}
)

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("spellbook")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// viper.AutomaticEnv()
	// HandleError(viper.BindEnv("SERVER_PORT"))
	// HandleError(viper.BindEnv("SERVER_ADDRESS"))
	// HandleError(viper.BindEnv("BLEVE_INDEX"))
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
		// log.Println("Using configuration file: ", viper.ConfigFileUsed())
	}

	index_file := viper.GetString("BLEVE_INDEX")
	log.Printf("bleve index: %s", index_file)
	if !Utils.FileExists(index_file) {
		log.Printf("No bleve index, creating one: %s", index_file)
		// return
	} else {
		index, err := bleve.Open(index_file)
		if err != nil {
			log.Fatalf("Error opening bleve index: %#v", err)
			os.Exit(1)
		}
		spellbook.Index = index
	}

}
