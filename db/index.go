package db

import (
	"math"

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

func Unique(strs []string) []string {
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
	df := len(Unique(term.Paths))
	idf := math.Log(float64(1+numDocuments) / float64(1+df))
	tfIdf := float64(tf) * idf
	return tfIdf
}
