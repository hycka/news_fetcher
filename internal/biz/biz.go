package biz

import (
	pb "github.com/hycka/news_fetcher/api/news_fetcher/api"
)

type FetcherCase struct {
	repo FetcherRepo
}

type FetcherRepo interface {
	List(*pb.ID) (*pb.Posts, error)
	Search(*pb.Keyword) (*pb.Posts, error)
	UpdateNews() error
	LoadExistNews() error
}

func NewFetcherCase(repo FetcherRepo) *FetcherCase {
	return &FetcherCase{repo: repo}
}

func (fc *FetcherCase) List(id *pb.ID) (*pb.Posts, error) {
	return fc.repo.List(id)
}

func (fc *FetcherCase) Search(nc *pb.Keyword) (*pb.Posts, error) {
	return fc.repo.Search(nc)
}

func (fc *FetcherCase) UpdateNews() error {
	return fc.repo.UpdateNews()
}
