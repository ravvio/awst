package cmd

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/ravvio/awst/ui/style"
	"github.com/ravvio/awst/ui/tlog"
	"github.com/ravvio/awst/utils"
	"github.com/spf13/cobra"
)

func init() {
	logsSearchCommand.Flags().StringP("pattern", "e", "", "pattern filter on log group name")
	logsSearchCommand.Flags().StringP("prefix", "p", "", "prefix filter on log group name")

	logsSearchCommand.Flags().Int32P("limit", "l", 1000, "limit number of log events to fetch")

	logsSearchCommand.Flags().StringP("filter", "f", "", "pattern filter on log events")
	logsSearchCommand.Flags().String("from", "1d", "moment in time to start the search, can be absolute or relative")
	logsSearchCommand.Flags().String("to", "", "")

	logsSearchCommand.Flags().BoolP("tail", "t", false, "start live tail")

	logsSearchCommand.MarkFlagsOneRequired("pattern", "prefix")
	logsSearchCommand.MarkFlagsMutuallyExclusive("pattern", "prefix")
}

var logsSearchCommand = &cobra.Command{
	Use: "search",
	Short: "Search for cloudwatch log groups matching given pattern or prefix and retrive logs",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := config.LoadDefaultConfig(context.TODO())
		utils.CheckErr(err)

		client := cloudwatchlogs.NewFromConfig(cfg)

		// Setup params for descibe operation
		describeParams := &cloudwatchlogs.DescribeLogGroupsInput{}

		pattern, err := cmd.Flags().GetString("pattern");
		utils.CheckErr(err)
		if pattern != "" {
			describeParams.LogGroupNamePattern = &pattern;
		}

		prefix, err := cmd.Flags().GetString("prefix");
		utils.CheckErr(err)
		if prefix != "" {
			describeParams.LogGroupNamePrefix = &prefix;
		}

		// Request describe
		logGroupsOutput, err := client.DescribeLogGroups(context.TODO(), describeParams)
		utils.CheckErr(err)

		logGroups := logGroupsOutput.LogGroups

		// TODO make limit configurable
		if len(logGroups) >= 10 {
			utils.CheckErr(fmt.Errorf("Log groups limit exceeded: found %d log groups > 10", len(logGroups)))
		}

		style.PrintInfo("%d groups found", len(logGroups))

		// Request logs
		filter, err := cmd.Flags().GetString("filter")
		utils.CheckErr(err)
		limit, err := cmd.Flags().GetInt32("limit")
		utils.CheckErr(err)
		tail, err := cmd.Flags().GetBool("tail")
		utils.CheckErr(err)

		fromDate, err := cmd.Flags().GetString("from")
		var fromUnix int64
		utils.CheckErr(err)
		if fromDate != "" {
			if t, err := time.Parse(time.RFC3339, fromDate); err != nil {
				m := t.UnixMilli()
				fromUnix = m
			} else if d, err := utils.ParseDuration(fromDate); err != nil {
				fromUnix = time.Now().UnixMilli() - d
			} else {
				utils.CheckErr(fmt.Errorf("Could not parse 'from' timestamp"))
			}
		}

		logs := []tlog.Log{}
		for _, group := range logGroups {
			params := &cloudwatchlogs.FilterLogEventsInput{
				FilterPattern: &filter,
				LogGroupName: group.LogGroupName,
				Limit: &limit,
				StartTime: &fromUnix,
			}

			eventsOutput, err := client.FilterLogEvents(context.TODO(), params)
			utils.CheckErr(err)

			for _, event := range eventsOutput.Events {
				logs = append(logs, tlog.Log{
					GroupName: *group.LogGroupName,
					Timestamp: time.UnixMilli(*event.Timestamp),
					Message: *event.Message,
				})
			}
		}

		sort.Slice(logs, func(i, j int) bool {
			return logs[i].Timestamp.Compare(logs[j].Timestamp) > 0
		})

		r := tlog.DefaultRenderer()
		for _, log := range logs {
			r.Render(&log)
		}

		if !tail {
			return
		}

		logStream := make(chan tlog.Log)

		// Create tail streams
		i := 0
		for {
			n := min(i + 10, len(logGroups))

			identifiers := []string{}
			for _, l := range logGroups[i:n] {
				identifiers = append(identifiers, *l.LogGroupArn)
			}


			tailParams := &cloudwatchlogs.StartLiveTailInput{
				LogEventFilterPattern: &filter,
				LogGroupIdentifiers: identifiers,
			}

			tailOutput, err := client.StartLiveTail(
				context.TODO(),
				tailParams,
			)
			utils.CheckErr(err)

			s := tailOutput.GetStream()
			handleStream(logStream, s)

			if (i + 10 >= len(logGroups)) {
				break;
			}
			i += 10;
		}

		for {
			log := <- logStream
			r.Render(&log)
		}

	},
}

func handleStream(
	logStream chan tlog.Log,
	eventStream *cloudwatchlogs.StartLiveTailEventStream,
) {
	for {
		if (logStream != nil) {
			event := <- eventStream.Events()

			switch e := event.(type) {
			case *types.StartLiveTailResponseStreamMemberSessionStart:
				style.PrintInfo("Session %d start", e.Value.SessionId)
			case *types.StartLiveTailResponseStreamMemberSessionUpdate:
				for _, logEvent := range e.Value.SessionResults {
					log := tlog.Log{
						GroupName: *logEvent.LogGroupIdentifier,
						Timestamp: time.UnixMilli(*logEvent.Timestamp),
						Message: *logEvent.Message,
					}
					logStream <- log
				}
			default:
				utils.CheckErr(eventStream.Err())
			}
		}
	}
}
