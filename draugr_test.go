package main

import (
	"testing"

	"github.com/bballant/draugr/db"
	"github.com/bballant/draugr/words"
)

func TestWords(t *testing.T) {

	var _db = db.NewMapIndex()
	db.IndexPathForExts(_db, "test_files", []string{".txt"})

	if len(_db.GetTerm("verrazano").Paths) != 17 {
		t.Error(`"verrazano" count should be 17`)
	}

	if db.SearchIndex(_db, words.Tokenize("wood fire"))[0].Path !=
		"test_files/essays/The-Hudson-River-And-Its-Early-Names-Susan-Fenimore-Cooper.txt" {
		t.Error(`"wood fire" should return "The Hudson ..." first`)

	}

	verneInTitle := db.SearchIndex(_db, words.Tokenize("Verne"))

	if verneInTitle[0].Count != 6 {
		t.Error(`"Verne" count should be 6 for one book because of title`)
	}
}
