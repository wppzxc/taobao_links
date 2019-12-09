package dataoke

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/wppzxc/taobao_links/pkg/features/dataoke/top"
	"strconv"
)

const (
	dataokeHost = "http://www.dataoke.com/"
	quanUrl     = "qlist/?px=sell&page="
	topUrl      = "top_sell"
)

type Dataoke struct {
	MainPage  *TabPage
	SubUrl    string
	GetBtn    *walk.PushButton
	LeiMu     int
	Links     *walk.TextEdit
	StartPage *walk.LineEdit
	EndPage   *walk.LineEdit
}

type LeiMu struct {
	Id   int
	Name string
}

func GetLeiMus() []LeiMu {
	return []LeiMu{
		{Name: "居家日用", Id: 4},
		{Name: "美食", Id: 6},
		{Name: "母婴", Id: 2},
		{Name: "美妆", Id: 3},
		{Name: "女装", Id: 1},
		{Name: "数码家电", Id: 8},
		{Name: "文娱车品", Id: 7},
		{Name: "内衣", Id: 10},
		{Name: "家装家纺", Id: 14},
		{Name: "鞋品", Id: 5},
		{Name: "男装", Id: 9},
		{Name: "配饰", Id: 12},
		{Name: "户外运动", Id: 13},
		{Name: "箱包", Id: 11},
	}
}

func GetDataokePage() *Dataoke {
	dataoke := &Dataoke{
		SubUrl:    "",
		Links:     &walk.TextEdit{},
		StartPage: &walk.LineEdit{},
		EndPage:   &walk.LineEdit{},
	}
	dataoke.MainPage = &TabPage{
		Title:  "大淘客",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource:  dataoke,
			AutoSubmit:  true,
			OnSubmitted: dataoke.ResetPage,
		},
		Children: []Widget{
			Composite{
				MaxSize: Size{0, 50},
				Layout:  HBox{},
				Children: []Widget{
					HSpacer{},
					Label{
						Text: "类目",
					},
					ComboBox{
						Value:         Bind("LeiMu"),
						BindingMember: "Id",
						DisplayMember: "Name",
						Model:         GetLeiMus(),
						Enabled:       Bind("Quan.Checked"),
					},
					RadioButtonGroup{
						DataMember: "SubUrl",
						Buttons: []RadioButton{{
							Name:  "Quan",
							Text:  "领券直播",
							Value: "quan",
						}, {
							Name:  "Top",
							Text:  "实时榜单",
							Value: "top",
						}},
					},
					PushButton{
						AssignTo:  &dataoke.GetBtn,
						Text:      "拉取",
						OnClicked: dataoke.GetLinks,
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					TextLabel{
						Text: "淘宝链接：",
					},
					TextEdit{
						ReadOnly: true,
						AssignTo: &dataoke.Links,
						VScroll:  true,
					},
				},
			},
			Composite{
				Layout: HBox{Margins: Margins{}},
				Children: []Widget{
					HSpacer{},
					TextLabel{
						Text: "起始页：",
					},
					LineEdit{
						AssignTo: &dataoke.StartPage,
						MaxSize:  Size{20, 0},
						Text:     "1",
						Enabled:  Bind("Quan.Checked"),
					},
					TextLabel{
						Text: "结束页：",
					},
					LineEdit{
						AssignTo: &dataoke.EndPage,
						MaxSize:  Size{20, 0},
						Text:     "1",
						Enabled:  Bind("Quan.Checked"),
					},
				},
			},
		},
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
	lm := ""
	lmUrl := ""
	if d.LeiMu > 0 {
		lm = strconv.Itoa(d.LeiMu)
	}
	if len(lm) > 0 {
		lmUrl = "&cid=" + lm
	}
	for i := start; i <= end; i++ {
		url := dataokeHost + quanUrl + strconv.Itoa(i) + lmUrl
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


func (d *Dataoke) GetMainPage() *TabPage {
	return d.MainPage
}
