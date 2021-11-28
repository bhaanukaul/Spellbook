package main

import (
	"Spellbook/internal/Spell"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/spell/:id", GetSpell)
	router.GET("/spells", GetAllSpells)
	router.POST("/spell", CreateSpell)
	router.PATCH("/spell/:id", UpdateSpell)
	router.DELETE("/spell/:id", DeleteSpell)

	router.Run("localhost:8081")
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

}

func DeleteSpell(c *gin.Context) {

}
