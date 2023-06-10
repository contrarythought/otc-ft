package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
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
	BASE_AUTHORITY = `www.otcmarkets.com`
	API_AUTHORITY  = `backend.otcmarkets.com`
)

func setHeaders(r *colly.Request, authority, path string) {
	r.Headers.Set(`authority`, authority)
	r.Headers.Set(`method`, r.Method)
	r.Headers.Set(`path`, path)
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
}

const (
	SCREENER_URL = `https://www.otcmarkets.com/research/stock-screener/api?page={{.PageNum}}&pageSize={{.PageSize}}`
	PATH         = `/research/stock-screener/api?page={{.PageNum}}&pageSize={{.PageSize}}`
)

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

func getPageData(pageNum, pageSize int, logger *log.Logger) (*PageData, error) {
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
		logger.Println(err)
		fmt.Println("err: status code ->", r.StatusCode, "msg ->", err)
	})

	c.OnRequest(func(r *colly.Request) {
		setHeaders(r, BASE_AUTHORITY, url.String()[len(`https://www.otcmarkets.com`):])
		fmt.Println("request:", r.URL)
	})

	var jsonData []byte
	c.OnResponse(func(r *colly.Response) {
		jsonData, err = jsonConverter(string(r.Body))
		if err = json.Unmarshal(jsonData, &pageData); err != nil {
			logger.Println(err)
		}
	})

	if err = c.Visit(url.String()); err != nil {
		logger.Println(err)
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

func Scrape(errLog *os.File) error {
	logger := log.New(errLog, "whenlambo", log.LstdFlags|log.Lshortfile)

	pageData, err := getPageData(0, 100, logger)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pageChan := make(chan int, 10)
	var wg sync.WaitGroup

	for i := 0; i < NUM_WORKERS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.Println(r)
				}
			}()
			for {
				select {
				case page := <-pageChan:
					if err := scrapePage(page, logger); err != nil {
						panic(err)
					}
				case <-ctx.Done():
					if len(pageChan) == 0 {
						return
					}
				}
			}
		}()
	}

	for i := 1; i < pageData.Pages; i++ {
		if (i+1)%10 == 0 {
			time.Sleep(time.Duration(rand.Intn(4) + 3))
		}

		pageChan <- i
	}

	close(pageChan)
	cancel()
	wg.Wait()

	return nil
}

func scrapePage(page int, logger *log.Logger) error {
	pageData, err := getPageData(page, 100, logger)
	if err != nil {
		logger.Println(err)
		return err
	}

	// loop through each stock in the page
	for _, stock := range pageData.Stocks {
		if err = scrapeReports(stock.Symbol, logger); err != nil {
			logger.Println(err)
			return fmt.Errorf("error scrape report %s: %s", stock.Symbol, err)
		}
		if err = scrapeNews(stock.Symbol, logger); err != nil {
			logger.Println(err)
			return fmt.Errorf("error scrape news %s: %s", stock.Symbol, err)
		}
	}

	return nil
}

const (
	ALL_NEWS_URL                  = `https://backend.otcmarkets.com/otcapi/company/{{.Symbol}}/dns/news?symbol={{.Symbol}}&page={{.PageNum}}&pageSize={{.PageSize}}&sortOn=releaseDate&sortDir=DESC`
	NEWS_URL                      = `https://www.otcmarkets.com/stock/{{.Symbol}}/news/{{.Title}}?id={{.ID}}`
	ALL_SEC_FILINGS               = `https://backend.otcmarkets.com/otcapi/company/sec-filings/AIMH?symbol=AIMH&page=1&pageSize=10`
	EXAMPLE_SEC_FILING            = `https://www.otcmarkets.com/filing/html?id=14305340&guid=2UT-kn10eYd-B3h`
	ALL_FINANCIAL_REPORTS_NOT_SEC = `https://backend.otcmarkets.com/otcapi/company/{{.Symbol}}/financial-report?symbol={{.Symbol}}&page={{.PageNum}}&pageSize={{.PageSize}}&statusId=A&sortOn=releaseDate&sortDir=DESC`
	FIN_REPORT_URL                = `https://www.otcmarkets.com/otcapi/company/financial-report/{{.ID}}/content`
	MAX_PAGE_SIZE                 = "50"
)

