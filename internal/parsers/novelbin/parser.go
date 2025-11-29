package novelbin

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"

	"github.com/crispyarty/novelparser/internal"
)

type ParseHtmlNobelBin struct {
	doc *goquery.Document
}

// type ParseHtmlNobelBin struct {
// 	html string
// }

func (parser *ParseHtmlNobelBin) Parse() internal.NovelData {
	fmt.Println("NobelBin Parse() implementation")

	return internal.NovelData{}
}

func (parser *ParseHtmlNobelBin) NextUrl() string {
	fmt.Println("NobelBin NextUrl() implementation")

	return "NextUrl"
}

func (parser *ParseHtmlNobelBin) Init(doc *goquery.Document) {
	parser.doc = doc
}
