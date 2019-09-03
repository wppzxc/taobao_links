package taokeyi

import (
	. "github.com/lxn/walk/declarative"
)

func GetTaokeyiPage() TabPage {
	return TabPage{
		Title:  "淘客易",
		Layout: VBox{},
		Children: []Widget{
			Composite{
				MaxSize: Size{0, 50},
				Layout:  HBox{},
				Children: []Widget{
					HSpacer{},
					RadioButtonGroup{
						Buttons: []RadioButton{{
							Text: "领券直播",
						}, {
							Text: "实时榜单",
						}},
					},
				},
			},
			Composite{
				Layout:             VBox{},
				AlwaysConsumeSpace: true,
				Children: []Widget{
					HSpacer{},
					TextLabel{
						Text: "淘宝链接：",
					},
					TextEdit{
						ReadOnly: true,
					},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text: "上一页",
					},
					LineEdit{
						MaxSize: Size{20, 0},
						Text:    "1",
					},
					PushButton{
						Text: "下一页",
					},
				},
			},
		},
	}
}
