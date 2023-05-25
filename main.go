package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	resource := `https://www.otcmarkets.com/otcapi/company/financial-report/206941/content`
	file, err := os.Create(`filings/test.pdf`)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	req, err := http.NewRequest(http.MethodGet, resource, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("user-agent", getUserAgent())

	client := http.Client{
		Timeout: time.Second * 3,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		log.Fatal(err)
	}
}
