package fetch_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ravvio/awst/fetch"
	"github.com/stretchr/testify/assert"
)

type TestFetchData = fetch.FetchData[string]

type TestFetcherClient struct {
	limit *int32
}

func (t *TestFetcherClient) Fetch(ctx context.Context) (TestFetchData, error) {
	data := []string{}
	i := 0
	for i = 0; i < int(*t.limit); i++ {
		data = append(data, fmt.Sprintf("hello %d", i))
	}
	next := "next"
	return TestFetchData{
		Data:      data,
		NextToken: &next,
	}, nil
}

func (t *TestFetcherClient) RequestLimit() *int32 {
	return t.limit
}

func (t *TestFetcherClient) SetRequestLimit(limit *int32) {
	t.limit = limit
}

func (t *TestFetcherClient) SetNextToken(_ *string) {}

type TestFetcher = fetch.Fetcher[*TestFetcherClient, string]

// --- //

func TestFetchLimit(t *testing.T) {
	const RLIMIT = 10
	const LIMIT = 125

	c := TestFetcherClient{}
	f := fetch.NewFetcher(context.Background(), &c, RLIMIT).WithLimit(LIMIT)

	assert.Equal(t, true, f.HasNextPage())

	r, e := f.All()
	assert.NoError(t, e)
	assert.Equal(t, LIMIT, len(r))

	assert.Equal(t, false, f.HasNextPage())
}

func TestPagination(t *testing.T) {
	const RLIMIT = 10
	const LIMIT = 15

	c := TestFetcherClient{}
	f := fetch.NewFetcher(context.Background(), &c, RLIMIT).WithLimit(LIMIT)

	assert.Equal(t, true, f.HasNextPage())

	r, e := f.NextPage()
	assert.NoError(t, e)
	assert.Equal(t, RLIMIT, len(r))

	r, e = f.NextPage()
	assert.NoError(t, e)
	assert.Equal(t, LIMIT-RLIMIT, len(r))

	assert.Equal(t, false, f.HasNextPage())
}
