package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ravvio/awst/fetch"
	"github.com/ravvio/awst/ui/style"
	"github.com/ravvio/awst/ui/tables"
	"github.com/ravvio/awst/utils"
	"github.com/spf13/cobra"
)

func init() {
	ddbListCommand.Flags().BoolP("all", "a", false, "fetch all log groups")
	ddbListCommand.Flags().Int32P("limit", "l", 50, "limit number of groups to fetch")

	ddbListCommand.MarkFlagsMutuallyExclusive("all", "limit")
}

var ddbListCommand = &cobra.Command{
	Use:   "list",
	Short: "List DynamoDB tables",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := loadAwsConfig(context.TODO())
		utils.CheckErr(err)

		// Setup params using flags
		params := &dynamodb.ListTablesInput{}

		limit, err := cmd.Flags().GetInt32("limit")
		utils.CheckErr(err)

		all, err := cmd.Flags().GetBool("all")
		utils.CheckErr(err)

		// Create client
		client := dynamodb.NewFromConfig(cfg)
		tablesFetcher := fetch.NewDDBTablesFetcher(
			context.TODO(),
			&fetch.DDBTablesFetcherClient{
				Client: client,
				Params: *params,
			},
		)
		if !all {
			tablesFetcher = tablesFetcher.WithLimit(limit)
		}

		tableList, err := tablesFetcher.All()
		utils.CheckErr(err)

		if len(tableList) == 0 {
			style.PrintInfo("No tables found")
			return
		}

		var (
			keyIndex = "index"
			keyName  = "name"
		)

		columns := []tables.Column{
			tables.NewColumn(keyIndex, "#", true).WithAlignment(tables.Right),
			tables.NewColumn(keyName, "Name", true),
		}

		rows := []tables.Row{}
		for index, table := range tableList {
			rows = append(rows, tables.Row{
				keyIndex: fmt.Sprintf("%d", index+1),
				keyName: table,
			})
		}

		table := tables.New(columns).WithRows(rows)
		fmt.Println(table.Render())
	},
}
