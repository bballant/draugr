package db

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bballant/draugr/words"
)

func removeAllTerms(index Index, path string) {
	allTerms := index.AllTerms()
	for _, term := range allTerms {
		index.RemoveTerm(term.Token, path)
	}
}

func indexFile(index Index, path string, info fs.FileInfo) error {
	if info.IsDir() {
		log.Printf("not indexing dir %s\n", path)
		return nil
	}
	bs, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	log.Printf("indexing %v\n", path)
	newTokens := words.Tokenize(string(bs))
	removeAllTerms(index, path)
	for _, tok := range newTokens {
		_, err := index.SaveTerm(tok, path)
		if err != nil {
			log.Printf("unable to save term %s", tok)
		}
	}
	return nil
}

func hasExtension(filename string, extensions []string) bool {
	ext := filepath.Ext(filename)
	for _, e := range extensions {
		if strings.EqualFold(e, ext) {
			return true
		}
	}
	return false
}

func IndexPathForExts(index Index, path string, extensions []string) error {
	var walkFile = func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		if info.IsDir() || !hasExtension(info.Name(), extensions) {
			return nil
		}
		return indexFile(index, path, info)
	}
	filepath.Walk(path, walkFile)
	return nil
}
