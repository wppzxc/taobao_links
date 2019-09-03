package main

import (
	"fmt"
	. "github.com/lxn/walk/declarative"
	"github.com/wpp/taobao_links/pkg/dataoke"
	haodanku "github.com/wpp/taobao_links/pkg/haodanku"
	"github.com/wpp/taobao_links/pkg/taokeyi"
)

func main() {
	// init dataoke
	dtk := dataoke.GetDataokePage()
	// init haodanku
	hdk := haodanku.GetHaodankuPage()
	// init taokeyi
	tky := taokeyi.GetTaokeyiPage()
	if _, err := (MainWindow{
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
							DataSource: dtk,
							AutoSubmit: true,
							OnSubmitted: dtk.ResetPage,
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
										AssignTo: &dtk.GetBtn,
										Text:      "拉取",
										OnClicked: dtk.GetLinks,
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
										AssignTo: &dtk.Links,
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
										AssignTo: &dtk.StartPage,
										MaxSize: Size{20, 0},
										Text: "1",
										Enabled: Bind("Quan.Checked"),
									},
									TextLabel{
										Text: "结束页：",
									},
									LineEdit{
										AssignTo: &dtk.EndPage,
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
							DataSource: hdk,
							AutoSubmit: true,
						},
						Children: []Widget{
							Composite{
								MaxSize: Size{0, 50},
								Layout:  HBox{},
								Children: []Widget{
									HSpacer{},
									PushButton{
										AssignTo: &hdk.GetBtn,
										Text:      "拉取",
										OnClicked: hdk.GetLinks,
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
										AssignTo: &hdk.Links,
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
										AssignTo: &hdk.StartPage,
										MaxSize: Size{20, 0},
										Text: "1",
									},
									TextLabel{
										Text: "结束页：",
									},
									LineEdit{
										AssignTo: &hdk.EndPage,
										MaxSize: Size{20, 0},
										Text: "1",
									},
								},
							},
						},
					},
					tky,
				},
			},
		},
	}).Run(); err != nil {
		fmt.Println(err)
	}
}