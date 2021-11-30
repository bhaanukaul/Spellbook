package main

import (
	"Spellbook/internal/Spell"
	"Spellbook/internal/Utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"gopkg.in/ini.v1"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "Spellbook",
		Usage:                "Create a reference for CLI commands so you don't have to remember them!",
		Commands: []*cli.Command{
			{
				Name:    "find",
				Aliases: []string{"f"},
				Usage:   "return all spells in your Spellbook",
				Action: func(c *cli.Context) error {
					tbl := GenerateTableHeader()

					spells, err := Spell.GetAllSpells()
					if err != nil {
						log.Fatalf("err: %s", err)
						return err
					}
					for _, spell := range spells {
						tbl.AddRow(spell.ID, spell.Description, spell.Contents, spell.Language, spell.Tags)
					}
					tbl.Print()

					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:    "tag",
						Aliases: []string{"t"},
						Usage:   "Find spells based on a tag.\nSpellbook find tag",
						Action: func(c *cli.Context) error {
							tbl := GenerateTableHeader()

							// fmt.Printf("args %s with len %d", c.Args(), c.Args().Len())

							spells, err := Spell.FindSpellsByTag(c.Args().First())
							if err != nil {
								log.Fatalf("%s", err)
								return err
							}
							for _, spell := range spells {
								tbl.AddRow(spell.ID, spell.Description, spell.Contents, spell.Language, spell.Tags)
							}
							tbl.Print()
							return nil
						},
					},
					{
						Name:    "id",
						Aliases: []string{},
						Usage:   "Return spell based on a ID in database.\nSpellbook find tag",
						Action: func(c *cli.Context) error {

							// fmt.Printf("args %s with len %d", c.Args(), c.Args().Len())
							tbl := GenerateTableHeader()

							int_id, err := strconv.Atoi(c.Args().First())
							if err != nil {
								log.Fatalf("%s", err)
								return err
							}
							spell, err := Spell.GetSpellByID(int_id)
							if err != nil {
								log.Fatalf("%s", err)
								return err
							}
							tbl.AddRow(spell.ID, spell.Description, spell.Contents, spell.Language, spell.Tags)
							tbl.Print()

							return nil
						},
					},
				},
			},
			{
				Name:    "add",
				Aliases: []string{},
				Usage:   "Add spell to spellbook",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "language",
						Aliases: []string{"l"},
						Usage:   "language of the spell",
					},
					&cli.StringFlag{
						Name:    "content",
						Aliases: []string{"c"},
						Usage:   "content of the spell",
					},
					&cli.StringFlag{
						Name:    "description",
						Aliases: []string{"d"},
						Usage:   "what does the spell do.",
					},
					&cli.StringFlag{
						Name:    "tags",
						Aliases: []string{"t"},
						Usage:   "associated tags for the spell",
					},
				},
				Action: func(c *cli.Context) error {
					log.Printf("Flags from add language: %s\ncontent: %s\ndescription: %s\ntags: %s",
						c.String("language"), c.String("content"), c.String("description"), c.String("tags"))
					tbl := GenerateTableHeader()

					newSpell, err := Spell.CreateSpell(c.String("language"), c.String("content"), c.String("description"), c.String("tags"))
					if err != nil {
						log.Fatalf("%s", err)
						return err
					}

					spell, err := Spell.GetSpellByID(newSpell.ID)
					if err != nil {
						log.Fatalf("%s", err)
						return err
					}
					tbl.AddRow(spell.ID, spell.Description, spell.Contents, spell.Language, spell.Tags)
					tbl.Print()

					return nil
				},
			},
			{
				Name:    "init",
				Aliases: []string{},
				Usage:   "initialize local spellbook database",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "Spellbook.db",
						Usage:   "Name of the Spellbook file. Default is Spellbook.db",
					},
				},
				Action: func(c *cli.Context) error {
					userHome, err := os.UserHomeDir()
					spellbookFileName := c.String("name")
					configFileName := "Spellbook.ini"
					if err != nil {
						log.Fatal(err)
					}
					spellbookConfigDir := filepath.Join(userHome, ".config", "Spellbook")
					err = os.MkdirAll(spellbookConfigDir, os.ModePerm)
					if err != nil {
						log.Fatal(err)
					}
					spellBookDB := spellbookConfigDir + "/" + spellbookFileName
					emptyDB, err := os.Create(spellBookDB)
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Created DB at %s\n", spellBookDB)
					emptyDB.Close()

					spellbookConfigFile := spellbookConfigDir + "/" + configFileName
					emptyConfig, err := os.Create(spellbookConfigFile)
					if err != nil {
						log.Fatal(err)
					}
					// fmt.Printf("Created config at %s\n", spellbookConfigFile)
					emptyConfig.Close()

					cfg, err := ini.Load(spellbookConfigFile)
					if err != nil {
						fmt.Printf("Fail to read file: %v", err)
						os.Exit(1)
					}
					cfg.Section("").Key("spellbookdb").SetValue(spellBookDB)
					cfg.SaveTo(spellbookConfigFile)

					db := Utils.GetDatabaseConnection()
					db.Migrator().CreateTable(&Spell.Spell{})
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateTableHeader() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "Description", "Contents", "Language", "Tags")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
