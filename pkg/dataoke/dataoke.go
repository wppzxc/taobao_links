package dataoke

import (
	"fmt"
	"github.com/lxn/walk"
	"github.com/wpp/taobao_links/pkg/dataoke/top"
	"strconv"
)

const (
	dataokeHost = "http://www.dataoke.com/"
	quanUrl     = "qlist/?px=sell&page="
	topUrl      = "top_sell"
)

type Dataoke struct {
	SubUrl      string
	GetBtn      *walk.PushButton
	Links       *walk.TextEdit
	StartPage   *walk.LineEdit
	EndPage     *walk.LineEdit
}

func GetDataokePage() *Dataoke {
	dataoke := &Dataoke{
		SubUrl:    "",
		Links:     &walk.TextEdit{},
		StartPage: &walk.LineEdit{},
		EndPage:   &walk.LineEdit{},
	}
	return dataoke
}

func (d *Dataoke) GetLinks() {
	d.SetUIEnable(false)
	switch d.SubUrl {
	case "top":
		fmt.Println("top")
		d.GetTopLinks()
	case "quan":
		startPage, _ := strconv.Atoi(d.StartPage.Text())
		endPage, _ := strconv.Atoi(d.EndPage.Text())
		if endPage <= startPage {
			startPage = endPage
		}
		fmt.Println("quan")
		fmt.Println(startPage, endPage)
		d.GetQuanLinks(startPage, endPage)
	}
}

func (d *Dataoke) GetTopLinks() {
	url := dataokeHost + topUrl
	fmt.Println(url)
	go func() {
		text := d.GetLinksText(url)
		d.UpdateLinks(text)
	}()
}

func (d *Dataoke) GetQuanLinks(start int, end int) {
	go func() {
		text := d.GetMutiPagesQuanLinks(start, end)
		d.UpdateLinks(text)
	}()
}

func (d *Dataoke) GetMutiPagesQuanLinks(start int, end int) string {
	text := ""
	for i := start ; i <= end ; i++ {
		url := dataokeHost + quanUrl + strconv.Itoa(i)
		fmt.Println(url)
		tmp := d.GetLinksText(url)
		text = text + "\n" + tmp
	}
	return text
}

func (d *Dataoke) GetLinksText(url string) string {
	items, err := top.GetTopItems(url)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	links, err := top.GetTblinks(items)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	return top.GetTextFromStrings(links)
}

func (d *Dataoke) UpdateLinks(text string) {
	defer d.SetUIEnable(true)
	d.Links.SetText(text)
}

func (d *Dataoke) SetUIEnable(enable bool) {
	d.GetBtn.SetEnabled(enable)
	d.StartPage.SetEnabled(enable)
	d.EndPage.SetEnabled(enable)
}

func (d *Dataoke) ResetPage() {
	d.StartPage.SetText("1")
	d.EndPage.SetText("1")
}
