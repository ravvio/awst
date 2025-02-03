package fetch

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const DEFAULT_DDB_TABLES_LIMIT = 100

type DDBTablesFetchData = FetchData[string]

type DDBTablesFetcherClient struct {
	Client *dynamodb.Client
	Params dynamodb.ListTablesInput
}

func (c *DDBTablesFetcherClient) Fetch(ctx context.Context) (DDBTablesFetchData, error) {
	res, err := c.Client.ListTables(ctx, &c.Params)
	if err != nil {
		return DDBTablesFetchData{}, err
	}

	data := DDBTablesFetchData{
		Data:      res.TableNames,
		NextToken: res.LastEvaluatedTableName,
	}
	return data, nil
}

func (c *DDBTablesFetcherClient) RequestLimit() *int32 {
	return c.Params.Limit
}

func (c *DDBTablesFetcherClient) SetRequestLimit(limit *int32) {
	c.Params.Limit = limit
}

func (c *DDBTablesFetcherClient) SetNextToken(token *string) {
	c.Params.ExclusiveStartTableName = token
}

type DDBTablesFetcher = Fetcher[*DDBTablesFetcherClient, string]

func NewDDBTablesFetcher(
	ctx context.Context,
	client *DDBTablesFetcherClient,
) DDBTablesFetcher {
	return NewFetcher(ctx, client, DEFAULT_DDB_TABLES_LIMIT)
}
