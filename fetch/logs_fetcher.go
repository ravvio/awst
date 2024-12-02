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

func (c *LogsFetcherClient) Fetch(ctx context.Context) (LogsFetchData, error) {
	res, err := c.Client.FilterLogEvents(ctx, &c.Params)
	if err != nil {
		return LogsFetchData{}, err
	}

	data := LogsFetchData{
		Data:      res.Events,
		NextToken: res.NextToken,
	}
	return data, nil
}

func (c *LogsFetcherClient) RequestLimit() *int32 {
	return c.Params.Limit
}

func (c *LogsFetcherClient) SetRequestLimit(limit *int32) {
	c.Params.Limit = limit
}

func (c *LogsFetcherClient) SetNextToken(token *string) {
	c.Params.NextToken = token
}

type LogsFetcher = Fetcher[*LogsFetcherClient, types.FilteredLogEvent]

func NewLogsFetcher(
	ctx context.Context,
	client *LogsFetcherClient,
) LogsFetcher {
	return NewFetcher(ctx, client, DEFAULT_LOG_LIMIT)
}
