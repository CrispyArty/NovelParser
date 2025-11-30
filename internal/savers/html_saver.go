package savers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/crispyarty/novelparser/internal"
)

const templatePath = "template.html"

type Content struct {
	Novels []*internal.NovelData
}

func (c *Content) Title() string {
	if len(c.Novels) == 1 {
		return c.Novels[0].Title
	}

	return fmt.Sprintf("Chapter %v - %v", c.Novels[0].ChapterNumber, c.Novels[len(c.Novels)-1].ChapterNumber)
}

// func GenHtml(data *Content) (buf *bytes.Buffer) {
// 	check := func(err error) {
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 	}

// 	tmpl, err := template.ParseFiles(templatePath)
// 	check(err)

// 	err = tmpl.Execute(buf, data)
// 	check(err)

// 	return
// }

func SaveNovel(name string, novels []*internal.NovelData) string {
	check := func(err error) {
		if err != nil {
			log.Panic(err)
		}
	}

	data := &Content{
		Novels: novels,
	}

	tmpl, err := template.ParseFiles(templatePath)
	check(err)

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	check(err)

	dir := fmt.Sprintf("uploads/%v", name)
	os.MkdirAll(dir, os.ModePerm)

	filename := fmt.Sprintf("%v/%v.html", dir, data.Title())

	file, err := os.Create(filename)
	check(err)

	defer file.Close()
	file.Write(buf.Bytes())

	return filename
}
