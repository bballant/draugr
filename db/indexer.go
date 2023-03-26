package db

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/bballant/draugr/words"
)

func addsAndRems(newTokens []string, oldTokens []string) (adds []string, rems []string) {
	sort.Strings(newTokens)
	sort.Strings(oldTokens)
	n := 0
	o := 0
	for {
		if n == len(newTokens) {
			rems = append(rems, oldTokens[o:]...)
			break
		}
		if o == len(oldTokens) {
			adds = append(adds, newTokens[n:]...)
			break
		}
		nstr := newTokens[n]
		ostr := oldTokens[o]
		if nstr > ostr {
			o++
			rems = append(rems, ostr)
		} else if ostr > nstr {
			n++
			adds = append(adds, ostr)
		} else {
			o++
			n++
		}
	}
	return adds, rems
}

func indexFile(db *DB, path string, info fs.FileInfo) error {
	prevModTime := db.PathIndex.GetModTime(path)
	currModTime := info.ModTime()
	if !currModTime.After(prevModTime) {
		log.Printf("'%v' <= '%v', not indexing %s\n", currModTime, prevModTime, path)
		return nil
	}
	if info.IsDir() {
		log.Printf("not indexing dir %s\n", path)
		return db.PathIndex.SetModTime(path, currModTime)
	}
	bs, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	log.Printf("'%v' > '%v', indexing %v\n", currModTime, prevModTime, path)
	newTokens := words.Tokenize(string(bs))
	oldTokens := db.PathTermIndex.GetTerms(path)
	addToks, remToks := addsAndRems(newTokens, oldTokens)
	for _, tok := range addToks {
		db.TermIndex.SaveTerm(tok, path)
	}
	for _, tok := range remToks {
		db.TermIndex.RemoveTerm(tok, path)
	}
	err = db.PathTermIndex.SetTerms(path, newTokens)
	if err != nil {
		return err
	}
	return db.PathIndex.SetModTime(path, currModTime)
}

func IndexPath(db *DB, path string) error {
	var walkFile func(string, fs.FileInfo, error) error
	walkFile = func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		return indexFile(db, path, info)
	}
	filepath.Walk(path, walkFile)
	return nil
}
