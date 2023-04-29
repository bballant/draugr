package main

import (
	"math"

	"github.com/bballant/draugr/db"
)

func termFrequency(termCount int, allTermsCount int) float64 {
	return float64(termCount) / float64(allTermsCount)
}

func TFIDF(N int, tf int, df int) float64 {
	idf := math.Log(float64(1+N) / float64(1+df))
	tfIdf := float64(tf) * idf
	return tfIdf
}

func calcTFIDF(term string, supah db.Index) float64 {
	termInfo := supah.GetTerm(term)
	paths := termInfo.Paths
	//for i, path := range paths {

	//}
	return float64(len(paths))
}
