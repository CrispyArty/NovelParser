package cmd

import (
	"fmt"

	"github.com/crispyarty/novelparser/internal"
	"github.com/crispyarty/novelparser/internal/mailer"
	"github.com/crispyarty/novelparser/internal/parsers"
	"github.com/spf13/cobra"
)

var TestCmd = &cobra.Command{
	Use:   "test [novel_name]",
	Short: "Test and debug",
	Long:  `This command will parse new batch of a novel`,
	Run: func(cmd *cobra.Command, args []string) {
		testParse()
	},
}

func init() {
	RootCmd.AddCommand(TestCmd)

}

func testEmail() {
	mailer.Validate()
	mailer.SendBook("uploads/my_simulated_road_to_immortality/test.epub")
}

func testParse() {
	url := "https://novelbin.com/b/my-simulated-road-to-immortality/chapter-1601-1357-the-immortal-worlds-dao-of-flames-tribulation-7k2"
	parserCreator := parsers.ParserFactory(url)
	doc := internal.Fetch(url)
	parser := parserCreator()
	parser.Init(doc)
	novel := parser.Parse()

	for p, v := range novel.Paragraphs {
		fmt.Println(p, v)
	}

}
