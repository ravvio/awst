package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

func init() {
	logsListCommad.Flags().Int32P("limit", "l", 50, "Maximum number of groups to fetch")
	logsListCommad.Flags().StringP("pattern", "e", "", "Pattern filter on log group name")
	logsListCommad.Flags().StringP("prefix", "p", "", "Prefix filter on log group name")

	logsListCommad.Flags().Bool("retention", false, "Show retention")
	logsListCommad.Flags().Bool("arn", false, "Show arn")

	logsListCommad.MarkFlagsMutuallyExclusive("pattern", "prefix")
}

var logsListCommad = &cobra.Command{
	Use:   "list",
	Short: "List log groups",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := config.LoadDefaultConfig(context.TODO())
		checkErr(err)

		// Load flag not used in params
		showRetention, err := cmd.Flags().GetBool("retention")
		checkErr(err)
		showArn, err := cmd.Flags().GetBool("arn")
		checkErr(err)

		// Setup params using flags
		params := &cloudwatchlogs.DescribeLogGroupsInput{}

		limit, err := cmd.Flags().GetInt32("limit")
		checkErr(err)
		params.Limit = &limit

		pattern, err := cmd.Flags().GetString("pattern");
		checkErr(err)
		if pattern != "" {
			params.LogGroupNamePattern = &pattern;
		}

		prefix, err := cmd.Flags().GetString("prefix");
		checkErr(err)
		if prefix != "" {
			params.LogGroupNamePrefix = &prefix;
		}

		// Request
		client := cloudwatchlogs.NewFromConfig(cfg)
		output, err := client.DescribeLogGroups(context.TODO(), params)
		checkErr(err)

		w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
		fields := []string{"#", "CreationDate", "Name"}
		if showRetention {
			fields = append(fields, "Retention")
		}
		if showArn {
			fields = append(fields, "ARN")
		}
		heading := strings.Join(fields, "\t")
		fmt.Fprintln(w, heading)

		for index, object := range output.LogGroups {
			fields := []string{
				fmt.Sprintf("%d", index+1),
				time.UnixMilli(*object.CreationTime).Format("2006-01-02"),
				*object.LogGroupName,
			}
			if showRetention {
				fields = append(
					fields,
					fmt.Sprintf("%d days", object.RetentionInDays),
				)
			}
			if showArn {
				fields = append(
					fields,
					*object.LogGroupArn,
				)
			}
			line := strings.Join(fields, "\t")
			fmt.Fprintln(w, line)
		}
		err = w.Flush()
		checkErr(err)
	},
}
