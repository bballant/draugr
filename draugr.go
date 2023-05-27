package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bballant/draugr/db"
	"github.com/bballant/draugr/words"
)

func init() {
	//log.SetOutput(ioutil.Discard)
}

type SearchResult struct {
	Path  string
	Count int
}

func SearchIndex(index db.Index, tokens []string) []SearchResult {
	pathTotals := map[string]int{}
	for _, token := range tokens {
		term := index.GetTerm(token)
		if term == nil {
			continue
		}
		for _, path := range db.Unique(term.Paths) {
			if _, ok := pathTotals[path]; !ok {
				pathTotals[path] = 0
			}
			basicScore := db.BasicScore(*index.GetIndexInfo(), *term, path)
			pathTotals[path] += basicScore
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
				// hack in filename match boost
				_, file := filepath.Split(path)
				fileMatches := words.Occurances(term, file)
				pathTotals[path] += 5 * fileMatches
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

func pathEnd(path string, lastN int) string {
	if len(path) <= lastN {
		return path
	}
	res := []rune(path[len(path)-lastN:])
	res[0] = '>'
	return string(res)
}

func main() {
	var helpFlag = flag.Bool("help", false, "Show help")
	var dirFlag = flag.String("dir", ".", "index dir")
	var searchFlag = flag.String("search", "", "search terms")
	var extensionFilterFlag = flag.String(
		"exts", ".txt .md .scala .go .hs .ts",
		"file extensions to filter for")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	if *searchFlag != "" {
		var _index = db.NewMapIndex()
		db.IndexIndexPathForExts(_index, *dirFlag, strings.Split(*extensionFilterFlag, " "))
		res := SearchIndex(_index, words.Tokenize(*searchFlag))
		for _, v := range res {
			fmt.Println(v)
		}
		return
	}

	var _db = db.NewMapDB()
	//db.IndexPathForExts(&_db, *dirFlag, strings.Split(*extensionFilterFlag, " "))
	var items []string

	_ /*doSearch*/ = func(term string) {
		res_ := Search(&_db, words.Tokenize(term))
		resText := ""
		max := len(res_)
		if max > 10 {
			max = 10
		}
		items = []string{}
		for _, v := range res_[:max] {
			// TODO - separate display from data
			cleanPath := pathEnd(v.Path, 1000)
			items = append(items, fmt.Sprintf("%s%s (%d)", resText, cleanPath, v.Count))
		}
	}

}
