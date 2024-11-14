package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "awst",
	Short: "A utility to manage AWS resources",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(s3command)
	s3command.AddCommand(s3listCommand)

	rootCmd.AddCommand(logsCommand)
	logsCommand.AddCommand(logsListCommad)
}

var s3command = &cobra.Command{
	Use: "s3",
	Short: "",
}
