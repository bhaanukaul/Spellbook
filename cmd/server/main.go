package main

import (
	"Spellbook/internal/Spell"
	"Spellbook/internal/Utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
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
				Usage:   "Initialize the Spellbook server.",
				Action: func(c *cli.Context) error {

					return nil
				},
			},
			{
				Name:    "start",
				Aliases: []string{},
				Usage:   "Start the Spellbook server.",
				Action: func(c *cli.Context) error {
					router := gin.Default()

					api := router.Group("/api")
					{
						api.GET("/spell/:id", GetSpell)
						api.GET("/spells", GetAllSpells)
						api.POST("/spell", CreateSpell)
						api.PATCH("/spell/:id", UpdateSpell)
					}
					router.GET("/ping", Ping)

					router.Run("localhost:8081")
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

func SetupServer() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configFileName := "Spellbook-server.ini"
	spellbookConfigDir := filepath.Join(userHome, ".config", "Spellbook")
	err = os.MkdirAll(spellbookConfigDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	spellbookConfigFile := spellbookConfigDir + "/" + configFileName
	emptyConfig, err := os.Create(spellbookConfigFile)

	if err != nil {
		log.Fatal(err)
	}
	emptyConfig.Close()

}

func Ping(c *gin.Context) {
	ping := Utils.SpellbookPing{
		Version: "1.0",
	}

	c.JSON(http.StatusOK, ping)
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
