package main

import (
	"Spellbook/internal/Constants"
	"Spellbook/internal/Utils"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

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

func AddRemoteServer(alias string, url string, description string) {

	section := "remotes." + alias
	configDir := GetServerConfig()

	checkSection := Utils.DoesSectionExist(configDir, section)
	log.Printf("%t", checkSection)
	Utils.AddKVToConfig(configDir, "url", url, section)
	Utils.AddKVToConfig(configDir, "alias", alias, section)
	Utils.AddKVToConfig(configDir, "description", description, section)

}

func SyncRemote(alias string) {
	config := GetServerConfig()
	cfg, err := ini.Load(config)
	if err != nil {
		Utils.Error("", err)
	}
	sectionExists := Utils.DoesSectionExist(config, alias)
	if !sectionExists {
		return
	}
	section, err := cfg.GetSection(alias)
	if err != nil {
		Utils.Error("", err)
	}
	remote := section.Key("url")
	log.Printf("remote: %s", remote)

}
