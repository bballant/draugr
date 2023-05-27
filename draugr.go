package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/bballant/draugr/db"
	"github.com/bballant/draugr/words"
)

const (
	socketPath = "/tmp/draugr.sock"
)

func init() {
	//log.SetOutput(ioutil.Discard)
}

func serveConnection(_index db.Index, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			return
		}
		res := db.SearchIndex(_index, words.Tokenize(string(buf[:n])))

		out := ""
		for _, v := range res {
			out = fmt.Sprintf("%s%s\n", out, v.Path)
		}

		_, err = conn.Write([]byte(out))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runServer(_index db.Index, dir string, exts []string) {
	db.IndexPathForExts(_index, dir, exts)
	os.Remove(socketPath)

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serveConnection(_index, conn)
	}
}

func runSearchClient(search string) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(search))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf[:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(buf[:n]))
}

func main() {
	var helpFlag = flag.Bool("help", false, "Show help")
	var debugFlag = flag.Bool("debug", false, "Print a lot of stuff")
	var serveFlag = flag.Bool("serve", false, "Run as service")
	var clientFlag = flag.Bool("client", false, "Run as client")
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

	if *serveFlag || (*searchFlag != "" && !*clientFlag) {
		var _index = db.NewMapIndex()
		var extensions = strings.Split(*extensionFilterFlag, " ")
		if *serveFlag {
			runServer(_index, *dirFlag, extensions)
		} else { // search index immediately
			db.IndexPathForExts(_index, *dirFlag, extensions)
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

	if *searchFlag != "" && *clientFlag {
		runSearchClient(*searchFlag)
	} else {
		fmt.Println("Invalid parameters")
		flag.Usage()
	}
}
