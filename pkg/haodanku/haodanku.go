package haodanku

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/wppzxc/taobao_links/pkg/haodanku/app"
	"strconv"
)

const (
	haodankuHost = "https://www.haodanku.com/"
	haodankuUrl  = "indexapi/get_allitem_list?sort=3&p="
)

type haodanku struct {
	MainPage  *TabPage
	GetBtn    *walk.PushButton
	Links     *walk.TextEdit
	StartPage *walk.LineEdit
	EndPage   *walk.LineEdit
}

func GetHaodankuPage() *haodanku {
	haodanku := &haodanku{
		GetBtn:    &walk.PushButton{},
		Links:     &walk.TextEdit{},
		StartPage: &walk.LineEdit{},
		EndPage:   &walk.LineEdit{},
	}
	haodanku.MainPage = &TabPage{
		Title: "好单库",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: haodanku,
			AutoSubmit: true,
		},
		Children: []Widget{
			Composite{
				MaxSize: Size{0, 50},
				Layout:  HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &haodanku.GetBtn,
						Text:      "拉取",
						OnClicked: haodanku.GetLinks,
					},
				},
			},
			Composite{
				Layout:             VBox{},
				Children: []Widget{
					HSpacer{},
					TextLabel{
						Text: "淘宝链接：",
					},
					TextEdit{
						ReadOnly: true,
						AssignTo: &haodanku.Links,
						VScroll: true,
					},
				},
			},
			Composite{
				Layout: HBox{Margins: Margins{}},
				Children: []Widget{
					HSpacer{},TextLabel{
						Text: "起始页：",
					},
					
					LineEdit{
						AssignTo: &haodanku.StartPage,
						MaxSize: Size{20, 0},
						Text: "1",
					},
					TextLabel{
						Text: "结束页：",
					},
					LineEdit{
						AssignTo: &haodanku.EndPage,
						MaxSize: Size{20, 0},
						Text: "1",
					},
				},
			},
		},
	}
	return haodanku
}

func (h *haodanku) GetLinks() {
	h.SetUIEnable(false)
	startPage, _ := strconv.Atoi(h.StartPage.Text())
	endPage, _ := strconv.Atoi(h.EndPage.Text())
	if startPage <= 1 {
		startPage = 1
	}
	if endPage <= startPage {
		endPage = startPage
	}
	h.ResetPage(startPage, endPage)
	fmt.Println("haodanku")
	fmt.Println(startPage, endPage)
	go func() {
		text := h.GetMutiPagesLinks(startPage, endPage)
		h.UpdateLinks(text)
	}()
}

func (h *haodanku) GetMutiPagesLinks(start int, end int) string {
	text := ""
	for i := start; i <= end; i++ {
		url := haodankuHost + haodankuUrl + strconv.Itoa(i)
		tmp := h.GetLinksText(url)
		text = text + "\n" + tmp
	}
	return text
}

func (h *haodanku) GetLinksText(url string) string {
	links, err := app.GetLinks(url)
	if err != nil {
		fmt.Println(err)
	}
	return GetTextFromLinks(links)
}

func GetTextFromLinks(links []string) string {
	text := ""
	for _, i := range links {
		text = text + "\n" + i
	}
	return text
}

func (h *haodanku) UpdateLinks(text string) {
	defer h.SetUIEnable(true)
	h.Links.SetText(text)
}

func (h *haodanku) SetUIEnable(enable bool) {
	h.GetBtn.SetEnabled(enable)
	h.StartPage.SetEnabled(enable)
	h.EndPage.SetEnabled(enable)
}

func (h *haodanku) ResetPage(startPage int, endPage int) {
	s := strconv.Itoa(startPage)
	e := strconv.Itoa(endPage)
	h.StartPage.SetText(s)
	h.EndPage.SetText(e)
}
