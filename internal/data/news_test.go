package data

import (
	"log"
	"testing"

	pb "github.com/hycka/news_fetcher/api/news_fetcher/api"
)

// go test -test.run=^TestGetNews$
func TestList(t *testing.T) {
	fr := NewFetcherRepo()
	fr.LoadExistNews()
	n, err := fr.List(&pb.ID{Id: "15cf90e5bdbf18a036d8a1c99c57ad0f"})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(n)
}

func TestSearch(t *testing.T) {
	fr := NewFetcherRepo()
	fr.LoadExistNews()
	n, err := fr.Search(&pb.Keyword{})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(n)
}
