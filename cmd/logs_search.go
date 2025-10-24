package cmd

import (
	"context"
	"sort"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/ravvio/awst/fetch"
	"github.com/ravvio/awst/ui/style"
	"github.com/ravvio/awst/ui/tlog"
	"github.com/ravvio/awst/utils"
	"github.com/spf13/cobra"
)

var logsSearchCommand = &cobra.Command{
	Use:   "search",
	Short: "Search for cloudwatch log groups matching given pattern or prefix and retrive logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Load config
		cfg, err := loadAwsConfig(ctx)
		if err != nil {
			return err
		}
		client := cloudwatchlogs.NewFromConfig(cfg)

		// Request describe
		groupsFetcher := fetch.NewGroupsFetcher(
			context.TODO(),
			&fetch.GroupsFetcherClient{
				Client: client,
				Params: cloudwatchlogs.DescribeLogGroupsInput{
					LogGroupNamePattern: patternGroupFlag,
					LogGroupNamePrefix: prefixGroupFlag,
				},
			},
		)
		if !allGroupsFlag {
			groupsFetcher = groupsFetcher.WithLimit(limitGroupsFlag)
		}
		logGroups, err := groupsFetcher.All()
		if err != nil {
			return err
		}

		// Render
		if len(logGroups) == 0 {
			style.PrintInfo("No groups found")
			return nil
		}

		style.PrintInfo("%d groups found", len(logGroups))

		// Set up sempaphore
		semaphore := make(chan struct{}, threadsFlag)
		var wg sync.WaitGroup = sync.WaitGroup{}

		logs := []tlog.Log{}
		for _, group := range logGroups {
			logsFetcher := fetch.NewLogsFetcher(
				context.TODO(),
				&fetch.LogsFetcherClient{
					Client: client,
					Params: cloudwatchlogs.FilterLogEventsInput{
						LogGroupName:  group.LogGroupName,
						StartTime:     sinceUnix,
						EndTime:       untilUnix,
						FilterPattern: filterLogsFlag,
					},
				},
			)
			if !allLogsFlag {
				logsFetcher = logsFetcher.WithLimit(limitLogsFlag)
			}

			// Acquire semaphore
			semaphore <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				res, err := logsFetcher.All()
				if err != nil {
					style.PrintError("could not fetch logs: %s", err)
					<-semaphore
				}
				for _, log := range res {
					logs = append(
						logs,
						utils.LogFromCloudwatchEvent(group.LogGroupName, &log),
					)
				}
				// Release semaphore
				<-semaphore
			}()
		}

		// Wait for all goroutines to finish fetching
		wg.Wait()

		sort.Slice(logs, func(i, j int) bool {
			return *logs[j].Timestamp > *logs[i].Timestamp
		})

		r := tlog.DefaultRenderer()
		for _, log := range logs {
			err = r.Render(&log)
			utils.CheckErr(err)
		}

		if !tailLogsFlag {
			return nil
		}

		logChan := make(chan tlog.Log)

		// Create tail streams
		i := 0
		for {
			n := min(i+10, len(logGroups))

			identifiers := []string{}
			for _, l := range logGroups[i:n] {
				identifiers = append(identifiers, *l.LogGroupArn)
			}

			tailParams := &cloudwatchlogs.StartLiveTailInput{
				LogEventFilterPattern: filterLogsFlag,
				LogGroupIdentifiers:   identifiers,
			}

			tailOutput, err := client.StartLiveTail(
				context.TODO(),
				tailParams,
			)
			utils.CheckErr(err)

			s := tailOutput.GetStream()
			go handleStream(logChan, s)

			i = n
			if i >= len(logGroups) {
				break
			}
		}

		for {
			log := <-logChan
			r.Render(&log)
		}

	},
}

func handleStream(
	logChan chan tlog.Log,
	eventStream *cloudwatchlogs.StartLiveTailEventStream,
) {
	eventsChan := eventStream.Events()
	for {
		if logChan != nil {
			event := <-eventsChan

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
					logChan <- log
				}
			default:
				utils.CheckErr(eventStream.Err())
			}
		}
	}
}
