package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/ravvio/awst/ui/style"
	"github.com/ravvio/awst/ui/tlog"
	"github.com/ravvio/awst/utils"
	"github.com/spf13/cobra"
)

func init() {
	logsGetCommand.Flags().Int32P("limit", "l", 1000, "limit number of log events to fetch")

	logsGetCommand.Flags().StringP("from", "f", "1d", "")

	logsGetCommand.Flags().BoolP("tail", "t", false, "start live tail")

	logsGetCommand.Flags().String("filter", "", "pattern filter on log events")
}

var logsGetCommand = &cobra.Command{
	Use: "get",
	Short: "Get cloudwatch logs of given log group",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := loadAwsConfig(context.TODO())
		utils.CheckErr(err)

		logGroupName := args[0]

		// Setup params
		params := &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName: &logGroupName,
		}

		filter, err := cmd.Flags().GetString("filter")
		utils.CheckErr(err)
		if filter != "" {
			params.FilterPattern = &filter
		}

		// Request
		client := cloudwatchlogs.NewFromConfig(cfg)

		logEventsOutput, err := client.FilterLogEvents(context.TODO(), params)
		utils.CheckErr(err)

		if len(logEventsOutput.Events) == 0 {
			style.PrintInfo("No events found")
			return
		}

		r := tlog.DefaultRenderer()
		for _, event := range logEventsOutput.Events {
			r.Render(&tlog.Log{
				GroupName: &logGroupName,
				Timestamp: event.Timestamp,
				Message: event.Message,
			})
		}

		// TODO Tail
	},
}
