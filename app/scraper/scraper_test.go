package scraper

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestMaxPage(t *testing.T) {
	pgdata, err := getPageData(0, 100)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("max page:", pgdata.Pages)
}

func TestJsonConverter(t *testing.T) {
	file, err := os.Open(`test2.txt`)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	byteData, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	data, err := jsonConverter(string(byteData))
	if err != nil {
		t.Error(err)
	}

	outFile, err := os.Create(`outfileTest.txt`)
	if err != nil {
		t.Error(err)
	}
	defer outFile.Close()

	if _, err = fmt.Fprintln(outFile, string(data)); err != nil {
		t.Error(err)
	}
}
