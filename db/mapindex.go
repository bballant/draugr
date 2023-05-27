package db

import (
	"errors"
	"sync"
)

var mmu sync.Mutex
var imu sync.Mutex

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
	imu.Lock()
	mi.indexInfo = info
	imu.Unlock()
	return nil
}

func (mi *MapIndex) AddPath(path string) *IndexInfo {
	imu.Lock()
	mi.indexInfo.Paths = append(mi.indexInfo.Paths, path)
	imu.Unlock()
	return mi.indexInfo
}

func (mi *MapIndex) SaveTerm(term string, path string) (*Term, error) {
	if term == "" || path == "" {
		return nil, errors.New("term or path cannot be empty")
	}

	mmu.Lock()
	t, ok := mi.termMap[term]
	if !ok {
		t = &Term{Token: term, Paths: []string{}}
		mi.termMap[term] = t
	}

	t.Paths = append(t.Paths, path)
	mmu.Unlock()
	return t, nil
}

func (mi *MapIndex) RemoveTerm(term string, path string) (*Term, error) {
	if term == "" || path == "" {
		return nil, errors.New("term or path cannot be empty")
	}

	mmu.Lock()
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
	mmu.Unlock()
	return t, nil
}

func (mi *MapIndex) GetTerm(term string) *Term {
	return mi.termMap[term]
}
