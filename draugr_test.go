package main

import (
	"fmt"
	"testing"

	"github.com/bballant/draugr/db"
	"github.com/bballant/draugr/words"
)

func TestWords(t *testing.T) {

	var _db = db.NewMapDB()
	db.IndexPathForExts(&_db, "test_files", []string{".txt"})

	fmt.Println(_db.TermIndex.GetTerm("name"))
	fmt.Println(_db.TermIndex.GetTerm("name").Count)

	if _db.TermIndex.GetTerm("verrazano").Count != 17 {
		t.Error(`"verrazano" count should be 17`)
	}

	if Search(&_db, words.Tokenize("wood fire"))[0].Path !=
		"test_files/essays/The-Hudson-River-And-Its-Early-Names-Susan-Fenimore-Cooper.txt" {
		t.Error(`"wood fire" should return "The Hudson ..." first`)

	}

	verneInTitle := Search(&_db, words.Tokenize("Verne"))

	if verneInTitle[0].Count != 6 {
		t.Error(`"Verne" count should be 6 for one book because of title`)
	}

	uglyPath := "/home/bballant/code/stm32/stm32f103c8t6/rtos/FreeRTOS-202112.00/FreeRTOS-Plus/Source/AWS/device-defender/docs/doxygen/include/size_table.md"
	if pathEnd(uglyPath, 10) != ">_table.md" {
		t.Error(`pathEnd not returning expected value`)
	}
}
