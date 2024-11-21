package fetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// TODO do not use params but custom options

type LogsFetcher struct {
	ctx    context.Context
	client *cloudwatchlogs.Client
	params cloudwatchlogs.FilterLogEventsInput
	limit  int32

	fetched    int32
	first_page bool
	next_token *string
}

func NewLogsFetcher(
	ctx context.Context,
	client *cloudwatchlogs.Client,
	params cloudwatchlogs.FilterLogEventsInput,
) LogsFetcher {
	return LogsFetcher{
		ctx:        ctx,
		client:     client,
		params:     params,
		limit:      -1,
		fetched:    0,
		first_page: true,
		next_token: nil,
	}
}

func (f LogsFetcher) WithLimit(limit int32) LogsFetcher {
	f.limit = limit
	return f
}

func (f *LogsFetcher) HasNextPage() bool {
	return f.first_page ||
	(f.next_token != nil && (f.limit < 0 || f.fetched < f.limit))
}

func (f *LogsFetcher) NextPage() ([]types.FilteredLogEvent, error) {
	if !f.HasNextPage() {
		return nil, fmt.Errorf("No next page")
	}
	f.params.NextToken = f.next_token

	if f.limit > 0 {
		newLimit := min(
			f.limit-int32(f.fetched),
			*f.params.Limit,
		)
		f.params.Limit = &newLimit
	}

	res, err := f.client.FilterLogEvents(f.ctx, &f.params)
	if err != nil {
		return nil, err
	}

	f.first_page = false
	f.fetched += int32(len(res.Events))
	f.next_token = res.NextToken
	return res.Events, nil
}

func (f *LogsFetcher) All() ([]types.FilteredLogEvent, error) {
	logs := []types.FilteredLogEvent{}

	for {
		l, err := f.NextPage()
		if err != nil {
			return nil, err
		}
		logs = append(logs, l...)

		if !f.HasNextPage() { break }
	}

	return logs, nil
}
