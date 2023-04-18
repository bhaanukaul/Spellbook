package main

import (
	"Spellbook/internal/Constants"
	"Spellbook/internal/Spellbook"
	"Spellbook/internal/Utils"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"
)

var spellbook Spellbook.Spellbook

func init() {
	rand.Seed(time.Now().Unix())
	index_file := "spellbook_test"
	var index bleve.Index
	var err error
	if !Utils.FileExists(index_file) {
		log.Print("No bleve index, creating test index.")
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(index_file, mapping)
		if err != nil {
			log.Fatalf("Error creating bleve index: %#v", err)
		}
	} else {
		index, err = bleve.Open(index_file)
		if err != nil {
			log.Fatalf("Error opening bleve index: %#v", err)
			os.Exit(1)
		}
	}

	log.Print("Bleve index exists.")
	spellbook.Index = index

}

func createRandomSpell() (string, string, string, string) {
	descriptions := []string{
		"description 1",
		"description 2",
		"description 3",
		"description 4",
		"description 5",
		"description 6",
	}

	contents := []string{
		"ls -lah",
		"ps -aux",
		"ps -ef",
		"netstat -tulpn",
		"kubctl exec -it $pod -- bash",
		"sleep 10",
	}

	tags := []string{
		"shell,k8s,container",
		"networking,shell",
		"sleep,shell",
		"processes,linux",
		"shell,linux,networking",
		"files,linux,directories",
	}

	languages := []string{
		"shell",
		"python",
		"bash",
		"kubernetes",
		"containers",
	}

	description := descriptions[rand.Intn(len(descriptions))]
	content := contents[rand.Intn(len(contents))]
	tag := tags[rand.Intn(len(tags))]
	language := languages[rand.Intn(len(languages))]
	return description, content, tag, language
}

func createTestSpells() {
	test_spells := []Spellbook.Spell{
		{ID: "1", Language: "containers", Contents: "ps -ef", Description: "description 1", Author: "", Tags: "processes,linux"},
		{ID: "2", Language: "containers", Contents: "sleep 10", Description: "description 2", Author: "", Tags: "shell,k8s,container"},
		{ID: "3", Language: "containers", Contents: "ps -aux", Description: "description 3", Author: "", Tags: "shell,linux,networking"},
		{ID: "4", Language: "shell", Contents: "sleep 10", Description: "description 4", Author: "", Tags: "shell,linux,networking"},
		{ID: "5", Language: "bash", Contents: "kubctl exec -it $pod -- bash", Description: "description 5", Author: "", Tags: "sleep,shell"},
		{ID: "6", Language: "shell", Contents: "ps -ef", Description: "description 6", Author: "", Tags: "sleep,shell"},
		{ID: "7", Language: "containers", Contents: "ps -ef", Description: "description 7", Author: "", Tags: "sleep,shell"},
		{ID: "8", Language: "bash", Contents: "cd ../", Description: "description 8", Author: "", Tags: "shell,k8s,container"},
		{ID: "9", Language: "bash", Contents: "docker exec -it $container /bin/bash", Description: "description 9", Author: "", Tags: "docker,container"},
		{ID: "10", Language: "shell", Contents: "docker ps -a", Description: "description 10", Author: "", Tags: "shell,docker,container"},
	}

	for _, spell := range test_spells {
		_, err := spellbook.CreateSpell(spell)
		if err != nil {
			log.Fatalf("error creating spell: %#v", err)
		}
	}
}

func TestCreateSpell(t *testing.T) {
	description, contents, tags, language := createRandomSpell()

	spellToCreate := Spellbook.Spell{
		Description: description, Language: language, Contents: contents, Tags: tags,
	}
	spell, err := spellbook.CreateSpell(spellToCreate)
	if err != nil {
		t.Errorf("error creating spell: %#v", err)
	}
	if spell.Description != description {
		t.Errorf("Spell description does not match")
	}
	if spell.Language != language {
		t.Errorf("Spell language does not match")
	}
	if spell.Contents != contents {
		t.Errorf("Spell contents does not match")
	}
	if spell.Tags != tags {
		t.Errorf("Spell Tags does not match")
	}

	spellbook.Index.Delete(spell.ID)

}

func TestGetAllSpells(t *testing.T) {
	createTestSpells()
	spells, err := spellbook.GetAllSpells()
	if err != nil {
		t.Errorf("error getting spells: %#v", err)
	}
	if len(spells) != 10 {
		t.Errorf("spell return count incorrect: %d", len(spells))
	}
}

func TestGetSpellByID(t *testing.T) {
	spell, err := spellbook.GetSpellByID("1")
	if err != nil {
		t.Errorf("error getting spells: %#v", err)
	}
	if spell.ID != "1" {
		t.Error("Spell id does not match")
	}

}

func TestGetSpellsByTags(t *testing.T) {
	spells, err := spellbook.FindSpellsByTag("shell", Constants.SearchResultLimit)
	if err != nil {
		t.Errorf("error getting spells: %#v", err)
	}
	if len(spells) != 8 {
		t.Errorf("Num spells found: %d", len(spells))
	}
}

func TestGetSpellsByTagsWithLimit(t *testing.T) {
	spells, err := spellbook.FindSpellsByTag("shell", 5)
	if err != nil {
		t.Errorf("error getting spells: %#v", err)
	}
	if len(spells) != 5 {
		t.Errorf("Num spells found: %d", len(spells))
	}
}

func TestGetSpellsByDescription(t *testing.T) {
	spells, err := spellbook.FindSpellsByDescription("10", Constants.SearchResultLimit)
	if err != nil {
		t.Errorf("error getting spells: %#v", err)
	}
	if len(spells) != 1 {
		t.Errorf("Num spells found: %d", len(spells))
	}
}

func TestGetSpellsByAuthor(t *testing.T) {

}

func TestUpdateSpell(t *testing.T) {
	spell, err := spellbook.GetSpellByID("2")
	if err != nil {
		t.Errorf("error getting spells: %#v", err)
	}
	spell.Language = "shell"
	updatedSpell, err := spellbook.UpdateSpell("2", *spell)
	if err != nil {
		t.Errorf("error getting spells: %#v", err)
	}
	if updatedSpell.Language != "shell" {
		t.Error("Spell not updated")
	}

}
