package internal

import (
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func buildRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicf("Error creating get request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	return req
}

func fetchWithAttempts(url string) (resp *http.Response, err error) {
	delay := 5 * time.Second
	attemts := 5
	req := buildRequest(url)

	for i := range attemts {
		log.Println("Fetcing content from:", url)
		client := &http.Client{}
		resp, err = client.Do(req)

		if err == nil && resp.StatusCode == 200 {
			return
		}

		if err != nil {
			log.Printf("Request failed with error (attempt %d/%d): %v. Retrying in %v...\n", i+1, attemts, err, delay)
		} else {
			log.Printf("Request failed with wrong status code (attempt %d/%d): %v. Retrying in %v...\n", i+1, attemts, resp.Status, delay)
		}

		if i+1 != attemts {
			time.Sleep(delay)
		}
	}

	return
}

func Fetch(url string) *goquery.Document {
	resp, err := fetchWithAttempts(url)
	if err != nil {
		log.Panicf("Error fetcing url %v\nurl:%v\n", err, url)
	}

	defer resp.Body.Close()

	htmlBody, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatalf("Error Parse response body: %v", err)
	}

	return htmlBody
}
