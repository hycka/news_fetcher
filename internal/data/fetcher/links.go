// Pacage links provides a link-extraction fuction.
package fetcher

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	htmldoc "github.com/hi20160616/exhtml"
	"github.com/hi20160616/gears"
	"github.com/pkg/errors"
)

// GetJsonLinks get links from a url that return json data.
func (f *Fetcher) GetLinksFromRaw() error {
	raw, _, err := htmldoc.GetRawAndDoc(f.Entrance, 1*time.Minute)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`"url":\s"(.*?)",`)
	rs := re.FindAllStringSubmatch(string(raw), -1)
	for _, item := range rs {
		f.Links = append(f.Links, "https://"+f.Entrance.Hostname()+item[1])
	}
	return nil
}

func (f *Fetcher) GetLinksFromNode() error {
	if f.Err != nil {
		return f.Err
	}
	if links, err := htmldoc.ExtractLinks(f.Entrance.String()); err != nil {
		f.Err = errors.WithMessage(err, "cannot extract links from "+f.Entrance.String())
	} else {
		f.Links = gears.StrSliceDeDupl(links)
	}
	return f.Err
}

// LinksInit fetch all links from entrance of f, save to f
func (f *Fetcher) GetLinks() error {
	hostname := f.Entrance.Hostname()
	switch hostname {
	case "www.dwnews.com":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		f.Links = LinksFilter(f.Links, `.*?/.*?/\d{8}/`)
		KickOutLinksMatchPath(&f.Links, "zone")
		KickOutLinksMatchPath(&f.Links, "/"+url.QueryEscape("视觉")+"/")
	case "www.voachinese.com":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		l1 := LinksFilter(f.Links, `.*?/a/\d*?.html`)
		l2 := LinksFilter(f.Links, `.*?/a/.*-.*.html`)
		f.Links = append(l1, l2...)
		KickOutLinksMatchPath(&f.Links, "voaweishi")
	case "www.zaobao.com":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		newsWorld := LinksFilter(f.Links, `.*?/news/world/.*`)
		newsChina := LinksFilter(f.Links, `.*?/news/china/.*`)
		realtimeWorld := LinksFilter(f.Links, `.*?/realtime/world/.*`)
		realtimeChina := LinksFilter(f.Links, `.*?/realtime/china/.*`)
		f.Links = append(append(append(newsWorld, newsChina...), realtimeWorld...), realtimeChina...)
	case "news.ltn.com.tw":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		f.Links = LinksFilter(f.Links, `https://news.*/news/.*`)
		KickOutLinksMatchPath(&f.Links, "/life/")
		KickOutLinksMatchPath(&f.Links, "/society/")
		KickOutLinksMatchPath(&f.Links, "/novelty/")
		KickOutLinksMatchPath(&f.Links, "/local/")
	case "www.cna.com.tw":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		newsFirst := LinksFilter(f.Links, `.*?/news/firstnews/.*`)
		newsWorld := LinksFilter(f.Links, `.*?/news/aopl/.*`)
		// TODO: ignore aipl and acn but this still fetch links?
		newsPolitical := LinksFilter(f.Links, `.*?/news/aipl/.*`)
		newsTW := LinksFilter(f.Links, `.*?/news/acn/.*`)
		f.Links = append(append(append(newsFirst, newsWorld...), newsPolitical...), newsTW...)
	case "www.bbc.com":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		f.Links = LinksFilter(f.Links, `.*?/zhongwen/simp/.*-\d*`)
		KickOutLinksMatchPath(&f.Links, "institutional")
	case "chinese.aljazeera.net":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		f.Links = LinksFilter(f.Links, `.*?\/[A-Za-z]+\/\d{4}\/\d{1,2}\/\d{1,2}\/.*`)
	case "cn.reuters.com":
		if err := f.GetLinksFromRaw(); err != nil {
			return err
		}
		f.Links = LinksFilter(f.Links, `.*?\/article/.*?-id\S*`)
	case "cn.kabar.kg":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		f.Links = LinksFilter(f.Links, `.*?\/news/.*?\/`)
	case "ucpnz.co.nz":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		f.Links = LinksFilter(f.Links, `.*?\/[0-9]{1,4}\/[0-9]{1,2}\/[0-9]{1,2}\/.*?\/`)
	case "www.dw.com":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		f.Links = LinksFilter(f.Links, `.*?\/zh\/.*`)
		KickOutLinksMatchPath(&f.Links, "s-")
		KickOutLinksMatchPath(&f.Links, "index-zh")
		KickOutLinksMatchPath(&f.Links, "av-")
	case "sputniknews.cn":
		if err := f.GetLinksFromNode(); err != nil {
			return err
		}
		f.Links = LinksFilter(f.Links, `http://sputniknews.cn/.*?\/[0-9]{10,25}\/`)
		KickOutLinksMatchPath(&f.Links, "entertainment")
	}
	return nil
}

// kickOutLinksMatchPath will kick out the links match the path,
func KickOutLinksMatchPath(links *[]string, path string) {
	tmp := []string{}
	// path = "/" + url.QueryEscape(path) + "/"
	// path = url.QueryEscape(path)
	for _, link := range *links {
		if !strings.Contains(link, path) {
			tmp = append(tmp, link)
		}
	}
	*links = tmp
}

// TODO: use point to impletement LinksFilter
// LinksFilter is support for SetLinks method
func LinksFilter(links []string, regex string) []string {
	flinks := []string{}
	re := regexp.MustCompile(regex)
	s := strings.Join(links, "\n")
	flinks = re.FindAllString(s, -1)
	return flinks
}
