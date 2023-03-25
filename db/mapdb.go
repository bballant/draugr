package db

import "time"

// PathCount Interface

type MapPathCount map[string]int

func (m MapPathCount) GetPaths() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func (m MapPathCount) GetCount(path string) int {
	return m[path]
}

func (m MapPathCount) SetCount(path string, count int) error {
	m[path] = count
	return nil
}

func (m MapPathCount) IncCount(path string) error {
	m[path]++
	return nil
}

// PathIndex Interface

type MapPathIndex map[string]time.Time

func (m MapPathIndex) SetModTime(path string, t time.Time) error {
	m[path] = t
	return nil
}

func (m MapPathIndex) GetModTime(path string) time.Time {
	return m[path]
}

// TermIndex Interface

type MapTermIndex map[string]*TermInfo

func (m MapTermIndex) GetTerm(term string) *TermInfo {
	return m[term]
}

func (m MapTermIndex) SaveTerm(term string, path string) (*TermInfo, error) {
	termInfo := m.GetTerm(term)
	if termInfo == nil {
		termInfo = &TermInfo{term, 0, MapPathCount{}}
	}
	termInfo.PathCount.IncCount(path)
	termInfo.Count++
	m[term] = termInfo
	return termInfo, nil
}

func (m MapTermIndex) RemoveTerm(term string, path string) (*TermInfo, error) {
	termInfo := m.GetTerm(term)
	if termInfo == nil {
		return termInfo, nil
	}
	currCount := termInfo.PathCount.GetCount(path)
	if currCount <= 0 {
		return termInfo, nil
	}
	termInfo.PathCount.SetCount(path, currCount-1)
	termInfo.Count--
	m[term] = termInfo
	return termInfo, nil
}

// PathTermIndex Interface

type MapPathTermIndex map[string][]string

func (m MapPathTermIndex) GetTerms(path string) []string {
	return m[path]
}

func (m MapPathTermIndex) SetTerms(path string, terms []string) error {
	m[path] = terms
	return nil
}

func NewMapDB() DB {
	pathIndex := MapPathIndex{}
	termIndex := MapTermIndex{}
	pathTermIndex := MapPathTermIndex{}
	return DB{pathIndex, termIndex, pathTermIndex}
}
