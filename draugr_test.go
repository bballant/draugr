package main

import (
	"fmt"
	"testing"

	"github.com/bballant/draugr/db"
	"github.com/bballant/draugr/words"
)

func TestWords(t *testing.T) {

	var db = db.NewMapDB()
	IndexPath(&db, "test_files")

	fmt.Println(db.TermIndex.GetTerm("name").Count)

	if db.TermIndex.GetTerm("verrazano").Count != 17 {
		t.Error(`"verrazano" count should be 17`)
	}

	fmt.Println(Search(&db, words.Tokenize("wood fire")))
}
