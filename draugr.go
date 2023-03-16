package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/bballant/draugr/db"
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

func indexFile(db *db.DB, path string, info fs.FileInfo) error {
	prevModTime := db.PathIndex.GetModTime(path)
	currModTime := info.ModTime()
	if !currModTime.After(prevModTime) {
		log.Printf("%v is not after %v, not indexing %s\n", currModTime, prevModTime, path)
		return nil
	}
	if info.IsDir() {
		log.Printf("not indexing dir %s\n", path)
		return nil
	}
	bs, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	log.Printf("%v is after %v, indexing %v\n", currModTime, prevModTime, path)
	newTokens := words.Tokenize(string(bs))
	oldTokens := db.PathTermIndex.GetTerms(path)
	fmt.Printf("new %v\n", newTokens)
	fmt.Printf("old %v\n", oldTokens)
	addToks, remToks := addsAndRems(newTokens, oldTokens)
	fmt.Printf("add %v\n", addToks)
	fmt.Printf("rem %v\n", remToks)
	// old terms - current terms = terms to delete
	// boost count for all add terms
	// drop count for all delete terms
	db.PathIndex.SetModTime(path, currModTime)
	return nil
}

func main() {
	var dirFlag = flag.String("dir", ".", "index current dir")
	var helpFlag = flag.Bool("help", false, "Show help")
	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}

	var db = db.NewMapDB()
	var walkFile func(string, fs.FileInfo, error) error
	walkFile = func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		indexFile(&db, path, info)
		return nil
	}
	filepath.Walk(*dirFlag, walkFile)
}
