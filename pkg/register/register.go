package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/wppzxc/taobao_links/pkg/license"
	"github.com/wppzxc/taobao_links/pkg/types"
	"strconv"
	"time"
)

type register struct {
	Selector    Selector
	Days        string
	Version     string
	RegisterBtn *walk.PushButton
	License     *walk.TextEdit
}

type Selector struct {
	Dataoke       bool
	Haodanku      bool
	Duoduojinbao  bool
	Taokeyi       bool
	GoodSearch    bool
	Yituike       bool
	PddUserNumber bool
	Taokouling    bool
	Coolq         bool
	Wechat        bool
}

// go build -ldflags="-H windowsgui"
func main() {
	reg := &register{
		RegisterBtn: &walk.PushButton{},
		License:     &walk.TextEdit{},
	}
	mw := &walk.MainWindow{}
	if _, err := (MainWindow{
		Title:    "淘宝客工具激活码生成器",
		AssignTo: &mw,
		//Icon: "./assets/img/icon.ico",
		Size:   Size{400, 500},
		Layout: VBox{},
		DataBinder: DataBinder{
			DataSource: reg,
			AutoSubmit: true,
		},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					TextLabel{
						Text: "功能模块：",
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							Composite{
								Layout: HBox{},
								Children: []Widget{
									CheckBox{
										Text:    "大淘客",
										Checked: Bind("Selector.Dataoke"),
									},
									CheckBox{
										Text:    "好单库",
										Checked: Bind("Selector.Haodanku"),
									},
									CheckBox{
										Text:    "多多进宝",
										Checked: Bind("Selector.Duoduojinbao"),
									},
								},
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									CheckBox{
										Text:    "淘客易",
										Checked: Bind("Selector.Taokeyi"),
									},
									CheckBox{
										Text:    "淘宝销量搜索",
										Checked: Bind("Selector.GoodSearch"),
									},
									CheckBox{
										Text:    "易推客发单",
										Checked: Bind("Selector.Yituike"),
									},
								},
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									CheckBox{
										Text:    "拼多多拼团人数查询",
										Checked: Bind("Selector.PddUserNumber"),
									},
									CheckBox{
										Text:    "淘口令提取",
										Checked: Bind("Selector.Taokouling"),
									},
									CheckBox{
										Text:    "酷Q消息转发",
										Checked: Bind("Selector.Coolq"),
									},
								},
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									CheckBox{
										Text:    "微信转发",
										Checked: Bind("Selector.Wechat"),
									},
								},
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text: "激活版本",
							},
							ComboBox{
								Value: Bind("Version"),
								Model: []string{"v1.x", "v2.x", "v3.x"},
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text: "激活天数",
							},
							ComboBox{
								Value: Bind("Days"),
								Model: []string{"30", "60", "90", "120", "150", "180"},
							},
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							PushButton{
								Text:      "生成",
								AssignTo:  &reg.RegisterBtn,
								OnClicked: reg.Register,
							},
							TextEdit{
								AssignTo: &reg.License,
								VScroll:  true,
								ReadOnly: true,
							},
						},
					},
				},
			},
		},
	}).Run(); err != nil {
		fmt.Println(err)
	}
}

func getSelectors() map[string]bool {
	return map[string]bool{
		"dataoke":       false,
		"haodanku":      false,
		"duoduojinbao":  false,
		"taokeyi":       false,
		"goodSearch":    false,
		"yituike":       false,
		"pddUserNumber": false,
		"taokouling":    false,
		"coolq":         false,
	}
}

func (r *register) Register() {
	features := getFeatures(r.Selector)
	now := time.Now()
	days, _ := strconv.Atoi(r.Days)
	expireTimestamp := now.Add(time.Duration(days) * 24 * time.Hour).Unix()
	lic := &types.License{
		MainVersion:     r.Version,
		Feature:         features,
		ExpireTimestamp: expireTimestamp,
	}
	data := license.EncodeLicense(lic)
	r.License.SetText(data)
}

func getFeatures(selector Selector) []string {
	features := []string{}
	if selector.Dataoke {
		features = append(features, "dataoke")
	}
	if selector.Haodanku {
		features = append(features, "haodanku")
	}
	if selector.Duoduojinbao {
		features = append(features, "duoduojinbao")
	}
	if selector.Taokeyi {
		features = append(features, "taokeyi")
	}
	if selector.GoodSearch {
		features = append(features, "goodSearch")
	}
	if selector.Yituike {
		features = append(features, "yituike")
	}
	if selector.PddUserNumber {
		features = append(features, "pddUserNumber")
	}
	if selector.Taokouling {
		features = append(features, "taokouling")
	}
	if selector.Coolq {
		features = append(features, "coolq")
	}
	if selector.Wechat {
		features = append(features, "wechat")
	}
	return features
}
