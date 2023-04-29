package db

import (
	"math"
	"os"
	"strings"
)

type Term struct {
	Token string
	Count int
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

func BasicScore(indexInfo IndexInfo, term Term, path string) int {
	tf := 0
	for _, s := range term.Paths {
		if s == path {
			tf++
		}
	}
	pathScore := 0
	for _, path := range term.Paths {
		dirs := strings.Split(path, string(os.PathSeparator))
		for _, dir := range dirs {
			if dir == term.Token {
				pathScore = pathScore + 2
			}
		}
	}
	return tf + pathScore
}

func countUnique(stringsSlice []string) int {
	uniqueStrings := make(map[string]struct{})

	for _, s := range stringsSlice {
		uniqueStrings[s] = struct{}{}
	}

	return len(uniqueStrings)
}

func TFIDFScore(indexInfo IndexInfo, term Term, path string) float64 {
	numDocuments := len(indexInfo.Paths)
	tf := 0
	for _, s := range term.Paths {
		if s == path {
			tf++
		}
	}
	df := countUnique(term.Paths)
	idf := math.Log(float64(1+numDocuments) / float64(1+df))
	tfIdf := float64(tf) * idf
	return tfIdf
}
