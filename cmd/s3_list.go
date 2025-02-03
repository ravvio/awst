package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ravvio/awst/ui/tables"
	"github.com/ravvio/awst/utils"
	"github.com/spf13/cobra"
)

func init() {
	s3listCommand.Flags().BoolP("all", "a", false, "Do not limit number of buckets to fetch")
	s3listCommand.Flags().Int32P("limit", "l", 50, "Maximum number of buckets to fetch")

	s3listCommand.MarkFlagsMutuallyExclusive("all", "limit")
}

var s3listCommand = &cobra.Command{
	Use:   "list",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := loadAwsConfig(context.TODO())
		utils.CheckErr(err)

		// Setup params
		params := &s3.ListBucketsInput{}

		region, err := cmd.Flags().GetString("region")
		utils.CheckErr(err)
		if region != "" {
			params.BucketRegion = &region
		}

		all, err := cmd.Flags().GetBool("all")
		utils.CheckErr(err)
		if !all {
			limit, err := cmd.Flags().GetInt32("limit")
			utils.CheckErr(err)
			params.MaxBuckets = &limit
		}

		// Request
		client := s3.NewFromConfig(cfg)
		output, err := client.ListBuckets(context.TODO(), params)
		utils.CheckErr(err)

		// Setup table
		var (
			keyIndex        = "index"
			keyCreationDate = "creation"
			keyName         = "name"
		)

		columns := []tables.Column{
			tables.NewColumn(keyIndex, "#", true).WithAlignment(tables.Right),
			tables.NewColumn(keyCreationDate, "Creation", true),
			tables.NewColumn(keyName, "Name", true),
		}

		rows := []tables.Row{}
		for index, bucket := range output.Buckets {
			rows = append(rows, tables.Row{
				keyIndex:        fmt.Sprintf("%d", index+1),
				keyCreationDate: bucket.CreationDate.Format("2006-01-02"),
				keyName:         *bucket.Name,
			})
		}

		table := tables.New(columns).WithRows(rows)

		// Render table
		fmt.Println(table.Render())
	},
}
