package main

import (
	"Spellbook/internal/Spell"
	"Spellbook/internal/Utils"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"

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
				Subcommands: []*cli.Command{
					{
						Name:  "remote",
						Usage: "initialize local Spellbook to sync with remote Spellbook server.",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "url",
								Aliases:  []string{"u"},
								Usage:    "Username for remote server",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "username",
								Aliases:  []string{},
								Usage:    "Username for remote server",
								Required: false,
							},
						},
						Action: func(c *cli.Context) error {
							// cfg.Section("").Key("remote").SetValue("")
							// Validate the remote url is a valid spellbook server
							remoteUrl := c.String("url")
							fmt.Println("Checking if remote server is up and running: " + remoteUrl)
							resp, err := http.Get(remoteUrl + "/ping")
							if err != nil {
								log.Fatal(err)
							}
							var ping Utils.SpellbookPing
							defer resp.Body.Close()
							body, err := io.ReadAll(resp.Body)
							if err != nil {
								log.Fatal(err)
							}

							err = json.Unmarshal(body, &ping)

							if resp.StatusCode != http.StatusOK || err != nil {
								log.Fatal("Error with remote Spellbook url: " + err.Error())
								return nil
							}

							if ping.Version != "1.0" {
								log.Fatal("Something not right")
								return nil
							}
							fmt.Println("Looks like we're good to go!")
							return nil
						},
					},
				},
				Action: func(c *cli.Context) error {
					userHome, err := os.UserHomeDir()
					if err != nil {
						log.Fatal(err)
					}

					spellbookFileName := c.String("name")
					configFileName := "Spellbook.ini"

					user, err := user.Current()
					if err != nil {
						log.Fatal(err)
					}
					username := user.Username
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

					now := time.Now()
					secs := now.Unix()
					userHashPlain := fmt.Sprintf("%s-%d", username, secs)
					userHash := fmt.Sprintf("%x", sha256.Sum256([]byte(userHashPlain)))
					Utils.AddKVToConfig(spellbookConfigFile, "spellbookdb", spellBookDB, "")
					Utils.AddKVToConfig(spellbookConfigFile, "username", username, "")
					Utils.AddKVToConfig(spellbookConfigFile, "user_hash", userHash, "")

					db := Utils.GetDatabaseConnection()
					db.Migrator().CreateTable(&Spell.Spell{})
					return nil
				},
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
				Action: func(c *cli.Context) error {
					tbl := GenerateTableHeader()

					var spellToUpdate Spell.Spell
					spell_language := c.String("language")
					// fmt.Printf("Language: %s", c.String("language"))
					spell_content := c.String("content")
					spell_description := c.String("description")
					spell_tags := c.String("tags")
					// int_id, err := strconv.Atoi(c.Args().First())
					spell_id := c.Int("id")

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
						spell, err := Spell.GetSpellByID(spell_id)
						if err != nil {
							log.Fatalf("%s", err)
							return err
						}
						newTags := spell.Tags + "," + spell_tags
						spellToUpdate.Tags = newTags
					}
					fmt.Printf("Updating %d with %#v\n", spell_id, spellToUpdate)

					_, err := Spell.UpdateSpell(spell_id, spellToUpdate)
					if err != nil {
						log.Fatalf("%s", err)
						return err
					}

					updatedSpell, err := Spell.GetSpellByID(spell_id)
					if err != nil {
						log.Fatalf("%s", err)
						return err
					}
					tbl.AddRow(updatedSpell.ID, updatedSpell.Description, updatedSpell.Contents,
						updatedSpell.Language, updatedSpell.Tags, updatedSpell.Author)
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

func GenerateTableHeader() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "Description", "Contents", "Language", "Tags")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
