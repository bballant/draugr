package main

import (
	"math"
	"time"
)

func termFrequency(termCount int, allTermsCount int) float64 {
	return float64(termCount) / float64(allTermsCount)
}

func TFIDF(N int, tf int, df int) float64 {
	idf := math.Log(float64(1+N) / float64(1+df))
	tfIdf := float64(tf) * idf
	return tfIdf
}

type TermInfo struct {
	Term  string
	Count int
	Paths []string
}

type DocInfo struct {
	Path    string
	Terms   []string
	ModTime time.Time
}

type IndexInfo struct {
	DocCount int
}

type Supah interface {
	GetIndexInfo() *IndexInfo
	SetIndexInfo(*IndexInfo) error
	GetDocInfo(path string) *DocInfo
	SetDocInfo(path string, docInfo *DocInfo) error
	SaveTerm(term string, path string) (*TermInfo, error)
	RemoveTerm(term string, path string) (*TermInfo, error)
	GetTerm(term string) *TermInfo
}

func calcTFIDF(term string, supah Supah) float64 {
	termInfo := supah.GetTerm(term)
	paths := termInfo.Paths
	//for i, path := range paths {

	//}
	return float64(len(paths))
}
