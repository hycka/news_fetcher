package service

import (
	"context"
	"log"
	"net"

	pb "github.com/hycka/news_fetcher/api/news_fetcher/api"
	"github.com/hycka/news_fetcher/internal/biz"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedNewsFetcherServer
	*grpc.Server
	address string
	fc      *biz.FetcherCase
}

type Options struct {
	Address string
}

func NewServer(opts Options) *Server {
	return &Server{Server: grpc.NewServer(), address: opts.Address}
}

func NewFetcherServer(fc *biz.FetcherCase) pb.NewsFetcherServer {
	return &Server{fc: fc}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	log.Printf("\ngrpc server start at: %s", s.address)
	return s.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	log.Printf("grpc server gracefully stopped.")
	return nil
}

// GetNews implements api.GetNews, it get News's info and set it back to `in`
func (s *Server) List(ctx context.Context, in *pb.ID) (*pb.Posts, error) {
	return s.fc.List(in)

}

// SearchNews, it search news from json file
func (s *Server) Search(ctx context.Context, in *pb.Keyword) (*pb.Posts, error) {
	return s.fc.Search(in)
}
