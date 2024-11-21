package cmd

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/ravvio/awst/fetch"
	"github.com/ravvio/awst/ui/style"
	"github.com/ravvio/awst/ui/tlog"
	"github.com/ravvio/awst/utils"
	"github.com/spf13/cobra"
)

func init() {
	logsSearchCommand.Flags().StringP("pattern", "e", "", "pattern filter on log group name")
	logsSearchCommand.Flags().StringP("prefix", "p", "", "prefix filter on log group name")
	logsSearchCommand.Flags().Bool("all-groups", false, "do not limit of log groups to use")
	logsSearchCommand.Flags().Int32("limit-groups", 50, "limit number of log groups to use")

	logsSearchCommand.Flags().BoolP("all", "a", false, "do not limit of log events to fetch from each group")
	logsSearchCommand.Flags().Int32P("limit", "l", 10000, "limit number of log events to fetch from each group")

	logsSearchCommand.Flags().StringP("filter", "f", "", "pattern filter on log events")
	logsSearchCommand.Flags().String("since", "1d", "moment in time to start the search, can be absolute or relative")
	logsSearchCommand.Flags().String("to", "", "")

	logsSearchCommand.Flags().BoolP("tail", "t", false, "start live tail")

	logsSearchCommand.MarkFlagsOneRequired("pattern", "prefix")

	logsSearchCommand.MarkFlagsMutuallyExclusive("pattern", "prefix")
	logsSearchCommand.MarkFlagsMutuallyExclusive("all", "limit")
}

var logsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search for cloudwatch log groups matching given pattern or prefix and retrive logs",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := loadAwsConfig(context.TODO())
		utils.CheckErr(err)

		client := cloudwatchlogs.NewFromConfig(cfg)

		// Setup params for descibe operation
		describeParams := &cloudwatchlogs.DescribeLogGroupsInput{}

		pattern, err := cmd.Flags().GetString("pattern")
		utils.CheckErr(err)
		if pattern != "" {
			describeParams.LogGroupNamePattern = &pattern
		}

		prefix, err := cmd.Flags().GetString("prefix")
		utils.CheckErr(err)
		if prefix != "" {
			describeParams.LogGroupNamePrefix = &prefix
		}

		allGroups, err := cmd.Flags().GetBool("all-groups")
		utils.CheckErr(err)
		limitGroups, err := cmd.Flags().GetInt32("limit-groups")
		utils.CheckErr(err)

		// Request describe
		logGroups := []types.LogGroup{}
		for {
			logGroupsOutput, err := client.DescribeLogGroups(context.TODO(), describeParams)
			utils.CheckErr(err)

			logGroups = append(logGroups, logGroupsOutput.LogGroups...)

			if !allGroups ||
				len(logGroups) >= int(limitGroups) ||
				logGroupsOutput.NextToken == nil {
				break
			}
			describeParams.NextToken = logGroupsOutput.NextToken
		}

		style.PrintInfo("%d groups found", len(logGroups))

		// Request logs
		filter, err := cmd.Flags().GetString("filter")
		utils.CheckErr(err)
		limitEvents, err := cmd.Flags().GetInt32("limit")
		utils.CheckErr(err)
		allEvents, err := cmd.Flags().GetBool("all")
		utils.CheckErr(err)
		tail, err := cmd.Flags().GetBool("tail")
		utils.CheckErr(err)

		sinceDate, err := cmd.Flags().GetString("since")
		var fromUnix int64
		utils.CheckErr(err)
		if sinceDate != "" {
			if t, err := utils.ParseDatetime(sinceDate); err == nil && t.UnixMilli() >= 0 {
				fromUnix = t.UnixMilli()
			} else if d, err := utils.ParseDuration(sinceDate); err == nil {
				fromUnix = time.Now().UnixMilli() - d
			} else {
				utils.CheckErr(fmt.Errorf("Could not parse 'from' timestamp"))
			}
		}

		var wg sync.WaitGroup

		logs := []tlog.Log{}
		for _, group := range logGroups {
			fetcher := fetch.NewLogsFetcher(
				context.TODO(),
				client,
				cloudwatchlogs.FilterLogEventsInput{
					FilterPattern: &filter,
					LogGroupName:  group.LogGroupName,
					Limit:         &limitEvents,
					StartTime:     &fromUnix,
				},
			)
			if !allEvents {
				fetcher = fetcher.WithLimit(limitEvents)
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				res, err := fetcher.All()
				utils.CheckErr(err)
				for _, log := range res {
					logs = append(
						logs,
						utils.LogFromCloudwatchEvent(group.LogGroupName, &log),
					)
				}
			}()
		}

		// Wait for all goroutines to finish fetching
		wg.Wait()

		sort.Slice(logs, func(i, j int) bool {
			return *logs[j].Timestamp > *logs[i].Timestamp
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
			n := min(i+10, len(logGroups))

			identifiers := []string{}
			for _, l := range logGroups[i:n] {
				identifiers = append(identifiers, *l.LogGroupArn)
			}

			tailParams := &cloudwatchlogs.StartLiveTailInput{
				LogEventFilterPattern: &filter,
				LogGroupIdentifiers:   identifiers,
			}

			tailOutput, err := client.StartLiveTail(
				context.TODO(),
				tailParams,
			)
			utils.CheckErr(err)

			s := tailOutput.GetStream()
			handleStream(logStream, s)

			if i+10 >= len(logGroups) {
				break
			}
			i += 10
		}

		for {
			log := <-logStream
			r.Render(&log)
		}

	},
}

func handleStream(
	logStream chan tlog.Log,
	eventStream *cloudwatchlogs.StartLiveTailEventStream,
) {
	for {
		if logStream != nil {
			event := <-eventStream.Events()

			switch e := event.(type) {
			case *types.StartLiveTailResponseStreamMemberSessionStart:
				style.PrintInfo("Session %d start", e.Value.SessionId)
			case *types.StartLiveTailResponseStreamMemberSessionUpdate:
				for _, logEvent := range e.Value.SessionResults {
					log := tlog.Log{
						GroupName: logEvent.LogGroupIdentifier,
						Timestamp: logEvent.Timestamp,
						Message:   logEvent.Message,
					}
					logStream <- log
				}
			default:
				utils.CheckErr(eventStream.Err())
			}
		}
	}
}
