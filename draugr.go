package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/bballant/draugr/db"
	"github.com/bballant/draugr/words"
)

func init() {
	//log.SetOutput(ioutil.Discard)
}

func main() {
	var helpFlag = flag.Bool("help", false, "Show help")
	var debugFlag = flag.Bool("debug", false, "Print a lot of stuff")
	var dirFlag = flag.String("dir", ".", "index dir")
	var searchFlag = flag.String("search", "", "search terms")
	var extensionFilterFlag = flag.String(
		"exts", ".txt .md .scala .go .hs .ts",
		"file extensions to filter for")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	if !*debugFlag {
		log.SetOutput(io.Discard)
	}

	if *searchFlag != "" {
		var _index = db.NewMapIndex()
		db.IndexPathForExts(_index, *dirFlag, strings.Split(*extensionFilterFlag, " "))
		res := db.SearchIndex(_index, words.Tokenize(*searchFlag))
		for _, v := range res {
			if *debugFlag {
				fmt.Println(v)
			} else {
				fmt.Println(v.Path)
			}
		}
	}

}
