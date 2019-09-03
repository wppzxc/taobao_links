package main

import (
	. "github.com/lxn/walk/declarative"
	"github.com/wpp/taobao_links/pkg/dataoke"
	haodanku "github.com/wpp/taobao_links/pkg/haodanku"
	"github.com/wpp/taobao_links/pkg/taokeyi"
)

func main() {
	dataoke := dataoke.GetDataokePage()
	haodanku := haodanku.GetHaodankuPage()
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
									HSpacer{},TextLabel{
										Text: "起始页：",
									},
									
									LineEdit{
										AssignTo: &dataoke.StartPage,
										MaxSize: Size{20, 0},
										Text: "1",
										Enabled: Bind("Quan.Checked"),
									},
									TextLabel{
										Text: "结束页：",
									},
									LineEdit{
										AssignTo: &dataoke.EndPage,
										MaxSize: Size{20, 0},
										Text: "1",
										Enabled: Bind("Quan.Checked"),
									},
								},
							},
						},
					},
					{
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
					},
					taokeyi,
				},
			},
		},
	}.Run()
}