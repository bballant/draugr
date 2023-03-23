package main

import (
	"testing"

	"github.com/bballant/draugr/db"
)

func TestWords(t *testing.T) {

	var db = db.NewMapDB()
	IndexPath(&db, ".")

	if db.TermIndex.GetTerm("files").Count != 6 {
		t.Error(`"files" count should be 6`)
	}
}
