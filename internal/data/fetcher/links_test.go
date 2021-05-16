package fetcher

import (
	"fmt"
	"net/url"
	"testing"
)

func TestGetLinksFromRaw(t *testing.T) {
	u, err := url.Parse("https://cn.reuters.com/assets/jsonWireNews?limit=100")
	if err != nil {
		t.Errorf("Url Parse fail!\n%s", err)
	}
	var f = &Fetcher{
		Entrance: u,
	}
	f.GetLinksFromRaw()
	fmt.Println(f.Links)
}

func TestKickOutLinksMatchPath(t *testing.T) {
	// link := "https://www.dwnews.com/%E8%A7%86%E8%A7%89/60202427/"
	beKick := []string{"https://www.dwnews.com/%E8%A7%86%E8%A7%89/60202427/"}
	// path := "/" + url.QueryEscape("视觉") + "/"
	path := url.QueryEscape("视觉")
	KickOutLinksMatchPath(&beKick, path)
	if len(beKick) != 0 {
		t.Errorf("want: len(beKick) == 0, got: len(beKick) == %d", len(beKick))
	}
}

func TestGetLinks(t *testing.T) {
	u, err := url.Parse("http://sputniknews.cn/archive/")
	assertLinks := []string{
		"/article/asia-financial-markets-0121-thur-idCNKBS29Q0CM",
	}
	// u, err := url.Parse("https://news.ltn.com.tw/list/breakingnews")
	// assertLinks := []string{
	//         "https://news.ltn.com.tw/news/society/breakingnews/3278253",
	//         "https://news.ltn.com.tw/news/society/breakingnews/3278250",
	//         "https://news.ltn.com.tw/news/politics/breakingnews/3278225",
	//         "https://news.ltn.com.tw/news/politics/breakingnews/3278170",
	// }
	// u, err := url.Parse("https://www.cna.com.tw/list/aall.aspx")
	// assertLinks := []string{
	//         "https://www.cna.com.tw/news/aopl/202009290075.aspx",
	//         "https://www.cna.com.tw/news/firstnews/202009290051.aspx",
	//         "https://www.cna.com.tw/news/acn/202009290063.aspx",
	//         "https://www.cna.com.tw/news/aipl/202009290055.aspx",
	// }
	// u, err := url.Parse("https://www.bbc.com/zhongwen/simp/topics/ck2l9z0em07t")
	// assertLinks := []string{
	//         "/zhongwen/simp/world-55655858",
	//         "/zhongwen/simp/world-55653976",
	//         "/zhongwen/simp/science-55632120",
	//         "/zhongwen/simp/chinese-news-55635500",
	// }
	// u, err := url.Parse("https://chinese.aljazeera.net/news")
	// assertLinks := []string{
	//         "https://chinese.aljazeera.net/news/2021/1/20/%e9%87%8d%e8%bf%94%e6%a0%b8%e5%8d%8f%e8%ae%ae%e5%b9%b6%e6%94%af%e6%8c%81%e4%b8%a4%e5%9b%bd%e6%96%b9%e6%a1%88%e5%b8%83%e6%9e%97%e8%82%af%e6%a6%82%e8%bf%b0%e6%8b%9c%e7%99%bb%e6%94%bf%e5%ba%9c%e7%9a%84",
	// }
	if err != nil {
		t.Errorf("Url Parse fail!\n%s", err)
	}
	var f = &Fetcher{
		Entrance: u,
	}
	f.GetLinks()
	shot := 0
	for _, link := range f.Links {
		for _, v := range assertLinks {
			if link == v {
				shot++
			}
		}
	}
	if shot != len(assertLinks) {
		t.Errorf("want: %v, got: %v", len(assertLinks), shot)
	}
}
