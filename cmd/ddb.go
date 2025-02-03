package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	loggCommand.AddCommand(loggListCommad)
}

var ddbCommand = &cobra.Command{
	Use: "ddb",
	Short: "Act on dynamo db tables",
}
