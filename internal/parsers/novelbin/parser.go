package novelbin

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/crispyarty/novelparser/internal"
)

type ParseHtmlNobelBin struct {
	doc *goquery.Document
}

func (parser *ParseHtmlNobelBin) Init(doc *goquery.Document) {
	parser.doc = doc
}

func parseTitle(doc *goquery.Document) (string, int) {
	title := doc.Find(".chr-title").First().Text()
	title = strings.TrimSpace(title)

	reg := regexp.MustCompile(`Chapter (\d+)[^a-zA-Z]*(.+)`)

	matches := reg.FindStringSubmatch(title)
	if len(matches) == 0 {
		return "", 0
	}

	numStr, str := matches[1], matches[2]

	// return reg.ReplaceAllString(title, "Chapter $1: $2")
	title = fmt.Sprintf("Chapter %v: %v", numStr, str)

	num, _ := strconv.Atoi(numStr)

	return title, num
}

func parseNextUrl(doc *goquery.Document) string {
	nextCpNode := doc.Find("#next_chap").First()
	url, _ := nextCpNode.Attr("href")

	return url
}

func parseParagraphs(doc *goquery.Document) (paragraphs []string) {
	// validTags := []string{"p", "h1", "h2", "h3", "h4", "h5", "h6"}
	validTags := []string{"p"}

	doc.Find("#chr-content").Children().Each(func(i int, s *goquery.Selection) {
		if !slices.Contains(validTags, goquery.NodeName(s)) {
			return
		}

		text := strings.TrimSpace(s.Text())

		if text == "" {
			return
		}

		paragraphs = append(paragraphs, text)
	})

	return
}

func (parser *ParseHtmlNobelBin) Parse() internal.NovelData {
	// fmt.Println("NobelBin Parse() implementation")

	title, number := parseTitle(parser.doc)

	return internal.NovelData{
		Title:         title,
		ChapterNumber: number,
		Paragraphs:    parseParagraphs(parser.doc),
		NextUrl:       parseNextUrl(parser.doc),
	}
}
