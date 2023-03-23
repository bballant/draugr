package db

import (
	"time"
)

type TermInfo struct {
	Term      string
	Count     int
	PathCount PathCount
}

var EmptyTerm = TermInfo{}

func (t *TermInfo) isEmpty() bool {
	return *t == EmptyTerm
}

type PathCount interface {
	GetCount(path string) int
	IncCount(path string) error
	SetCount(path string, count int) error
}

type PathIndex interface {
	SetModTime(path string, t time.Time) error
	GetModTime(path string) time.Time
}

type TermIndex interface {
	SaveTerm(term string, path string) (TermInfo, error)
	RemoveTerm(term string, path string) (TermInfo, error)
	GetTerm(term string) TermInfo
}

type PathTermIndex interface {
	SetTerms(path string, terms []string) error
	GetTerms(path string) []string
}

type DB struct {
	PathIndex     PathIndex
	TermIndex     TermIndex
	PathTermIndex PathTermIndex
}
