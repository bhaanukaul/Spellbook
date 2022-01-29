package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "Spellbook-Server",
		Usage:                "Create a reference for CLI commands so you don't have to remember them!",
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "8081",
						Usage:   "Port to run the spellbook server.",
					},
				},
				Usage:  "Initialize the Spellbook server. Creates config file in directory the command was ran.",
				Action: SetupServer,
			},
			{
				Name:    "add",
				Aliases: []string{},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "alias",
						Aliases:  []string{"a"},
						Required: true,
						Usage:    "One word alias for the remote spellbook.",
					},
					&cli.StringFlag{
						Name:     "url",
						Aliases:  []string{"u"},
						Required: true,
						Usage:    "Url of the remote spellbook.",
					},
					&cli.StringFlag{
						Name:     "description",
						Aliases:  []string{"d"},
						Required: true,
						Usage:    "Description for the remote spellbook.",
					},
				},
				Usage:  "Add a remote server.",
				Action: AddRemoteSpellbookCli,
			},
			{
				Name:    "sync-remote",
				Aliases: []string{},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "alias",
						Aliases: []string{"a"},
						Usage:   "Alias to sync to local server.",
					},
				},
				Usage:  "Sync a remote Spellbook file to the local server. If no alias is provided, then command will sync all aliases in config.",
				Action: SyncRemoteSpellbook,
			},
			{
				Name:    "start",
				Aliases: []string{},
				Usage:   "Start the Spellbook server.",
				Action:  StartServer,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
