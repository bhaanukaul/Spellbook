package main

import (
	"Spellbook/internal/Constants"
	"Spellbook/internal/Spell"
	"Spellbook/internal/Utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
)

type RemoteSection struct {
	Alias       string `json:"alias"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

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
				Action: AddRemoteServer,
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

/*
CLI Functions
*/
func StartServer(c *cli.Context) error {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/spell/:id", GetSpell)
		api.GET("/spells", GetAllSpells)
		api.POST("/spell", CreateSpell)
		api.PATCH("/spell/:id", UpdateSpell)
		api.POST("/spellbook", AddSpellBook)
		api.GET("/spellbooks", GetAllSpellbooks)
	}
	router.GET("/ping", Ping)
	configFile := GetServerConfig()
	port := Utils.GetKVFromConfig(configFile, "http_port", "server")
	router.Run("0.0.0.0:" + port)
	return nil
}

func AddRemoteServer(c *cli.Context) error {
	if c.Args().Len() != 2 {
		Utils.Error("\"spellbook-server add\" requires 2 arguments.", nil)
	}
	alias := c.Args().Get(0)
	spellbookUrl := c.Args().Get(1)
	splitSpellbook := strings.Split(spellbookUrl, ".")
	extension := splitSpellbook[len(splitSpellbook)-1]
	if extension != "json" {
		Utils.Error("Remote needs to be a json file.", nil)
	}

	section := "remotes." + alias

	configDir := GetServerConfig()
	Utils.AddKVToConfig(configDir, "url", spellbookUrl, section)
	Utils.AddKVToConfig(configDir, "alias", alias, section)

	return nil
}

func SetupServer(c *cli.Context) error {
	var serverHome string
	var err error
	serverHome, err = os.Getwd()
	port := c.String("port")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	configFileName := "spellbook-server.ini"
	spellbookConfigDir := GetServerConfigDir(serverHome)
	err = os.MkdirAll(spellbookConfigDir, os.ModePerm)
	if err != nil {
		Utils.Error("Error", err)
	}

	spellbookRemotesDir := GetServerRemotesDir(serverHome)
	err = os.MkdirAll(spellbookRemotesDir, os.ModePerm)
	if err != nil {
		Utils.Error("Error", err)
	}

	spellbookConfigFile := spellbookConfigDir + "/" + configFileName
	emptyConfig, err := os.Create(spellbookConfigFile)

	if err != nil {
		Utils.Error("Error", err)
	}
	emptyConfig.Close()
	Utils.AddKVToConfig(spellbookConfigFile, "http_port", port, "server")
	return nil
}

/*
API Functions
*/

func GetAllSpellbooks(c *gin.Context) {
	cd := GetServerConfig()

	cfg, err := ini.Load(cd)
	if err != nil {
		Utils.Error("Cannot load config.", err)
	}

	names := cfg.ChildSections("remotes")
	var returnSb []RemoteSection
	for _, n := range names {
		toAdd := RemoteSection{
			Alias: n.Key("alias").String(), Description: n.Key("description").String(), Url: n.Key("url").String(),
		}
		returnSb = append(returnSb, toAdd)
	}
	c.JSON(http.StatusOK, returnSb)
}

func AddSpellBook(c *gin.Context) {

}

func UpdateSpell(c *gin.Context) {
	var spellToUpdate Spell.Spell
	if err := c.BindJSON(&spellToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Request Body")
		return
	}
	spell_id_param := c.Param("id")
	spell_id, err := strconv.Atoi(spell_id_param)
	if err != nil {
		panic(err)
	}

	Spell.UpdateSpell(spell_id, spellToUpdate)
}

func CreateSpell(c *gin.Context) {
	// language := c.PostForm("language")
	// content := c.PostForm("content")
	// description := c.PostForm("description")
	// tags := c.PostForm("tags")

	var newSpell Spell.Spell
	if err := c.BindJSON(&newSpell); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Request Body")
		return
	}

	fmt.Printf("From POST CreateSpell: %s, %s, %s, %s", newSpell.Language, newSpell.Contents, newSpell.Description, newSpell.Tags)

	spell, err := Spell.CreateSpell(newSpell.Language, newSpell.Contents, newSpell.Description, newSpell.Tags)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, spell)

}

func GetSpell(c *gin.Context) {
	spell_id_param := c.Param("id")
	spell_id, err := strconv.Atoi(spell_id_param)
	if err != nil {
		panic(err)
	}
	spell, err := Spell.GetSpellByID(spell_id)
	if err != nil {
		panic(err)
	}
	if spell.ID == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(http.StatusOK, spell)
}

func GetAllSpells(c *gin.Context) {
	spells, err := Spell.GetAllSpells()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, spells)
}

func Ping(c *gin.Context) {
	ping := Utils.SpellbookPing{
		Version: "1.0",
	}

	c.JSON(http.StatusOK, ping)
}

/*
Util Functions
*/

func GetServerConfig() string {
	var serverHome string
	var err error
	serverHome, err = os.Getwd()

	if err != nil {
		Utils.Error("Error getting working directory", err)
	}
	spellbookConfigDir := GetServerConfigDir(serverHome)

	spellbookConfigFile := spellbookConfigDir + "/" + Constants.ConfigFile
	return spellbookConfigFile
}

func GetServerConfigDir(basedir string) string {
	spellbookConfigDir := filepath.Join(basedir, "config")
	return spellbookConfigDir
}

func GetServerRemotesDir(basedir string) string {
	spellbookRemotesDir := filepath.Join(basedir, "remotes")
	return spellbookRemotesDir
}
