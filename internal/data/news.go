package data

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	pb "github.com/hycka/news_fetcher/api/news_fetcher/api"
	fetcher "github.com/hycka/news_fetcher/internal/data/fetcher"
	db "github.com/hycka/news_fetcher/internal/pkg/db/json"
)

// load aready exist news in json file
func (fr *fetcherRepo) LoadExistNews() error {
	return db.LoadLocalNews()
}

// GetNews get News info if it's Id is currect
func (fr *fetcherRepo) List(n *pb.ID) (*pb.Posts, error) {
	if n.Id == "" {
		// return all news
		return db.GetAllNews(), nil
	}
	return getNews(n)
}

// getNews get News info
func getNews(n *pb.ID) (*pb.Posts, error) {
	ns, err := db.SelectNews(n)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

// SearchNewss select all Newss match the keywords,
// If no vs.Keywords provided, it'll return all the Newss in table of database.
func (fr *fetcherRepo) Search(keyword *pb.Keyword) (*pb.Posts, error) {
	return db.SearchNews(keyword)
}

//judege whether the news' publishdate <= d days
func valuableNews(ns *pb.Post, d int) bool {
	a := time.Now().AddDate(0, 0, -d).UnixNano()
	if ns.UpdateTime >= a {
		return true
	}
	return false
}

//data maintenance
func maintenanceDB(d int) error {
	//update newsCollection in cache，delete news > 5 days
	db.UpdateCache(d)
	//delete storaged links in file, the probability is 1/100
	a := rand.Intn(100)
	if a == 1 {
		if err := db.TruncLinks(); err != nil {
			return err
		}
	}
	return nil
}

//update news
func (fr *fetcherRepo) UpdateNews() error {
	// get news links
	links, err := fetcher.GetNewsLinks()
	if err != nil {
		return err
	}
	//filter links
	linksWithDate := []string{}
	curLinks, err := db.FilterCrawledLinks(links)
	if err != nil {
		return err
	}
	newsCollection := pb.Posts{}
	//fetcher news and storage
	for _, link := range *curLinks {
		post, err := fetcher.CrawlNews(link)
		if err == nil {
			ns := pb.Post{Id: fmt.Sprintf("%x", md5.Sum([]byte(post.URL.String()))), Title: post.Title, Content: post.Body, UpdateTime: post.Date, WebsiteId: fmt.Sprintf("%x", md5.Sum([]byte("sputniknews.cn"))), WebsiteTitle: "新西兰联合报"}
			//append links and date to varible linksWithDate
			linksWithDate = append(linksWithDate, link+"\t"+strconv.FormatInt(post.Date, 10))
			//judege whether the news' publishdate <= 3 days
			if valuableNews(&ns, 5) {
				newsCollection.Posts = append(newsCollection.Posts, &ns)
			}
		} else {
			continue
			// return err
		}
	}
	//maintenance data
	if err = maintenanceDB(5); err != nil {
		return err
	}
	//save NewsCollection and linksWithDate
	if len(newsCollection.Posts) > 0 {
		//save NewsCollection
		if save_err := db.SaveNewsCollection(&newsCollection); save_err != nil {
			return err
		}
		//save linksWithDate to file
		db.SaveNewsLinks(linksWithDate)
	}
	return nil
}
