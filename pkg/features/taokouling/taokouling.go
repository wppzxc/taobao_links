package taokouling

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/mvdan/xurls"
	"github.com/wppzxc/taobao_links/pkg/features/taokouling/app"
	"strings"
)

type taokouling struct {
	MainPage *TabPage
	GetBtn   *walk.PushButton
	Inputs   *walk.TextEdit
	Links    *walk.TextEdit
}

func GetTaokoulingPage() *taokouling {
	taokouling := &taokouling{
		GetBtn: &walk.PushButton{},
		Inputs: &walk.TextEdit{},
		Links:  &walk.TextEdit{},
	}
	taokouling.MainPage = &TabPage{
		Title:  "淘口令提取",
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: taokouling,
			AutoSubmit: true,
		},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					TextLabel{
						Text: "输入：",
					},
					TextEdit{
						AssignTo: &taokouling.Inputs,
						VScroll:  true,
					},
				},
			},
			Composite{
				MaxSize: Size{0, 50},
				Layout:  HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo:  &taokouling.GetBtn,
						Text:      "解析",
						OnClicked: taokouling.GetLinks,
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					HSpacer{},
					TextLabel{
						Text: "淘口令：",
					},
					TextEdit{
						AssignTo: &taokouling.Links,
						VScroll:  true,
						ReadOnly: true,
					},
				},
			},
		},
	}
	return taokouling
}

func (t *taokouling) GetLinks() {
	inputs := t.Inputs.Text()
	t.SetUIEnable(false)
	fmt.Printf("input is '%s' ", inputs)
	if len(inputs) == 0 {
		fmt.Println("must provide inputs !")
		return
	}
	go func() {
		links := t.ResolveInputs(inputs)
		t.UpdateLinks(links)
	}()

}

func (t *taokouling) ResolveInputs(inputs string) string {
	fmt.Print(inputs)
	outputs := inputs
	sLinks := xurls.Relaxed.FindAllString(inputs, -1)
	for _, l := range sLinks {
		outputs = strings.ReplaceAll(outputs, l, "%s")
		kl := app.GetKouling(l)
		outputs = fmt.Sprintf(outputs, kl)
	}
	return outputs
}

func (t *taokouling) UpdateLinks(links string) {
	defer t.SetUIEnable(true)
	t.Links.SetText(links)
}

func (t *taokouling) SetUIEnable(enable bool) {
	t.Inputs.SetEnabled(enable)
	t.GetBtn.SetEnabled(enable)
	t.Links.SetEnabled(enable)
}

func (t *taokouling) GetMainPage() *TabPage {
	return t.MainPage
}
