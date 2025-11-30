package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "novelparser [command]",
	Short: "novelparser",
	Long:  `This application will parse various novels and compact them in batches`,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
}
