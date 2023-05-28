package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/gocolly/colly/v2"
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

func TestAimhDisclosure(t *testing.T) {
	var data TotalFinancialReports

	urlTemp := template.New("urlTemp")
	urlTemp, err := urlTemp.Parse(ALL_FINANCIAL_REPORTS_NOT_SEC)
	if err != nil {
		t.Error(err)
	}
	var url strings.Builder
	if err = urlTemp.Execute(&url, struct {
		Symbol   string
		PageNum  string
		PageSize string
	}{
		Symbol:   "AIMH",
		PageNum:  "1",
		PageSize: "10",
	}); err != nil {
		t.Error(err)
	}

	c := colly.NewCollector(colly.UserAgent(getUserAgent()))

	c.OnError(func(r *colly.Response, err error) {
		t.Error(err)
	})

	c.OnRequest(func(r *colly.Request) {
		if err := setHeaders(r, API_AUTHORITY, url.String()[len(`https://backend.otcmarkets.com`):]); err != nil {
			t.Error(err)
		}
		fmt.Println("request:", r.URL, "path:", r.Headers.Get("path"))
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("response:", r.StatusCode)

		file, err := os.Create("forbidden.html")
		if err != nil {
			t.Error(err)
		}
		defer file.Close()

		htmlData, err := io.ReadAll(file)
		if err != nil {
			t.Error(err)
		}

		fmt.Fprintln(file, string(htmlData))

		if err := json.Unmarshal(r.Body, &data); err != nil {
			t.Error(err)
		}
	})

	linkHits, err := os.Create("disclosure_links.txt")
	if err != nil {
		t.Error(err)
	}
	defer linkHits.Close()

	for _, r := range data.Records {
		fmt.Fprintln(linkHits, r)
	}

	if err := c.Visit(url.String()); err != nil {
		t.Error(err)
	}
}
