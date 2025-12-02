package parsers

import (
	"github.com/PuerkitoBio/goquery"

	"github.com/crispyarty/novelparser/internal"
)

type ParseHtml interface {
	Init(*goquery.Document)
	Parse() internal.NovelData
}
