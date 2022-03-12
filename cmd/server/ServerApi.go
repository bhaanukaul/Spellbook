package main

import (
	"Spellbook/internal/Spell"
	"Spellbook/internal/Utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

/*
API Functions
*/

type RemoteSection struct {
	Alias       string `json:"alias"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

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

func AddRemoteSpellbookApi(c *gin.Context) {
	var remoteSection RemoteSection
	if err := c.BindJSON(&remoteSection); err != nil {
		Utils.Error("", err)
		return
	}

	AddRemoteServer(remoteSection.Alias, remoteSection.Url, remoteSection.Description)
	c.JSON(http.StatusCreated, remoteSection)
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