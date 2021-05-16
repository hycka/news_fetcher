package json

import (
	"log"
	"testing"

	pb "github.com/hycka/news_fetcher/api/news_fetcher/api"
)

// go test -test.run=^TestGetNews$
func TestLoadLocalNews(t *testing.T) {
	err := LoadLocalNews()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewsExist(t *testing.T) {
	LoadLocalNews()
	b := NewsExist("15cf90e5bdbf18a036d8a1c99c57ad0f")
	log.Println(b)
}

func TestSearchNews(t *testing.T) {
	LoadLocalNews()
	p, err := SearchNews(&pb.Keyword{Keyword: ""})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(p)
}

func TestSelectNews(t *testing.T) {
	LoadLocalNews()
	p, err := SelectNews(&pb.ID{Id: "15cf90e5bdbf18a036d8a1c99c57ad0f"})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(p)
}

func TestSaveNewsCollection(t *testing.T) {
	tmp := pb.Post{Id: "1", Title: "ssssssssssssssss", Content: "kkkkkkkkkkkkkkkkkkkk", UpdateTime: 1621096770000000000, WebsiteId: "b1c062a456dd250b0ccd24d310d1b5d3", WebsiteTitle: "新西兰联合报"}
	tmp2 := []*pb.Post{&tmp}
	err := SaveNewsCollection(&pb.Posts{Posts: tmp2})
	if err != nil {
		t.Fatal(err)
	}
	LoadLocalNews()
	p, _ := SelectNews(&pb.ID{Id: "1"})
	log.Println(p)
}
