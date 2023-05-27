package db

import (
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bballant/draugr/words"
)

type Term struct {
	Token string
	Paths []string
}

type IndexInfo struct {
	Paths []string
}

type Index interface {
	GetIndexInfo() *IndexInfo
	SetIndexInfo(*IndexInfo) error
	AllTerms() []*Term
	SaveTerm(token string, path string) (*Term, error)
	RemoveTerm(token string, path string) (*Term, error)
	GetTerm(token string) *Term
}

func IndexPathForExts(index Index, path string, extensions []string) error {
	var walkFile = func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		// node modules are the worst
		if strings.Contains(path, "node_modules") ||
			strings.Contains(path, "gogol") ||
			strings.Contains(path, "amazonka") {
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

type SearchResult struct {
	Path  string
	Count int
}

func SearchIndex(index Index, tokens []string) []SearchResult {
	pathTotals := map[string]int{}
	for _, token := range tokens {
		term := index.GetTerm(token)
		if term == nil {
			continue
		}
		for _, path := range unique(term.Paths) {
			if _, ok := pathTotals[path]; !ok {
				pathTotals[path] = 0
			}
			basicScore := BasicScore(*index.GetIndexInfo(), *term, path)
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

func BasicScore(indexInfo IndexInfo, term Term, path string) int {
	tf := 0
	for _, s := range term.Paths {
		if s == path {
			tf++
		}
	}
	pathScore := 0
	dirs := words.Tokenize(path)
	for _, dir := range dirs {
		if dir == term.Token {
			pathScore = pathScore + 2
		}
	}
	return tf + pathScore
}

func TFIDFScore(indexInfo IndexInfo, term Term, path string) float64 {
	numDocuments := len(indexInfo.Paths)
	tf := 0
	for _, s := range term.Paths {
		if s == path {
			tf++
		}
	}
	df := len(unique(term.Paths))
	idf := math.Log(float64(1+numDocuments) / float64(1+df))
	tfIdf := float64(tf) * idf
	return tfIdf
}

func unique(strs []string) []string {
	strMap := make(map[string]struct{})
	uniqueStrings := []string{}
	for _, s := range strs {
		_, ok := strMap[s]
		if !ok {
			strMap[s] = struct{}{}
			uniqueStrings = append(uniqueStrings, s)
		}
	}
	return uniqueStrings
}

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
