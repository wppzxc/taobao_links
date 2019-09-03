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
	PageUpBtn   *walk.PushButton
	PageDownBtn *walk.PushButton
	Links       *walk.TextEdit
	Page        *walk.LineEdit
	PageNum     int
}

func GetDataokePage() *Dataoke {
	dataoke := &Dataoke{
		SubUrl:  "",
		Links:   new(walk.TextEdit),
		Page:    &walk.LineEdit{},
		PageNum: 3,
	}
	dataoke.Page.SetText("3")
	return dataoke
}

func (d *Dataoke) GetLinks() {
	d.SetBtnEnable(false)
	p, _ := strconv.Atoi(d.Page.Text())
	if p <= 0 {
		p = 1
	}
	d.PageNum = p
	switch d.SubUrl {
	case "top":
		fmt.Println("top")
		d.GetTopLinks()
	case "quan":
		fmt.Println("quan")
		d.GetQuanLinks()
	}
}

func (d *Dataoke) GetTopLinks() {
	url := dataokeHost + topUrl
	fmt.Println(url)
	go d.GetAndUpdateLinks(url)
}

func (d *Dataoke) GetQuanLinks() {
	url := dataokeHost + quanUrl + d.Page.Text()
	fmt.Println(url)
	go d.GetAndUpdateLinks(url)
}

func (d *Dataoke) GetAndUpdateLinks(url string) {
	defer d.SetBtnEnable(true)
	items, err := top.GetTopItems(url)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	links, err := top.GetTblinks(items)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	text := top.GetTextFromStrings(links)
	d.Links.SetText(text)
}

func (d *Dataoke) GetNextPageQuanLinks() {
	d.PageNum ++
	d.Page.SetText(strconv.Itoa(d.PageNum))
	d.GetQuanLinks()
}

func (d *Dataoke) GetBackPageQuanLinks() {
	if d.PageNum - 1 <= 1 {
		d.PageNum = 1
	} else {
		d.PageNum --
		d.Page.SetText(strconv.Itoa(d.PageNum))
	}
	d.GetQuanLinks()
}

func (d *Dataoke) SetBtnEnable(enable bool) {
	d.GetBtn.SetEnabled(enable)
	d.PageUpBtn.SetEnabled(enable)
	d.PageDownBtn.SetEnabled(enable)
}

func (d *Dataoke) ResetPage() {
	d.PageNum = 1
	d.Page.SetText("1")
}
