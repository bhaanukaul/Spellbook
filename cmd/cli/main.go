package main

import (
	"Spellbook/internal/Constants"
	"log"
	"os"

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
				Action:  SpellbookFind,
				Subcommands: []*cli.Command{
					{
						Name:    "tag",
						Aliases: []string{"t"},
						Usage:   "Find spells based on a tag.\nSpellbook find tag",
						Action:  SpellbookFinByTag,
					},
					{
						Name:    "id",
						Aliases: []string{},
						Usage:   "Return spell based on a ID in database.\nSpellbook find tag",
						Action:  SpellbookFindById,
					},
				},
			},
			{
				Name:    "add",
				Aliases: []string{},
				Usage:   "Add spell to spellbook",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "language",
						Aliases:  []string{"l"},
						Usage:    "language of the spell",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "content",
						Aliases:  []string{"c"},
						Usage:    "content of the spell",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "description",
						Aliases:  []string{"d"},
						Usage:    "what does the spell do.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "tags",
						Aliases:  []string{"t"},
						Usage:    "associated tags for the spell",
						Required: true,
					},
				},
				Action: SpellbookAddSpell,
			},
			{
				Name:    "init",
				Aliases: []string{},
				Usage:   "initialize local spellbook database",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   Constants.CliConfigBleveIndexFileName,
						Usage:   "Name of the Spellbook file. Default is Spellbook.bleve",
					},
				},
				Action: SpellbookInit,
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "updates spell in spellbook",
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
					&cli.IntFlag{
						Name:     "id",
						Usage:    "ID of the spell to update",
						Required: true,
					},
				},
				Action: SpellbookUpdateSpell,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
