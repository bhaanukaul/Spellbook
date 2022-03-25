package main

import (
	"Spellbook/internal/Constants"
	"Spellbook/internal/Utils"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func GetBleveIndex() bleve.Index {
	config := GetCliConfig()
	index_name := Utils.GetKVFromConfig(config, Constants.CliConfigBleveIndexKey, "")
	index, _ := bleve.Open(index_name)
	return index
}

func GetCliConfig() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	spellbookConfigDir := filepath.Join(userHome, ".config", "Spellbook")
	spellbookConfigFile := fmt.Sprintf("%s/%s", spellbookConfigDir, Constants.CliConfigFileName)
	return spellbookConfigFile
}

func GenerateTableHeader() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "Description", "Contents", "Language", "Tags")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
