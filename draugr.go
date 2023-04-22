package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/bballant/draugr/db"
	"github.com/bballant/draugr/words"
)

func init() {
	//log.SetOutput(ioutil.Discard)
}

type SearchResult struct {
	Path  string
	Count int
}

func Search(db_ *db.DB, terms []string) []SearchResult {
	pathTotals := map[string]int{}
	for _, term := range terms {
		inf := db_.TermIndex.GetTerm(term)
		if inf != nil {
			for _, path := range inf.PathCount.GetPaths() {
				if _, ok := pathTotals[path]; !ok {
					pathTotals[path] = 0
				}
				pathTotals[path] += inf.PathCount.GetCount(path)
				// hack in filename match boost
				_, file := filepath.Split(path)
				fileMatches := words.Occurances(term, file)
				pathTotals[path] += 5 * fileMatches
			}
		}
	}
	results := make([]SearchResult, len(pathTotals))
	i := 0
	for k, v := range pathTotals {
		results[i] = SearchResult{k, v}
		i++
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})
	return results
}

func pathEnd(path string, lastN int) string {
	if len(path) <= lastN {
		return path
	}
	res := []rune(path[len(path)-lastN:])
	res[0] = '>'
	return string(res)
}

func openTerminalAt(path string) error {
	dirStr, _ := filepath.Split(path)
	cmd := exec.Command("gnome-terminal", "--working-directory="+dirStr)
	return cmd.Run()
}

func showi3Icon() {
	// create a new i3bar output
	bar := i3.NewOutput()

	// create a simple text block
	helloBlock := i3.NewBlock("Hello, World!")
	bar.AddBlock(helloBlock)

	// create an icon block
	iconBlock := i3.NewBlock("")
	iconBlock.FullText = "" // Unicode icon code for a gear icon
	iconBlock.Color = "#ffffff"
	bar.AddBlock(iconBlock)

	// start the i3bar output loop
	bar.Start()

	// update the icon block every second
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		// set the icon block's full text to the current time
		iconBlock.FullText = time.Now().Format("15:04:05")

		// marshal the i3bar output to JSON and write it to stdout
		jsonData, err := json.Marshal(bar)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			continue
		}
		fmt.Println(string(jsonData))
	}
}

func main() {
	var helpFlag = flag.Bool("help", false, "Show help")
	var dirFlag = flag.String("dir", ".", "index dir")
	var searchFlag = flag.String("search", "", "search terms")
	var extensionFilterFlag = flag.String(
		"exts", ".txt .md .scala .go .hs",
		"file extensions to filter for")
	flag.Parse()
	if *helpFlag {
		flag.Usage()
		return
	}
	var _db = db.NewMapDB()
	db.IndexPathForExts(&_db, *dirFlag, strings.Split(*extensionFilterFlag, " "))

	if *searchFlag != "" {
		res := Search(&_db, words.Tokenize(*searchFlag))
		for _, v := range res {
			fmt.Println(v.Path)
		}
		return
	}

	myApp := app.New()
	myWin := myApp.NewWindow("Search App")
	input := widget.NewEntry()
	items := []string{}
	list := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(items[id])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		openTerminalAt(items[id])
	}
	doSearch := func(term string) {
		res_ := Search(&_db, words.Tokenize(term))
		resText := ""
		max := len(res_)
		if max > 10 {
			max = 10
		}
		items = []string{}
		for _, v := range res_[:max] {
			// TODO - separate display from data
			cleanPath := pathEnd(v.Path, 1000)
			items = append(items, fmt.Sprintf("%s%s (%d)", resText, cleanPath, v.Count))
		}
		list.Refresh()
	}
	input.OnSubmitted = doSearch
	searchButton := widget.NewButton("Search", func() { doSearch(input.Text) })
	listContainer := container.NewScroll(list)
	listContainer.SetMinSize(fyne.NewSize(200, 300))
	content := container.NewVBox(
		input,
		searchButton,
		listContainer,
	)
	myWin.SetContent(content)
	myWin.ShowAndRun()
}
