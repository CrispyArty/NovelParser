package savers

import (
	"fmt"
	"strings"
	"time"

	"github.com/crispyarty/novelparser/internal"
)

type Content struct {
	NovelName  string
	Novels     []*internal.NovelData
	Translator string

	identifier string
}

func (c *Content) Title() string {
	if len(c.Novels) == 1 {
		return c.Novels[0].Title
	}

	return fmt.Sprintf("Chapter %v - %v", c.Novels[0].ChapterNumber, c.Novels[len(c.Novels)-1].ChapterNumber)
}

func (c *Content) Date() string {
	return time.Now().Format("2006-01-02")
}

func (c *Content) Identifier() string {
	if c.identifier != "" {
		return c.identifier
	}

	if len(c.Novels) == 1 {
		return fmt.Sprintf("chapter-%v", c.Novels[0].ChapterNumber)
	}

	return fmt.Sprintf("chapter-%v-%v", c.Novels[0].ChapterNumber, c.Novels[len(c.Novels)-1].ChapterNumber)
}

func (c *Content) Author() string {
	human := strings.ReplaceAll(c.NovelName, "_", " ")

	return strings.ToUpper(string(human[0])) + strings.ToLower(human[1:])
}

func newContent(name string, novels []*internal.NovelData) *Content {
	return &Content{
		NovelName:  name,
		Novels:     novels,
		Translator: "",
	}
}
