package main

import (
	. "github.com/lxn/walk/declarative"
	"github.com/wpp/taobao_links/pkg/dataoke"
	"github.com/wpp/taobao_links/pkg/taokeyi"
)

func main() {
	dataoke := dataoke.GetDataokePage()
	taokeyi := taokeyi.GetTaokeyiPage()
	MainWindow{
		Title:   "淘宝链接获取工具",
		Size: Size{400, 600},
		Layout:  VBox{},
		Children: []Widget{
			TabWidget{
				Pages: []TabPage{
					{
						Title:  "大淘客",
						Layout: VBox{},
						DataBinder: DataBinder{
							DataSource: dataoke,
							AutoSubmit: true,
							OnSubmitted: dataoke.ResetPage,
						},
						Children: []Widget{
							Composite{
								MaxSize: Size{0, 50},
								Layout:  HBox{},
								Children: []Widget{
									HSpacer{},
									RadioButtonGroup{
										DataMember: "SubUrl",
										Buttons: []RadioButton{{
											Name: "Quan",
											Text:  "领券直播",
											Value: "quan",
										}, {
											Name: "Top",
											Text:  "实时榜单",
											Value: "top",
										}},
									},
									PushButton{
										AssignTo: &dataoke.GetBtn,
										Text:      "拉取",
										OnClicked: dataoke.GetLinks,
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
										AssignTo: &dataoke.Links,
										VScroll: true,
									},
								},
							},
							Composite{
								Layout: HBox{Margins: Margins{}},
								Children: []Widget{
									HSpacer{},
									PushButton{
										Text: "上一页",
										AssignTo: &dataoke.PageUpBtn,
										OnClicked: dataoke.GetBackPageQuanLinks,
										Enabled: Bind("Quan.Checked"),
									},
									LineEdit{
										AssignTo: &dataoke.Page,
										MaxSize: Size{20, 0},
										Text: "1",
									},
									PushButton{
										Text: "下一页",
										AssignTo: &dataoke.PageDownBtn,
										OnClicked: dataoke.GetNextPageQuanLinks,
										Enabled: Bind("Quan.Checked"),
									},
								},
							},
						},
					},
					taokeyi,
				},
			},
		},
	}.Run()
}