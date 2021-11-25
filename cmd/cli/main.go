package main

import (
	"Spellbook/internal/Spell"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/rodaine/table"

	"github.com/urfave/cli/v2"
)

func main() {

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

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
					tbl := table.New("ID", "Description", "Contents", "Language", "Tags")
					tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

					spells, err := Spell.GetAllSpells()
					if err != nil {
						fmt.Errorf("%s", err)
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
							tbl := table.New("ID", "Description", "Contents", "Language", "Tags")
							tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
							fmt.Printf("args %s with len %d", c.Args(), c.Args().Len())

							spells, err := Spell.FindSpellsByTag(c.Args().First())
							if err != nil {
								fmt.Errorf("%s", err)
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

							tbl := table.New("ID", "Description", "Contents", "Language", "Tags")
							tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
							fmt.Printf("args %s with len %d", c.Args(), c.Args().Len())

							int_id, err := strconv.Atoi(c.Args().First())
							if err != nil {
								fmt.Errorf("%s", err)
								return err
							}
							spell, err := Spell.GetSpellByID(int_id)
							if err != nil {
								fmt.Errorf("%s", err)
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
					fmt.Printf("Flags from add language: %s\ncontent: %s\ndescription: %s\ntags: %s",
						c.String("language"), c.String("content"), c.String("description"), c.String("tags"))
					tbl := table.New("ID", "Description", "Contents", "Language", "Tags")
					tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
					spell_id, err := Spell.CreateSpell(c.String("language"), c.String("content"), c.String("description"), c.String("tags"))
					if err != nil {
						fmt.Errorf("%s", err)
						return err
					}

					spell, err := Spell.GetSpellByID(spell_id)
					if err != nil {
						fmt.Errorf("%s", err)
						return err
					}
					tbl.AddRow(spell.ID, spell.Description, spell.Contents, spell.Language, spell.Tags)
					tbl.Print()

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

func GenerateTableHeader() (table.Table, error) {

}
