package main

import (
	"fmt"
	"testing"

	"github.com/bballant/draugr/db"
	"github.com/bballant/draugr/words"
)

func TestWords(t *testing.T) {

	var _db = db.NewMapDB()
	db.IndexPath(&_db, "test_files")

	fmt.Println(_db.TermIndex.GetTerm("name").Count)

	if _db.TermIndex.GetTerm("verrazano").Count != 17 {
		t.Error(`"verrazano" count should be 17`)
	}

	fmt.Println(Search(&_db, words.Tokenize("wood fire")))
}
