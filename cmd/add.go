package cmd

import (
	"fmt"
	"net/url"

	"github.com/crispyarty/novelparser/internal/config"
	"github.com/spf13/cobra"
)

var (
	batchSize int
)

func validateUrl(argPos int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		arg := args[argPos-1]

		u, err := url.Parse(arg)

		if err != nil || u.Host == "" {
			return fmt.Errorf("argument %d must be valid url, received \"%v\"", argPos, arg)
		}

		return nil
	}
}

var AddCmd = &cobra.Command{
	Use:   "add [novel name] [first chapter url] --batch-size",
	Short: "Add new novel to config with first chapter",
	Args:  cobra.MatchAll(cobra.ExactArgs(2), validateUrl(2)),
	Long:  `This command will add novel with first chapter link, and you can then use this name for "parse" command `,
	Run: func(cmd *cobra.Command, args []string) {
		config.UpdateLastChapter(args[0], args[1])

		fmt.Printf("New novel in config: %v\n", args[0])
	},
}

func init() {
	RootCmd.AddCommand(AddCmd)
	// AddCmd.Flags().StringVarP(&novelName, "name", "n", "", "Name on the novel to use in `parser` command")
	// AddCmd.MarkFlagRequired("name")
	AddCmd.Flags().IntVarP(&batchSize, "batch-size", "b", 10, "Number of novels saved as one")
}

// add https://nobellink.com/noveltitle/chapter-1 --name noveltitle
