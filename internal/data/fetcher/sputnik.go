package fetcher

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
	"github.com/hi20160616/gears"
)

func SetPost(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	p.Err = setDate(p)
	p.Err = setTitle(p)
	p.Err = setBody(p)
	// p.Err = transform(p)
	return p.Err
}

func setDate(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	if p.DOC == nil {
		return fmt.Errorf("p.DOC is nil")
	}
	metas := htmldoc.MetasByName(p.DOC, "pubdate")
	cs := []string{}
	for _, meta := range metas {
		for _, a := range meta.Attr {
			if a.Key == "content" {
				cs = append(cs, a.Val)
			}
		}
	}
	if len(cs) <= 0 {
		return fmt.Errorf("setData got nothing.")
	}
	t, err := time.Parse(time.RFC3339, cs[0])
	t, _ = reduce8Hour(t)
	if err != nil {
		return fmt.Errorf("setData got nothing.")
	}
	p.Date = t.UnixNano()
	return nil
}

//UTC + 8H
func reduce8Hour(t time.Time) (time.Time, error) {
	h, _ := time.ParseDuration("-1h")
	h1 := t.Add(8 * h)
	return h1, nil
}

func setTitle(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	if p.DOC == nil {
		return fmt.Errorf("p.DOC is nil")
	}
	doc := htmldoc.ElementsByTag(p.DOC, "title")
	if doc == nil {
		return fmt.Errorf("there is no element <title>")
	}
	title := doc[0].FirstChild.Data
	title = strings.TrimSpace(title)
	gears.ReplaceIllegalChar(&title)
	p.Title = title
	return nil
}

func setBody(p *Post) error {
	if p.Err != nil {
		return p.Err
	}
	if p.DOC == nil {
		return fmt.Errorf("p.DOC is nil")
	}
	b, err := sputnik(p)
	if err != nil {
		return err
	}
	t := time.Unix(p.Date/1e9, 0)
	if err != nil {
		return err
	}
	h1 := fmt.Sprintf("# [%02d.%02d][%02d%02dH] %s", t.Month(), t.Day(), t.Hour(), t.Minute(), p.Title)
	p.Body = h1 + "\n\n" + b + "\n\n原地址：" + p.URL.String()
	return nil
}

func sputnik(p *Post) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	if p.Raw == nil {
		return "", fmt.Errorf("p.Raw is nil")
	}
	raw := p.Raw
	//crawl article lead
	r := htmldoc.DivWithAttr2(raw, "class", "b-article__lead")
	ps := [][]byte{}
	b := bytes.Buffer{}
	re := regexp.MustCompile(`(?s)<p.*?>(.*?)</p>`)
	for _, v := range re.FindAllSubmatch(r, -1) {
		ps = append(ps, v[1])
	}
	//crawl article
	// r = htmldoc.DivWithAttr2(raw, "class", "b-article__text")
	// b = bytes.Buffer{}
	// re = regexp.MustCompile(`(?s)<p.*?>(.*?)</p>`)
	// for _, v := range re.FindAllSubmatch(r, -1) {
	// 	ps = append(ps, v[1])
	// }

	if len(ps) == 0 {
		return "", fmt.Errorf("no <p> matched")
	}
	for _, p := range ps {
		b.Write(p)
		b.Write([]byte("  \n"))
	}
	body := b.String()
	re = regexp.MustCompile(`「`)
	body = re.ReplaceAllString(body, "“")
	re = regexp.MustCompile(`」`)
	body = re.ReplaceAllString(body, "”")
	re = regexp.MustCompile(`<a.*?>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`</a>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`<script.*?</script>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`<blockquote.*?</blockquote>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`<iframe.*?</iframe>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`<strong.*?</strong>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`(?s)<span.*?</span>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`<img.*?>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`<div.*?</div>`)
	body = re.ReplaceAllString(body, "")
	re = regexp.MustCompile(`</div>`)
	body = re.ReplaceAllString(body, "")
	return body, nil
}

//transform HANZI
// func transform(p *Post) error {
// 	tw2s, err := gocc.New("hk2s")
// 	if err != nil {
// 		p.Err = err
// 		return err
// 	}
// 	//transform title
// 	in := p.Title
// 	out, err := tw2s.Convert(in)
// 	if err != nil {
// 		p.Err = err
// 		return err
// 	}
// 	p.Title = out
// 	//transform body
// 	in = p.Body
// 	out, err = tw2s.Convert(in)
// 	if err != nil {
// 		p.Err = err
// 		return err
// 	}
// 	p.Body = out
// 	return p.Err
// }
