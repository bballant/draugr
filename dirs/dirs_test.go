package dirs

import (
	"fmt"
	"testing"
)

func TestDirs(t *testing.T) {

	fileInfos, err := FileInfo("../")
	for _, f := range fileInfos {
		fmt.Println(sPrintFileInfo(f))
	}

	if err != nil {
		t.Error(`unable to get fileinfo for .`)
	}

	if len(fileInfos) != 2 {
		t.Error(`There should be 2 files here`)
	}

}
