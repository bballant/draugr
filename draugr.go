package main

import (
	"flag"
	"fmt"
	"sort"

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
	var dirFlag = flag.String("dir", ".", "index dir")
	var searchFlag = flag.String("search", "", "search terms")
	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}
	if *searchFlag == "" {
		flag.Usage()
		return
	}

	var _db = db.NewMapDB()
	db.IndexPath(&_db, *dirFlag)
	res := Search(&_db, words.Tokenize(*searchFlag))
	for _, v := range res {
		fmt.Println(v.Path)
	}
}
