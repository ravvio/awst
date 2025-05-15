package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/ravvio/awst/fetch"
	"github.com/ravvio/awst/ui/style"
	"github.com/ravvio/awst/ui/tlog"
	"github.com/ravvio/awst/utils"
	"github.com/spf13/cobra"
)

func init() {
	logsGetCommand.Flags().BoolP("all", "a", false, "do not limit of log events to fetch from each group")
	logsGetCommand.Flags().Int32P("limit", "l", 10000, "limit number of log events to fetch from each group")

	logsGetCommand.Flags().StringP("filter", "f", "", "pattern filter on log events")
	logsGetCommand.Flags().String("since", "1d", "moment in time to start the search, can be absolute or relative")
	logsGetCommand.Flags().String("until", "0s", "moment in time to end the search, can be absolute or relative")

	logsGetCommand.Flags().BoolP("tail", "t", false, "start live tail")
}

var logsGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get cloudwatch logs of given log group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := loadAwsConfig(context.TODO())
		utils.CheckErr(err)

		now := time.Now()

		logGroupName := args[0]

		// Setup params
		filter, err := cmd.Flags().GetString("filter")
		utils.CheckErr(err)
		limitEvents, err := cmd.Flags().GetInt32("limit")
		utils.CheckErr(err)
		allEvents, err := cmd.Flags().GetBool("all")
		utils.CheckErr(err)
		// tail, err := cmd.Flags().GetBool("tail")
		// utils.CheckErr(err)

		sinceDate, err := cmd.Flags().GetString("since")
		var sinceUnix int64
		utils.CheckErr(err)
		if sinceDate != "" {
			if t, err := utils.ParseDatetime(sinceDate); err == nil && t.UnixMilli() >= 0 {
				sinceUnix = t.UnixMilli()
			} else if d, err := utils.ParseDuration(sinceDate); err == nil {
				sinceUnix = now.UnixMilli() - d
			} else {
				utils.CheckErr(fmt.Errorf("Could not parse 'since' timestamp"))
			}
		}

		untilDate, err := cmd.Flags().GetString("until")
		var untilUnix int64
		utils.CheckErr(err)
		if untilDate != "" {
			if t, err := utils.ParseDatetime(untilDate); err == nil && t.UnixMilli() >= 0 {
				untilUnix = t.UnixMilli()
			} else if d, err := utils.ParseDuration(untilDate); err == nil {
				untilUnix = now.UnixMilli() - d
			} else {
				utils.CheckErr(fmt.Errorf("Could not parse 'until' timestamp"))
			}
		}

		// Request
		client := cloudwatchlogs.NewFromConfig(cfg)
		logFetcher := fetch.NewLogsFetcher(
			context.TODO(),
			&fetch.LogsFetcherClient{
				Client: client,
				Params: cloudwatchlogs.FilterLogEventsInput{
					LogGroupName:  &logGroupName,
					StartTime:     &sinceUnix,
					EndTime:       &untilUnix,
					FilterPattern: &filter,
				},
			},
		)
		if !allEvents {
			logFetcher = logFetcher.WithLimit(limitEvents)
		}

		logEvents, err := logFetcher.All()
		utils.CheckErr(err)

		if len(logEvents) == 0 {
			style.PrintInfo("No events found")
			return
		}

		r := tlog.DefaultRenderer()
		for _, event := range logEvents {
			err = r.Render(&tlog.Log{
				GroupName: &logGroupName,
				Timestamp: event.Timestamp,
				Message:   event.Message,
			})
			utils.CheckErr(err)
		}

		// TODO Tail
	},
}
