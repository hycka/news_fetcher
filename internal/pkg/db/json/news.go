package json

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hi20160616/gears"
	pb "github.com/hycka/news_fetcher/api/news_fetcher/api"
)

var (
	// currentNews []*pb.News
	// newsPath  string = "../../../internal/pkg/db/json/sputnik.json"      //for test
	// linksPath string = "../../../internal/pkg/db/json/sputnik_links.txt" //for test
	// newsPath  string = "../../internal/pkg/db/json/sputnik.json"         //for data/news.go test
	// newsPath string = "sputnik.json" //for json/news.go test
	newsPath    string = "./internal/pkg/db/json/sputnik.json"      //for run
	linksPath   string = "./internal/pkg/db/json/sputnik_links.txt" //for run
	currentNews        = make(map[string]*pb.Post)
)

//Load json file to global varible 'currentNews'
func LoadLocalNews() error {
	if gears.Exists(newsPath) {
		b, err := os.ReadFile(newsPath)
		if err != nil {
			return err
		}
		if len(b) == 0 {
			return nil
		}
		if err = json.Unmarshal(b, &currentNews); err != nil {
			return err
		}
	}
	return nil
}

//return all news in cache
func GetAllNews() *pb.Posts {
	return &pb.Posts{Posts: getValues(currentNews)}
}

// judge whether News exist
func NewsExist(id string) bool {
	for _, news := range currentNews {
		if news.Id == id {
			return true
		}
	}
	return false
}

// Search just search single keyword contained in title or content
func SearchNews(keywords *pb.Keyword) (*pb.Posts, error) {
	retPost := pb.Posts{}
	keys := strings.Split(keywords.Keyword, ",")
	for _, keyword := range keys {
		for _, v := range currentNews {
			if strings.Contains(v.Title, keyword) || strings.Contains(v.Content, keyword) {
				retPost.Posts = append(retPost.Posts, v)
			}
		}
	}
	return &retPost, nil
}

//select news
func SelectNews(n *pb.ID) (*pb.Posts, error) {
	retPost := pb.Posts{}
	newsIds := strings.Split(n.Id, ",")
	for _, id := range newsIds {
		// if exists, return ns
		if _, ok := currentNews[id]; ok {
			retPost.Posts = append(retPost.Posts, currentNews[id])
		}
	}
	return &retPost, nil
}

//update currentNews in cache, delete news > 5 days
func UpdateCache(d int) {
	a := time.Now().AddDate(0, 0, -d).UnixNano()
	for k, v := range currentNews {
		if v.UpdateTime < a {
			delete(currentNews, k)
		}
	}
}

//update links file
func TruncLinks() error {
	np, err := os.OpenFile(linksPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer np.Close()
	return nil
}

//get map values
func getValues(m map[string]*pb.Post) []*pb.Post {
	values := make([]*pb.Post, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

//save all news
func SaveNewsCollection(ns *pb.Posts) error {
	// save newscolletion to varible currentNews
	for _, news := range ns.Posts {
		currentNews[news.Id] = news
	}

	asJson, err := json.Marshal(currentNews)
	if err != nil {
		return err
	}
	// storage to newsPath
	np, err := os.OpenFile(newsPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer np.Close()
	_, np_err := np.Write(asJson)
	if np_err != nil {
		return np_err
	}

	return nil
}

//read already crawled news's links
func readExistLinks() (map[string]int64, error) {
	// existLinks := []string{}
	existLinks := make(map[string]int64) // New empty set
	// if linksPath not exist, return empty map
	if !gears.Exists(linksPath) {
		return existLinks, nil
	}
	lp, err := os.Open(linksPath)
	if err != nil {
		return nil, err
	}
	defer lp.Close()

	br := bufio.NewReader(lp)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		tmp := strings.Split(string(a), "\t")
		tmp2, _ := strconv.ParseInt(tmp[1], 10, 64)
		existLinks[string(tmp[0])] = tmp2
	}
	return existLinks, nil
}

//filter path and write crawled links to linksPath
func FilterCrawledLinks(currentlinks *[]string) (*[]string, error) {
	retLinks := []string{}
	existLinks, err := readExistLinks()
	if err != nil {
		return nil, err
	}
	for _, link := range *currentlinks {
		if _, exist := existLinks[link]; !exist {
			//append link to retLinks
			retLinks = append(retLinks, link)
		} else { //if exist, judge the published date, crawl if publish_date < 30min
			m, _ := time.ParseDuration("-0.5h")
			a := time.Now().UTC().Add(m).UnixNano()
			if existLinks[link] > a {
				retLinks = append(retLinks, link)
			}
		}
	}
	return &retLinks, nil
}

//storage Links
func SaveNewsLinks(links []string) error {
	np, err := os.OpenFile(linksPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer np.Close()

	for _, link := range links { // Loop
		_, np_err := np.WriteString(strings.TrimSpace(link) + "\n")
		if np_err != nil {
			return np_err
		}
	}
	return nil
}
