package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/hycka/news_fetcher/api/news_fetcher/api"
	"github.com/hycka/news_fetcher/internal/biz"
	"github.com/hycka/news_fetcher/internal/data"
	"github.com/hycka/news_fetcher/internal/service"
	"golang.org/x/sync/errgroup"
)

func InitFetcherCase() *biz.FetcherCase {
	fetcherRepo := data.NewFetcherRepo()
	fetcherCase := biz.NewFetcherCase(fetcherRepo)
	return fetcherCase
}

func main() {
	opts := service.Options{Address: ":10001"}

	fc := InitFetcherCase()

	fservice := service.NewFetcherServer(fc)

	s := service.NewServer(opts)
	pb.RegisterNewsFetcherServer(s, fservice)

	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer cancel()
		return s.Start(ctx)
	})
	g.Go(func() error {
		defer cancel()
		return LoadExistNews(ctx)
	})
	g.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		select {
		case sig := <-sigs:
			fmt.Println()
			log.Printf("signal caught: %s, ready to quit...", sig.String())
			defer cancel()
			s.Stop(ctx)
			os.Exit(1)
		case <-ctx.Done():
			defer cancel()
			s.Stop(ctx)
			return ctx.Err()
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Printf("news_fetcher server main error: %v", err)
	}
}

//Update news collection in cache
func LoadExistNews(ctx context.Context) error {
	fr := data.NewFetcherRepo()
	for {
		err := fr.LoadExistNews()
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Minute)
	}
}
