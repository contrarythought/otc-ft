package scraper

import (
	"fmt"
	"testing"
)

func TestMaxPage(t *testing.T) {
	max, err := getMaxPage()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("max page:", max)
}
