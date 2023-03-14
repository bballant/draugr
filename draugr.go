package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/bballant/draugr/db"
)

func indexFile(db *db.DB, path string, info fs.FileInfo) error {
	prevModTime := db.PathIndex.GetModTime(path)
	currModTime := info.ModTime()
	if !currModTime.After(prevModTime) {
		log.Printf("%v is not after %v, not indexing %v\n", currModTime, prevModTime, path)
		return nil
	}
	db.PathIndex.SetModTime(path, currModTime)
	log.Printf("%v is after %v, indexing %v\n", currModTime, prevModTime, path)
	return nil
}

func main() {
	var dirFlag = flag.String("dir", ".", "index current dir")
	var helpFlag = flag.Bool("help", false, "Show help")
	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}

	var db = db.NewMapDB()
	var walkFile func(string, fs.FileInfo, error) error
	walkFile = func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		indexFile(&db, path, info)
		return nil
	}
	filepath.Walk(*dirFlag, walkFile)
	filepath.Walk(*dirFlag, walkFile)
	fmt.Println(db.PathIndex)
}
