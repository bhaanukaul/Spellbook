package cmd

import (
	"Spellbook/internal/Spellbook"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var (
	spellbook Spellbook.Spellbook

	rootCmd = &cobra.Command{
		Use:   "spellbook",
		Short: "A collection of personal code snippets",
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			defer spellbook.Index.Close()
		},
	}
	sugar *zap.SugaredLogger
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

	if err := viper.ReadInConfig(); err != nil {
		sugar.Fatalf("fatal error config file: %w", err)
	}

	log_level, err := zap.ParseAtomicLevel(viper.GetString("LOG_LEVEL"))
	if err != nil {
		sugar.Fatalf("Error parsing log level from viper config: %w", err)
	}

	logger := zap.Must(zap.Config{
		Level:            zap.NewAtomicLevelAt(log_level.Level()),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build())
	defer logger.Sync() // flushes buffer, if any

	sugar = logger.Sugar()
}
