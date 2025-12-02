package savers

import (
	"fmt"
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
	t := time.Now()
	return t.Format("2006-01-02")
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

func newContent(name string, novels []*internal.NovelData) *Content {
	return &Content{
		NovelName:  name,
		Novels:     novels,
		Translator: "",
	}
}
