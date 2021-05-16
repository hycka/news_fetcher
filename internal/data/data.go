package data

import (
	"github.com/hycka/news_fetcher/internal/biz"
)

var _ biz.FetcherRepo = new(fetcherRepo)

type fetcherRepo struct {
}

func NewFetcherRepo() biz.FetcherRepo {
	return &fetcherRepo{}
}