func setHeadersAPI(r *colly.Request) {
	r.Headers.Set(`authority`, API_AUTHORITY)
	r.Headers.Set(`method`, r.Method)
	r.Headers.Set(`path`, r.URL.String()[len(API_AUTHORITY):])
	r.Headers.Set(`scheme`, `https`)
	r.Headers.Set(`accept`, `application/json, text/plain, */*`)
	r.Headers.Set(`Accept-Encoding`, `gzip, deflate, br`)
	r.Headers.Set(`Accept-Language`, `en,en-US;q=0.9,zh-TW;q=0.8,zh;q=0.7`)
	r.Headers.Set(`origin`, `https://www.otcmarkets.com`)
	r.Headers.Set(`referer`, `https://www.otcmarkets.com/`)
	r.Headers.Set(`sec-ch-ua`, `"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`)
	r.Headers.Set(`Sec-Ch-Ua-Mobile`, `?0`)
	r.Headers.Set(`Sec-Ch-Ua-Platform`, `"Windows"`)
	r.Headers.Set(`Sec-Fetch-Dest`, `empty`)
	r.Headers.Set(`Sec-Fetch-Mode`, `cors`)
	r.Headers.Set(`Sec-Fetch-Site`, `same-site`)
}

// downloads reports and puts them in server directory
func scrapeReports(symbol string, logger *log.Logger) error {
	// reports will be unmarshaled into data
	var data TotalFinancialReports

	// forming the initial url to call the API with to obtain total pages
	urlTemp := template.New("urlTemp")
	urlTemp, err := urlTemp.Parse(ALL_FINANCIAL_REPORTS_NOT_SEC)
	if err != nil {
		return err
	}
	var url strings.Builder
	if err = urlTemp.Execute(&url, struct {
		Symbol   string
		PageNum  string
		PageSize string
	}{
		Symbol:   symbol,
		PageNum:  "1",
		PageSize: MAX_PAGE_SIZE,
	}); err != nil {
		return err
	}

	// find out how many pages of records there are
	// I will iterate across each page (MAX_PAGE_SIZE records/page)
	totalPages, err := getTotalRecordPages(url.String())
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	// iterate across each page
	for i := 1; i <= totalPages; i++ {
		url.Reset()
		if err = urlTemp.Execute(&url, struct {
			Symbol   string
			PageNum  string
			PageSize string
		}{
			Symbol:   symbol,
			PageNum:  strconv.Itoa(i),
			PageSize: MAX_PAGE_SIZE,
		}); err != nil {
			return err
		}

		c := colly.NewCollector(colly.UserAgent(getUserAgent()))

		c.OnRequest(func(r *colly.Request) {
			setHeadersAPI(r)
		})

		c.OnResponse(func(r *colly.Response) {
			if err = json.Unmarshal(r.Body, &data); err != nil {
				logger.Println(err)
			}
		})

		if err = c.Visit(url.String()); err != nil {
			logger.Println(err)
			return err
		}

		// download each report and put them in server directory
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		for _, r := range data.Records {
			fmt.Println(r.ID, symbol, r.TypeID)
			if err = downloadRecord(r.ID, symbol, r.TypeID); err != nil {
				return err
			}
			time.Sleep(time.Duration(r1.Intn(6)+1) + time.Second)
		}
	}

	return err
}

const (
	SERVER_PATH = `C:\Users\athor\go\otc_ft\app\server`
)

// builds url used to fetch pdf documents
func buildResourceURL(id int) (string, error) {
	var url strings.Builder
	urlTemp := template.New("urlTemp")
	urlTemp, err := urlTemp.Parse(FIN_REPORT_URL)
	if err != nil {
		return "", err
	}
	if err = urlTemp.Execute(&url, struct {
		ID string
	}{
		ID: strconv.Itoa(id),
	}); err != nil {
		return "", err
	}
	return url.String(), nil
}

