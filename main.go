package main

import (
	"fmt"
	"reflect"

	// "io"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/crispyarty/novelparser/internal/parsers"
	// "rsc.io/quote"
)

func fetch(url string) *goquery.Document {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatalf("Error creating get request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error making HTTP request: %v, status: %v", err, resp)
	}

	if resp.StatusCode == 429 {
		// TODO: Try again in 5 seconds
		log.Panicf("!!!, %v", resp.Status)
		time.Sleep(5 * time.Second)
		// fetch(url)
	}

	// htmlBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("Error reading response body: %v", err)
	// }

	// log.Println(resp.Status)
	defer resp.Body.Close() // Ensure the response body is closed

	// htmlBody, err := html.Parse(resp.Body)
	htmlBody, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatalf("Error Parse response body: %v", err)
	}

	return htmlBody
}

func main() {
	// go say("world")
	// say("hello")
	// fmt.Println(internal.Qwe())
	url := "https://novelbin.com/b/my-simulated-road-to-immortality/chapter-201-201-187-the-sacred-one-with-a-strange-aura"
	// content := fetch(url)

	parserCreator := parsers.ParserFactory(url)

	parser := parserCreator()
	// parser.Init(content)

	// parser.Parse()

	fmt.Printf("%T, %v\n", parser, parser)
	fmt.Println(reflect.TypeOf(parser))
	fmt.Println("Hello!")

	// for i := range 50 {
	// fetch(url)
	// fmt.Printf("%v - %v\n", i, url)
	// }

	// html.Parse(content)
	// fetch(url)
	// fmt.Println(content)

	// titleNode := content.Find(".chr-title").First()

	// fmt.Println(titleNode)
	// fmt.Println(titleNode.Text())

	// nextCpNode := content.Find("#next_chap").First()

	// fmt.Println(nextCpNode.Attr("href"))
	// parserFunc := parserFactory
	// parser := parserFunc(url)

	// fmt.Println(parser)
	// fmt.Println(parser.(*novelbin.ParseHtmlNobelBin).Url)
	// parser.NextUrl()

	// str, error := parseAndSaveNovel(parser, 10)

	// fmt.Println(str)
	// fmt.Println(error)
}
