package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	loggCommand.AddCommand(loggListCommad)
}

var loggCommand = &cobra.Command{
	Use: "logg",
	Short: "Retrieve log groups",
}
