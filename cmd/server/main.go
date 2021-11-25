package main

import (
	"Spellbook/internal/Spell"
	"net/http"

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

}

func GetAllSpells(c *gin.Context) {
	spells, err := Spell.GetAllSpells()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, spells)
}

func CreateSpell(c *gin.Context) {

}

func UpdateSpell(c *gin.Context) {

}

func DeleteSpell(c *gin.Context) {

}
