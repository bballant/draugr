package dirs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func sPrintFileInfo(f fs.FileInfo) string {
	return fmt.Sprintf("%s %v %v\n", f.Name(), f.ModTime(), f.IsDir())
}

func walkFile(path string, info fs.FileInfo, err error) error {
	fmt.Printf("%s: %v \n", path, info.ModTime())
	return nil
}

func FileInfo(dir string) ([]fs.FileInfo, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	filepath.Walk(dir, walkFile)

	return fileInfo, nil
}
