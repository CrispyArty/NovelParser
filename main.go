package main

import (
	"log"

	"github.com/crispyarty/novelparser/internal"
	"github.com/crispyarty/novelparser/internal/config"
	"github.com/crispyarty/novelparser/internal/parsers"
	"github.com/crispyarty/novelparser/internal/savers"
	// "github.com/crispyarty/novelparser/internal/savers"
)

func singleParse() {
	novelConf := config.Novel("my_simulated_road_to_immortality")
	parserCreator := parsers.ParserFactory(novelConf.LastChapterUrl)
	parser := parserCreator()

	doc := internal.Fetch(novelConf.LastChapterUrl)
	parser.Init(doc)
	novel := parser.Parse()

	// fmt.Println(novel)

	novels := []*internal.NovelData{&novel}

	savers.SaveNovel("my_simulated_road_to_immortality", novels)
}

func batchParse(novelName string) {
	novelConf := config.Novel(novelName)
	parserCreator := parsers.ParserFactory(novelConf.LastChapterUrl)

	var novels []*internal.NovelData

	nextUrl := novelConf.LastChapterUrl

	for range novelConf.BatchSize {
		doc := internal.Fetch(nextUrl)
		parser := parserCreator()
		parser.Init(doc)
		novel := parser.Parse()
		novels = append(novels, &novel)

		nextUrl = novel.NextUrl
	}

	filename := savers.SaveNovel(novelName, novels)

	log.Println("Saved!", filename)

	config.UpdateLastChapter(novelName, nextUrl)
}

func main() {
	config.Init()
	defer config.Save()

	batchParse("my_simulated_road_to_immortality")

}
