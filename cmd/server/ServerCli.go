package main

import (
	"Spellbook/internal/Constants"
	"Spellbook/internal/Utils"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
)

/*
CLI Functions
*/

func SyncRemoteSpellbook(c *cli.Context) error {
	alias := c.String("alias")
	log.Printf("alias: %s", alias)
	if alias != "" {
		SyncRemote(alias)
		return nil
	}
	config := GetServerConfig()
	cfg, err := ini.Load(config)
	if err != nil {
		Utils.Error("", err)
	}
	parent, err := cfg.GetSection(Constants.RemotesParent)
	if err != nil {
		Utils.Error("", err)
	}
	sections := parent.ChildSections()
	for _, s := range sections {
		SyncRemote(s.Name())
	}
	return nil
}

func StartServer(c *cli.Context) error {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/spell/:id", GetSpell)
		api.GET("/spells", GetAllSpells)
		api.POST("/spell", CreateSpell)
		api.PATCH("/spell/:id", UpdateSpell)
		// api.POST("/spellbook", AddRemoteSpellbookApi)
		api.GET("/spellbooks", GetAllSpellbooks)
	}
	router.GET("/ping", Ping)
	configFile := GetServerConfig()
	port := Utils.GetKVFromConfig(configFile, "http_port", "server")
	router.Run("0.0.0.0:" + port)
	return nil
}

func AddRemoteSpellbookCli(c *cli.Context) error {
	if c.Args().Len() != 3 {
		Utils.Error("\"spellbook-server add\" requires 2 arguments.", nil)
	}
	alias := c.String("alias")
	spellbookUrl := c.String("url")
	description := c.String("description")
	splitSpellbook := strings.Split(spellbookUrl, ".")
	extension := splitSpellbook[len(splitSpellbook)-1]
	if extension != "json" {
		Utils.Error("Remote needs to be a json file.", nil)
	}

	AddRemoteServer(alias, spellbookUrl, description)
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
