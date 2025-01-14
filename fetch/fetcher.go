package fetch

import (
	"context"
	"fmt"
)

type FetchData[T any] struct {
	Data      []T
	NextToken *string
}

type FetcherClient[T any] interface {
	Fetch(context.Context) (FetchData[T], error)
	RequestLimit() *int32
	SetRequestLimit(*int32)
	SetNextToken(*string)
}

type Fetcher[C FetcherClient[T], T any] struct {
	ctx    context.Context
	client C

	limit      int32
	fetched    int32
	first_page bool
	next_token *string
}

func NewFetcher[C FetcherClient[T], T any](
	ctx context.Context,
	c C,
	defaultRequestLimit int32,
) Fetcher[C, T] {
	f := Fetcher[C, T]{
		ctx:    ctx,
		client: c,

		limit:      -1,
		fetched:    0,
		first_page: true,
		next_token: nil,
	}

	if f.client.RequestLimit() == nil {
		f.client.SetRequestLimit(&defaultRequestLimit)
	}
	return f
}

func (f Fetcher[C, T]) WithLimit(limit int32) Fetcher[C, T] {
	f.limit = limit
	return f
}

func (f *Fetcher[C, T]) HasNextPage() bool {
	return f.first_page ||
		(f.next_token != nil && (f.limit < 0 || f.fetched < f.limit))
}

func (f *Fetcher[C, T]) NextPage() ([]T, error) {
	if !f.HasNextPage() {
		return nil, fmt.Errorf("no next page")
	}
	f.client.SetNextToken(f.next_token)

	if f.limit > 0 {
		newLimit := min(
			f.limit-int32(f.fetched),
			*f.client.RequestLimit(),
		)
		f.client.SetRequestLimit(&newLimit)
	}

	res, err := f.client.Fetch(f.ctx)
	if err != nil {
		return nil, err
	}

	f.first_page = false
	f.fetched += int32(len(res.Data))
	f.next_token = res.NextToken
	return res.Data, nil
}

func (f *Fetcher[C, T]) All() ([]T, error) {
	res, err := f.NextPage()
	if err != nil {
		return nil, err
	}

	for {
		r, err := f.NextPage()
		if err != nil {
			break
		}

		res = append(res, r...)

		if !f.HasNextPage() {
			break
		}
	}

	return res, nil
}
