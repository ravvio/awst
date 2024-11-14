package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ravvio/awst/ui/style"
	"github.com/ravvio/awst/ui/tables"
	"github.com/spf13/cobra"
)

const (
	colKeyIndex     = "index"
	colKeyDate      = "date"
	colKeyName      = "name"
	colKeyArn       = "arn"
	colKeyRetention = "retention"
)

func init() {
	logsCommand.Flags().Int32P("limit", "l", 50, "Maximum number of groups to fetch")
	logsCommand.Flags().StringP("pattern", "e", "", "Pattern filter on log group name")
	logsCommand.Flags().StringP("prefix", "p", "", "Prefix filter on log group name")
}

var logsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Interact with cloudwatch log groups",
	Run: func(cmd *cobra.Command, args []string) {
		// Load aws config
		cfg, err := config.LoadDefaultConfig(context.TODO())
		checkErr(err)

		// Setup params using flags
		params := &cloudwatchlogs.DescribeLogGroupsInput{}

		limit, err := cmd.Flags().GetInt32("limit")
		checkErr(err)
		params.Limit = &limit

		pattern, err := cmd.Flags().GetString("pattern")
		checkErr(err)
		if pattern != "" {
			params.LogGroupNamePattern = &pattern
		}

		prefix, err := cmd.Flags().GetString("prefix")
		checkErr(err)
		if prefix != "" {
			params.LogGroupNamePrefix = &prefix
		}

		// Request
		client := cloudwatchlogs.NewFromConfig(cfg)
		output, err := client.DescribeLogGroups(context.TODO(), params)
		checkErr(err)

		// Check empty
		if len(output.LogGroups) == 0 {
			style.PrintInfo(
				"No log groups found for current filters",
			)
			return
		}

		// Create Table
		columns := []table.Column{
			{Title: "#", Width: 3},
			{Title: "Creation", Width: 10},
			{Title: "Name", Width: 50},
			{Title: "Arn", Width: 50},
		}

		rows := []table.Row{}
		for index, object := range output.LogGroups {
			rows = append(rows, table.Row{
				fmt.Sprintf("%d", index+1),
				time.UnixMilli(*object.CreationTime).Format("2006-01-02"),
				*object.LogGroupName,
				*object.LogGroupArn,
			})
		}

		var quit bool;
		table := tables.NewInteractive(columns, &quit).WithRows(rows)

		tp := tea.NewProgram(table)
		_, err = tp.Run()
		checkErr(err)

		if quit {
			tp.ReleaseTerminal()
			return
		}

		selected := table.SelectedRow()
		selectedLog := selected[3]

		tailParams := &cloudwatchlogs.StartLiveTailInput{
			LogGroupIdentifiers: []string{selectedLog},
		}
		logTail, err := client.StartLiveTail(
			context.TODO(),
			tailParams,
		)
		checkErr(err)

		logTailStream := logTail.GetStream()
		defer logTailStream.Close()

		eventsChan := logTailStream.Events()

		for {
			event := <-eventsChan
			switch e := event.(type) {
			case *types.StartLiveTailResponseStreamMemberSessionStart:
				fmt.Println(style.TitleStyle.Render("\nSessionStart"))
			case *types.StartLiveTailResponseStreamMemberSessionUpdate:
				for _, logEvent := range e.Value.SessionResults {
					fmt.Println(*logEvent.Message)
				}
			default:
				checkErr(err)
			}
		}
	},
}
