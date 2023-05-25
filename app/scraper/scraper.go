package scraper

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gocolly/colly/v2"
)

func getUserAgent() string {
	userAgents := []string{
		`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36`,
		`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36`,
		`Mozilla/5.0 (Macintosh; Intel Mac OS X 13_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15`,
		`Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36`,
		`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36`,
	}
	idx := rand.Intn(len(userAgents))
	return userAgents[idx]
}

func downloadFile(path string, body []byte) error {
	if strings.Index(path, `filings\`) == -1 {
		path = `filings\` + path
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err = io.Copy(outFile, bytes.NewReader(body)); err != nil {
		return err
	}

	return nil
}

const (
	AUTHORITY = `www.otcmarkets.com`
)

func setHeaders(r *colly.Request, pageNum, pageSize int) error {
	pathTemp := template.New("pathTemp")
	pathTemp, err := pathTemp.Parse(PATH)
	if err != nil {
		return err
	}
	var path strings.Builder
	if err = pathTemp.Execute(&path, struct {
		PageNum  int
		PageSize int
	}{
		PageNum:  pageNum,
		PageSize: pageSize,
	}); err != nil {
		return err
	}

	r.Headers.Set(`authority`, AUTHORITY)
	r.Headers.Set(`method`, r.Method)
	r.Headers.Set(`path`, path.String())
	
	return nil
}

const (
	SCREENER_URL = `https://www.otcmarkets.com/research/stock-screener/api?page={{.PageNum}}&pageSize={{.PageSize}}`
	PATH         = `/research/stock-screener/api?page={{.PageNum}}&pageSize={{.PageSize}}`
)

type PressReleaseData struct {
	Content     []string  `json:"content"`
	Date        time.Time `json:"date"`
	CompanyName string    `json:"companyName"`
	Ticker      string    `json:"ticker"`
}

type FilingData struct {
	Content     []byte    `json:"content"`
	Date        time.Time `json:"date"`
	CompanyName string    `json:"companyName"`
	Ticker      string    `json:"ticker"`
}

func getMaxPage() (int, error) {
	var maxPages int
	var err error
	urlTemp := template.New("urlTemp")
	urlTemp, err = urlTemp.Parse(SCREENER_URL)
	if err != nil {
		return -1, err
	}

	var url strings.Builder
	if err = urlTemp.Execute(&url, struct {
		PageNum  int
		PageSize int
	}{
		PageNum:  0,
		PageSize: 100,
	}); err != nil {
		return -1, err
	}

	c := colly.NewCollector(colly.UserAgent(getUserAgent()))

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("err: status code ->", r.StatusCode, "msg ->", err)
	})

	c.OnHTML(`#pagination > ul > li:nth-child(9) > a`, func(h *colly.HTMLElement) {
		maxPages, err = strconv.Atoi(h.Text)
		if err != nil {
			log.Fatal(err)
		}
	})

	if err = c.Visit(url.String()); err != nil {
		return -1, err
	}

	return maxPages, nil
}

// TODO: visit each company link in the stock screener
// download each pdf filling and store in server\filings
// visit each press release link, download the text and store in server\press_releases
func Scrape() error {
	maxPage, err := getMaxPage()
	if err != nil {
		return err
	}
	return nil
}