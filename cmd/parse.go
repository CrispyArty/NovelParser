package cmd

import (
	"log"

	"github.com/crispyarty/novelparser/internal"
	"github.com/crispyarty/novelparser/internal/config"
	"github.com/crispyarty/novelparser/internal/parsers"
	"github.com/crispyarty/novelparser/internal/savers"
	"github.com/spf13/cobra"
)

// add https://nobellink.com/noveltitle/chapter-1 --name noveltitle

var (
	batchCount int
)

var ParseCmd = &cobra.Command{
	Use:   "parse [novel name]",
	Short: "Parse batch-size chapters of novel by name",
	Long:  `This command will parse new batch of a novel`,
	Args:  cobra.ExactArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		for i := range batchCount {
			log.Printf("Downloading batch %v/%v\n", i+1, batchCount)
			batchParse(args[0])
		}
	},
}

func init() {
	RootCmd.AddCommand(ParseCmd)
	ParseCmd.Flags().IntVarP(&batchCount, "count", "c", 1, "Number of batches of chapters to parse")
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
