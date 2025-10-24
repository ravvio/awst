package cmd

import (
	"time"

	"github.com/ravvio/awst/utils"
	"github.com/spf13/cobra"
)

var (
	patternGroupFlag *string
	prefixGroupFlag  *string
	allGroupsFlag    bool
	limitGroupsFlag  int32

	allLogsFlag    bool
	limitLogsFlag  int32
	filterLogsFlag *string
	tailLogsFlag   bool

	sinceLogsFlag *string
	untilLogsFlag *string

	threadsFlag int
)

var (
	now int64
	sinceUnix *int64
	untilUnix *int64
)

func init() {
	logsCommand.AddCommand(logsSearchCommand)
	logsCommand.AddCommand(logsGetCommand)

	logsCommand.Flags().BoolVar(&allLogsFlag, "all-logs", false, "do not limit of log events to fetch from each group")
	logsCommand.Flags().Int32Var(&limitLogsFlag, "limit-logs", 10000, "limit number of log events to fetch from each group")
	logsCommand.MarkFlagsMutuallyExclusive("all-logs", "limit-logs")

	logsCommand.Flags().StringVarP(filterLogsFlag, "filter-logs", "f", "", "pattern filter on log events")
	logsCommand.Flags().StringVar(sinceLogsFlag, "since", "1d", "moment in time to start the search, can be absolute or relative")
	logsCommand.Flags().StringVar(untilLogsFlag, "until", "0s", "moment in time to end the search, can be absolute or relative")

	logsCommand.Flags().BoolP("tail", "t", false, "start live tail")

	logsCommand.Flags().IntVar(&threadsFlag, "threads", 5, "maximum parallelization for fetching")

	// Search
	logsSearchCommand.Flags().StringVarP(patternGroupFlag, "pattern", "e", "", "pattern filter on log group name")
	logsSearchCommand.Flags().StringVarP(prefixGroupFlag, "prefix", "p", "", "prefix filter on log group name")
	logsSearchCommand.Parent().MarkFlagsOneRequired("pattern", "prefix")
	logsSearchCommand.MarkFlagsMutuallyExclusive("pattern", "prefix")

	logsSearchCommand.Flags().BoolVar(&allGroupsFlag, "all-groups", false, "do not limit of log groups to use")
	logsSearchCommand.Flags().Int32Var(&limitGroupsFlag, "limit-groups", 50, "limit number of log groups to use")
	logsCommand.MarkFlagsMutuallyExclusive("all-groups", "limit-groups")
}

var logsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve and tail logs",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		now := time.Now().UnixMilli()
		// Timestamps
		sinceUnix, err = utils.ParseTimestamp(now, sinceLogsFlag);
		if err != nil {
			return err
		}
		untilUnix, err = utils.ParseTimestamp(now, sinceLogsFlag);
		if err != nil {
			return err
		}
		return nil
	},
}
