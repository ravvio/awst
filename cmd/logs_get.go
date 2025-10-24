package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/ravvio/awst/fetch"
	"github.com/ravvio/awst/ui/style"
	"github.com/ravvio/awst/ui/tlog"
	"github.com/spf13/cobra"
)

var logsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get cloudwatch logs of given log group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Load config
		cfg, err := loadAwsConfig(ctx)
		if err != nil {
			return err
		}

		// Log group name name
		logGroupName := args[0]

		// Request
		client := cloudwatchlogs.NewFromConfig(cfg)
		logFetcher := fetch.NewLogsFetcher(
			ctx,
			&fetch.LogsFetcherClient{
				Client: client,
				Params: cloudwatchlogs.FilterLogEventsInput{
					LogGroupName:  &logGroupName,
					StartTime:     sinceUnix,
					EndTime:       untilUnix,
					FilterPattern: filterLogsFlag,
				},
			},
		)
		if !allLogsFlag {
			logFetcher = logFetcher.WithLimit(limitLogsFlag)
		}

		logEvents, err := logFetcher.All()
		if err != nil {
			return err
		}

		// Render results
		if len(logEvents) == 0 {
			style.PrintInfo("No events found")
			return nil
		}

		r := tlog.DefaultRenderer()
		for _, event := range logEvents {
			err = r.Render(&tlog.Log{
				GroupName: &logGroupName,
				Timestamp: event.Timestamp,
				Message:   event.Message,
			})
			if err != nil {
				return err
			}
		}

		// TODO Tail

		return nil
	},
}
