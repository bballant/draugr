package main

import (
	"flag"
	"fmt"
	"io"
	"log"
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

func main() {
	var helpFlag = flag.Bool("help", false, "Show help")
	var debugFlag = flag.Bool("debug", false, "Print a lot of stuff")
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

	if !*debugFlag {
		log.SetOutput(io.Discard)
	}

	if *searchFlag != "" {
		var _index = db.NewMapIndex()
		db.IndexPathForExts(_index, *dirFlag, strings.Split(*extensionFilterFlag, " "))
		res := SearchIndex(_index, words.Tokenize(*searchFlag))
		for _, v := range res {
			if *debugFlag {
				fmt.Println(v)
			} else {
				fmt.Println(v.Path)
			}
		}
	}

}
