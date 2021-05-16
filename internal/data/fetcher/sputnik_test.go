package fetcher

import (
	"fmt"
	"log"
	"testing"
)

func TestSetDate(t *testing.T) {
	p := NewPost("http://sputniknews.cn/politics/202105141033691405/")
	p.PostInit()
	err := setDate(p)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(p.Date)
}

func TestSetTitle(t *testing.T) {
	p := NewPost("http://sputniknews.cn/politics/202105141033691405/")
	p.PostInit()
	err := setTitle(p)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(p.Title)
}

func TestSetBody(t *testing.T) {
	p := NewPost("http://sputniknews.cn/covid-2019/202105141033687799/")
	p.PostInit()
	setTitle(p)
	setDate(p)
	err := setBody(p)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(p.Body)
}
