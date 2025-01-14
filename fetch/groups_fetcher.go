package fetch

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

const DEFAULT_GROUPS_LIMIT = 20

type GroupsFetchData = FetchData[types.LogGroup]

type GroupsFetcherClient struct {
	Client *cloudwatchlogs.Client
	Params cloudwatchlogs.DescribeLogGroupsInput
}

func (c *GroupsFetcherClient) Fetch(ctx context.Context) (GroupsFetchData, error) {
	res, err := c.Client.DescribeLogGroups(ctx, &c.Params)
	if err != nil {
		return GroupsFetchData{}, err
	}

	data := GroupsFetchData{
		Data:      res.LogGroups,
		NextToken: res.NextToken,
	}
	return data, nil
}

func (c *GroupsFetcherClient) RequestLimit() *int32 {
	return c.Params.Limit
}

func (c *GroupsFetcherClient) SetRequestLimit(limit *int32) {
	c.Params.Limit = limit
}

func (c *GroupsFetcherClient) SetNextToken(token *string) {
	c.Params.NextToken = token
}

type GroupsFetcher = Fetcher[*GroupsFetcherClient, types.LogGroup]

func NewGroupsFetcher(
	ctx context.Context,
	client *GroupsFetcherClient,
) GroupsFetcher {
	return NewFetcher(ctx, client, DEFAULT_GROUPS_LIMIT)
}
