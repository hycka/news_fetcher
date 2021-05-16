package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/hycka/news_fetcher/api/news_fetcher/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:10001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNewsFetcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// search news: pass test
	// keywords := &pb.Keyword{Keyword: "以色列"}
	// ns, err := c.Search(ctx, keywords)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(len(ns.Posts))

	// get posts: pass test
	// v, err := c.List(ctx, &pb.ID{})
	// if err != nil {
	// 	log.Printf("c.GetPosts err: %+v", err)
	// }
	// fmt.Println(v)

	// get post pass test
	v, err := c.List(ctx, &pb.ID{Id: "c263ded8cfbb5820e13e14b878680e90,3a7997ec004dad7d0befe251600d22e6"})
	if err != nil {
		log.Printf("c.GetPosts err: %+v", err)
	}
	fmt.Println(v)

}
