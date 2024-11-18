package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
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

var (
	keyIndex = "index"
	keyCreationDate = "creation"
	keyName = "name"
	keyArn = "arn"
	keyRetention = "retention"
	keyStreams = "streams"
)

var logsListCommad = &cobra.Command{
	Use:   "list",
	Short: "List cloudwatch log groups",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := config.LoadDefaultConfig(context.TODO())
		utils.CheckErr(err)

		// Setup params using flags
		params := &cloudwatchlogs.DescribeLogGroupsInput{}

		limit, err := cmd.Flags().GetInt32("limit")
		utils.CheckErr(err)
		if limit < 0 || limit > 50 {
			utils.CheckErr(fmt.Errorf("Limit must have a value between 0 and 50, was %d", limit))
		}
		params.Limit = &limit

		pattern, err := cmd.Flags().GetString("pattern");
		utils.CheckErr(err)
		if pattern != "" {
			params.LogGroupNamePattern = &pattern;
		}

		prefix, err := cmd.Flags().GetString("prefix");
		utils.CheckErr(err)
		if prefix != "" {
			params.LogGroupNamePrefix = &prefix;
		}

		// Request
		allPages, err := cmd.Flags().GetBool("all")
		utils.CheckErr(err)

		logGroups := []types.LogGroup{}
		client := cloudwatchlogs.NewFromConfig(cfg)
		for {
			output, err := client.DescribeLogGroups(context.TODO(), params)
			utils.CheckErr(err)

			logGroups = append(logGroups, output.LogGroups...)

			if !allPages || output.NextToken == nil {
				break
			}
			params.NextToken = output.NextToken
		}

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
				logStreams, err := client.DescribeLogStreams( context.TODO(), params)
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
		columns := []tables.Column{
			tables.NewColumn(keyIndex, "#", true),
			tables.NewColumn(keyCreationDate, "Creation", true),
			tables.NewColumn(keyName, "Name", true),
			tables.NewColumn(keyArn, "ARN", showArn),
			tables.NewColumn(keyRetention, "Retention", showRetention),
			tables.NewColumn(keyStreams, "Streams", showStreams),
		}

		rows := []tables.Row{}
		for index, group := range logGroups {
			rows = append(rows, tables.Row{
				keyIndex: fmt.Sprintf("%d", index+1),
				keyCreationDate: time.UnixMilli(*group.CreationTime).Format("2006-01-02"),
				keyName: *group.LogGroupName,
				keyArn: *group.LogGroupArn,
				keyRetention: fmt.Sprintf("%d days", group.RetentionInDays),
				keyStreams: strings.Join(streams[*group.LogGroupName], ", "),
			})
		}

		table := tables.New(columns).WithRows(rows)

		// Render Table
		fmt.Println(table.Render())
	},
}
