package internal

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/crispyarty/novelparser/internal/config"
)

func buildRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicf("Error creating get request: %v", err)
	}

	httpHeaders := config.Config().HttpHeaders

	req.Header.Set("User-Agent", httpHeaders.UserAgent)

	for name, value := range httpHeaders.Cookies {
		cookie := &http.Cookie{Name: name, Value: value}
		req.AddCookie(cookie)
	}

	return req
}

func fetchWithAttempts(url string, attemts int) (resp *http.Response, err error) {
	delay := 7 * time.Second
	req := buildRequest(url)

	for i := range attemts {
		log.Println("Fetcing content from:", url)
		client := &http.Client{}
		resp, err = client.Do(req)

		if err == nil && resp.StatusCode == 200 {
			return
		}

		if err != nil {
			log.Printf("Request failed with error (attempt %d/%d): %v", i+1, attemts, err)
		} else {
			log.Printf("Request failed with wrong status code (attempt %d/%d): %v", i+1, attemts, resp.Status)
		}

		if i+1 == attemts {
			err = errors.New("too many attempts")
			return
		}

		log.Printf("Retrying in %v...\n", delay)
		time.Sleep(delay)
	}

	return
}

func Fetch(url string) *goquery.Document {
	resp, err := fetchWithAttempts(url, 1)
	if err != nil {
		log.Panicf("Error fetcing url %v\n", err)
	}

	defer resp.Body.Close()

	htmlBody, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatalf("Error Parse response body: %v", err)
	}

	return htmlBody
}
