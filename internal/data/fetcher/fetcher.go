package fetcher

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
)

var worklist = []string{
	"http://sputniknews.cn/archive/",
}

type Fetcher struct {
	Entrance *url.URL
	Links    []string
	Err      error
}

func NewFetcher(site string) *Fetcher {
	u, err := url.Parse(site)
	if err != nil {
		log.Printf("url parse err: %s", err)
	}
	return &Fetcher{
		Entrance: u,
	}
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once  for each item.
// breadthFirst(crawl, os.Args[1:])
func BreadthFirst(f func(item string), worklist []string) {
	for _, item := range worklist {
		f(item)
	}
}

//get all news links
func GetNewsLinks() (*[]string, error) {
	allLinks := []string{}
	for _, item := range worklist {
		tmp, err := extractLinks(item)
		if err != nil {
			log.Println(err)
			return nil, err
		} else {
			for _, link := range tmp {
				allLinks = append(allLinks, link)
			}
		}
	}
	return &allLinks, nil
}

// extract news links
func extractLinks(_url string) ([]string, error) {
	defer func() {
		if err := recover(); err != nil {
			e := err.(error)
			log.Println(e)
			PanicLog(e)
		}
	}()
	f := NewFetcher(_url)
	log.Printf("[*] Deal with: [%s]\n", _url)
	log.Println("[*] Fetch links ...")
	if err := f.GetLinks(); err != nil { // f.Links update to the _url website is.
		log.Println(err)
		// if links cannot fetch sleep 1 minute then continue
		time.Sleep(1 * time.Minute)
		// continue // only useful by goroutine
		return nil, err
	}
	return f.Links, nil
}

//fecher news
func CrawlNews(link string) (*Post, error) {
	post := NewPost(link)
	fmt.Println("crawl url:" + link)
	if err := post.TreatPost(); err != nil {
		errMsg := "TreatPost error occur on: " + link
		log.Println(errMsg)
		log.Println(err)
		return nil, err
	} else {
		return post, nil
	}
}

func PanicLog(_err error) error {
	filePath := "./PanicLog.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString("[" + time.Now().Format(time.RFC3339) + "]--------------------------------------\n")
	write.WriteString(_err.Error() + "\n")
	write.Flush()
	return nil
}
