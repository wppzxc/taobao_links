package taokeyi

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strconv"
)

type taokeyi struct {
	MainPage  *TabPage
	GetBtn    *walk.PushButton
	Links     *walk.TextEdit
	StartPage *walk.LineEdit
	EndPage   *walk.LineEdit
}

func GetTaokeyiPage() *taokeyi {
	taokeyi := &taokeyi{
		GetBtn:    &walk.PushButton{},
		Links:     &walk.TextEdit{},
		StartPage: &walk.LineEdit{},
		EndPage:   &walk.LineEdit{},
	}
	taokeyi.MainPage = &TabPage{
		Title:  "淘客易",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: taokeyi,
			AutoSubmit: true,
		},
		Children: []Widget{
			Composite{
				MaxSize: Size{0, 50},
				Layout:  HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo:  &taokeyi.GetBtn,
						Text:      "拉取",
						OnClicked: taokeyi.GetLinks,
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
						AssignTo: &taokeyi.Links,
						VScroll:  true,
					},
				},
			},
			Composite{
				Layout: HBox{Margins: Margins{}},
				Children: []Widget{
					HSpacer{}, TextLabel{
						Text: "起始页：",
					},

					LineEdit{
						AssignTo: &taokeyi.StartPage,
						MaxSize:  Size{20, 0},
						Text:     "1",
					},
					TextLabel{
						Text: "结束页：",
					},
					LineEdit{
						AssignTo: &taokeyi.EndPage,
						MaxSize:  Size{20, 0},
						Text:     "1",
					},
				},
			},
		},
	}
	return taokeyi
}

func (t *taokeyi) GetLinks() {
	t.SetUIEnable(false)
	startPage, _ := strconv.Atoi(t.StartPage.Text())
	endPage, _ := strconv.Atoi(t.EndPage.Text())
	if startPage <= 1 {
		startPage = 1
	}
	if endPage <= startPage {
		endPage = startPage
	}
	t.ResetPage(startPage, endPage)
	fmt.Println("taokeyi")
	fmt.Println(startPage, endPage)
	go func() {
		text := t.GetMutiPagesLinks(startPage, endPage)
		t.UpdateLinks(text)
	}()
}

func (t *taokeyi) SetUIEnable(enable bool) {
	t.GetBtn.SetEnabled(enable)
	t.StartPage.SetEnabled(enable)
	t.EndPage.SetEnabled(enable)
}

func (t *taokeyi) ResetPage(startPage int, endPage int) {
	s := strconv.Itoa(startPage)
	e := strconv.Itoa(endPage)
	t.StartPage.SetText(s)
	t.EndPage.SetText(e)
}

func (t *taokeyi) GetMutiPagesLinks(startPage int, endPage int) string {
	text := ""
	// TODO: get links

	return text
}

func (t *taokeyi) UpdateLinks(text string) {
	defer t.SetUIEnable(true)
	t.Links.SetText(text)
}
