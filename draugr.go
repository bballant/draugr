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

func IndexPath(db *db.DB, path string) error {
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

type SearchResult struct {
	Path  string
	Count int
}

func Search(db_ *db.DB, terms []string) []SearchResult {
	pathTotals := map[string]int{}
	for _, term := range terms {
		inf := db_.TermIndex.GetTerm(term)
		if inf != nil {
			for _, path := range inf.PathCount.GetPaths() {
				if _, ok := pathTotals[path]; !ok {
					pathTotals[path] = 0
				}
				pathTotals[path] += inf.PathCount.GetCount(path)
			}
		}
	}
	results := make([]SearchResult, len(pathTotals))
	i := 0
	for k, v := range pathTotals {
		results[i] = SearchResult{k, v}
		i++
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})
	return results
}

func main() {
	var helpFlag = flag.Bool("help", false, "Show help")
	var dirFlag = flag.String("dir", ".", "index current dir")
	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}

	var db = db.NewMapDB()
	IndexPath(&db, *dirFlag)
	fmt.Printf("%v\n", db.TermIndex.GetTerm("files"))
}
