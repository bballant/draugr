package db

import (
	"errors"
)

type MapIndex struct {
	indexInfo *IndexInfo
	termMap   map[string]*Term
}

func NewMapIndex() *MapIndex {
	return &MapIndex{
		indexInfo: &IndexInfo{},
		termMap:   make(map[string]*Term),
	}
}

func (mi *MapIndex) GetIndexInfo() *IndexInfo {
	return mi.indexInfo
}

func (mi *MapIndex) SetIndexInfo(info *IndexInfo) error {
	if info == nil {
		return errors.New("nil IndexInfo provided")
	}
	mi.indexInfo = info
	return nil
}

func (mi *MapIndex) SaveTerm(term string, path string) (*Term, error) {
	if term == "" || path == "" {
		return nil, errors.New("term or path cannot be empty")
	}

	t, ok := mi.termMap[term]
	if !ok {
		t = &Term{Token: term, Paths: []string{}}
		mi.termMap[term] = t
	}

	t.Paths = append(t.Paths, path)
	return t, nil
}

func (mi *MapIndex) RemoveTerm(term string, path string) (*Term, error) {
	if term == "" || path == "" {
		return nil, errors.New("term or path cannot be empty")
	}

	t, ok := mi.termMap[term]
	if !ok {
		return nil, errors.New("term not found")
	}

	for i, p := range t.Paths {
		if p == path {
			t.Paths = append(t.Paths[:i], t.Paths[i+1:]...)
			break
		}
	}

	if len(t.Paths) == 0 {
		delete(mi.termMap, term)
	}
	return t, nil
}

func (mi *MapIndex) GetTerm(term string) *Term {
	return mi.termMap[term]
}