func downloadRecord(id int, symbol, typeID string) error {
	var err error

	// file to output the report into
	outFile, err := os.Create(SERVER_PATH + `\` + symbol + strconv.Itoa(id) + typeID + ".pdf")
	if err != nil {
		return err
	}
	defer outFile.Close()

	resourceURL, err := buildResourceURL(id)
	if err != nil {
		return err
	}

	c := colly.NewCollector(colly.UserAgent(getUserAgent()))

	c.OnRequest(func(r *colly.Request) {
		setHeaders(r, BASE_AUTHORITY, resourceURL[len(BASE_AUTHORITY):])
		fmt.Println("request:", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		respData, err := io.ReadAll(bytes.NewReader(r.Body))
		if err != nil {
			// need to log this error
			fmt.Println(err)
		}
		_, err = io.Copy(outFile, bytes.NewReader(respData))
	})

	if err = c.Visit(resourceURL); err != nil {
		return err
	}

	return err
}

// gets the total amount of financial records for a stock
func getTotalRecordPages(urlAPI string) (int, error) {
	var data TotalFinancialReports
	var totalPages int
	var err error

	c := colly.NewCollector(colly.UserAgent(getUserAgent()))

	c.OnRequest(func(r *colly.Request) {
		setHeadersAPI(r)
	})

	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &data)
		totalPages = data.Pages
	})

	if err = c.Visit(urlAPI); err != nil {
		return -1, err
	}

	return totalPages, err
}

func getTotalNewsPages(urlAPI string, logger *log.Logger) (int, error) {
	var data TotalNews
	var totalPages int
	var err error

	c := colly.NewCollector(colly.UserAgent(getUserAgent()))

	c.OnRequest(func(r *colly.Request) {
		setHeadersAPI(r)
	})

	c.OnResponse(func(r *colly.Response) {
		if err = json.Unmarshal(r.Body, &data); err != nil {
			logger.Println(err)
		}
		totalPages = data.Pages
	})

	if err = c.Visit(urlAPI); err != nil {
		return -1, err
	}

	return totalPages, err
}

func scrapeNews(symbol string, logger *log.Logger) error {
	// all news will be unmarshaled into data
	var data TotalNews

	// form url for initial API request to get total records
	urlTemp := template.New("urlTemp")
	urlTemp, err := urlTemp.Parse(ALL_NEWS_URL)
	if err != nil {
		return err
	}
	var url strings.Builder
	if err = urlTemp.Execute(&url, struct {
		Symbol   string
		PageNum  string
		PageSize string
	}{
		Symbol:   symbol,
		PageNum:  "1",
		PageSize: MAX_PAGE_SIZE,
	}); err != nil {
		return err
	}

	// get total amount of pages to scrape
	totalPages, err := getTotalNewsPages(url.String(), logger)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	// grab records of each page
	for i := 1; i <= totalPages; i++ {
		url.Reset()
		if err = urlTemp.Execute(&url, struct {
			Symbol   string
			PageNum  string
			PageSize string
		}{
			Symbol:   symbol,
			PageNum:  strconv.Itoa(i),
			PageSize: MAX_PAGE_SIZE,
		}); err != nil {
			return err
		}

		c := colly.NewCollector(colly.UserAgent(getUserAgent()))

		c.OnRequest(func(r *colly.Request) {
			setHeadersAPI(r)
		})

		c.OnResponse(func(r *colly.Response) {
			if err = json.Unmarshal(r.Body, &data); err != nil {
				logger.Println(err)
			}
		})

		if err = c.Visit(url.String()); err != nil {
			logger.Println(err)
			return err
		}

		// download news reports
		randSource := rand.NewSource(time.Now().UnixNano())
		randGen := rand.New(randSource)

		for _, r := range data.Records {
			if err = downloadNews(symbol, r.Title, strconv.Itoa(r.ID)); err != nil {
				return err
			}

			time.Sleep(time.Duration(randGen.Intn(4)+1) * time.Second)
		}
	}

	return err
}

// store file name and url of article in db
func downloadNews(symbol, title, id string) error {
	// create txt file to contain news article text
	outFile, err := os.Create(SERVER_PATH + `\` + symbol + id + title + ".txt")
	if err != nil {
		return err
	}
	defer outFile.Close()

	// create url to fetch the HTML of the news article
	urlTemp := template.New("urlTemp")
	urlTemp, err = urlTemp.Parse(NEWS_URL)
	if err != nil {
		return err
	}
	var url strings.Builder
	if err = urlTemp.Execute(&url, struct {
		Symbol string
		Title  string
		ID     string
	}{
		Symbol: symbol,
		Title:  title,
		ID:     id,
	}); err != nil {
		return err
	}

	// execute the python script
	cmd := exec.Command("py", "script.py", url.String(), getUserAgent(), outFile.Name())
	if _, err := cmd.Output(); err != nil {
		return err
	}

	return err
}

type TotalFinancialReports struct {
	TotalRecords int    `json:"totalRecords"`
	Pages        int    `json:"pages"`
	CurrentPage  int    `json:"currentPage"`
	PageSize     int    `json:"pageSize"`
	SortOn       string `json:"sortOn"`
	SortDir      string `json:"sortDir"`
	Records      []struct {
		ID               int    `json:"id"`
		CompanyID        int    `json:"companyId"`
		UserID           int    `json:"userId"`
		Title            string `json:"title"`
		TypeID           string `json:"typeId"`
		StatusID         string `json:"statusId"`
		PeriodDate       int64  `json:"periodDate"`
		IsImmediate      bool   `json:"isImmediate"`
		CreatedDate      int64  `json:"createdDate"`
		LastModifiedDate int64  `json:"lastModifiedDate"`
		ReleaseDate      int64  `json:"releaseDate"`
		CanDistribute    bool   `json:"canDistribute"`
		WasDistributed   bool   `json:"wasDistributed"`
		CompanyName      string `json:"companyName"`
		ReportType       string `json:"reportType"`
		Name             string `json:"name"`
		StatusDescript   string `json:"statusDescript"`
		Symbol           string `json:"symbol"`
		PrimarySymbol    string `json:"primarySymbol"`
		IsCaveatEmptor   bool   `json:"isCaveatEmptor"`
		EdgarSECFiling   bool   `json:"edgarSECFiling"`
		TierCode         string `json:"tierCode"`
	} `json:"records"`
	Singular  string `json:"singular"`
	Plural    string `json:"plural"`
	CompanyID int    `json:"companyId"`
	StatusID  string `json:"statusId"`
	Empty     bool   `json:"empty"`
}

type TotalNews struct {
	TotalRecords int    `json:"totalRecords"`
	Pages        int    `json:"pages"`
	CurrentPage  int    `json:"currentPage"`
	PageSize     int    `json:"pageSize"`
	SortOn       string `json:"sortOn"`
	SortDir      string `json:"sortDir"`
	Records      []struct {
		ID                         int    `json:"id"`
		CompanyID                  int    `json:"companyId"`
		UserID                     int    `json:"userId"`
		Title                      string `json:"title"`
		TypeID                     string `json:"typeId"`
		StatusID                   string `json:"statusId"`
		Location                   string `json:"location"`
		IsImmediate                bool   `json:"isImmediate"`
		CreatedDate                int64  `json:"createdDate"`
		LastModifiedDate           int64  `json:"lastModifiedDate"`
		ReleaseDate                int64  `json:"releaseDate"`
		CanDistribute              bool   `json:"canDistribute"`
		WasDistributed             bool   `json:"wasDistributed"`
		NewsTypeDescript           string `json:"newsTypeDescript"`
		StatusDescript             string `json:"statusDescript"`
		Symbol                     string `json:"symbol"`
		IsCaveatEmptor             bool   `json:"isCaveatEmptor"`
		SourceID                   string `json:"sourceId"`
		DisplayDateTime            string `json:"displayDateTime"`
		Display                    bool   `json:"display"`
		TierCode                   string `json:"tierCode"`
		IsItemFromExternalProvider bool   `json:"isItemFromExternalProvider"`
		Immediate                  bool   `json:"immediate"`
	} `json:"records"`
	Singular string `json:"singular"`
	Plural   string `json:"plural"`
	Empty    bool   `json:"empty"`
}
