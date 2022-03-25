package main

import (
	"Spellbook/internal/Constants"
	"Spellbook/internal/Spell"
	"Spellbook/internal/Utils"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/urfave/cli/v2"
)

func SpellbookInit(c *cli.Context) error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Name of the config file that will store the spellbook cli configuration
	configFileName := Constants.CliConfigFileName

	// Creating a path in /home/<user>/.config/Spellbook (*nix)
	// /Users/<user>/.config/Spellbook (MacOs)
	// C:\Users\<user>\.config\Spellbook (Win)
	spellbookConfigDir := filepath.Join(userHome, ".config", "Spellbook")
	err = os.MkdirAll(spellbookConfigDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Create the initial empty config file
	spellbookConfigFile := fmt.Sprintf("%s/%s", spellbookConfigDir, configFileName)
	emptyConfig, err := os.Create(spellbookConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	emptyConfig.Close()

	// Name of the index bleve will use to store data
	spellbookFileName := c.String("name")
	// Creating the empty bleve index
	spellbookIndex := fmt.Sprintf("%s/%s", spellbookConfigDir, spellbookFileName)
	mapping := bleve.NewIndexMapping()
	_, err = bleve.New(spellbookIndex, mapping)
	if err != nil {
		Utils.Error("Error", err)
	}
	// Add the bleve index to the config for future use.
	Utils.AddKVToConfig(spellbookConfigFile, Constants.CliConfigBleveIndexKey, spellbookIndex, "")

	return nil
}

func SpellbookFind(c *cli.Context) error {
	tbl := GenerateTableHeader()
	index := GetBleveIndex()
	searchResult, err := Spell.GetAllSpells(index)
	if err != nil {
		log.Fatalf("err: %s", err)
		return err
	}
	for _, spell := range searchResult.Hits {
		tbl.AddRow(spell.Fields["ID"], spell.Fields["Description"], spell.Fields["Contents"], spell.Fields["Language"], spell.Fields["Tags"])
	}
	tbl.Print()
	index.Close()
	return nil
}

func SpellbookFinByTag(c *cli.Context) error {
	tbl := GenerateTableHeader()
	index := GetBleveIndex()
	// fmt.Printf("args %s with len %d", c.Args(), c.Args().Len())

	searchResult, err := Spell.FindSpellsByTag(c.Args().First(), index)
	if err != nil {
		log.Fatalf("%s", err)
		return err
	}
	for _, spell := range searchResult.Hits {
		tbl.AddRow(spell.ID, spell.Fields["Description"], spell.Fields["Contents"], spell.Fields["Language"], spell.Fields["Tags"])
	}
	tbl.Print()
	index.Close()
	return nil
}

func SpellbookFindById(c *cli.Context) error {
	// fmt.Printf("args %s with len %d", c.Args(), c.Args().Len())
	tbl := GenerateTableHeader()

	// int_id, err := strconv.Atoi(c.Args().First())
	id := c.Args().First()
	index := GetBleveIndex()
	spell, err := Spell.GetSpellByID(id, index)
	if err != nil {
		log.Fatalf("%s", err)
		return err
	}
	if spell.Hits.Len() == 0 {
		fmt.Printf("Spell with id %s cannot be found.", id)
	} else {
		foundSpell := spell.Hits[0]
		tbl.AddRow(foundSpell.Fields["ID"], foundSpell.Fields["Description"], foundSpell.Fields["Contents"], foundSpell.Fields["Language"], foundSpell.Fields["Tags"])
		tbl.Print()
	}
	index.Close()
	return nil
}

func SpellbookAddSpell(c *cli.Context) error {
	log.Printf("Flags from add language: %s\ncontent: %s\ndescription: %s\ntags: %s",
		c.String("language"), c.String("content"), c.String("description"), c.String("tags"))
	tbl := GenerateTableHeader()
	index := GetBleveIndex()
	newSpell := Spell.Spell{
		Description: c.String("description"), Language: c.String("language"), Contents: c.String("content"), Tags: c.String("tags"),
	}
	newSpell, err := Spell.CreateSpell(newSpell, index)
	if err != nil {
		log.Fatalf("%s", err)
		return err
	}

	spell, err := Spell.GetSpellByID(newSpell.ID, index)
	if err != nil {
		log.Fatalf("%s", err)
		return err
	}
	if spell.Hits.Len() == 0 {
		fmt.Printf("Spell with id %s cannot be found.", newSpell.ID)
	} else {
		foundSpell := spell.Hits[0]
		tbl.AddRow(foundSpell.Fields["ID"], foundSpell.Fields["Description"], foundSpell.Fields["Contents"], foundSpell.Fields["Language"], foundSpell.Fields["Tags"])
		tbl.Print()
	}

	index.Close()
	return nil
}

func SpellbookUpdateSpell(c *cli.Context) error {
	spell_id := c.String("id")
	index := GetBleveIndex()
	searchResult, err := Spell.GetSpellByID(spell_id, index)
	if err != nil {
		log.Fatalf("%s", err)
		return err
	}
	// var spell bleve.DocumentMatch
	// var foundSpell *search.DocumentMatch
	if searchResult.Hits.Len() == 0 {
		fmt.Printf("Spell with id %s cannot be found to update.", spell_id)
		return nil
	}
	foundSpell := searchResult.Hits[0]

	tbl := GenerateTableHeader()

	var spellToUpdate Spell.Spell
	spell_language := c.String("language")
	// fmt.Printf("Language: %s", c.String("language"))
	spell_content := c.String("content")
	spell_description := c.String("description")
	spell_tags := c.String("tags")
	// int_id, err := strconv.Atoi(c.Args().First())

	if spell_language != "" {
		spellToUpdate.Language = spell_language
	}

	if spell_content != "" {
		spellToUpdate.Contents = spell_content
	}

	if spell_description != "" {
		spellToUpdate.Description = spell_description
	}

	if spell_tags != "" {
		old_tags := foundSpell.Fields["Tags"]
		newTags := fmt.Sprintf("%s,%s", old_tags, spell_tags)
		spellToUpdate.Tags = newTags
	}
	fmt.Printf("Updating %s with %#v\n", spell_id, spellToUpdate)

	_, err = Spell.UpdateSpell(spell_id, spellToUpdate, index)
	if err != nil {
		log.Fatalf("%s", err)
		return err
	}

	updatedSpell, err := Spell.GetSpellByID(spell_id, index)
	if searchResult.Hits.Len() == 0 {
		fmt.Printf("Spell with id %s cannot be found after update.", spell_id)
		return nil
	}
	foundSpell = updatedSpell.Hits[0]
	if err != nil {
		log.Fatalf("%s", err)
		return err
	}

	tbl.AddRow(foundSpell.Fields["ID"], foundSpell.Fields["Description"], foundSpell.Fields["Contents"], foundSpell.Fields["Language"], foundSpell.Fields["Tags"])

	tbl.Print()
	index.Close()
	return nil
}
