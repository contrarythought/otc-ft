package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
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

// TODO
func jsonConverter(in string) ([]byte, error) {
	var out strings.Builder
	for _, c := range in {
		if c != '\\' {
			if _, err := out.WriteRune(c); err != nil {
				return nil, err
			}
		}
	}

	return []byte(out.String()[1 : len(out.String())-2]), nil
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
	r.Headers.Set(`scheme`, `https`)
	r.Headers.Set(`accept`, `*/*`)
	r.Headers.Set(`accept-encoding`, `gzip, deflate, br`)
	r.Headers.Set(`accept-language`, `en,en-US;q=0.9,zh-TW;q=0.8,zh;q=0.7`)
	r.Headers.Set(`referer`, `https://www.otcmarkets.com/research/stock-screener`)
	r.Headers.Set(`sec-ch-ua`, `"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`)
	r.Headers.Set(`sec-ch-ua-mobile`, `?0`)
	r.Headers.Set(`sec-ch-ua-platform`, `"Windows"`)
	r.Headers.Set(`sec-fetch-dest`, `empty`)
	r.Headers.Set(`sec-fetch-mode`, `cors`)
	r.Headers.Set(`sec-fetch-site`, `same-origin`)
	r.Headers.Set(`x-requested-with`, `XMLHttpRequest`)
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

type PageData struct {
	Count  int `json:"count"`
	Pages  int `json:"pages"`
	Stocks []struct {
		SecurityID           int     `json:"securityId"`
		ReportDate           string  `json:"reportDate"`
		Symbol               string  `json:"symbol"`
		SecurityName         string  `json:"securityName"`
		Market               string  `json:"market"`
		MarketID             int     `json:"marketId"`
		SecurityType         string  `json:"securityType"`
		Country              string  `json:"country,omitempty"`
		State                string  `json:"state,omitempty"`
		ForexCountry         string  `json:"forexCountry,omitempty"`
		CaveatEmptor         bool    `json:"caveatEmptor"`
		IndustryID           int     `json:"industryId,omitempty"`
		Industry             string  `json:"industry,omitempty"`
		Volume               int     `json:"volume"`
		VolumeChange         float64 `json:"volumeChange,omitempty"`
		DividendYield        float64 `json:"dividendYield"`
		DividendPayer        bool    `json:"dividendPayer"`
		Penny                bool    `json:"penny"`
		Price                float64 `json:"price"`
		ShortInterestRatio   float64 `json:"shortInterestRatio"`
		IsBank               string  `json:"isBank"`
		ShortInterest        int     `json:"shortInterest,omitempty"`
		ShortInterestPercent float64 `json:"shortInterestPercent,omitempty"`
		Pct1Day              float64 `json:"pct1Day,omitempty"`
		Pct5Day              float64 `json:"pct5Day,omitempty"`
		Pct4Weeks            float64 `json:"pct4Weeks,omitempty"`
		Pct13Weeks           float64 `json:"pct13Weeks,omitempty"`
		Pct52Weeks           float64 `json:"pct52Weeks,omitempty"`
		PerfQxComp4Weeks     float64 `json:"perfQxComp4Weeks,omitempty"`
		PerfQxComp13Weeks    float64 `json:"perfQxComp13Weeks,omitempty"`
		PerfQxComp52Weeks    float64 `json:"perfQxComp52Weeks,omitempty"`
		PerfQxBillion4Weeks  float64 `json:"perfQxBillion4Weeks,omitempty"`
		PerfQxBillion13Weeks float64 `json:"perfQxBillion13Weeks,omitempty"`
		PerfQxBillion52Weeks float64 `json:"perfQxBillion52Weeks,omitempty"`
		PerfQxBanks4Weeks    float64 `json:"perfQxBanks4Weeks,omitempty"`
		PerfQxBanks13Weeks   float64 `json:"perfQxBanks13Weeks,omitempty"`
		PerfQxBanks52Weeks   float64 `json:"perfQxBanks52Weeks,omitempty"`
		PerfQxIntl4Weeks     float64 `json:"perfQxIntl4Weeks,omitempty"`
		PerfQxIntl13Weeks    float64 `json:"perfQxIntl13Weeks,omitempty"`
		PerfQxIntl52Weeks    float64 `json:"perfQxIntl52Weeks,omitempty"`
		PerfQxUs4Weeks       float64 `json:"perfQxUs4Weeks,omitempty"`
		PerfQxUs13Weeks      float64 `json:"perfQxUs13Weeks,omitempty"`
		PerfQxUs52Weeks      float64 `json:"perfQxUs52Weeks,omitempty"`
		PerfQb4Weeks         float64 `json:"perfQb4Weeks,omitempty"`
		PerfQb13Weeks        float64 `json:"perfQb13Weeks,omitempty"`
		PerfQb52Weeks        float64 `json:"perfQb52Weeks,omitempty"`
		PerfSp4Weeks         float64 `json:"perfSp4Weeks,omitempty"`
		PerfSp13Weeks        float64 `json:"perfSp13Weeks,omitempty"`
		PerfSp52Weeks        float64 `json:"perfSp52Weeks,omitempty"`
		PerfQxDiv4Weeks      float64 `json:"perfQxDiv4Weeks,omitempty"`
		PerfQxDiv13Weeks     float64 `json:"perfQxDiv13Weeks,omitempty"`
		PerfQxDiv52Weeks     float64 `json:"perfQxDiv52Weeks,omitempty"`
		PerfQxCan4Weeks      float64 `json:"perfQxCan4Weeks,omitempty"`
		PerfQxCan13Weeks     float64 `json:"perfQxCan13Weeks,omitempty"`
		PerfQxCan52Weeks     float64 `json:"perfQxCan52Weeks,omitempty"`
		MorningStarRating    int     `json:"morningStarRating,omitempty"`
	} `json:"stocks"`
}

func getPageData(pageNum, pageSize int) (*PageData, error) {
	var pageData PageData
	var err error = nil
	urlTemp := template.New("urlTemp")
	urlTemp, err = urlTemp.Parse(SCREENER_URL)
	if err != nil {
		return nil, err
	}

	var url strings.Builder
	if err = urlTemp.Execute(&url, struct {
		PageNum  int
		PageSize int
	}{
		PageNum:  pageNum,
		PageSize: pageSize,
	}); err != nil {
		return nil, err
	}

	c := colly.NewCollector(colly.UserAgent(getUserAgent()))

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("err: status code ->", r.StatusCode, "msg ->", err)
	})

	c.OnRequest(func(r *colly.Request) {
		err = setHeaders(r, pageNum, pageSize)
		fmt.Println("request:", r.URL)
	})

	var jsonData []byte
	c.OnResponse(func(r *colly.Response) {
		jsonData, err = jsonConverter(string(r.Body))
		err = json.Unmarshal(jsonData, &pageData)
	})

	if err = c.Visit(url.String()); err != nil {
		return nil, err
	}

	return &pageData, err
}

const (
	NUM_WORKERS = 10
)

// TODO:
// 1. iterate through each page in the screener and grab the JSON data of the page
// 2. iterate through each stock's details site and download pdf filing and press release info

func Scrape() error {
	pageData, err := getPageData(0, 100)
	if err != nil {
		return err
	}

	for i := 0; i < NUM_WORKERS; i++ {

	}

	for i := 0; i < pageData.Pages; i++ {
		if (i+1)%10 == 0 {
			time.Sleep(time.Duration(rand.Intn(4) + 3))
		}
	}

	return nil
}
