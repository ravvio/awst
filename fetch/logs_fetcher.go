package fetch

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

const DEFAULT_LOG_LIMIT = 10_000

type LogsFetchData = FetchData[types.FilteredLogEvent]

type LogsFetcherClient struct {
	Client *cloudwatchlogs.Client
	Params cloudwatchlogs.FilterLogEventsInput
}

func (l *LogsFetcherClient) Fetch(ctx context.Context) (LogsFetchData, error) {
	res, err := l.Client.FilterLogEvents(ctx, &l.Params)
	if err != nil {
		return LogsFetchData{}, err
	}

	data := LogsFetchData{
		Data:      res.Events,
		NextToken: res.NextToken,
	}
	return data, nil
}

func (l *LogsFetcherClient) RequestLimit() *int32 {
	return l.Params.Limit
}

func (l *LogsFetcherClient) SetRequestLimit(limit *int32) {
	l.Params.Limit = limit
}

func (l *LogsFetcherClient) SetNextToken(token *string) {
	l.Params.NextToken = token
}

type LogsFetcher = Fetcher[*LogsFetcherClient, types.FilteredLogEvent]

func NewLogsFetcher(
	ctx context.Context,
	client *LogsFetcherClient,
) LogsFetcher {
	return NewFetcher(ctx, client, DEFAULT_LOG_LIMIT)
}
