package fetcher

import (
	"fmt"
	"net/url"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
	"golang.org/x/net/html"
)

type Post struct {
	Domain   string
	URL      *url.URL
	DOC      *html.Node
	Raw      []byte
	Title    string
	Body     string
	Date     int64
	Filename string
	Err      error
}

type Paragraph struct {
	Type    string
	Content string
}

func NewPost(rawurl string) *Post {
	p := &Post{}
	p.URL, p.Err = url.Parse(rawurl)
	p.Domain = p.URL.Hostname()
	return p
}

// TODO: use func init
// PostInit open url and get raw and doc
func (p *Post) PostInit() error {
	if p.Err != nil {
		return p.Err
	}
	p.Raw, p.DOC, p.Err = htmldoc.GetRawAndDoc(p.URL, 1*time.Minute)
	return p.Err
}

// RoutePost will switch post to the right dealer.
func (p *Post) RoutePost() error {
	if p.Err != nil {
		return p.Err
	}
	switch p.Domain {
	case "sputniknews.cn":
		post := Post(*p)
		p.Err = SetPost(&post)
		*p = Post(post)
	default:
		return fmt.Errorf("switch no case on: %s", p.Domain)
	}
	return p.Err
}

// TreatPost get post things and set to `p` then save it.
func (p *Post) TreatPost() error {
	// Post prepare
	if p.Err = p.PostInit(); p.Err != nil {
		return p.Err
	}
	if p.Err = p.RoutePost(); p.Err != nil {
		return p.Err
	}

	// Post storage
	// if p.Err = p.setFilename(); p.Err != nil {
	// 	return p.Err
	// }
	//p.Err = p.savePost()

	return p.Err
}
