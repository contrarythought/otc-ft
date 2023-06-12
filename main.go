package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"os"
	"otc_ft/app/scraper"
)

type Creds struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func getCreds() (*Creds, error) {
	file, err := os.Open("creds.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var creds Creds
	if err = json.Unmarshal(jsonData, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

func findErrLog() (bool, error) {
	entries, err := os.ReadDir(`app`)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.Name() == `err_log.txt` {
			return true, nil
		}
	}

	return false, nil
}

func main() {
	creds, err := getCreds()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", "postgres://"+creds.User+":"+creds.Pass+"@localhost/otc_fts?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	var file *os.File
	haveFile, err := findErrLog()
	if err != nil {
		log.Fatal(err)
	}

	if !haveFile {
		file, err = os.Create(`err_log.txt`)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		file, err = os.OpenFile(`err_log.txt`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()

	if err = scraper.Scrape(file, db); err != nil {
		log.Fatal(err)
	}
}
