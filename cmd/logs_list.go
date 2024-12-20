package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/ravvio/awst/fetch"
	"github.com/ravvio/awst/ui/style"
	"github.com/ravvio/awst/ui/tables"
	"github.com/ravvio/awst/utils"
	"github.com/spf13/cobra"
)

func init() {
	logsListCommad.Flags().BoolP("all", "a", false, "fetch all log groups")
	logsListCommad.Flags().Int32P("limit", "l", 50, "limit number of groups to fetch")

	logsListCommad.Flags().StringP("pattern", "e", "", "pattern filter on log group name")
	logsListCommad.Flags().StringP("prefix", "p", "", "prefix filter on log group name")

	logsListCommad.Flags().Bool("retention", false, "show log groups retention")
	logsListCommad.Flags().Bool("arn", false, "show log groups arn")
	logsListCommad.Flags().Bool("streams", false, "show log groups streams")

	logsListCommad.MarkFlagsMutuallyExclusive("pattern", "prefix")
	logsListCommad.MarkFlagsMutuallyExclusive("all", "limit")
}

var logsListCommad = &cobra.Command{
	Use:   "list",
	Short: "List cloudwatch log groups",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := loadAwsConfig(context.TODO())
		utils.CheckErr(err)

		// Setup params using flags
		params := &cloudwatchlogs.DescribeLogGroupsInput{}

		limit, err := cmd.Flags().GetInt32("limit")
		utils.CheckErr(err)

		pattern, err := cmd.Flags().GetString("pattern")
		utils.CheckErr(err)
		if pattern != "" {
			params.LogGroupNamePattern = &pattern
		}

		prefix, err := cmd.Flags().GetString("prefix")
		utils.CheckErr(err)
		if prefix != "" {
			params.LogGroupNamePrefix = &prefix
		}

		// Request
		all, err := cmd.Flags().GetBool("all")
		utils.CheckErr(err)

		client := cloudwatchlogs.NewFromConfig(cfg)

		groupsFetcher := fetch.NewGroupsFetcher(
			context.TODO(),
			&fetch.GroupsFetcherClient{
				Client: client,
				Params: *params,
			},
		)
		if !all {
			groupsFetcher = groupsFetcher.WithLimit(limit)
		}
		logGroups, err := groupsFetcher.All()
		utils.CheckErr(err)

		showStreams, err := cmd.Flags().GetBool("streams")
		utils.CheckErr(err)
		showRetention, err := cmd.Flags().GetBool("retention")
		utils.CheckErr(err)
		showArn, err := cmd.Flags().GetBool("arn")
		utils.CheckErr(err)

		// If streams are requested recover them
		var streams = map[string][]string{}
		if showStreams {
			for _, group := range logGroups {
				params := &cloudwatchlogs.DescribeLogStreamsInput{
					LogGroupName: group.LogGroupName,
				}
				logStreams, err := client.DescribeLogStreams(context.TODO(), params)
				utils.CheckErr(err)

				var names = []string{}
				for _, stream := range logStreams.LogStreams {
					names = append(names, *stream.LogStreamName)
				}
				streams[*group.LogGroupName] = names
			}
		}

		if len(logGroups) == 0 {
			style.PrintInfo("No groups found")
			return
		}

		// Setup table
		var (
			keyIndex        = "index"
			keyCreationDate = "creation"
			keyName         = "name"
			keyArn          = "arn"
			keyRetention    = "retention"
			keyStreams      = "streams"
		)

		columns := []tables.Column{
			tables.NewColumn(keyIndex, "#", true).WithAlignment(tables.Right),
			tables.NewColumn(keyCreationDate, "Creation", true),
			tables.NewColumn(keyName, "Name", true),
			tables.NewColumn(keyArn, "Arn", showArn),
			tables.NewColumn(keyRetention, "Retention", showRetention).WithAlignment(tables.Right),
			tables.NewColumn(keyStreams, "Streams", showStreams),
		}


		rows := []tables.Row{}
		for index, group := range logGroups {
			var retention string
			if group.RetentionInDays != nil {
				retention = fmt.Sprintf("%d days", *group.RetentionInDays)
			} else {
				retention = "-"
			}
			rows = append(rows, tables.Row{
				keyIndex:        fmt.Sprintf("%d", index+1),
				keyCreationDate: time.UnixMilli(*group.CreationTime).Format("2006-01-02"),
				keyName:         *group.LogGroupName,
				keyArn:          *group.LogGroupArn,
				keyRetention:    retention,
				keyStreams:      strings.Join(streams[*group.LogGroupName], ", "),
			})
		}

		table := tables.New(columns).WithRows(rows)

		// Render Table
		fmt.Println(table.Render())
	},
}
