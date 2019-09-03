package top

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"time"
)

const (
	goodsListId = ".goods-list"
	topId       = "id"
	getTblUrl   = "http://www.dataoke.com/gettpl?gid=%s&_=%s"
)

func GetTopItems(url string) ([]string, error) {
	var items []string
	fmt.Println(url)
	dom, err := goquery.NewDocument(url)
	if err != nil {
		return nil, fmt.Errorf("Error in get html : %s ", err)
	}
	dom.Find(goodsListId).Children().Each(func(index int, content *goquery.Selection) {
		id, ok := content.Attr(topId)
		if ok {
			id = strings.Replace(id, "goods-items_", "", 1)
			items = append(items, id)
		}
	})
	return items, nil
}

func GetTblinks(items []string) ([]string, error) {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	var links []string
	for _, i := range items {
		l, err := GetTbl(i, now)
		if err != nil {
			continue
		}
		links = append(links, l)
	}
	if len(links) == 0 {
		return nil, fmt.Errorf("Error : the links is null ")
	}
	return links, nil
}

func GetTbl (id string, timestamp string) (string, error) {
	url := fmt.Sprintf(getTblUrl, id, timestamp)
	dom, err := goquery.NewDocument(url)
	if err != nil {
		return "", err
	}
	a := dom.Find("a")
	link := a.Nodes[1].Attr[1].Val
	return link, nil
}

func GetTextFromStrings(links []string) string {
	var text string
	for _, s := range links {
		text = text + "\n" + s
	}
	return text
}