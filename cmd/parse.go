package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/crispyarty/novelparser/internal"
	"github.com/crispyarty/novelparser/internal/config"
	"github.com/crispyarty/novelparser/internal/mailer"
	"github.com/crispyarty/novelparser/internal/parsers"
	"github.com/crispyarty/novelparser/internal/savers"
	"github.com/spf13/cobra"
)

var (
	batchCount int
	sendEmails bool
)

var ParseCmd = &cobra.Command{
	Use:   "parse [novel_name]",
	Short: "Parse batch-size chapters of novel by name",
	Long:  `This command will parse new batch of a novel`,
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		novels := config.Config().Novels
		names := make([]string, 0, len(novels))

		if len(args) >= 1 {
			// Stop further completion (including file completion)
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		for name := range novels {
			names = append(names, name)
			// names = append(names, fmt.Sprintf("%v\t%v\n", name, novel.LastChapterUrl))
		}

		return names, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		ch := make(chan string)

		if sendEmails {
			mailer.Validate()
		}

		for i := range batchCount {
			log.Printf("Downloading batch %v/%v\n", i+1, batchCount)
			filename := batchParse(args[0])

			if sendEmails {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							ch <- `Error sending email: "` + filepath.Base(filename) + `"`
						}
					}()
					mailer.SendBook(filename)
					ch <- `Book sent via email: "` + filepath.Base(filename) + `"`
				}()
			}
		}

		if sendEmails {
			for range batchCount {
				fmt.Println(<-ch)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(ParseCmd)
	ParseCmd.Flags().IntVarP(&batchCount, "count", "c", 1, "Number of batches of chapters to parse")
	ParseCmd.Flags().BoolVarP(&sendEmails, "email", "e", false, "Send books to email")
}

func batchParse(novelName string) string {
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

	return filename
}
