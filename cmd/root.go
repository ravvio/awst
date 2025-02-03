package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	region  string
	profile string
)

func init() {
	rootCmd.AddCommand(s3command)

	rootCmd.PersistentFlags().StringVar(&region, "region", "", "Specify AWS region")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "Specify AWS profile")

	s3command.AddCommand(s3listCommand)

	rootCmd.AddCommand(logsCommand)
	logsCommand.AddCommand(logsListCommad)
	logsCommand.AddCommand(logsGetCommand)
	logsCommand.AddCommand(logsSearchCommand)

	rootCmd.AddCommand(ddbCommand)
	ddbCommand.AddCommand(ddbListCommand)
}

var rootCmd = &cobra.Command{
	Use:   "awst",
	Short: "A utility to manage AWS resources",
}

var s3command = &cobra.Command{
	Use:   "s3",
	Short: "",
}

var logsCommand = &cobra.Command{
	Use:   "logs",
	Short: "",
}

var ddbCommand = &cobra.Command{
	Use:   "ddb",
	Short: "",
}
